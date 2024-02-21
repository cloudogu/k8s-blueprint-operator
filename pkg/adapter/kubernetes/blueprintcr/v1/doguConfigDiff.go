package v1

type CombinedDoguConfigDiff struct {
	DoguConfigDiff          DoguConfigDiff          `json:"doguConfigDiff,omitempty"`
	SensitiveDoguConfigDiff SensitiveDoguConfigDiff `json:"sensitiveDoguConfigDiff,omitempty"`
}

type DoguConfigValueState ConfigValueState

type DoguConfigDiff []DoguConfigEntryDiff
type DoguConfigEntryDiff struct {
	Key          string               `json:"key,omitempty"`
	Actual       DoguConfigValueState `json:"actual,omitempty"`
	Expected     DoguConfigValueState `json:"expected,omitempty"`
	NeededAction ConfigAction         `json:"neededAction,omitempty"`
}

type SensitiveDoguConfigDiff []SensitiveDoguConfigEntryDiff
type SensitiveDoguConfigEntryDiff struct {
	Key          string               `json:"key,omitempty"`
	Actual       DoguConfigValueState `json:"actual,omitempty"`
	Expected     DoguConfigValueState `json:"expected,omitempty"`
	NeededAction ConfigAction         `json:"neededAction,omitempty"`
}
