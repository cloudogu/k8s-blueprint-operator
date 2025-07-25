package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/Masterminds/semver/v3"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/cloudogu/cesapp-lib/core"
)

const (
	StageDevelopment = "development"
	StageProduction  = "production"
	StageEnvVar      = "STAGE"
	namespaceEnvVar  = "NAMESPACE"
	logLevelEnvVar   = "LOG_LEVEL"
)

const (
	doguRegistryEndpointEnvVar  = "DOGU_REGISTRY_ENDPOINT"
	doguRegistryUsernameEnvVar  = "DOGU_REGISTRY_USERNAME"
	doguRegistryPasswordEnvVar  = "DOGU_REGISTRY_PASSWORD"
	doguRegistryURLSchemaEnvVar = "DOGU_REGISTRY_URLSCHEMA"
)

const registryCacheDir = "/tmp/dogu-registry-cache"

var log = ctrl.Log.WithName("config")
var Stage = StageProduction

// OperatorConfig contains all configurable values for the dogu operator.
type OperatorConfig struct {
	// Version contains the current version of the operator
	Version *semver.Version
	// Namespace specifies the namespace that the operator is deployed to.
	Namespace string
}

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

func GetLogLevel() (string, error) {
	logLevel, err := getEnvVar(logLevelEnvVar)
	if err != nil {
		return "", fmt.Errorf("failed to get env var [%s]: %w", logLevelEnvVar, err)
	}

	return logLevel, nil
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

// GetRemoteConfiguration creates a remote configuration with the configured values.
func GetRemoteConfiguration() (*core.Remote, error) {
	// We can safely ignore this error since the url schema variable is optional.
	urlSchema, _ := getEnvVar(doguRegistryURLSchemaEnvVar)

	if urlSchema != "index" {
		log.Info("URLSchema is not index. Setting it to default.")
		urlSchema = "default"
	}

	endpoint, err := getEnvVar(doguRegistryEndpointEnvVar)
	if err != nil {
		return nil, err
	}

	if urlSchema == "default" {
		// trim suffix 'dogus' or 'dogus/' to provide maximum compatibility with the old remote configuration of the operator
		endpoint = strings.TrimSuffix(endpoint, "dogus/")
		endpoint = strings.TrimSuffix(endpoint, "dogus")
	}

	proxyURL, b := os.LookupEnv("PROXY_URL")

	proxySettings := core.ProxySettings{}
	if b && len(proxyURL) > 0 {
		var err error
		if proxySettings, err = configureProxySettings(proxyURL); err != nil {
			return nil, err
		}
	}

	return &core.Remote{
		Endpoint:      endpoint,
		CacheDir:      registryCacheDir,
		URLSchema:     urlSchema,
		ProxySettings: proxySettings,
	}, nil
}

func configureProxySettings(proxyURL string) (core.ProxySettings, error) {
	parsedURL, err := url.Parse(proxyURL)
	if err != nil {
		return core.ProxySettings{}, fmt.Errorf("invalid proxy url: %w", err)
	}

	proxySettings := core.ProxySettings{}
	proxySettings.Enabled = true
	if parsedURL.User != nil {
		proxySettings.Username = parsedURL.User.Username()
		if password, set := parsedURL.User.Password(); set {
			proxySettings.Password = password
		}
	}

	proxySettings.Server = parsedURL.Hostname()

	port, err := strconv.Atoi(parsedURL.Port())
	if err != nil {
		return core.ProxySettings{}, fmt.Errorf("invalid port %s: %w", parsedURL.Port(), err)
	}
	proxySettings.Port = port

	return proxySettings, nil
}

// GetRemoteCredentials creates a remote credential pair with the configured values.
func GetRemoteCredentials() (*core.Credentials, error) {
	username, err := getEnvVar(doguRegistryUsernameEnvVar)
	if err != nil {
		return nil, err
	}

	password, err := getEnvVar(doguRegistryPasswordEnvVar)
	if err != nil {
		return nil, err
	}

	return &core.Credentials{
		Username: username,
		Password: password,
	}, nil
}
