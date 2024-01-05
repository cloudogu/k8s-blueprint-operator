package maintenance

import (
	"encoding/json"
	"fmt"
)

type defaultLock struct {
	globalConfig globalConfig
}

func (l *defaultLock) isActiveAndOurs() (bool, bool, error) {
	isActive, err := l.globalConfig.Exists(registryKeyMaintenance)
	if err != nil {
		return false, false,
			fmt.Errorf("failed to check if maintenance mode registry key exists: %w", err)
	}

	if !isActive {
		return false, false, nil
	}

	isOurs, err := l.isOurs()
	if err != nil {
		return false, false, err
	}

	return true, isOurs, nil
}

func (l *defaultLock) isOurs() (bool, error) {
	rawValue, err := l.globalConfig.Get(registryKeyMaintenance)
	if err != nil {
		return false, fmt.Errorf("failed to get maintenance mode from configuration registry: %w", err)
	}

	var value maintenanceRegistryObject
	err = json.Unmarshal([]byte(rawValue), &value)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal json of maintenance mode object: %w", err)
	}

	return value.Holder == blueprintOperatorHolder, nil
}
