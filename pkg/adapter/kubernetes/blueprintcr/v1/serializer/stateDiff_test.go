package serializer

import (
	"cmp"
	"github.com/Masterminds/semver/v3"
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
	crd "github.com/cloudogu/k8s-blueprint-lib/api/v1"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/stretchr/testify/assert"
	"reflect"
	"slices"
	"testing"
)

const testComponentName = "my-component"

var (
	testVersionLowRaw  = "1.2.3"
	testVersionLow     = semver.MustParse(testVersionLowRaw)
	testVersionHighRaw = "2.3.4"
	testVersionHigh    = semver.MustParse(testVersionHighRaw)
	testDogu           = cescommons.SimpleName("testDogu")
	testDogu2          = cescommons.SimpleName("testDogu2")
	testDoguKey1       = common.DoguConfigKey{DoguName: testDogu, Key: "key1"}
	testDoguKey2       = common.DoguConfigKey{DoguName: testDogu2, Key: "key2"}
)

func TestConvertToDTO(t *testing.T) {
	tests := []struct {
		name        string
		domainModel domain.StateDiff
		want        crd.StateDiff
	}{
		{
			name: "should convert single dogu diff",
			domainModel: domain.StateDiff{DoguDiffs: []domain.DoguDiff{{
				DoguName: "ldap",
				Actual: domain.DoguDiffState{
					Namespace:         "official",
					Version:           mustParseVersion("1.1.1-1"),
					InstallationState: domain.TargetStatePresent,
				},
				Expected: domain.DoguDiffState{
					Namespace:         "official",
					Version:           mustParseVersion("1.2.3-1"),
					InstallationState: domain.TargetStatePresent,
				},
				NeededActions: []domain.Action{domain.ActionUpgrade},
			}}},
			want: crd.StateDiff{DoguDiffs: map[string]crd.DoguDiff{
				"ldap": {
					Actual: crd.DoguDiffState{
						Namespace:         "official",
						Version:           "1.1.1-1",
						InstallationState: "present",
					},
					Expected: crd.DoguDiffState{
						Namespace:         "official",
						Version:           "1.2.3-1",
						InstallationState: "present",
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
						Namespace:         "official",
						InstallationState: domain.TargetStateAbsent,
					},
					Expected: domain.DoguDiffState{
						Namespace:         "official",
						Version:           mustParseVersion("1.2.3-1"),
						InstallationState: domain.TargetStatePresent,
					},
					NeededActions: []domain.Action{domain.ActionInstall},
				},
				{
					DoguName: "nginx-ingress",
					Actual: domain.DoguDiffState{
						Namespace:         "k8s",
						Version:           mustParseVersion("8.2.3-2"),
						InstallationState: domain.TargetStatePresent,
					},
					Expected: domain.DoguDiffState{
						Namespace:         "k8s",
						InstallationState: domain.TargetStateAbsent,
					},
					NeededActions: []domain.Action{domain.ActionUninstall},
				},
			}},
			want: crd.StateDiff{DoguDiffs: map[string]crd.DoguDiff{
				"ldap": {
					Actual: crd.DoguDiffState{
						Namespace:         "official",
						InstallationState: "absent",
					},
					Expected: crd.DoguDiffState{
						Namespace:         "official",
						Version:           "1.2.3-1",
						InstallationState: "present",
					},
					NeededActions: []crd.DoguAction{"install"},
				},
				"nginx-ingress": {
					Actual: crd.DoguDiffState{
						Namespace:         "k8s",
						Version:           "8.2.3-2",
						InstallationState: "present",
					},
					Expected: crd.DoguDiffState{
						Namespace:         "k8s",
						InstallationState: "absent",
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
						Actual:        domain.ComponentDiffState{Version: testVersionLow, InstallationState: domain.TargetStatePresent},
						Expected:      domain.ComponentDiffState{Version: testVersionHigh, InstallationState: domain.TargetStatePresent},
						NeededActions: []domain.Action{domain.ActionUpgrade, domain.ActionSwitchComponentNamespace},
					},
					{
						Name:          "my-component-2",
						Actual:        domain.ComponentDiffState{Version: testVersionHigh, InstallationState: domain.TargetStatePresent},
						Expected:      domain.ComponentDiffState{InstallationState: domain.TargetStateAbsent},
						NeededActions: []domain.Action{domain.ActionUninstall},
					},
				}},
			want: crd.StateDiff{
				DoguDiffs: map[string]crd.DoguDiff{},
				ComponentDiffs: map[string]crd.ComponentDiff{
					testComponentName: {
						Actual:        crd.ComponentDiffState{Version: testVersionLowRaw, InstallationState: "present"},
						Expected:      crd.ComponentDiffState{Version: testVersionHighRaw, InstallationState: "present"},
						NeededActions: []crd.ComponentAction{"upgrade", "component namespace switch"},
					},
					"my-component-2": {
						Actual:        crd.ComponentDiffState{Version: testVersionHighRaw, InstallationState: "present"},
						Expected:      crd.ComponentDiffState{InstallationState: "absent"},
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
			want: crd.StateDiff{
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
						Value:  "ces1.example.com",
						Exists: true,
					},
					Expected: domain.GlobalConfigValueState{
						Value:  "ces2.example.com",
						Exists: true,
					},
					NeededAction: domain.ConfigActionSet,
				}},
			},
			want: crd.StateDiff{
				DoguDiffs:      map[string]crd.DoguDiff{},
				ComponentDiffs: map[string]crd.ComponentDiff{},
				GlobalConfigDiff: []crd.GlobalConfigEntryDiff{{
					Key: "fqdn",
					Actual: crd.GlobalConfigValueState{
						Value:  "ces1.example.com",
						Exists: true,
					},
					Expected: crd.GlobalConfigValueState{
						Value:  "ces2.example.com",
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
					Namespace:         "official",
					Version:           mustParseVersion("1.1.1-1"),
					InstallationState: domain.TargetStatePresent,
					AdditionalMounts: []ecosystem.AdditionalMount{
						{
							SourceType: ecosystem.DataSourceConfigMap,
							Name:       "configmap",
							Volume:     "volume",
							Subfolder:  "different_subfolder",
						},
					},
				},
				Expected: domain.DoguDiffState{
					Namespace:         "official",
					Version:           mustParseVersion("1.2.3-1"),
					InstallationState: domain.TargetStatePresent,
					AdditionalMounts: []ecosystem.AdditionalMount{
						{
							SourceType: ecosystem.DataSourceConfigMap,
							Name:       "secret",
							Volume:     "volume2",
							Subfolder:  "different_subfolder2",
						},
					},
				},
				NeededActions: []domain.Action{domain.ActionUpgrade},
			}}},
			want: crd.StateDiff{
				ComponentDiffs: map[string]crd.ComponentDiff{},
				DoguDiffs: map[string]crd.DoguDiff{
					"ldap": {
						Actual: crd.DoguDiffState{
							Namespace:         "official",
							Version:           "1.1.1-1",
							InstallationState: "present",
							AdditionalMounts: []crd.AdditionalMount{
								{
									SourceType: crd.DataSourceConfigMap,
									Name:       "configmap",
									Volume:     "volume",
									Subfolder:  "different_subfolder",
								},
							},
						},
						Expected: crd.DoguDiffState{
							Namespace:         "official",
							Version:           "1.2.3-1",
							InstallationState: "present",
							AdditionalMounts: []crd.AdditionalMount{
								{
									SourceType: crd.DataSourceConfigMap,
									Name:       "secret",
									Volume:     "volume2",
									Subfolder:  "different_subfolder2",
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
				assert.Equal(t, tt.want, got)
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

func TestConvertToDomainModel(t *testing.T) {
	tests := []struct {
		name    string
		dto     crd.StateDiff
		want    domain.StateDiff
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "fail to parse actual version of single dogu diff",
			dto: crd.StateDiff{
				DoguDiffs: map[string]crd.DoguDiff{
					"ldap": {
						Actual: crd.DoguDiffState{
							Namespace:         "official",
							Version:           "a.b.c-d",
							InstallationState: "present",
						},
						Expected: crd.DoguDiffState{
							Namespace:         "official",
							Version:           "1.2.3-4",
							InstallationState: "present",
						},
						NeededActions: []crd.DoguAction{"upgrade"},
					},
				},
			},
			want: domain.StateDiff{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "failed to parse version \"a.b.c-d\"") &&
					assert.ErrorContains(t, err, "failed to convert dogu diff dto \"ldap\" to domain model")
			},
		},
		{
			name: "fail to parse expected version of single dogu diff",
			dto: crd.StateDiff{
				DoguDiffs: map[string]crd.DoguDiff{
					"ldap": {
						Actual: crd.DoguDiffState{
							Namespace:         "official",
							Version:           "1.2.3-4",
							InstallationState: "present",
						},
						Expected: crd.DoguDiffState{
							Namespace:         "official",
							Version:           "a.b.c-d",
							InstallationState: "present",
						},
						NeededActions: []crd.DoguAction{"downgrade"},
					},
				},
			},
			want: domain.StateDiff{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "failed to parse version \"a.b.c-d\"") &&
					assert.ErrorContains(t, err, "failed to convert dogu diff dto \"ldap\" to domain model")
			},
		},
		{
			name: "fail to parse actual installation state of single dogu diff",
			dto: crd.StateDiff{
				DoguDiffs: map[string]crd.DoguDiff{
					"ldap": {
						Actual: crd.DoguDiffState{
							Namespace:         "official",
							Version:           "1.2.3-4",
							InstallationState: "invalid",
						},
						Expected: crd.DoguDiffState{
							Namespace:         "official",
							Version:           "2.3.4-5",
							InstallationState: "present",
						},
						NeededActions: []crd.DoguAction{"install"},
					},
				},
			},
			want: domain.StateDiff{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "failed to parse installation state \"invalid\"") &&
					assert.ErrorContains(t, err, "failed to convert dogu diff dto \"ldap\" to domain model")
			},
		},
		{
			name: "fail to parse expected installation state of single dogu diff",
			dto: crd.StateDiff{
				DoguDiffs: map[string]crd.DoguDiff{
					"ldap": {
						Actual: crd.DoguDiffState{
							Namespace:         "official",
							Version:           "1.2.3-4",
							InstallationState: "present",
						},
						Expected: crd.DoguDiffState{
							Namespace:         "official",
							Version:           "2.3.4-5",
							InstallationState: "invalid",
						},
						NeededActions: []crd.DoguAction{"upgrade"},
					},
				},
			},
			want: domain.StateDiff{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "failed to parse installation state \"invalid\"") &&
					assert.ErrorContains(t, err, "failed to convert dogu diff dto \"ldap\" to domain model")
			},
		},
		{
			name: "fail with multiple errors in single dogu diff",
			dto: crd.StateDiff{
				DoguDiffs: map[string]crd.DoguDiff{
					"ldap": {
						Actual: crd.DoguDiffState{
							Namespace:         "official",
							Version:           "a.b.c-d",
							InstallationState: "invalid",
						},
						Expected: crd.DoguDiffState{
							Namespace:         "official",
							Version:           "a.b.c-d",
							InstallationState: "invalid",
						},
						NeededActions: []crd.DoguAction{"none"},
					},
				},
			},
			want: domain.StateDiff{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "failed to parse version \"a.b.c-d\"") &&
					assert.ErrorContains(t, err, "failed to parse version \"a.b.c-d\"") &&
					assert.ErrorContains(t, err, "failed to parse installation state \"invalid\"") &&
					assert.ErrorContains(t, err, "failed to parse installation state \"invalid\"") &&
					assert.ErrorContains(t, err, "failed to convert dogu diff dto \"ldap\" to domain model")
			},
		},
		{
			name: "fail for one of multiple dogu diffs",
			dto: crd.StateDiff{
				DoguDiffs: map[string]crd.DoguDiff{
					"postfix": {
						Actual: crd.DoguDiffState{
							Namespace:         "official",
							Version:           "1.2.3-4",
							InstallationState: "present",
						},
						Expected: crd.DoguDiffState{
							Namespace:         "official",
							Version:           "2.3.4-5",
							InstallationState: "present",
						},
						NeededActions: []crd.DoguAction{"upgrade"},
					},
					"ldap": {
						Actual: crd.DoguDiffState{
							Namespace:         "official",
							Version:           "1.2.3-4",
							InstallationState: "present",
						},
						Expected: crd.DoguDiffState{
							Namespace:         "official",
							Version:           "2.3.4-5",
							InstallationState: "invalid",
						},
						NeededActions: []crd.DoguAction{"upgrade"},
					},
				},
			},
			want: domain.StateDiff{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "failed to parse installation state \"invalid\"") &&
					assert.ErrorContains(t, err, "failed to convert dogu diff dto \"ldap\" to domain model")
			},
		},
		{
			name: "fail for multiple dogu diffs",
			dto: crd.StateDiff{
				DoguDiffs: map[string]crd.DoguDiff{
					"postfix": {
						Actual: crd.DoguDiffState{
							Namespace:         "official",
							Version:           "a.b.c-d",
							InstallationState: "present",
						},
						Expected: crd.DoguDiffState{
							Namespace:         "official",
							Version:           "2.3.4-5",
							InstallationState: "present",
						},
						NeededActions: []crd.DoguAction{"none"},
					},
					"ldap": {
						Actual: crd.DoguDiffState{
							Namespace:         "official",
							Version:           "1.2.3-4",
							InstallationState: "present",
						},
						Expected: crd.DoguDiffState{
							Namespace:         "official",
							Version:           "2.3.4-5",
							InstallationState: "invalid",
						},
						NeededActions: []crd.DoguAction{"upgrade"},
					},
				},
			},
			want: domain.StateDiff{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "failed to parse version \"a.b.c-d\"") &&
					assert.ErrorContains(t, err, "failed to parse installation state \"invalid\"") &&
					assert.ErrorContains(t, err, "failed to convert dogu diff dto \"ldap\" to domain model") &&
					assert.ErrorContains(t, err, "failed to convert dogu diff dto \"postfix\" to domain model")
			},
		},
		{
			name: "succeed for multiple dogu diffs",
			dto: crd.StateDiff{
				DoguDiffs: map[string]crd.DoguDiff{
					"postfix": {
						Actual: crd.DoguDiffState{
							Namespace:         "official",
							Version:           "1.2.3-4",
							InstallationState: "present",
						},
						Expected: crd.DoguDiffState{
							Namespace:         "official",
							Version:           "2.3.4-5",
							InstallationState: "present",
						},
						NeededActions: []crd.DoguAction{"upgrade"},
					},
					"ldap": {
						Actual: crd.DoguDiffState{
							Namespace:         "official",
							Version:           "1.2.3-4",
							InstallationState: "present",
						},
						Expected: crd.DoguDiffState{
							Namespace:         "official",
							InstallationState: "absent",
						},
						NeededActions: []crd.DoguAction{"uninstall"},
					},
				},
			},
			want: domain.StateDiff{DoguDiffs: []domain.DoguDiff{
				{
					DoguName: "ldap",
					Actual: domain.DoguDiffState{
						Namespace:         "official",
						Version:           mustParseVersion("1.2.3-4"),
						InstallationState: domain.TargetStatePresent,
					},
					Expected: domain.DoguDiffState{
						Namespace:         "official",
						InstallationState: domain.TargetStateAbsent,
					}, NeededActions: []domain.Action{domain.ActionUninstall},
				},
				{
					DoguName: "postfix",
					Actual: domain.DoguDiffState{
						Namespace:         "official",
						Version:           mustParseVersion("1.2.3-4"),
						InstallationState: domain.TargetStatePresent,
					},
					Expected: domain.DoguDiffState{
						Namespace:         "official",
						Version:           mustParseVersion("2.3.4-5"),
						InstallationState: domain.TargetStatePresent,
					},
					NeededActions: []domain.Action{domain.ActionUpgrade},
				},
			}},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
		{
			name: "succeed for multiple dogu config diffs",
			dto: crd.StateDiff{
				DoguConfigDiffs: map[string]crd.CombinedDoguConfigDiff{
					"ldap": {
						DoguConfigDiff:          crd.DoguConfigDiff{},
						SensitiveDoguConfigDiff: crd.SensitiveDoguConfigDiff{},
					},
					"postfix": {
						DoguConfigDiff:          crd.DoguConfigDiff{},
						SensitiveDoguConfigDiff: crd.SensitiveDoguConfigDiff{},
					},
				},
			},
			want: domain.StateDiff{
				DoguConfigDiffs: map[cescommons.SimpleName]domain.DoguConfigDiffs{
					"ldap":    nil,
					"postfix": nil,
				},
				SensitiveDoguConfigDiffs: map[cescommons.SimpleName]domain.SensitiveDoguConfigDiffs{
					"ldap":    nil,
					"postfix": nil,
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
		{
			name: "succeed for global config diffs",
			dto: crd.StateDiff{
				GlobalConfigDiff: crd.GlobalConfigDiff{{
					Key: "fqdn",
					Actual: crd.GlobalConfigValueState{
						Value:  "ces1.example.com",
						Exists: true,
					},
					Expected: crd.GlobalConfigValueState{
						Value:  "ces2.example.com",
						Exists: true,
					},
					NeededAction: "set",
				}},
			},
			want: domain.StateDiff{
				GlobalConfigDiffs: []domain.GlobalConfigEntryDiff{{
					Key: "fqdn",
					Actual: domain.GlobalConfigValueState{
						Value:  "ces1.example.com",
						Exists: true,
					},
					Expected: domain.GlobalConfigValueState{
						Value:  "ces2.example.com",
						Exists: true,
					},
					NeededAction: domain.ConfigActionSet,
				}},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
		{
			name: "succeed for multiple component diffs",
			dto: crd.StateDiff{
				ComponentDiffs: map[string]crd.ComponentDiff{
					testComponentName: {
						Actual:        crd.ComponentDiffState{Version: testVersionLowRaw, InstallationState: "present"},
						Expected:      crd.ComponentDiffState{Version: testVersionHighRaw, InstallationState: "present"},
						NeededActions: []crd.ComponentAction{"upgrade", "component namespace switch"},
					},
					"my-component-2": {
						Actual:        crd.ComponentDiffState{Version: testVersionHighRaw, InstallationState: "present"},
						Expected:      crd.ComponentDiffState{InstallationState: "absent"},
						NeededActions: []crd.ComponentAction{"uninstall"},
					},
				},
			},
			want: domain.StateDiff{ComponentDiffs: []domain.ComponentDiff{
				{
					Name:          testComponentName,
					Actual:        domain.ComponentDiffState{Version: testVersionLow, InstallationState: domain.TargetStatePresent},
					Expected:      domain.ComponentDiffState{Version: testVersionHigh, InstallationState: domain.TargetStatePresent},
					NeededActions: []domain.Action{domain.ActionUpgrade, domain.ActionSwitchComponentNamespace},
				},
				{
					Name:          "my-component-2",
					Actual:        domain.ComponentDiffState{Version: testVersionHigh, InstallationState: domain.TargetStatePresent},
					Expected:      domain.ComponentDiffState{InstallationState: domain.TargetStateAbsent},
					NeededActions: []domain.Action{domain.ActionUninstall},
				},
			}},
			wantErr: assert.NoError,
		},
		{
			name: "fail for multiple component diffs",
			dto: crd.StateDiff{
				ComponentDiffs: map[string]crd.ComponentDiff{
					testComponentName: {
						Actual: crd.ComponentDiffState{
							Version:           "a.b.c-d",
							InstallationState: "present",
						},
						Expected: crd.ComponentDiffState{
							Version:           "2.3.4-5",
							InstallationState: "present",
						},
						NeededActions: []crd.ComponentAction{"none"},
					},
					"my-component-2": {
						Actual: crd.ComponentDiffState{
							Version:           "1.2.3-4",
							InstallationState: "present",
						},
						Expected: crd.ComponentDiffState{
							Version:           "2.3.4-5",
							InstallationState: "invalid",
						},
						NeededActions: []crd.ComponentAction{"upgrade"},
					},
				},
			},
			want: domain.StateDiff{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "failed to parse actual version \"a.b.c-d\"") &&
					assert.ErrorContains(t, err, "failed to parse expected installation state \"invalid\"") &&
					assert.ErrorContains(t, err, "failed to convert component diff dto \"my-component\" to domain model") &&
					assert.ErrorContains(t, err, "failed to convert component diff dto \"my-component-2\" to domain model")
			},
		},

		{
			name: "should convert additional mounts",
			dto: crd.StateDiff{
				DoguDiffs: map[string]crd.DoguDiff{
					"ldap": {
						Actual: crd.DoguDiffState{
							Namespace:         "official",
							InstallationState: "present",
							AdditionalMounts: []crd.AdditionalMount{
								{
									SourceType: crd.DataSourceConfigMap,
									Name:       "config",
									Volume:     "volume",
									Subfolder:  "subfolder",
								},
							},
						},
						Expected: crd.DoguDiffState{
							Namespace:         "official",
							InstallationState: "present",
							AdditionalMounts: []crd.AdditionalMount{
								{
									SourceType: crd.DataSourceConfigMap,
									Name:       "config-different",
									Volume:     "volume",
									Subfolder:  "subfolder",
								},
							},
						},
						NeededActions: []crd.DoguAction{"update additional mounts"},
					},
				},
			},
			want: domain.StateDiff{DoguDiffs: []domain.DoguDiff{
				{
					DoguName: "ldap",
					Actual: domain.DoguDiffState{
						Namespace:         "official",
						InstallationState: domain.TargetStatePresent,
						AdditionalMounts: []ecosystem.AdditionalMount{
							{
								SourceType: ecosystem.DataSourceConfigMap,
								Name:       "config",
								Volume:     "volume",
								Subfolder:  "subfolder",
							},
						},
					},
					Expected: domain.DoguDiffState{
						Namespace:         "official",
						InstallationState: domain.TargetStatePresent,
						AdditionalMounts: []ecosystem.AdditionalMount{
							{
								SourceType: ecosystem.DataSourceConfigMap,
								Name:       "config-different",
								Volume:     "volume",
								Subfolder:  "subfolder",
							},
						},
					},
					NeededActions: []domain.Action{domain.ActionUpdateAdditionalMounts},
				},
			}},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertToStateDiffDomain(tt.dto)
			tt.wantErr(t, err)
			// sort to avoid flaky tests
			slices.SortFunc(got.DoguDiffs, func(a, b domain.DoguDiff) int {
				return cmp.Compare(a.DoguName, b.DoguName)
			})
			assert.Equalf(t, tt.want, got, "ConvertToStateDiffDomain(%v)", tt.dto)
		})
	}
}

func TestConvertToStateDiffDTO(t *testing.T) {

	tests := []struct {
		name  string
		model domain.StateDiff
		want  crd.StateDiff
	}{
		{
			name: "ok",
			model: domain.StateDiff{
				DoguDiffs:       nil,
				ComponentDiffs:  nil,
				DoguConfigDiffs: map[cescommons.SimpleName]domain.DoguConfigDiffs{},
				SensitiveDoguConfigDiffs: map[cescommons.SimpleName]domain.SensitiveDoguConfigDiffs{
					testDogu: {
						{
							Key: testDoguKey1,
							Actual: domain.DoguConfigValueState{
								Value:  "1",
								Exists: true,
							},
							Expected: domain.DoguConfigValueState{
								Value:  "123",
								Exists: true,
							},
							NeededAction: domain.ConfigActionSet,
						},
					},
					testDogu2: {
						{
							Key: testDoguKey2,
							Actual: domain.DoguConfigValueState{
								Value:  "",
								Exists: false,
							},
							Expected: domain.DoguConfigValueState{
								Value:  "123",
								Exists: true,
							},
							NeededAction: domain.ConfigActionSet,
						},
					},
				},
				GlobalConfigDiffs: nil,
			},
			want: crd.StateDiff{
				DoguDiffs:      map[string]crd.DoguDiff{},
				ComponentDiffs: map[string]crd.ComponentDiff{},
				DoguConfigDiffs: map[string]crd.CombinedDoguConfigDiff{
					testDogu.String(): {
						DoguConfigDiff: crd.DoguConfigDiff(nil),
						SensitiveDoguConfigDiff: crd.SensitiveDoguConfigDiff{
							crd.DoguConfigEntryDiff{
								Key: testDoguKey1.Key.String(),
								Actual: crd.DoguConfigValueState{
									Value:  "1",
									Exists: true,
								},
								Expected: crd.DoguConfigValueState{
									Value:  "123",
									Exists: true,
								},
								NeededAction: crd.ConfigAction("set"),
							},
						},
					},
					testDogu2.String(): {
						SensitiveDoguConfigDiff: crd.SensitiveDoguConfigDiff{
							crd.DoguConfigEntryDiff{
								Key: testDoguKey2.Key.String(),
								Actual: crd.DoguConfigValueState{
									Value:  "",
									Exists: false,
								},
								Expected: crd.DoguConfigValueState{
									Value:  "123",
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
