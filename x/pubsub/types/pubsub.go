package types

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net/url"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	MaxDomainLength         = 256
	MaxUrlLength            = 512
	MaxCertLength           = 4096
	MaxSubscriptionIdLength = 128
	MaxAllowedSubscribers   = 256
)

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

	if err := ValidateProofUrl(p.ProofUrl, p.Domain, p.Address); err != nil {
		return fmt.Errorf("invalid proof URL: %s", err.Error())
	}

	if err := ValidateCaCertificateBase64(p.CaCert); err != nil {
		return fmt.Errorf("invalid CA certificate: %s", err.Error())
	}

	return nil
}

func (s *Subscriber) ValidateBasic() error {
	if err := ValidateAddress(s.Address); err != nil {
		return fmt.Errorf("invalid address: %s", err.Error())
	}

	// if a subcsriber does not provide a domain, ca_cert, and proof_url, they will not be able to
	// subscribe to any publisher intents that use the PUSH method

	if s.Domain != "" {
		if err := ValidateDomain(s.Domain); err != nil {
			return fmt.Errorf("invalid domain: %s", err.Error())
		}
	}

	if s.ProofUrl != "" {
		if err := ValidateProofUrl(s.ProofUrl, s.Domain, s.Address); err != nil {
			return fmt.Errorf("invalid proof URL: %s", err.Error())
		}
	}

	if s.CaCert != "" {
		if err := ValidateCaCertificateBase64(s.CaCert); err != nil {
			return fmt.Errorf("invalid CA certificate: %s", err.Error())
		}
	}

	return nil
}

func (pi *PublisherIntent) ValidateBasic() error {
	if err := ValidateSubscriptionId(pi.SubscriptionId); err != nil {
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

		if err := ValidateGenericUrl(pi.PullUrl); err != nil {
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
				return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress,
					fmt.Sprintf("allowed address entry %s is invalid: %s", allowedAddress, err.Error()))
			}
		}
	}

	return nil
}

func (si *SubscriberIntent) ValidateBasic() error {
	if err := ValidateSubscriptionId(si.SubscriptionId); err != nil {
		return fmt.Errorf("invalid subscription ID: %s", err.Error())
	}

	if err := ValidateAddress(si.SubscriberAddress); err != nil {
		return fmt.Errorf("invalid subscriber address: %s", err.Error())
	}

	if err := ValidateDomain(si.PublisherDomain); err != nil {
		return fmt.Errorf("invalid publisher domain: %s", err.Error())
	}

	// PushUrl is optional, but the SubscriberIntent will be rejected if the publisher intent is using
	// the PUSH method
	if si.PushUrl != "" {
		if err := ValidateGenericUrl(si.PushUrl); err != nil {
			return fmt.Errorf("invalid push URL: %s", err.Error())
		}
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
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("invalid address: %s", err.Error()))
	}

	return nil
}

func ValidateProofUrl(proofUrl string, domain string, address string) error {
	if proofUrl == "" {
		return fmt.Errorf("empty proof URL")
	}

	if err := ValidateGenericUrl(proofUrl); err != nil {
		return err
	}

	validProofUrl := fmt.Sprintf("https://%s/%s/cacert.pem", domain, address)
	if proofUrl != validProofUrl {
		return fmt.Errorf("invalid proof URL format, should be: %s", validProofUrl)
	}

	return nil
}

func ValidateGenericUrl(urlString string) error {
	if urlString == "" {
		return fmt.Errorf("empty URL")
	}

	if len(urlString) > MaxUrlLength {
		return fmt.Errorf("URL over max length of %d: %d", MaxUrlLength, len(urlString))
	}

	if _, err := url.Parse(urlString); err != nil {
		return fmt.Errorf("invalid URL format: %s", err.Error())
	}

	return nil
}

func ValidateCaCertificateBase64(certBase64 string) error {
	if certBase64 == "" {
		return fmt.Errorf("empty CA certificate")
	}

	if len(certBase64) > MaxCertLength {
		return fmt.Errorf("CA cert over max length of %d: %d", MaxCertLength, len(certBase64))
	}

	certBytes, err := base64.StdEncoding.DecodeString(certBase64)
	if err != nil {
		return fmt.Errorf("invalid base64 encoding for CA cert")
	}

	block, rest := pem.Decode(certBytes)
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

func ValidateSubscriptionId(subscriptionId string) error {
	if subscriptionId == "" {
		return fmt.Errorf("empty subscription ID")
	}

	if len(subscriptionId) > MaxSubscriptionIdLength {
		return fmt.Errorf("subscription ID over max length of %d: %d", MaxSubscriptionIdLength, len(subscriptionId))
	}

	// TODO(bolten): any other character limitations we should add here?
	if strings.Contains(subscriptionId, "|") {
		return fmt.Errorf("subscription IDs may not contain the pipe character '|'")
	}

	return nil
}
