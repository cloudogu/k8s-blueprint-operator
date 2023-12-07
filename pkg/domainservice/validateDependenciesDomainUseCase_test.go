package domainservice

import (
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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

func TestValidateDependenciesDomainUseCase_ValidateDependenciesForAllDogus(t *testing.T) {
	type args struct {
		effectiveBlueprint domain.EffectiveBlueprint
	}
	tests := []struct {
		name          string
		args          args
		wantErr       bool
		errorContains string
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
			wantErr:       true,
			errorContains: "dependencies for dogu 'redmine' are not satisfied blueprint: dependency 'postgres' in version '1.0.0-1' is not a present dogu in the effective blueprint",
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
			wantErr:       true,
			errorContains: "remote dogu registry has no dogu specification for at least one wanted dogu: dogu official/redmine2 in version 1.0.0-1 not found",
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
			err := useCase.ValidateDependenciesForAllDogus(tt.args.effectiveBlueprint)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("ValidateDependenciesForAllDogus() error = %v, wantErr %v", err, tt.wantErr)
				}
				assert.ErrorContains(t, err, tt.errorContains)
			}
		})
	}
}

func TestValidateDependenciesDomainUseCase_ValidateDependenciesForAllDogus_internalError(t *testing.T) {
	//given
	RegistryMock := NewMockRemoteDoguRegistry(t)
	useCase := NewValidateDependenciesDomainUseCase(RegistryMock)

	RegistryMock.EXPECT().GetDogus(mock.Anything).Return(nil, &InternalError{Message: "my error"})
	//when
	err := useCase.ValidateDependenciesForAllDogus(domain.EffectiveBlueprint{})
	//then
	require.Error(t, err)
	var internalError *InternalError
	assert.ErrorAs(t, err, &internalError)
}
