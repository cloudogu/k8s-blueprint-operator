package domainservice

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"testing"
)

func Test_checkDependencyVersion(t *testing.T) {
	type args struct {
		doguInBlueprint domain.TargetDogu
		expectedVersion string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "exact version", args: args{doguInBlueprint: domain.TargetDogu{Name: "nginx", Version: "2.0.0-1"}, expectedVersion: "2.0.0-1"}, wantErr: false},
		{name: "has lower version", args: args{doguInBlueprint: domain.TargetDogu{Name: "nginx", Version: "2.0.0-1"}, expectedVersion: ">=2.0.0-2"}, wantErr: true},
		{name: "has higher version", args: args{doguInBlueprint: domain.TargetDogu{Name: "nginx", Version: "2.0.0-3"}, expectedVersion: ">=2.0.0-2"}, wantErr: false},
		{name: "needs lower version", args: args{doguInBlueprint: domain.TargetDogu{Name: "nginx", Version: "2.0.0-3"}, expectedVersion: "<=2.0.0-2"}, wantErr: true},
		{name: "needs higher version", args: args{doguInBlueprint: domain.TargetDogu{Name: "nginx", Version: "2.0.0-1"}, expectedVersion: ">2.0.0-1"}, wantErr: true},
		{name: "no constraint", args: args{doguInBlueprint: domain.TargetDogu{Name: "nginx", Version: "2.0.0-1"}, expectedVersion: ""}, wantErr: false},
		{name: "not parsable expected version", args: args{doguInBlueprint: domain.TargetDogu{Name: "nginx", Version: "2.0.0-1"}, expectedVersion: "abc"}, wantErr: true},
		{name: "not parsable actual version", args: args{doguInBlueprint: domain.TargetDogu{Name: "nginx", Version: "abc"}, expectedVersion: "2.0.0-1"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkDependencyVersion(tt.args.doguInBlueprint, tt.args.expectedVersion); (err != nil) != tt.wantErr {
				t.Errorf("checkDependencyVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBlueprintSpecDomainUseCase_ValidateDependenciesForAllDogus(t *testing.T) {

	type fields struct {
		remoteDoguRegistry RemoteDoguRegistry
	}
	type args struct {
		effectiveBlueprint domain.EffectiveBlueprint
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "all ok",
			args: args{
				effectiveBlueprint: domain.EffectiveBlueprint{
					Dogus: []domain.TargetDogu{
						{Namespace: "official", Name: "redmine", Version: "1.0.0-1"},
						{Namespace: "official", Name: "postgres", Version: "1.0.0-1"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "namespace change", // redmine has dependency to postgres, but postgres is usually in official
			args: args{
				effectiveBlueprint: domain.EffectiveBlueprint{
					Dogus: []domain.TargetDogu{
						{Namespace: "official", Name: "redmine", Version: "1.0.0-1"},
						{Namespace: "premium", Name: "postgres", Version: "1.0.0-1"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing dependency",
			args: args{
				effectiveBlueprint: domain.EffectiveBlueprint{
					Dogus: []domain.TargetDogu{
						{Namespace: "official", Name: "redmine", Version: "1.0.0-1"},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "unknown dogu",
			args: args{
				effectiveBlueprint: domain.EffectiveBlueprint{
					Dogus: []domain.TargetDogu{
						{Namespace: "official", Name: "redmine2", Version: "1.0.0-1"},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "package dependency",
			args: args{
				effectiveBlueprint: domain.EffectiveBlueprint{
					Dogus: []domain.TargetDogu{
						{Namespace: "official", Name: "k8s-ces-control", Version: "1.0.0-1"},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useCase := NewValidateDependenciesDomainUseCase(testDataDoguRegistry)
			if err := useCase.ValidateDependenciesForAllDogus(tt.args.effectiveBlueprint); (err != nil) != tt.wantErr {
				t.Errorf("ValidateDependenciesForAllDogus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
