package serializer

import (
	"reflect"
	"testing"

	"github.com/Masterminds/semver/v3"
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
	crd "github.com/cloudogu/k8s-blueprint-lib/v2/api/v2"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	googlecmp "github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

const testComponentName = "my-component"

var (
	testSemverVersionLowRaw  = "1.2.3"
	testSemverVersionLow     = semver.MustParse(testSemverVersionLowRaw)
	testSemverVersionHighRaw = "2.3.4"
	testSemverVersionHigh    = semver.MustParse(testSemverVersionHighRaw)
	testCoreVersionLow       = mustParseVersion("1.1.1-1")
	testCoreVersionLowStr    = testCoreVersionLow.String()
	testCoreVersionHigh      = mustParseVersion("1.2.3-1")
	testCoreVersionHighStr   = testCoreVersionHigh.String()
	testDogu                 = cescommons.SimpleName("testDogu")
	testDogu2                = cescommons.SimpleName("testDogu2")
	testDoguKey1             = common.DoguConfigKey{DoguName: testDogu, Key: "key1"}
	testDoguKey2             = common.DoguConfigKey{DoguName: testDogu2, Key: "key2"}
	testFqdn1                = "ces1.example.com"
	testFqdn2                = "ces2.example.com"
	testSubfolderStr         = "subfolder"
	testSubfolderStr2        = "different_subfolder"
)

func TestConvertToDTO(t *testing.T) {
	tests := []struct {
		name        string
		domainModel domain.StateDiff
		want        *crd.StateDiff
	}{
		{
			name: "should convert single dogu diff",
			domainModel: domain.StateDiff{DoguDiffs: []domain.DoguDiff{{
				DoguName: "ldap",
				Actual: domain.DoguDiffState{
					Namespace: "official",
					Version:   &testCoreVersionLow,
					Absent:    false,
				},
				Expected: domain.DoguDiffState{
					Namespace: "official",
					Version:   &testCoreVersionHigh,
					Absent:    false,
				},
				NeededActions: []domain.Action{domain.ActionUpgrade},
			}}},
			want: &crd.StateDiff{DoguDiffs: map[string]crd.DoguDiff{
				"ldap": {
					Actual: crd.DoguDiffState{
						Namespace: "official",
						Version:   &testCoreVersionLowStr,
						Absent:    false,
					},
					Expected: crd.DoguDiffState{
						Namespace: "official",
						Version:   &testCoreVersionHighStr,
						Absent:    false,
					},
					NeededActions: []crd.DoguAction{"upgrade"},
				},
			}, ComponentDiffs: map[string]crd.ComponentDiff{}},
		},
		{
			name: "should convert multiple dogu diffs",
			domainModel: domain.StateDiff{DoguDiffs: []domain.DoguDiff{
				{
					DoguName: "ldap",
					Actual: domain.DoguDiffState{
						Namespace: "official",
						Absent:    true,
					},
					Expected: domain.DoguDiffState{
						Namespace: "official",
						Version:   &testCoreVersionHigh,
						Absent:    false,
					},
					NeededActions: []domain.Action{domain.ActionInstall},
				},
				{
					DoguName: "nginx-ingress",
					Actual: domain.DoguDiffState{
						Namespace: "k8s",
						Version:   &testCoreVersionLow,
						Absent:    false,
					},
					Expected: domain.DoguDiffState{
						Namespace: "k8s",
						Absent:    true,
					},
					NeededActions: []domain.Action{domain.ActionUninstall},
				},
			}},
			want: &crd.StateDiff{DoguDiffs: map[string]crd.DoguDiff{
				"ldap": {
					Actual: crd.DoguDiffState{
						Namespace: "official",
						Absent:    true,
					},
					Expected: crd.DoguDiffState{
						Namespace: "official",
						Version:   &testCoreVersionHighStr,
						Absent:    false,
					},
					NeededActions: []crd.DoguAction{"install"},
				},
				"nginx-ingress": {
					Actual: crd.DoguDiffState{
						Namespace: "k8s",
						Version:   &testCoreVersionLowStr,
						Absent:    false,
					},
					Expected: crd.DoguDiffState{
						Namespace: "k8s",
						Absent:    true,
					},
					NeededActions: []crd.DoguAction{"uninstall"},
				},
			}, ComponentDiffs: map[string]crd.ComponentDiff{}},
		},
		{
			name: "should convert multiple component diffs",
			domainModel: domain.StateDiff{
				DoguDiffs: domain.DoguDiffs{},
				ComponentDiffs: []domain.ComponentDiff{
					{
						Name:          testComponentName,
						Actual:        domain.ComponentDiffState{Version: testSemverVersionLow, Absent: false},
						Expected:      domain.ComponentDiffState{Version: testSemverVersionHigh, Absent: false},
						NeededActions: []domain.Action{domain.ActionUpgrade, domain.ActionSwitchComponentNamespace},
					},
					{
						Name:          "my-component-2",
						Actual:        domain.ComponentDiffState{Version: testSemverVersionHigh, Absent: false},
						Expected:      domain.ComponentDiffState{Absent: true},
						NeededActions: []domain.Action{domain.ActionUninstall},
					},
				}},
			want: &crd.StateDiff{
				DoguDiffs: map[string]crd.DoguDiff{},
				ComponentDiffs: map[string]crd.ComponentDiff{
					testComponentName: {
						Actual:        crd.ComponentDiffState{Version: &testSemverVersionLowRaw, Absent: false},
						Expected:      crd.ComponentDiffState{Version: &testSemverVersionHighRaw, Absent: false},
						NeededActions: []crd.ComponentAction{"upgrade", "component namespace switch"},
					},
					"my-component-2": {
						Actual:        crd.ComponentDiffState{Version: &testSemverVersionHighRaw, Absent: false},
						Expected:      crd.ComponentDiffState{Absent: true},
						NeededActions: []crd.ComponentAction{"uninstall"},
					},
				}},
		},
		{
			name: "should convert multiple dogu config diffs",
			domainModel: domain.StateDiff{
				DoguConfigDiffs: map[cescommons.SimpleName]domain.DoguConfigDiffs{
					"ldap":    {},
					"postfix": {},
				},
				SensitiveDoguConfigDiffs: map[cescommons.SimpleName]domain.SensitiveDoguConfigDiffs{
					"ldap":    {},
					"postfix": {},
				},
			},
			want: &crd.StateDiff{
				DoguDiffs:      map[string]crd.DoguDiff{},
				ComponentDiffs: map[string]crd.ComponentDiff{},
				DoguConfigDiffs: map[string]crd.CombinedDoguConfigDiff{
					"ldap":    {},
					"postfix": {},
				},
			},
		},
		{
			name: "should convert global config diff",
			domainModel: domain.StateDiff{
				GlobalConfigDiffs: []domain.GlobalConfigEntryDiff{{
					Key: "fqdn",
					Actual: domain.GlobalConfigValueState{
						Value:  &testFqdn1,
						Exists: true,
					},
					Expected: domain.GlobalConfigValueState{
						Value:  &testFqdn2,
						Exists: true,
					},
					NeededAction: domain.ConfigActionSet,
				}},
			},
			want: &crd.StateDiff{
				DoguDiffs:      map[string]crd.DoguDiff{},
				ComponentDiffs: map[string]crd.ComponentDiff{},
				GlobalConfigDiff: []crd.ConfigEntryDiff{{
					Key: "fqdn",
					Actual: crd.ConfigValueState{
						Value:  &testFqdn1,
						Exists: true,
					},
					Expected: crd.ConfigValueState{
						Value:  &testFqdn2,
						Exists: true,
					},
					NeededAction: "set",
				}},
			},
		},
		{
			name: "should convert additional mounts in diff",
			domainModel: domain.StateDiff{DoguDiffs: []domain.DoguDiff{{
				DoguName: "ldap",
				Actual: domain.DoguDiffState{
					Namespace: "official",
					Version:   &testCoreVersionLow,
					Absent:    false,
					AdditionalMounts: []ecosystem.AdditionalMount{
						{
							SourceType: ecosystem.DataSourceConfigMap,
							Name:       "configmap",
							Volume:     "volume",
							Subfolder:  &testSubfolderStr,
						},
					},
				},
				Expected: domain.DoguDiffState{
					Namespace: "official",
					Version:   &testCoreVersionHigh,
					Absent:    false,
					AdditionalMounts: []ecosystem.AdditionalMount{
						{
							SourceType: ecosystem.DataSourceConfigMap,
							Name:       "secret",
							Volume:     "volume2",
							Subfolder:  &testSubfolderStr2,
						},
					},
				},
				NeededActions: []domain.Action{domain.ActionUpgrade},
			}}},
			want: &crd.StateDiff{
				ComponentDiffs: map[string]crd.ComponentDiff{},
				DoguDiffs: map[string]crd.DoguDiff{
					"ldap": {
						Actual: crd.DoguDiffState{
							Namespace: "official",
							Version:   &testCoreVersionLowStr,
							Absent:    false,
							AdditionalMounts: []crd.AdditionalMount{
								{
									SourceType: crd.DataSourceConfigMap,
									Name:       "configmap",
									Volume:     "volume",
									Subfolder:  &testSubfolderStr,
								},
							},
						},
						Expected: crd.DoguDiffState{
							Namespace: "official",
							Version:   &testCoreVersionHighStr,
							Absent:    false,
							AdditionalMounts: []crd.AdditionalMount{
								{
									SourceType: crd.DataSourceConfigMap,
									Name:       "secret",
									Volume:     "volume2",
									Subfolder:  &testSubfolderStr2,
								},
							},
						},
						NeededActions: []crd.DoguAction{"upgrade"},
					},
				},
			},
		}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertToStateDiffDTO(tt.domainModel); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertToStateDiffDTO() = %v, want %v", got, tt.want)
				assert.Empty(t, googlecmp.Diff(tt.want, got))
			}
		})
	}
}

func mustParseVersion(raw string) core.Version {
	version, err := core.ParseVersion(raw)
	if err != nil {
		panic(err)
	}

	return version
}

func TestConvertToStateDiffDTO(t *testing.T) {
	value1 := "1"
	value123 := "123"
	tests := []struct {
		name  string
		model domain.StateDiff
		want  *crd.StateDiff
	}{
		{
			name: "normal dogu config",
			model: domain.StateDiff{
				DoguDiffs:      nil,
				ComponentDiffs: nil,
				DoguConfigDiffs: map[cescommons.SimpleName]domain.DoguConfigDiffs{
					testDogu: {
						{
							Key: testDoguKey1,
							Actual: domain.DoguConfigValueState{
								Value:  &value1,
								Exists: true,
							},
							Expected: domain.DoguConfigValueState{
								Value:  &value123,
								Exists: true,
							},
							NeededAction: domain.ConfigActionSet,
						},
					},
					testDogu2: {
						{
							Key: testDoguKey2,
							Actual: domain.DoguConfigValueState{
								Value:  nil,
								Exists: false,
							},
							Expected: domain.DoguConfigValueState{
								Value:  &value123,
								Exists: true,
							},
							NeededAction: domain.ConfigActionSet,
						},
					},
				},
				GlobalConfigDiffs: nil,
			},
			want: &crd.StateDiff{
				DoguDiffs:      map[string]crd.DoguDiff{},
				ComponentDiffs: map[string]crd.ComponentDiff{},
				DoguConfigDiffs: map[string]crd.CombinedDoguConfigDiff{
					testDogu.String(): {
						DoguConfigDiff: crd.DoguConfigDiff{
							{
								Key: testDoguKey1.Key.String(),
								Actual: crd.ConfigValueState{
									Value:  &value1,
									Exists: true,
								},
								Expected: crd.ConfigValueState{
									Value:  &value123,
									Exists: true,
								},
								NeededAction: crd.ConfigAction("set"),
							},
						},
					},
					testDogu2.String(): {
						DoguConfigDiff: crd.DoguConfigDiff{
							{
								Key: testDoguKey2.Key.String(),
								Actual: crd.ConfigValueState{
									Value:  nil,
									Exists: false,
								},
								Expected: crd.ConfigValueState{
									Value:  &value123,
									Exists: true,
								},
								NeededAction: crd.ConfigAction("set"),
							},
						},
					},
				},
				GlobalConfigDiff: nil,
			},
		},
		{
			name: "censor sensitive config",
			model: domain.StateDiff{
				DoguDiffs:       nil,
				ComponentDiffs:  nil,
				DoguConfigDiffs: map[cescommons.SimpleName]domain.DoguConfigDiffs{},
				SensitiveDoguConfigDiffs: map[cescommons.SimpleName]domain.SensitiveDoguConfigDiffs{
					testDogu: {
						{
							Key: testDoguKey1,
							Actual: domain.DoguConfigValueState{
								Value:  &value1,
								Exists: true,
							},
							Expected: domain.DoguConfigValueState{
								Value:  &value123,
								Exists: true,
							},
							NeededAction: domain.ConfigActionSet,
						},
					},
					testDogu2: {
						{
							Key: testDoguKey2,
							Actual: domain.DoguConfigValueState{
								Value:  nil,
								Exists: false,
							},
							Expected: domain.DoguConfigValueState{
								Value:  &value123,
								Exists: true,
							},
							NeededAction: domain.ConfigActionSet,
						},
					},
				},
				GlobalConfigDiffs: nil,
			},
			want: &crd.StateDiff{
				DoguDiffs:      map[string]crd.DoguDiff{},
				ComponentDiffs: map[string]crd.ComponentDiff{},
				DoguConfigDiffs: map[string]crd.CombinedDoguConfigDiff{
					testDogu.String(): {
						SensitiveDoguConfigDiff: crd.SensitiveDoguConfigDiff{
							{
								Key: testDoguKey1.Key.String(),
								Actual: crd.ConfigValueState{
									Value:  nil,
									Exists: true,
								},
								Expected: crd.ConfigValueState{
									Value:  nil,
									Exists: true,
								},
								NeededAction: crd.ConfigAction("set"),
							},
						},
					},
					testDogu2.String(): {
						SensitiveDoguConfigDiff: crd.SensitiveDoguConfigDiff{
							{
								Key: testDoguKey2.Key.String(),
								Actual: crd.ConfigValueState{
									Value:  nil,
									Exists: false,
								},
								Expected: crd.ConfigValueState{
									Value:  nil,
									Exists: true,
								},
								NeededAction: crd.ConfigAction("set"),
							},
						},
					},
				},
				GlobalConfigDiff: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, ConvertToStateDiffDTO(tt.model), "ConvertToStateDiffDTO(%v)", tt.model)
		})
	}
}
