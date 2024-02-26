package v1

type ConfigAction string

const (
	ConfigActionNone         ConfigAction = "none"
	ConfigActionSet          ConfigAction = "set"
	ConfigActionSetToEncrypt ConfigAction = "setToEncrypt"
	//TODO: add new actions
	ConfigActionRemove ConfigAction = "remove"
)

type ConfigValueState struct {
	Value  string `json:"value,omitempty"`
	Exists bool   `json:"exists,omitempty"`
}
