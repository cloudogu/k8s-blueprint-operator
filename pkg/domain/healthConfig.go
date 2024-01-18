package domain

type HealthConfig struct {
	ComponentHealthConfig `yaml:"components" json:"components"`
}

type ComponentHealthConfig struct {
	RequiredComponents []string `yaml:"required" json:"required"`
}
