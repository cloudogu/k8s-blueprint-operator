package domainservice

import (
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	version1_0_0_1, _ = core.ParseVersion("1.0.0-1")
	version2_0_0_1, _ = core.ParseVersion("2.0.0-1")
	version2_0_0_3, _ = core.ParseVersion("2.0.0-3")
)

func Test_checkDependencyVersion(t *testing.T) {
	type args struct {
		doguInBlueprint domain.Dogu
		expectedVersion string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "exact version", args: args{doguInBlueprint: domain.Dogu{Name: "nginx", Version: version2_0_0_1}, expectedVersion: "2.0.0-1"}, wantErr: false},
		{name: "has lower version", args: args{doguInBlueprint: domain.Dogu{Name: "nginx", Version: version2_0_0_1}, expectedVersion: ">=2.0.0-2"}, wantErr: true},
		{name: "has higher version", args: args{doguInBlueprint: domain.Dogu{Name: "nginx", Version: version2_0_0_3}, expectedVersion: ">=2.0.0-2"}, wantErr: false},
		{name: "needs lower version", args: args{doguInBlueprint: domain.Dogu{Name: "nginx", Version: version2_0_0_3}, expectedVersion: "<=2.0.0-2"}, wantErr: true},
		{name: "needs higher version", args: args{doguInBlueprint: domain.Dogu{Name: "nginx", Version: version2_0_0_1}, expectedVersion: ">2.0.0-1"}, wantErr: true},
		{name: "no constraint", args: args{doguInBlueprint: domain.Dogu{Name: "nginx", Version: version2_0_0_1}, expectedVersion: ""}, wantErr: false},
		{name: "not parsable expected version", args: args{doguInBlueprint: domain.Dogu{Name: "nginx", Version: version2_0_0_1}, expectedVersion: "abc"}, wantErr: true},
		{name: "not parsable actual version", args: args{doguInBlueprint: domain.Dogu{Name: "nginx", Version: core.Version{Raw: "abc"}}, expectedVersion: "2.0.0-1"}, wantErr: true},
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
					Dogus: []domain.Dogu{
						{Namespace: "official", Name: "redmine", Version: version1_0_0_1},
						{Namespace: "official", Name: "postgres", Version: version1_0_0_1},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "namespace change", // redmine has dependency to postgres, but postgres is usually in official
			args: args{
				effectiveBlueprint: domain.EffectiveBlueprint{
					Dogus: []domain.Dogu{
						{Namespace: "official", Name: "redmine", Version: version1_0_0_1},
						{Namespace: "premium", Name: "postgres", Version: version1_0_0_1},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing dependency",
			args: args{
				effectiveBlueprint: domain.EffectiveBlueprint{
					Dogus: []domain.Dogu{
						{Namespace: "official", Name: "redmine", Version: version1_0_0_1},
					},
				},
			},
			wantErr:       true,
			errorContains: "dependencies for dogu 'redmine' are not satisfied in blueprint: dependency 'postgres' in version '1.0.0-1' is not a present dogu in the effective blueprint",
		},
		{
			name: "unknown dogu",
			args: args{
				effectiveBlueprint: domain.EffectiveBlueprint{
					Dogus: []domain.Dogu{
						{Namespace: "official", Name: "redmine2", Version: version1_0_0_1},
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
					Dogus: []domain.Dogu{
						{Namespace: "official", Name: "k8s-ces-control", Version: version1_0_0_1},
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

func TestValidateDependenciesDomainUseCase_ValidateDependenciesForAllDogus_collectDependencyErrors(t *testing.T) {
	//given
	useCase := NewValidateDependenciesDomainUseCase(testDataDoguRegistry)
	//when
	err := useCase.ValidateDependenciesForAllDogus(domain.EffectiveBlueprint{
		Dogus: []domain.Dogu{
			{Namespace: "official", Name: "redmine", Version: version1_0_0_1},
			{Namespace: "helloworld", Name: "bluespice", Version: version1_0_0_1},
		},
	})
	//then
	require.Error(t, err)
	var expectedErrorType *domain.InvalidBlueprintError
	require.ErrorAs(t, err, &expectedErrorType)
	assert.ErrorContains(t, err, "dependencies are not satisfied in effective blueprint")
	assert.ErrorContains(t, err, "dependencies for dogu 'redmine' are not satisfied in blueprint: dependency 'postgres' in version '1.0.0-1' is not a present dogu in the effective blueprint")
	assert.ErrorContains(t, err, "dependencies for dogu 'bluespice' are not satisfied in blueprint: dependency 'official/mysql' in version '1.0.0-1' is not a present dogu in the effective blueprint")
}
