package stateDiffV1

import (
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestConvertToDTO(t *testing.T) {
	tests := []struct {
		name        string
		domainModel domain.StateDiff
		want        StateDiffV1
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
				NeededAction: domain.ActionUpgrade,
			}}},
			want: StateDiffV1{DoguDiffs: []DoguDiffV1{{
				DoguName: "ldap",
				Actual: DoguDiffV1State{
					Namespace:         "official",
					Version:           "1.1.1-1",
					InstallationState: "present",
				},
				Expected: DoguDiffV1State{
					Namespace:         "official",
					Version:           "1.2.3-1",
					InstallationState: "present",
				},
				NeededAction: "upgrade",
			}}},
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
					NeededAction: domain.ActionInstall,
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
					NeededAction: domain.ActionUninstall,
				},
			}},
			want: StateDiffV1{DoguDiffs: []DoguDiffV1{
				{
					DoguName: "ldap",
					Actual: DoguDiffV1State{
						Namespace:         "official",
						InstallationState: "absent",
					},
					Expected: DoguDiffV1State{
						Namespace:         "official",
						Version:           "1.2.3-1",
						InstallationState: "present",
					},
					NeededAction: "install",
				},
				{
					DoguName: "nginx-ingress",
					Actual: DoguDiffV1State{
						Namespace:         "k8s",
						Version:           "8.2.3-2",
						InstallationState: "present",
					},
					Expected: DoguDiffV1State{
						Namespace:         "k8s",
						InstallationState: "absent",
					},
					NeededAction: "uninstall",
				},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertToDTO(tt.domainModel); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertToDTO() = %v, want %v", got, tt.want)
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
		dto     StateDiffV1
		want    domain.StateDiff
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "fail to parse actual version of single dogu diff",
			dto: StateDiffV1{
				DoguDiffs: []DoguDiffV1{{
					DoguName: "ldap",
					Actual: DoguDiffV1State{
						Namespace:         "official",
						Version:           "a.b.c-d",
						InstallationState: "present",
					},
					Expected: DoguDiffV1State{
						Namespace:         "official",
						Version:           "1.2.3-4",
						InstallationState: "present",
					},
					NeededAction: "upgrade",
				}},
			},
			want: domain.StateDiff{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "failed to parse actual version \"a.b.c-d\"") &&
					assert.ErrorContains(t, err, "failed to convert dogu diff dto \"ldap\" to domain model")
			},
		},
		{
			name: "fail to parse expected version of single dogu diff",
			dto: StateDiffV1{
				DoguDiffs: []DoguDiffV1{{
					DoguName: "ldap",
					Actual: DoguDiffV1State{
						Namespace:         "official",
						Version:           "1.2.3-4",
						InstallationState: "present",
					},
					Expected: DoguDiffV1State{
						Namespace:         "official",
						Version:           "a.b.c-d",
						InstallationState: "present",
					},
					NeededAction: "downgrade",
				}},
			},
			want: domain.StateDiff{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "failed to parse expected version \"a.b.c-d\"") &&
					assert.ErrorContains(t, err, "failed to convert dogu diff dto \"ldap\" to domain model")
			},
		},
		{
			name: "fail to parse actual installation state of single dogu diff",
			dto: StateDiffV1{
				DoguDiffs: []DoguDiffV1{{
					DoguName: "ldap",
					Actual: DoguDiffV1State{
						Namespace:         "official",
						Version:           "1.2.3-4",
						InstallationState: "invalid",
					},
					Expected: DoguDiffV1State{
						Namespace:         "official",
						Version:           "2.3.4-5",
						InstallationState: "present",
					},
					NeededAction: "install",
				}},
			},
			want: domain.StateDiff{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "failed to parse actual installation state \"invalid\"") &&
					assert.ErrorContains(t, err, "failed to convert dogu diff dto \"ldap\" to domain model")
			},
		},
		{
			name: "fail to parse expected installation state of single dogu diff",
			dto: StateDiffV1{
				DoguDiffs: []DoguDiffV1{{
					DoguName: "ldap",
					Actual: DoguDiffV1State{
						Namespace:         "official",
						Version:           "1.2.3-4",
						InstallationState: "present",
					},
					Expected: DoguDiffV1State{
						Namespace:         "official",
						Version:           "2.3.4-5",
						InstallationState: "invalid",
					},
					NeededAction: "upgrade",
				}},
			},
			want: domain.StateDiff{},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorContains(t, err, "failed to parse expected installation state \"invalid\"") &&
					assert.ErrorContains(t, err, "failed to convert dogu diff dto \"ldap\" to domain model")
			},
		},
		{
			name: "fail with multiple errors in single dogu diff",
			dto: StateDiffV1{
				DoguDiffs: []DoguDiffV1{{
					DoguName: "ldap",
					Actual: DoguDiffV1State{
						Namespace:         "official",
						Version:           "a.b.c-d",
						InstallationState: "invalid",
					},
					Expected: DoguDiffV1State{
						Namespace:         "official",
						Version:           "a.b.c-d",
						InstallationState: "invalid",
					},
					NeededAction: "none",
				}},
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
			dto: StateDiffV1{
				DoguDiffs: []DoguDiffV1{
					{
						DoguName: "postfix",
						Actual: DoguDiffV1State{
							Namespace:         "official",
							Version:           "1.2.3-4",
							InstallationState: "present",
						},
						Expected: DoguDiffV1State{
							Namespace:         "official",
							Version:           "2.3.4-5",
							InstallationState: "present",
						},
						NeededAction: "upgrade",
					},
					{
						DoguName: "ldap",
						Actual: DoguDiffV1State{
							Namespace:         "official",
							Version:           "1.2.3-4",
							InstallationState: "present",
						},
						Expected: DoguDiffV1State{
							Namespace:         "official",
							Version:           "2.3.4-5",
							InstallationState: "invalid",
						},
						NeededAction: "upgrade",
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
			dto: StateDiffV1{
				DoguDiffs: []DoguDiffV1{
					{
						DoguName: "postfix",
						Actual: DoguDiffV1State{
							Namespace:         "official",
							Version:           "a.b.c-d",
							InstallationState: "present",
						},
						Expected: DoguDiffV1State{
							Namespace:         "official",
							Version:           "2.3.4-5",
							InstallationState: "present",
						},
						NeededAction: "none",
					},
					{
						DoguName: "ldap",
						Actual: DoguDiffV1State{
							Namespace:         "official",
							Version:           "1.2.3-4",
							InstallationState: "present",
						},
						Expected: DoguDiffV1State{
							Namespace:         "official",
							Version:           "2.3.4-5",
							InstallationState: "invalid",
						},
						NeededAction: "upgrade",
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
			dto: StateDiffV1{
				DoguDiffs: []DoguDiffV1{
					{
						DoguName: "postfix",
						Actual: DoguDiffV1State{
							Namespace:         "official",
							Version:           "1.2.3-4",
							InstallationState: "present",
						},
						Expected: DoguDiffV1State{
							Namespace:         "official",
							Version:           "2.3.4-5",
							InstallationState: "present",
						},
						NeededAction: "upgrade",
					},
					{
						DoguName: "ldap",
						Actual: DoguDiffV1State{
							Namespace:         "official",
							Version:           "1.2.3-4",
							InstallationState: "present",
						},
						Expected: DoguDiffV1State{
							Namespace:         "official",
							InstallationState: "absent",
						},
						NeededAction: "uninstall",
					},
				},
			},
			want: domain.StateDiff{DoguDiffs: []domain.DoguDiff{
				{
					DoguName: "postfix",
					Actual: domain.DoguDiffState{Namespace: "official",
						Version:           mustParseVersion("1.2.3-4"),
						InstallationState: domain.TargetStatePresent,
					},
					Expected: domain.DoguDiffState{
						Namespace:         "official",
						Version:           mustParseVersion("2.3.4-5"),
						InstallationState: domain.TargetStatePresent,
					},
					NeededAction: domain.ActionUpgrade,
				},
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
					}, NeededAction: domain.ActionUninstall,
				},
			}},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NoError(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertToDomainModel(tt.dto)
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
