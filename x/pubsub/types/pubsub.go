package types

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net/url"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (p *Publisher) ValidateBasic() error {
	if p.Domain == "" {
		return fmt.Errorf("empty domain")
	}

	if err := ValidateDomain(p.Domain); err != nil {
		return fmt.Errorf("invalid domain: %s", err.Error())
	}

	if p.CaCert == "" {
		return fmt.Errorf("empty CA cert")
	}

	if err := ValidateCaCertificateBase64(p.CaCert); err != nil {
		return fmt.Errorf("invalid CA certificate: %s", err.Error())
	}

	if p.Address == "" {
		return fmt.Errorf("empty address")
	}

	if _, err := sdk.AccAddressFromBech32(p.Address); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress,
			fmt.Sprintf("address %s is invalid: %s", p.Address, err.Error()))
	}

	if p.ProofUrl == "" {
		return fmt.Errorf("empty proof URL")
	}

	if err := ValidateProofUrl(p.ProofUrl, p.Domain, p.Address); err != nil {
		return fmt.Errorf("invalid proof URL: %s", err.Error())
	}

	return nil
}

func (s *Subscriber) ValidateBasic() error {
	if s.Address == "" {
		return fmt.Errorf("empty address")
	}

	if _, err := sdk.AccAddressFromBech32(s.Address); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress,
			fmt.Sprintf("address %s is invalid: %s", s.Address, err.Error()))
	}

	// if a subcsriber does not provide a domain, ca_cert, and proof_url, they will not be able to
	// subscribe to any publisher intents that use the PUSH method

	if s.Domain != "" {
		if err := ValidateDomain(s.Domain); err != nil {
			return fmt.Errorf("invalid domain: %s", err.Error())
		}
	}

	if s.CaCert != "" {
		if err := ValidateCaCertificateBase64(s.CaCert); err != nil {
			return fmt.Errorf("invalid CA certificate: %s", err.Error())
		}
	}

	if s.ProofUrl != "" {
		if err := ValidateProofUrl(s.ProofUrl, s.Domain, s.Address); err != nil {
			return fmt.Errorf("invalid proof URL: %s", err.Error())
		}
	}

	return nil
}

func (pi *PublisherIntent) ValidateBasic() error {
	if pi.SubscriptionId == "" {
		return fmt.Errorf("empty subscription ID")
	}

	if err := ValidateSubscriptionId(pi.SubscriptionId); err != nil {
		return fmt.Errorf("invalid subscription ID: %s", err.Error())
	}

	if pi.PublisherDomain == "" {
		return fmt.Errorf("empty publisher domain")
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

		for _, allowedAddress := range pi.AllowedAddresses {
			if _, err := sdk.AccAddressFromBech32(allowedAddress); err != nil {
				return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress,
					fmt.Sprintf("allowed address entry %s is invalid: %s", allowedAddress, err.Error()))
			}
		}
	}

	return nil
}

func (si *SubscriberIntent) ValidateBasic() error {
	if si.SubscriptionId == "" {
		return fmt.Errorf("empty subscription ID")
	}

	if err := ValidateSubscriptionId(si.SubscriptionId); err != nil {
		return fmt.Errorf("invalid subscription ID: %s", err.Error())
	}

	if si.SubscriberAddress == "" {
		return fmt.Errorf("empty subscriber address")
	}

	if _, err := sdk.AccAddressFromBech32(si.SubscriberAddress); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress,
			fmt.Sprintf("subscriber address %s is invalid: %s", si.SubscriberAddress, err.Error()))
	}

	if si.PublisherDomain == "" {
		return fmt.Errorf("empty publisher domain")
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

func ValidateDomain(domain string) error {
	if len(domain) > 256 {
		return fmt.Errorf("domain over max length of 256: %d", len(domain))
	}

	if _, err := url.Parse(fmt.Sprintf("https://%s", domain)); err != nil {
		return fmt.Errorf("invalid URL format: %s", err.Error())
	}

	return nil
}

func ValidateProofUrl(proofUrl string, domain string, address string) error {
	if err := ValidateGenericUrl(proofUrl); err != nil {
		return err
	}

	validProofUrl := fmt.Sprintf("https://%s/%s/cacert.pem", domain, address)
	if proofUrl != validProofUrl {
		return fmt.Errorf("Invalid proof URL format, should be: %s", validProofUrl)
	}

	return nil
}

func ValidateGenericUrl(urlString string) error {
	if len(urlString) > 512 {
		return fmt.Errorf("URL over max length of 512: %d", len(urlString))
	}

	if _, err := url.Parse(urlString); err != nil {
		return fmt.Errorf("invalid URL format: %s", err.Error())
	}

	return nil
}

func ValidateCaCertificateBase64(certBase64 string) error {
	if len(certBase64) > 4096 {
		return fmt.Errorf("CA cert over max length of 4096: %d", len(certBase64))
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
	if len(subscriptionId) > 128 {
		return fmt.Errorf("subscription ID over max length of 128: %d", len(subscriptionId))
	}

	return nil
}
