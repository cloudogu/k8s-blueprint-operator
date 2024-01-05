package maintenance

import (
	"encoding/json"

	"github.com/cloudogu/cesapp-lib/registry"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
)

const registryKeyMaintenance = "maintenance"
const blueprintOperatorHolder = "k8s-blueprint-operator"

type Switch struct {
	globalConfig globalConfig
}

func NewSwitch(globalConfig registry.ConfigurationContext) *Switch {
	return &Switch{globalConfig: globalConfig}
}

func (m *Switch) GetLock() (domainservice.MaintenanceLock, error) {
	active, err := m.isActive()
	if err != nil {
		return nil, err
	}

	var ours bool
	if active {
		ours, err = m.isOurs()
		if err != nil {
			return nil, err
		}
	}

	return lock{
		isActive: active,
		isOurs:   ours,
	}, nil
}

func (m *Switch) Activate(content domainservice.MaintenancePageModel) error {
	value := maintenanceRegistryObject{
		Title:  content.Title,
		Text:   content.Text,
		Holder: blueprintOperatorHolder,
	}

	marshal, err := json.Marshal(value)
	if err != nil {
		return &domainservice.InternalError{
			WrappedError: err,
			Message:      "failed to marshal maintenance page model",
		}
	}

	err = m.globalConfig.Set(registryKeyMaintenance, string(marshal))
	if err != nil {
		return &domainservice.InternalError{
			WrappedError: err,
			Message:      "failed to set maintenance mode registry key",
		}
	}

	return nil
}

func (m *Switch) Deactivate() error {
	err := m.globalConfig.Delete(registryKeyMaintenance)
	if err != nil {
		return &domainservice.InternalError{
			WrappedError: err,
			Message:      "failed to delete maintenance mode registry key",
		}
	}

	return nil
}

func (m *Switch) isActive() (bool, error) {
	exists, err := m.globalConfig.Exists(registryKeyMaintenance)
	if err != nil {
		return false, &domainservice.InternalError{
			WrappedError: err,
			Message:      "failed to check if maintenance mode registry key exists",
		}
	}

	return exists, nil
}

func (m *Switch) isOurs() (bool, error) {
	rawValue, err := m.globalConfig.Get(registryKeyMaintenance)
	if err != nil {
		return false, &domainservice.InternalError{
			WrappedError: err,
			Message:      "failed to get maintenance mode from configuration registry",
		}
	}

	var value maintenanceRegistryObject
	err = json.Unmarshal([]byte(rawValue), &value)
	if err != nil {
		return false, &domainservice.InternalError{
			WrappedError: err,
			Message:      "failed to unmarshal json of maintenance mode object",
		}
	}

	return value.Holder == blueprintOperatorHolder, nil
}

type lock struct {
	isActive bool
	isOurs   bool
}

func (l lock) IsActive() bool {
	return l.isActive
}

func (l lock) IsOurs() bool {
	return l.isOurs
}

type maintenanceRegistryObject struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	Holder string `json:"holder,omitempty"`
}
