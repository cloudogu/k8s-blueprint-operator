package maintenance

import (
	"encoding/json"
	"fmt"

	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
)

const registryKeyMaintenance = "maintenance"
const blueprintOperatorHolder = "k8s-blueprint-operator"

type defaultSwitcher struct {
	globalConfig globalConfig
}

// activate enables the maintenance mode.
func (m *defaultSwitcher) activate(content domainservice.MaintenancePageModel) error {
	value := maintenanceRegistryObject{
		Title:  content.Title,
		Text:   content.Text,
		Holder: blueprintOperatorHolder,
	}

	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to serialize maintenance mode object: %w", err)
	}

	err = m.globalConfig.Set(registryKeyMaintenance, string(jsonBytes))
	if err != nil {
		return fmt.Errorf("failed to set maintenance mode registry key: %w", err)
	}

	return nil
}

// Deactivate disables the maintenance mode.
func (m *defaultSwitcher) deactivate() error {
	err := m.globalConfig.Delete(registryKeyMaintenance)
	if err != nil {
		return fmt.Errorf("failed to delete maintenance mode registry key: %w", err)
	}

	return nil
}

type maintenanceRegistryObject struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	Holder string `json:"holder,omitempty"`
}
