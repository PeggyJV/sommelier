package types

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/url"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	MaxDomainLength         = 256
	MaxURLLength            = 512
	MaxCertLength           = 4096
	MaxSubscriptionIDLength = 128
	MaxAllowedSubscribers   = 256
)

func StringHash(inputString string) []byte {
	return crypto.Keccak256Hash([]byte(inputString)).Bytes()
}

///////////////////
// ValidateBasic //
///////////////////

func (p *Publisher) ValidateBasic() error {
	if err := ValidateDomain(p.Domain); err != nil {
		return fmt.Errorf("invalid domain: %s", err.Error())
	}

	if err := ValidateAddress(p.Address); err != nil {
		return fmt.Errorf("invalid address: %s", err.Error())
	}

	if err := ValidateCaCertificate(p.CaCert); err != nil {
		return fmt.Errorf("invalid CA certificate: %s", err.Error())
	}

	return nil
}

func (s *Subscriber) ValidateBasic() error {
	if err := ValidateAddress(s.Address); err != nil {
		return fmt.Errorf("invalid address: %s", err.Error())
	}

	// if a subcsriber does not provide a CA cert and push URL, they will not be able to
	// subscribe to any publisher intents that use the PUSH method

	if s.CaCert != "" {
		if err := ValidateCaCertificate(s.CaCert); err != nil {
			return fmt.Errorf("invalid CA certificate: %s", err.Error())
		}
	}

	// PushUrl is optional, but SubscriberIntents will be rejected if the PublisherIntent is using
	// the PUSH method and this is missing
	if s.PushUrl != "" {
		if err := ValidateGenericURL(s.PushUrl); err != nil {
			return fmt.Errorf("invalid push URL: %s", err.Error())
		}
	}

	return nil
}

func (pi *PublisherIntent) ValidateBasic() error {
	if err := ValidateSubscriptionID(pi.SubscriptionId); err != nil {
		return fmt.Errorf("invalid subscription ID: %s", err.Error())
	}

	if err := ValidateDomain(pi.PublisherDomain); err != nil {
		return fmt.Errorf("invalid publisher domain: %s", err.Error())
	}

	if pi.Method != PublishMethod_PULL &&
		pi.Method != PublishMethod_PUSH {
		return fmt.Errorf("invalid enum value for method: %d", pi.Method)
	}

	if pi.Method == PublishMethod_PULL {
		if pi.PullUrl == "" {
			return fmt.Errorf("empty pull URL when it is required by the PULL method")
		}

		if err := ValidateGenericURL(pi.PullUrl); err != nil {
			return fmt.Errorf("invalid pull URL: %s", err.Error())
		}
	}

	if pi.AllowedSubscribers != AllowedSubscribers_ANY &&
		pi.AllowedSubscribers != AllowedSubscribers_VALIDATORS &&
		pi.AllowedSubscribers != AllowedSubscribers_LIST {
		return fmt.Errorf("invalid enum value for allowed subscribers: %d", pi.AllowedSubscribers)
	}

	if pi.AllowedSubscribers == AllowedSubscribers_LIST {
		if len(pi.AllowedAddresses) == 0 {
			return fmt.Errorf("empty list of allowed addresses when it is required by the allowed subscribers LIST value")
		}

		// TODO(bolten): currently set to 256, how many addresses should we allow in this list given that it is heavy on
		// storage and must be iterated over?
		if len(pi.AllowedAddresses) > MaxAllowedSubscribers {
			return fmt.Errorf("allowed address list over maximum length of %d: %d", MaxAllowedSubscribers, len(pi.AllowedAddresses))
		}

		for _, allowedAddress := range pi.AllowedAddresses {
			if err := ValidateAddress(allowedAddress); err != nil {
				return errorsmod.Wrap(sdkerrors.ErrInvalidAddress,
					fmt.Sprintf("allowed address entry %s is invalid: %s", allowedAddress, err.Error()))
			}
		}
	}

	return nil
}

func (si *SubscriberIntent) ValidateBasic() error {
	if err := ValidateSubscriptionID(si.SubscriptionId); err != nil {
		return fmt.Errorf("invalid subscription ID: %s", err.Error())
	}

	if err := ValidateAddress(si.SubscriberAddress); err != nil {
		return fmt.Errorf("invalid subscriber address: %s", err.Error())
	}

	if err := ValidateDomain(si.PublisherDomain); err != nil {
		return fmt.Errorf("invalid publisher domain: %s", err.Error())
	}

	return nil
}

func (ds *DefaultSubscription) ValidateBasic() error {
	if err := ValidateSubscriptionID(ds.SubscriptionId); err != nil {
		return fmt.Errorf("invalid subscription ID: %s", err.Error())
	}

	if err := ValidateDomain(ds.PublisherDomain); err != nil {
		return fmt.Errorf("invalid publisher domain: %s", err.Error())
	}

	return nil
}

///////////////////////
// Field Validations //
///////////////////////

func ValidateDomain(domain string) error {
	if domain == "" {
		return fmt.Errorf("empty domain")
	}

	if len(domain) > MaxDomainLength {
		return fmt.Errorf("domain over max length of %d: %d", MaxDomainLength, len(domain))
	}

	// TODO(bolten): perhaps we should add a regex to ensure a limited character set for the domain
	if _, err := url.Parse(fmt.Sprintf("https://%s", domain)); err != nil {
		return fmt.Errorf("invalid URL format: %s", err.Error())
	}

	return nil
}

func ValidateAddress(address string) error {
	if address == "" {
		return fmt.Errorf("empty address")
	}

	if _, err := sdk.AccAddressFromBech32(address); err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err.Error()))
	}

	return nil
}

func ValidateCaCertificate(certPem string) error {
	if certPem == "" {
		return fmt.Errorf("empty CA certificate")
	}

	if len(certPem) > MaxCertLength {
		return fmt.Errorf("CA cert over max length of %d: %d", MaxCertLength, len(certPem))
	}

	block, rest := pem.Decode([]byte(certPem))
	if block == nil || len(rest) > 0 {
		return fmt.Errorf("invalid PEM formating, expecting only certificate block")
	}

	if block.Type != "CERTIFICATE" {
		return fmt.Errorf("invalid PEM certificate block, must be CERTIFICATE")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("could not parse x509 certificate")
	}

	if !cert.IsCA {
		return fmt.Errorf("must be a CA cert")
	}

	if cert.PublicKeyAlgorithm != x509.ECDSA {
		return fmt.Errorf("must be an ECDSA cert")
	}

	return nil
}

func ValidateSubscriptionID(subscriptionID string) error {
	if subscriptionID == "" {
		return fmt.Errorf("empty subscription ID")
	}

	if len(subscriptionID) > MaxSubscriptionIDLength {
		return fmt.Errorf("subscription ID over max length of %d: %d", MaxSubscriptionIDLength, len(subscriptionID))
	}

	return nil
}
