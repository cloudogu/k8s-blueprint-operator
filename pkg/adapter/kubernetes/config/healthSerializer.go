package config

import (
	"encoding/json"
	"fmt"
	"time"
)

type healthConfig struct {
	// Components contains configuration concerning the health-check of components.
	Components componentHealthConfig `yaml:"components,omitempty" json:"components,omitempty"`
	// Wait contains configuration concerning how the stand-by-period for the ecosystem to become healthy.
	Wait waitHealthConfig `yaml:"wait,omitempty" json:"wait,omitempty"`
}

type componentHealthConfig struct {
	// Required is a list of components that have to be installed for the health-check to succeed.
	Required []requiredComponent `yaml:"required,omitempty" json:"required,omitempty"`
}

type requiredComponent struct {
	// Name identifies the component.
	Name string `yaml:"name" json:"name"`
}

type waitHealthConfig struct {
	// Timeout is the maximum time to wait for the ecosystem to become healthy.
	Timeout duration `yaml:"timeout,omitempty" json:"timeout,omitempty"`
	// Interval is the time to wait between health checks.
	Interval duration `yaml:"interval,omitempty" json:"interval,omitempty"`
}

// duration is a wrapper for time.Duration so that it can be marshalled and unmarshalled to JSON
// using time.Duration's string formatting.
type duration struct {
	time.Duration
}

// UnmarshalJSON expects either a JSON string formatted to be parseable by time.ParseDuration or
// a JSON integer number that represents the duration in nanoseconds.
func (d *duration) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' {
		sd := string(b[1 : len(b)-1])
		d.Duration, err = time.ParseDuration(sd)
		return err
	}

	var id int64
	id, err = json.Number(b).Int64()
	d.Duration = time.Duration(id)

	return err
}

// MarshalJSON returns the duration as a JSON string formatted with time.Duration's string formatting
func (d *duration) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, d.String())), nil
}
