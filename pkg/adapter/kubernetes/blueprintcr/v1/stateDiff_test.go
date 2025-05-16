package v1

import (
	"cmp"
	"github.com/Masterminds/semver/v3"
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
	. "github.com/cloudogu/k8s-blueprint-lib/api/v1"
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
		want        StateDiff
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
			want: StateDiff{DoguDiffs: map[string]DoguDiff{
				"ldap": {
					Actual: DoguDiffState{
						Namespace:         "official",
						Version:           "1.1.1-1",
						InstallationState: "present",
					},
					Expected: DoguDiffState{
						Namespace:         "official",
						Version:           "1.2.3-1",
						InstallationState: "present",
					},
					NeededActions: []DoguAction{"upgrade"},
				},
			}, ComponentDiffs: map[string]ComponentDiff{}},
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
			want: StateDiff{DoguDiffs: map[string]DoguDiff{
				"ldap": {
					Actual: DoguDiffState{
						Namespace:         "official",
						InstallationState: "absent",
					},
					Expected: DoguDiffState{
						Namespace:         "official",
						Version:           "1.2.3-1",
						InstallationState: "present",
					},
					NeededActions: []DoguAction{"install"},
				},
				"nginx-ingress": {
					Actual: DoguDiffState{
						Namespace:         "k8s",
						Version:           "8.2.3-2",
						InstallationState: "present",
					},
					Expected: DoguDiffState{
						Namespace:         "k8s",
						InstallationState: "absent",
					},
					NeededActions: []DoguAction{"uninstall"},
				},
			}, ComponentDiffs: map[string]ComponentDiff{}},
		}, {
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
			want: StateDiff{
				DoguDiffs: map[string]DoguDiff{},
				ComponentDiffs: map[string]ComponentDiff{
					testComponentName: {
						Actual:        ComponentDiffState{Version: testVersionLowRaw, InstallationState: "present"},
						Expected:      ComponentDiffState{Version: testVersionHighRaw, InstallationState: "present"},
						NeededActions: []ComponentAction{"upgrade", "component namespace switch"},
					},
					"my-component-2": {
						Actual:        ComponentDiffState{Version: testVersionHighRaw, InstallationState: "present"},
						Expected:      ComponentDiffState{InstallationState: "absent"},
						NeededActions: []ComponentAction{"uninstall"},
					},
				}},
		}, {
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
			want: StateDiff{
				DoguDiffs:      map[string]DoguDiff{},
				ComponentDiffs: map[string]ComponentDiff{},
				DoguConfigDiffs: map[string]CombinedDoguConfigDiff{
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
			want: StateDiff{
				DoguDiffs:      map[string]DoguDiff{},
				ComponentDiffs: map[string]ComponentDiff{},
				GlobalConfigDiff: []GlobalConfigEntryDiff{{
					Key: "fqdn",
					Actual: GlobalConfigValueState{
						Value:  "ces1.example.com",
						Exists: true,
					},
					Expected: GlobalConfigValueState{
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
							Name:       "configmap",
							Volume:     "volume",
							Subfolder:  "different_subfolder",
						},
					},
				},
				NeededActions: []domain.Action{domain.ActionUpgrade},
			}}},
			want: StateDiff{
				ComponentDiffs: map[string]ComponentDiff{},
				DoguDiffs: map[string]DoguDiff{
					"ldap": {
						Actual: DoguDiffState{
							Namespace:         "official",
							Version:           "1.1.1-1",
							InstallationState: "present",
							AdditionalMounts: []AdditionalMount{
								{
									SourceType: DataSourceConfigMap,
									Name:       "configmap",
									Volume:     "volume",
									Subfolder:  "different_subfolder",
								},
							},
						},
						Expected: DoguDiffState{
							Namespace:         "official",
							Version:           "1.2.3-1",
							InstallationState: "present",
							AdditionalMounts: []AdditionalMount{
								{
									SourceType: DataSourceConfigMap,
									Name:       "configmap",
									Volume:     "volume",
									Subfolder:  "different_subfolder",
								},
							},
						},
						NeededActions: []DoguAction{"upgrade"},
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
		dto     StateDiff
		want    domain.StateDiff
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "fail to parse actual version of single dogu diff",
			dto: StateDiff{
				DoguDiffs: map[string]DoguDiff{
					"ldap": {
						Actual: DoguDiffState{
							Namespace:         "official",
							Version:           "a.b.c-d",
							InstallationState: "present",
						},
						Expected: DoguDiffState{
							Namespace:         "official",
							Version:           "1.2.3-4",
							InstallationState: "present",
						},
						NeededActions: []DoguAction{"upgrade"},
					},
				},
			},
			want: domain.StateDiff{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "failed to parse actual version \"a.b.c-d\"") &&
					assert.ErrorContains(t, err, "failed to convert dogu diff dto \"ldap\" to domain model")
			},
		},
		{
			name: "fail to parse expected version of single dogu diff",
			dto: StateDiff{
				DoguDiffs: map[string]DoguDiff{
					"ldap": {
						Actual: DoguDiffState{
							Namespace:         "official",
							Version:           "1.2.3-4",
							InstallationState: "present",
						},
						Expected: DoguDiffState{
							Namespace:         "official",
							Version:           "a.b.c-d",
							InstallationState: "present",
						},
						NeededActions: []DoguAction{"downgrade"},
					},
				},
			},
			want: domain.StateDiff{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "failed to parse expected version \"a.b.c-d\"") &&
					assert.ErrorContains(t, err, "failed to convert dogu diff dto \"ldap\" to domain model")
			},
		},
		{
			name: "fail to parse actual installation state of single dogu diff",
			dto: StateDiff{
				DoguDiffs: map[string]DoguDiff{
					"ldap": {
						Actual: DoguDiffState{
							Namespace:         "official",
							Version:           "1.2.3-4",
							InstallationState: "invalid",
						},
						Expected: DoguDiffState{
							Namespace:         "official",
							Version:           "2.3.4-5",
							InstallationState: "present",
						},
						NeededActions: []DoguAction{"install"},
					},
				},
			},
			want: domain.StateDiff{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "failed to parse actual installation state \"invalid\"") &&
					assert.ErrorContains(t, err, "failed to convert dogu diff dto \"ldap\" to domain model")
			},
		},
		{
			name: "fail to parse expected installation state of single dogu diff",
			dto: StateDiff{
				DoguDiffs: map[string]DoguDiff{
					"ldap": {
						Actual: DoguDiffState{
							Namespace:         "official",
							Version:           "1.2.3-4",
							InstallationState: "present",
						},
						Expected: DoguDiffState{
							Namespace:         "official",
							Version:           "2.3.4-5",
							InstallationState: "invalid",
						},
						NeededActions: []DoguAction{"upgrade"},
					},
				},
			},
			want: domain.StateDiff{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "failed to parse expected installation state \"invalid\"") &&
					assert.ErrorContains(t, err, "failed to convert dogu diff dto \"ldap\" to domain model")
			},
		},
		{
			name: "fail with multiple errors in single dogu diff",
			dto: StateDiff{
				DoguDiffs: map[string]DoguDiff{
					"ldap": {
						Actual: DoguDiffState{
							Namespace:         "official",
							Version:           "a.b.c-d",
							InstallationState: "invalid",
						},
						Expected: DoguDiffState{
							Namespace:         "official",
							Version:           "a.b.c-d",
							InstallationState: "invalid",
						},
						NeededActions: []DoguAction{"none"},
					},
				},
			},
			want: domain.StateDiff{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "failed to parse actual version \"a.b.c-d\"") &&
					assert.ErrorContains(t, err, "failed to parse expected version \"a.b.c-d\"") &&
					assert.ErrorContains(t, err, "failed to parse actual installation state \"invalid\"") &&
					assert.ErrorContains(t, err, "failed to parse expected installation state \"invalid\"") &&
					assert.ErrorContains(t, err, "failed to convert dogu diff dto \"ldap\" to domain model")
			},
		},
		{
			name: "fail for one of multiple dogu diffs",
			dto: StateDiff{
				DoguDiffs: map[string]DoguDiff{
					"postfix": {
						Actual: DoguDiffState{
							Namespace:         "official",
							Version:           "1.2.3-4",
							InstallationState: "present",
						},
						Expected: DoguDiffState{
							Namespace:         "official",
							Version:           "2.3.4-5",
							InstallationState: "present",
						},
						NeededActions: []DoguAction{"upgrade"},
					},
					"ldap": {
						Actual: DoguDiffState{
							Namespace:         "official",
							Version:           "1.2.3-4",
							InstallationState: "present",
						},
						Expected: DoguDiffState{
							Namespace:         "official",
							Version:           "2.3.4-5",
							InstallationState: "invalid",
						},
						NeededActions: []DoguAction{"upgrade"},
					},
				},
			},
			want: domain.StateDiff{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "failed to parse expected installation state \"invalid\"") &&
					assert.ErrorContains(t, err, "failed to convert dogu diff dto \"ldap\" to domain model")
			},
		},
		{
			name: "fail for multiple dogu diffs",
			dto: StateDiff{
				DoguDiffs: map[string]DoguDiff{
					"postfix": {
						Actual: DoguDiffState{
							Namespace:         "official",
							Version:           "a.b.c-d",
							InstallationState: "present",
						},
						Expected: DoguDiffState{
							Namespace:         "official",
							Version:           "2.3.4-5",
							InstallationState: "present",
						},
						NeededActions: []DoguAction{"none"},
					},
					"ldap": {
						Actual: DoguDiffState{
							Namespace:         "official",
							Version:           "1.2.3-4",
							InstallationState: "present",
						},
						Expected: DoguDiffState{
							Namespace:         "official",
							Version:           "2.3.4-5",
							InstallationState: "invalid",
						},
						NeededActions: []DoguAction{"upgrade"},
					},
				},
			},
			want: domain.StateDiff{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "failed to parse actual version \"a.b.c-d\"") &&
					assert.ErrorContains(t, err, "failed to parse expected installation state \"invalid\"") &&
					assert.ErrorContains(t, err, "failed to convert dogu diff dto \"ldap\" to domain model") &&
					assert.ErrorContains(t, err, "failed to convert dogu diff dto \"postfix\" to domain model")
			},
		},
		{
			name: "succeed for multiple dogu diffs",
			dto: StateDiff{
				DoguDiffs: map[string]DoguDiff{
					"postfix": {
						Actual: DoguDiffState{
							Namespace:         "official",
							Version:           "1.2.3-4",
							InstallationState: "present",
						},
						Expected: DoguDiffState{
							Namespace:         "official",
							Version:           "2.3.4-5",
							InstallationState: "present",
						},
						NeededActions: []DoguAction{"upgrade"},
					},
					"ldap": {
						Actual: DoguDiffState{
							Namespace:         "official",
							Version:           "1.2.3-4",
							InstallationState: "present",
						},
						Expected: DoguDiffState{
							Namespace:         "official",
							InstallationState: "absent",
						},
						NeededActions: []DoguAction{"uninstall"},
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
			dto: StateDiff{
				DoguConfigDiffs: map[string]CombinedDoguConfigDiff{
					"ldap": {
						DoguConfigDiff:          DoguConfigDiff{},
						SensitiveDoguConfigDiff: SensitiveDoguConfigDiff{},
					},
					"postfix": {
						DoguConfigDiff:          DoguConfigDiff{},
						SensitiveDoguConfigDiff: SensitiveDoguConfigDiff{},
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
			dto: StateDiff{
				GlobalConfigDiff: GlobalConfigDiff{{
					Key: "fqdn",
					Actual: GlobalConfigValueState{
						Value:  "ces1.example.com",
						Exists: true,
					},
					Expected: GlobalConfigValueState{
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
			dto: StateDiff{
				ComponentDiffs: map[string]ComponentDiff{
					testComponentName: {
						Actual:        ComponentDiffState{Version: testVersionLowRaw, InstallationState: "present"},
						Expected:      ComponentDiffState{Version: testVersionHighRaw, InstallationState: "present"},
						NeededActions: []ComponentAction{"upgrade", "component namespace switch"},
					},
					"my-component-2": {
						Actual:        ComponentDiffState{Version: testVersionHighRaw, InstallationState: "present"},
						Expected:      ComponentDiffState{InstallationState: "absent"},
						NeededActions: []ComponentAction{"uninstall"},
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
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
		{
			name: "fail for multiple component diffs",
			dto: StateDiff{
				ComponentDiffs: map[string]ComponentDiff{
					testComponentName: {
						Actual: ComponentDiffState{
							Version:           "a.b.c-d",
							InstallationState: "present",
						},
						Expected: ComponentDiffState{
							Version:           "2.3.4-5",
							InstallationState: "present",
						},
						NeededActions: []ComponentAction{"none"},
					},
					"my-component-2": {
						Actual: ComponentDiffState{
							Version:           "1.2.3-4",
							InstallationState: "present",
						},
						Expected: ComponentDiffState{
							Version:           "2.3.4-5",
							InstallationState: "invalid",
						},
						NeededActions: []ComponentAction{"upgrade"},
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
			dto: StateDiff{
				DoguDiffs: map[string]DoguDiff{
					"ldap": {
						Actual: DoguDiffState{
							Namespace:         "official",
							InstallationState: "present",
							AdditionalMounts: []AdditionalMount{
								{
									SourceType: DataSourceConfigMap,
									Name:       "config",
									Volume:     "volume",
									Subfolder:  "subfolder",
								},
							},
						},
						Expected: DoguDiffState{
							Namespace:         "official",
							InstallationState: "present",
							AdditionalMounts: []AdditionalMount{
								{
									SourceType: DataSourceConfigMap,
									Name:       "config-different",
									Volume:     "volume",
									Subfolder:  "subfolder",
								},
							},
						},
						NeededActions: []DoguAction{"update additional mounts"},
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
		want  StateDiff
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
			want: StateDiff{
				DoguDiffs:      map[string]DoguDiff{},
				ComponentDiffs: map[string]ComponentDiff{},
				DoguConfigDiffs: map[string]CombinedDoguConfigDiff{
					testDogu.String(): {
						DoguConfigDiff: DoguConfigDiff(nil),
						SensitiveDoguConfigDiff: SensitiveDoguConfigDiff{
							DoguConfigEntryDiff{
								Key: testDoguKey1.Key.String(),
								Actual: DoguConfigValueState{
									Value:  "1",
									Exists: true,
								},
								Expected: DoguConfigValueState{
									Value:  "123",
									Exists: true,
								},
								NeededAction: ConfigAction("set"),
							},
						},
					},
					testDogu2.String(): {
						SensitiveDoguConfigDiff: SensitiveDoguConfigDiff{
							DoguConfigEntryDiff{
								Key: testDoguKey2.Key.String(),
								Actual: DoguConfigValueState{
									Value:  "",
									Exists: false,
								},
								Expected: DoguConfigValueState{
									Value:  "123",
									Exists: true,
								},
								NeededAction: ConfigAction("set"),
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
