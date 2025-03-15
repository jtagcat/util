package webauthn

import (
	"fmt"
	"strings"
)

// Code generated by cdproto-gen. DO NOT EDIT.

// AuthenticatorID [no description].
//
// See: https://chromedevtools.github.io/devtools-protocol/tot/WebAuthn#type-AuthenticatorId
type AuthenticatorID string

// String returns the AuthenticatorID as string value.
func (t AuthenticatorID) String() string {
	return string(t)
}

// AuthenticatorProtocol [no description].
//
// See: https://chromedevtools.github.io/devtools-protocol/tot/WebAuthn#type-AuthenticatorProtocol
type AuthenticatorProtocol string

// String returns the AuthenticatorProtocol as string value.
func (t AuthenticatorProtocol) String() string {
	return string(t)
}

// AuthenticatorProtocol values.
const (
	AuthenticatorProtocolU2f   AuthenticatorProtocol = "u2f"
	AuthenticatorProtocolCtap2 AuthenticatorProtocol = "ctap2"
)

// UnmarshalJSON satisfies [json.Unmarshaler].
func (t *AuthenticatorProtocol) UnmarshalJSON(buf []byte) error {
	s := string(buf)
	s = strings.TrimSuffix(strings.TrimPrefix(s, `"`), `"`)

	switch AuthenticatorProtocol(s) {
	case AuthenticatorProtocolU2f:
		*t = AuthenticatorProtocolU2f
	case AuthenticatorProtocolCtap2:
		*t = AuthenticatorProtocolCtap2
	default:
		return fmt.Errorf("unknown AuthenticatorProtocol value: %v", s)
	}
	return nil
}

// Ctap2version [no description].
//
// See: https://chromedevtools.github.io/devtools-protocol/tot/WebAuthn#type-Ctap2Version
type Ctap2version string

// String returns the Ctap2version as string value.
func (t Ctap2version) String() string {
	return string(t)
}

// Ctap2version values.
const (
	Ctap2versionCtap20 Ctap2version = "ctap2_0"
	Ctap2versionCtap21 Ctap2version = "ctap2_1"
)

// UnmarshalJSON satisfies [json.Unmarshaler].
func (t *Ctap2version) UnmarshalJSON(buf []byte) error {
	s := string(buf)
	s = strings.TrimSuffix(strings.TrimPrefix(s, `"`), `"`)

	switch Ctap2version(s) {
	case Ctap2versionCtap20:
		*t = Ctap2versionCtap20
	case Ctap2versionCtap21:
		*t = Ctap2versionCtap21
	default:
		return fmt.Errorf("unknown Ctap2version value: %v", s)
	}
	return nil
}

// AuthenticatorTransport [no description].
//
// See: https://chromedevtools.github.io/devtools-protocol/tot/WebAuthn#type-AuthenticatorTransport
type AuthenticatorTransport string

// String returns the AuthenticatorTransport as string value.
func (t AuthenticatorTransport) String() string {
	return string(t)
}

// AuthenticatorTransport values.
const (
	AuthenticatorTransportUsb      AuthenticatorTransport = "usb"
	AuthenticatorTransportNfc      AuthenticatorTransport = "nfc"
	AuthenticatorTransportBle      AuthenticatorTransport = "ble"
	AuthenticatorTransportCable    AuthenticatorTransport = "cable"
	AuthenticatorTransportInternal AuthenticatorTransport = "internal"
)

// UnmarshalJSON satisfies [json.Unmarshaler].
func (t *AuthenticatorTransport) UnmarshalJSON(buf []byte) error {
	s := string(buf)
	s = strings.TrimSuffix(strings.TrimPrefix(s, `"`), `"`)

	switch AuthenticatorTransport(s) {
	case AuthenticatorTransportUsb:
		*t = AuthenticatorTransportUsb
	case AuthenticatorTransportNfc:
		*t = AuthenticatorTransportNfc
	case AuthenticatorTransportBle:
		*t = AuthenticatorTransportBle
	case AuthenticatorTransportCable:
		*t = AuthenticatorTransportCable
	case AuthenticatorTransportInternal:
		*t = AuthenticatorTransportInternal
	default:
		return fmt.Errorf("unknown AuthenticatorTransport value: %v", s)
	}
	return nil
}

// VirtualAuthenticatorOptions [no description].
//
// See: https://chromedevtools.github.io/devtools-protocol/tot/WebAuthn#type-VirtualAuthenticatorOptions
type VirtualAuthenticatorOptions struct {
	Protocol                    AuthenticatorProtocol  `json:"protocol"`
	Ctap2version                Ctap2version           `json:"ctap2Version,omitempty,omitzero"` // Defaults to ctap2_0. Ignored if |protocol| == u2f.
	Transport                   AuthenticatorTransport `json:"transport"`
	HasResidentKey              bool                   `json:"hasResidentKey,omitempty,omitzero"`              // Defaults to false.
	HasUserVerification         bool                   `json:"hasUserVerification,omitempty,omitzero"`         // Defaults to false.
	HasLargeBlob                bool                   `json:"hasLargeBlob,omitempty,omitzero"`                // If set to true, the authenticator will support the largeBlob extension. https://w3c.github.io/webauthn#largeBlob Defaults to false.
	HasCredBlob                 bool                   `json:"hasCredBlob,omitempty,omitzero"`                 // If set to true, the authenticator will support the credBlob extension. https://fidoalliance.org/specs/fido-v2.1-rd-20201208/fido-client-to-authenticator-protocol-v2.1-rd-20201208.html#sctn-credBlob-extension Defaults to false.
	HasMinPinLength             bool                   `json:"hasMinPinLength,omitempty,omitzero"`             // If set to true, the authenticator will support the minPinLength extension. https://fidoalliance.org/specs/fido-v2.1-ps-20210615/fido-client-to-authenticator-protocol-v2.1-ps-20210615.html#sctn-minpinlength-extension Defaults to false.
	HasPrf                      bool                   `json:"hasPrf,omitempty,omitzero"`                      // If set to true, the authenticator will support the prf extension. https://w3c.github.io/webauthn/#prf-extension Defaults to false.
	AutomaticPresenceSimulation bool                   `json:"automaticPresenceSimulation,omitempty,omitzero"` // If set to true, tests of user presence will succeed immediately. Otherwise, they will not be resolved. Defaults to true.
	IsUserVerified              bool                   `json:"isUserVerified,omitempty,omitzero"`              // Sets whether User Verification succeeds or fails for an authenticator. Defaults to false.
	DefaultBackupEligibility    bool                   `json:"defaultBackupEligibility,omitempty,omitzero"`    // Credentials created by this authenticator will have the backup eligibility (BE) flag set to this value. Defaults to false. https://w3c.github.io/webauthn/#sctn-credential-backup
	DefaultBackupState          bool                   `json:"defaultBackupState,omitempty,omitzero"`          // Credentials created by this authenticator will have the backup state (BS) flag set to this value. Defaults to false. https://w3c.github.io/webauthn/#sctn-credential-backup
}

// Credential [no description].
//
// See: https://chromedevtools.github.io/devtools-protocol/tot/WebAuthn#type-Credential
type Credential struct {
	CredentialID         string `json:"credentialId"`
	IsResidentCredential bool   `json:"isResidentCredential"`
	RpID                 string `json:"rpId,omitempty,omitzero"`              // Relying Party ID the credential is scoped to. Must be set when adding a credential.
	PrivateKey           string `json:"privateKey"`                           // The ECDSA P-256 private key in PKCS#8 format.
	UserHandle           string `json:"userHandle,omitempty,omitzero"`        // An opaque byte sequence with a maximum size of 64 bytes mapping the credential to a specific user.
	SignCount            int64  `json:"signCount"`                            // Signature counter. This is incremented by one for each successful assertion. See https://w3c.github.io/webauthn/#signature-counter
	LargeBlob            string `json:"largeBlob,omitempty,omitzero"`         // The large blob associated with the credential. See https://w3c.github.io/webauthn/#sctn-large-blob-extension
	BackupEligibility    bool   `json:"backupEligibility,omitempty,omitzero"` // Assertions returned by this credential will have the backup eligibility (BE) flag set to this value. Defaults to the authenticator's defaultBackupEligibility value.
	BackupState          bool   `json:"backupState,omitempty,omitzero"`       // Assertions returned by this credential will have the backup state (BS) flag set to this value. Defaults to the authenticator's defaultBackupState value.
	UserName             string `json:"userName,omitempty,omitzero"`          // The credential's user.name property. Equivalent to empty if not set. https://w3c.github.io/webauthn/#dom-publickeycredentialentity-name
	UserDisplayName      string `json:"userDisplayName,omitempty,omitzero"`   // The credential's user.displayName property. Equivalent to empty if not set. https://w3c.github.io/webauthn/#dom-publickeycredentialuserentity-displayname
}
