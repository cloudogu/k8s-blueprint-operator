package domain

type RequiredComponent struct {
	Name string `yaml:"name" json:"name"`
}

type HealthConfig struct {
	ComponentHealthConfig `yaml:"components" json:"components"`
}

type ComponentHealthConfig struct {
	RequiredComponents []RequiredComponent `yaml:"required" json:"required"`
}
