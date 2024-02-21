package v1

type GlobalConfigDiff []GlobalConfigEntryDiff

type GlobalConfigValueState ConfigValueState
type GlobalConfigEntryDiff struct {
	Key          string                 `json:"key,omitempty"`
	Actual       GlobalConfigValueState `json:"actual,omitempty"`
	Expected     GlobalConfigValueState `json:"expected,omitempty"`
	NeededAction ConfigAction           `json:"neededAction,omitempty"`
}
