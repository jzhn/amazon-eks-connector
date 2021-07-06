// Package state provides types and functions for eks connector state management
package state

import "encoding/json"

type State struct {
	FingerPrint           string
	InstanceID            string
	PrivateKey            string
	PrivateKeyType        string
	PrivateKeyCreatedDate string
	Region                string
}

func (state *State) Serialize() (serializedState SerializedState, err error) {
	serializedState = SerializedState{}

	serializedState[FileManifest], err = state.serializeManifest()
	if err != nil {
		return
	}

	serializedState[FileRegistrationKey], err = state.serializeRegistrationKey()
	if err != nil {
		return
	}

	serializedState[FileInstanceFingerprint], err = state.serializeInstanceFingerprint()
	if err != nil {
		return
	}

	return serializedState, nil
}

func (state *State) serializeManifest() (string, error) {
	// Manifest is a static file.
	manifest := &manifestState{
		InstanceFingerprint: "/var/lib/amazon/ssm/Vault/Store/InstanceFingerprint",
		RegistrationKey:     "/var/lib/amazon/ssm/Vault/Store/RegistrationKey",
	}
	return marshal(manifest)
}

func (state *State) serializeInstanceFingerprint() (string, error) {
	fingerprint := &instanceFingerprintState{}
	fingerprint.Fingerprint = state.FingerPrint
	fingerprint.HardwareHash = make(map[string]string)
	fingerprint.SimilarityThreshold = -1

	return marshal(fingerprint)
}

func (state *State) serializeRegistrationKey() (string, error) {
	registrationKey := &registrationKeyState{}

	registrationKey.PrivateKey = state.PrivateKey
	registrationKey.PrivateKeyType = state.PrivateKeyType
	registrationKey.Region = state.Region
	registrationKey.PrivateKeyCreatedDate = state.PrivateKeyCreatedDate
	registrationKey.InstanceID = state.InstanceID
	registrationKey.AvailabilityZone = ""
	registrationKey.InstanceType = ""

	return marshal(registrationKey)
}

func marshal(v interface{}) (string, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return "", nil
	}
	return string(data), nil
}
