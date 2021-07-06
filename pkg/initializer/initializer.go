// Package initializer contains init container related functionalities.
package initializer

import (
	"k8s.io/klog/v2"

	"golang.a2z.com/eks-connector/pkg/agent"
	"golang.a2z.com/eks-connector/pkg/state"
)

type Initializer interface {
	Initialize() error
}

func NewInitializer(secretPersistence state.Persistence,
	fsPersistence state.Persistence,
	registration agent.Registration) Initializer {
	return &ssmInitializer{
		secretPersistence: secretPersistence,
		fsPersistence:     fsPersistence,
		registration:      registration,
	}
}

type ssmInitializer struct {
	secretPersistence state.Persistence
	fsPersistence     state.Persistence
	registration      agent.Registration
}

func (i *ssmInitializer) Initialize() error {
	klog.Infof("eks-connector initializer starts...")

	klog.Infof("loading persisted state from secrets...")
	serializedSecret, err := i.secretPersistence.Load()
	if err != nil {
		return err
	}

	if serializedSecret != nil {
		klog.Infof("state information is available at eks-connector secret state store, skipping registration")
	} else {
		klog.Infof("state information is not available at eks-connector secret state store, registering as new instance")
		connectorState, err := i.registration.Register()
		if err != nil {
			return err
		}

		klog.Infof("serializing state information...")
		serializedSecret, err = connectorState.Serialize()
		if err != nil {
			return err
		}

		klog.Infof("persisting state information to secrets...")
		err = i.secretPersistence.Save(serializedSecret)
		if err != nil {
			return err
		}
	}

	klog.Infof("persisting state information to filesystem...")
	err = i.fsPersistence.Save(serializedSecret)

	return err
}
