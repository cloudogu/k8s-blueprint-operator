//go:build !ignore_autogenerated

/*
This file was generated with "make generate".
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1

import (
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/adapter/serializer"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Blueprint) DeepCopyInto(out *Blueprint) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Blueprint.
func (in *Blueprint) DeepCopy() *Blueprint {
	if in == nil {
		return nil
	}
	out := new(Blueprint)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Blueprint) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BlueprintList) DeepCopyInto(out *BlueprintList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Blueprint, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BlueprintList.
func (in *BlueprintList) DeepCopy() *BlueprintList {
	if in == nil {
		return nil
	}
	out := new(BlueprintList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BlueprintList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BlueprintSpec) DeepCopyInto(out *BlueprintSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BlueprintSpec.
func (in *BlueprintSpec) DeepCopy() *BlueprintSpec {
	if in == nil {
		return nil
	}
	out := new(BlueprintSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BlueprintStatus) DeepCopyInto(out *BlueprintStatus) {
	*out = *in
	in.EffectiveBlueprint.DeepCopyInto(&out.EffectiveBlueprint)
	in.StateDiff.DeepCopyInto(&out.StateDiff)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BlueprintStatus.
func (in *BlueprintStatus) DeepCopy() *BlueprintStatus {
	if in == nil {
		return nil
	}
	out := new(BlueprintStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CombinedDoguConfig) DeepCopyInto(out *CombinedDoguConfig) {
	*out = *in
	in.Config.DeepCopyInto(&out.Config)
	in.SensitiveConfig.DeepCopyInto(&out.SensitiveConfig)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CombinedDoguConfig.
func (in *CombinedDoguConfig) DeepCopy() *CombinedDoguConfig {
	if in == nil {
		return nil
	}
	out := new(CombinedDoguConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CombinedDoguConfigDiff) DeepCopyInto(out *CombinedDoguConfigDiff) {
	*out = *in
	if in.DoguConfigDiff != nil {
		in, out := &in.DoguConfigDiff, &out.DoguConfigDiff
		*out = make(DoguConfigDiff, len(*in))
		copy(*out, *in)
	}
	if in.SensitiveDoguConfigDiff != nil {
		in, out := &in.SensitiveDoguConfigDiff, &out.SensitiveDoguConfigDiff
		*out = make(SensitiveDoguConfigDiff, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CombinedDoguConfigDiff.
func (in *CombinedDoguConfigDiff) DeepCopy() *CombinedDoguConfigDiff {
	if in == nil {
		return nil
	}
	out := new(CombinedDoguConfigDiff)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ComponentDiff) DeepCopyInto(out *ComponentDiff) {
	*out = *in
	out.Actual = in.Actual
	out.Expected = in.Expected
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ComponentDiff.
func (in *ComponentDiff) DeepCopy() *ComponentDiff {
	if in == nil {
		return nil
	}
	out := new(ComponentDiff)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ComponentDiffState) DeepCopyInto(out *ComponentDiffState) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ComponentDiffState.
func (in *ComponentDiffState) DeepCopy() *ComponentDiffState {
	if in == nil {
		return nil
	}
	out := new(ComponentDiffState)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Config) DeepCopyInto(out *Config) {
	*out = *in
	if in.Dogus != nil {
		in, out := &in.Dogus, &out.Dogus
		*out = make(map[string]CombinedDoguConfig, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
	}
	in.Global.DeepCopyInto(&out.Global)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Config.
func (in *Config) DeepCopy() *Config {
	if in == nil {
		return nil
	}
	out := new(Config)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConfigValueState) DeepCopyInto(out *ConfigValueState) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConfigValueState.
func (in *ConfigValueState) DeepCopy() *ConfigValueState {
	if in == nil {
		return nil
	}
	out := new(ConfigValueState)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DoguConfig) DeepCopyInto(out *DoguConfig) {
	*out = *in
	if in.Present != nil {
		in, out := &in.Present, &out.Present
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Absent != nil {
		in, out := &in.Absent, &out.Absent
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DoguConfig.
func (in *DoguConfig) DeepCopy() *DoguConfig {
	if in == nil {
		return nil
	}
	out := new(DoguConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in DoguConfigDiff) DeepCopyInto(out *DoguConfigDiff) {
	{
		in := &in
		*out = make(DoguConfigDiff, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DoguConfigDiff.
func (in DoguConfigDiff) DeepCopy() DoguConfigDiff {
	if in == nil {
		return nil
	}
	out := new(DoguConfigDiff)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DoguConfigEntryDiff) DeepCopyInto(out *DoguConfigEntryDiff) {
	*out = *in
	out.Actual = in.Actual
	out.Expected = in.Expected
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DoguConfigEntryDiff.
func (in *DoguConfigEntryDiff) DeepCopy() *DoguConfigEntryDiff {
	if in == nil {
		return nil
	}
	out := new(DoguConfigEntryDiff)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DoguConfigValueState) DeepCopyInto(out *DoguConfigValueState) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DoguConfigValueState.
func (in *DoguConfigValueState) DeepCopy() *DoguConfigValueState {
	if in == nil {
		return nil
	}
	out := new(DoguConfigValueState)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DoguDiff) DeepCopyInto(out *DoguDiff) {
	*out = *in
	out.Actual = in.Actual
	out.Expected = in.Expected
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DoguDiff.
func (in *DoguDiff) DeepCopy() *DoguDiff {
	if in == nil {
		return nil
	}
	out := new(DoguDiff)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DoguDiffState) DeepCopyInto(out *DoguDiffState) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DoguDiffState.
func (in *DoguDiffState) DeepCopy() *DoguDiffState {
	if in == nil {
		return nil
	}
	out := new(DoguDiffState)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EffectiveBlueprint) DeepCopyInto(out *EffectiveBlueprint) {
	*out = *in
	if in.Dogus != nil {
		in, out := &in.Dogus, &out.Dogus
		*out = make([]serializer.TargetDogu, len(*in))
		copy(*out, *in)
	}
	if in.Components != nil {
		in, out := &in.Components, &out.Components
		*out = make([]serializer.TargetComponent, len(*in))
		copy(*out, *in)
	}
	in.Config.DeepCopyInto(&out.Config)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EffectiveBlueprint.
func (in *EffectiveBlueprint) DeepCopy() *EffectiveBlueprint {
	if in == nil {
		return nil
	}
	out := new(EffectiveBlueprint)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GlobalConfig) DeepCopyInto(out *GlobalConfig) {
	*out = *in
	if in.Present != nil {
		in, out := &in.Present, &out.Present
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Absent != nil {
		in, out := &in.Absent, &out.Absent
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GlobalConfig.
func (in *GlobalConfig) DeepCopy() *GlobalConfig {
	if in == nil {
		return nil
	}
	out := new(GlobalConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in GlobalConfigDiff) DeepCopyInto(out *GlobalConfigDiff) {
	{
		in := &in
		*out = make(GlobalConfigDiff, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GlobalConfigDiff.
func (in GlobalConfigDiff) DeepCopy() GlobalConfigDiff {
	if in == nil {
		return nil
	}
	out := new(GlobalConfigDiff)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GlobalConfigEntryDiff) DeepCopyInto(out *GlobalConfigEntryDiff) {
	*out = *in
	out.Actual = in.Actual
	out.Expected = in.Expected
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GlobalConfigEntryDiff.
func (in *GlobalConfigEntryDiff) DeepCopy() *GlobalConfigEntryDiff {
	if in == nil {
		return nil
	}
	out := new(GlobalConfigEntryDiff)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GlobalConfigValueState) DeepCopyInto(out *GlobalConfigValueState) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GlobalConfigValueState.
func (in *GlobalConfigValueState) DeepCopy() *GlobalConfigValueState {
	if in == nil {
		return nil
	}
	out := new(GlobalConfigValueState)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SensitiveDoguConfig) DeepCopyInto(out *SensitiveDoguConfig) {
	*out = *in
	if in.Present != nil {
		in, out := &in.Present, &out.Present
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Absent != nil {
		in, out := &in.Absent, &out.Absent
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SensitiveDoguConfig.
func (in *SensitiveDoguConfig) DeepCopy() *SensitiveDoguConfig {
	if in == nil {
		return nil
	}
	out := new(SensitiveDoguConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in SensitiveDoguConfigDiff) DeepCopyInto(out *SensitiveDoguConfigDiff) {
	{
		in := &in
		*out = make(SensitiveDoguConfigDiff, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SensitiveDoguConfigDiff.
func (in SensitiveDoguConfigDiff) DeepCopy() SensitiveDoguConfigDiff {
	if in == nil {
		return nil
	}
	out := new(SensitiveDoguConfigDiff)
	in.DeepCopyInto(out)
	return *out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SensitiveDoguConfigEntryDiff) DeepCopyInto(out *SensitiveDoguConfigEntryDiff) {
	*out = *in
	out.Actual = in.Actual
	out.Expected = in.Expected
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SensitiveDoguConfigEntryDiff.
func (in *SensitiveDoguConfigEntryDiff) DeepCopy() *SensitiveDoguConfigEntryDiff {
	if in == nil {
		return nil
	}
	out := new(SensitiveDoguConfigEntryDiff)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StateDiff) DeepCopyInto(out *StateDiff) {
	*out = *in
	if in.DoguDiffs != nil {
		in, out := &in.DoguDiffs, &out.DoguDiffs
		*out = make(map[string]DoguDiff, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.ComponentDiffs != nil {
		in, out := &in.ComponentDiffs, &out.ComponentDiffs
		*out = make(map[string]ComponentDiff, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.DoguConfigDiffs != nil {
		in, out := &in.DoguConfigDiffs, &out.DoguConfigDiffs
		*out = make(map[string]CombinedDoguConfigDiff, len(*in))
		for key, val := range *in {
			(*out)[key] = *val.DeepCopy()
		}
	}
	if in.GlobalConfigDiff != nil {
		in, out := &in.GlobalConfigDiff, &out.GlobalConfigDiff
		*out = make(GlobalConfigDiff, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StateDiff.
func (in *StateDiff) DeepCopy() *StateDiff {
	if in == nil {
		return nil
	}
	out := new(StateDiff)
	in.DeepCopyInto(out)
	return out
}
