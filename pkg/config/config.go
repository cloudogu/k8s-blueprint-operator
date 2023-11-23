package config

import (
	"fmt"
	"os"

	"github.com/Masterminds/semver/v3"
	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	StageDevelopment = "development"
	StageProduction  = "production"
	StageEnvVar      = "STAGE"
	namespaceEnvVar  = "NAMESPACE"
)

var log = ctrl.Log.WithName("config")

// OperatorConfig contains all configurable values for the dogu operator.
type OperatorConfig struct {
	// Version contains the current version of the operator
	Version *semver.Version
	// Namespace specifies the namespace that the operator is deployed to.
	Namespace string
}

var Stage = StageProduction

func IsStageDevelopment() bool {
	return Stage == StageDevelopment
}

// NewOperatorConfig creates a new operator config by reading values from the environment variables
func NewOperatorConfig(version string) (*OperatorConfig, error) {
	configureStage()

	parsedVersion, err := semver.NewVersion(version)
	if err != nil {
		return nil, fmt.Errorf("failed to parse version: %w", err)
	}
	log.Info(fmt.Sprintf("Version: [%s]", version))

	namespace, err := GetNamespace()
	if err != nil {
		return nil, fmt.Errorf("failed to read namespace: %w", err)
	}
	log.Info(fmt.Sprintf("Deploying the k8s dogu operator in namespace %s", namespace))

	return &OperatorConfig{
		Version:   parsedVersion,
		Namespace: namespace,
	}, nil
}

func configureStage() {
	var err error
	Stage, err = getEnvVar(StageEnvVar)
	if err != nil {
		log.Error(err, "Error reading stage environment variable. Use stage production")
	}

	if IsStageDevelopment() {
		log.Info("Starting in development mode! This is not recommended for production!")
	}
}

func GetNamespace() (string, error) {
	namespace, err := getEnvVar(namespaceEnvVar)
	if err != nil {
		return "", fmt.Errorf("failed to get env var [%s]: %w", namespaceEnvVar, err)
	}

	return namespace, nil
}

func getEnvVar(name string) (string, error) {
	env, found := os.LookupEnv(name)
	if !found {
		return "", fmt.Errorf("environment variable %s must be set", name)
	}
	return env, nil
}

// used for mocks

//nolint:unused
//goland:noinspection GoUnusedType
type logSink interface {
	logr.LogSink
}
