package config

import (
	"encoding/json"
	"fmt"
	"time"
)

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
