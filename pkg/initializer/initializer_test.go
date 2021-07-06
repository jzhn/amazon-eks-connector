package initializer

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"

	"golang.a2z.com/eks-connector/pkg/agent"
	"golang.a2z.com/eks-connector/pkg/state"
)

const (
	testFingerPrint           = "7eb32fa4-5b75-4431-866e-c2af92f5440b"
	testInstanceID            = "9e629669-d5f6-47e6-97e0-3f5c57c01d82"
	testPrivateKey            = "greedisgood"
	testPrivateKeyType        = "war3"
	testPrivateKeyCreatedDate = "2021-07-30 00:00:00.999999999 -0700 PDT"
	testRegion                = "mars-northeast-2"
)

func TestInitializerSuite(t *testing.T) {
	suite.Run(t, new(InitializerSuite))
}

type InitializerSuite struct {
	suite.Suite
	secretPersistence *state.MockPersistence
	fsPersistence     *state.MockPersistence
	registration      *agent.MockRegistration

	initializer Initializer
}

func (suite *InitializerSuite) SetupTest() {
	suite.secretPersistence = &state.MockPersistence{}
	suite.fsPersistence = &state.MockPersistence{}
	suite.registration = &agent.MockRegistration{}
	suite.initializer = NewInitializer(suite.secretPersistence, suite.fsPersistence, suite.registration)
}

func (suite *InitializerSuite) TestInitializeNoSavedStateHappyCase() {
	// prepare
	state := NewTestState()
	serializedState, err := state.Serialize()
	suite.NoError(err)
	suite.secretPersistence.On("Load").Return(nil, nil)
	suite.registration.On("Register").Return(state, nil)
	suite.secretPersistence.On("Save", serializedState).Return(nil)
	suite.fsPersistence.On("Save", serializedState).Return(nil)

	// test
	actualErr := suite.initializer.Initialize()

	// verify
	suite.NoError(actualErr)
	suite.secretPersistence.AssertExpectations(suite.T())
	suite.fsPersistence.AssertExpectations(suite.T())
	suite.registration.AssertExpectations(suite.T())
}

func (suite *InitializerSuite) TestInitializeNoSavedStateFailedRegistration() {
	// prepare
	err := errors.New("failed registration")
	suite.secretPersistence.On("Load").Return(nil, nil)
	suite.registration.On("Register").Return(nil, err)

	// test
	actualErr := suite.initializer.Initialize()

	// verify
	suite.ErrorIs(actualErr, err)
	suite.secretPersistence.AssertExpectations(suite.T())
	suite.fsPersistence.AssertExpectations(suite.T())
	suite.registration.AssertExpectations(suite.T())
}

func (suite *InitializerSuite) TestInitializeNoSavedStateFailedSecretPersistence() {
	// prepare
	state := NewTestState()
	serializedState, err := state.Serialize()
	suite.NoError(err)
	err = errors.New("failed persistence")
	suite.secretPersistence.On("Load").Return(nil, nil)
	suite.registration.On("Register").Return(state, nil)
	suite.secretPersistence.On("Save", serializedState).Return(err)

	// test
	actualErr := suite.initializer.Initialize()

	// verify
	suite.ErrorIs(actualErr, err)
	suite.secretPersistence.AssertExpectations(suite.T())
	suite.fsPersistence.AssertExpectations(suite.T())
	suite.registration.AssertExpectations(suite.T())
}

func (suite *InitializerSuite) TestInitializeSavedStateHappyCase() {
	// prepare
	state := NewTestState()
	serializedState, err := state.Serialize()
	suite.NoError(err)
	suite.secretPersistence.On("Load").Return(serializedState, nil)
	suite.fsPersistence.On("Save", serializedState).Return(nil)

	// test
	actualErr := suite.initializer.Initialize()

	// verify
	suite.NoError(actualErr)
	suite.secretPersistence.AssertExpectations(suite.T())
	suite.fsPersistence.AssertExpectations(suite.T())
	suite.registration.AssertExpectations(suite.T())
}

func (suite *InitializerSuite) TestInitializeSavedStateFailedFSPersistence() {
	// prepare
	state := NewTestState()
	serializedState, err := state.Serialize()
	suite.NoError(err)
	err = errors.New("failed to persist")
	suite.secretPersistence.On("Load").Return(serializedState, nil)
	suite.fsPersistence.On("Save", serializedState).Return(err)

	// test
	actualErr := suite.initializer.Initialize()

	// verify
	suite.ErrorIs(actualErr, err)
	suite.secretPersistence.AssertExpectations(suite.T())
	suite.fsPersistence.AssertExpectations(suite.T())
	suite.registration.AssertExpectations(suite.T())
}

func NewTestState() *state.State {
	return &state.State{
		FingerPrint:           testFingerPrint,
		InstanceID:            testInstanceID,
		PrivateKey:            testPrivateKey,
		PrivateKeyType:        testPrivateKeyType,
		PrivateKeyCreatedDate: testPrivateKeyCreatedDate,
		Region:                testRegion,
	}
}
