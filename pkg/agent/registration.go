// Package agent is provides types and functions to perform registration with SSM
package agent

import (
	"time"

	"k8s.io/klog/v2"

	"golang.a2z.com/eks-connector/pkg/config"
	"golang.a2z.com/eks-connector/pkg/ssm"
	"golang.a2z.com/eks-connector/pkg/state"
)

type Registration interface {
	Register() (*state.State, error)
}

type ssmRegistration struct {
	ssm              ssm.Client
	activationConfig *config.ActivationConfig
}

func NewRegistration(ssmService ssm.Client, activationConfig *config.ActivationConfig) Registration {
	return &ssmRegistration{
		ssm:              ssmService,
		activationConfig: activationConfig,
	}
}

func (r *ssmRegistration) Register() (*state.State, error) {
	state := &state.State{}

	klog.Infof("creating %s keypair...", KeyType)
	keyPair, err := createKeypair()
	if err != nil {
		return nil, err
	}

	klog.Infof("encoding %s keypair...", KeyType)
	privateKey := keyPair.encodePrivateKey()

	publicKey, err := keyPair.encodePublicKey()
	if err != nil {
		return nil, err
	}

	klog.Infof("generating fingerprint...")
	fingerPrint, err := createFingerPrint()
	if err != nil {
		return nil, err
	}
	klog.Infof("fingerprint %s generated", fingerPrint)

	klog.Infof("registering at SSM with activationId %s...", r.activationConfig.ID)
	instanceID, err := r.ssm.RegisterManagedInstance(
		r.activationConfig.ID,
		r.activationConfig.Code,
		publicKey,
		KeyType,
		fingerPrint,
	)
	if err != nil {
		return nil, err
	}
	klog.Infof("successfully registered at SSM with InstanceID %s", instanceID)

	state.PrivateKey = privateKey
	state.PrivateKeyType = KeyType
	state.PrivateKeyCreatedDate = time.Now().String()

	state.FingerPrint = fingerPrint
	state.InstanceID = instanceID
	state.Region = r.ssm.Region()

	return state, nil
}
