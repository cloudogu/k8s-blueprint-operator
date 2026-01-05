package domainservice

import (
	"fmt"
	"testing"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/stretchr/testify/assert"
)

func TestValidateStorageClassDomainUseCase_ValidateDoguStorageClass(t *testing.T) {
	storageClassName := "test-storage-class"
	anotherStorageClassName := "another-test-storage-class"
	tests := []struct {
		name               string
		doguRepositoryFn   func(t *testing.T) DoguInstallationRepository
		effectiveBlueprint domain.EffectiveBlueprint
		wantErr            assert.ErrorAssertionFunc
	}{
		{
			name: "fail to load installed dogus",
			doguRepositoryFn: func(t *testing.T) DoguInstallationRepository {
				m := NewMockDoguInstallationRepository(t)
				m.EXPECT().GetAll(ctx).Return(nil, assert.AnError)
				return m
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				var internalErr *InternalError
				return assert.ErrorIs(t, err, assert.AnError) &&
					assert.ErrorAs(t, err, &internalErr) &&
					assert.ErrorContains(t, err, "cannot get installed dogus for storage class validation")
			},
		},
		{
			name: "fail validation for multiple dogus",
			doguRepositoryFn: func(t *testing.T) DoguInstallationRepository {
				m := NewMockDoguInstallationRepository(t)
				m.EXPECT().GetAll(ctx).Return(map[cescommons.SimpleName]*ecosystem.DoguInstallation{
					"redmine":  {StorageClassName: nil},
					"postgres": {StorageClassName: &storageClassName},
					"scm":      {StorageClassName: &storageClassName},
				}, nil)
				return m
			},
			effectiveBlueprint: domain.EffectiveBlueprint{
				Dogus: []domain.Dogu{
					{Name: officialRedmine, StorageClassName: &storageClassName},
					{Name: officialPostgres, StorageClassName: nil},
					{Name: officialScm, StorageClassName: &anotherStorageClassName},
				},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				var invalidErr *domain.InvalidBlueprintError
				return assert.ErrorAs(t, err, &invalidErr) &&
					assert.ErrorContains(t, err, "storage classes are invalid in effective blueprint") &&
					assert.ErrorContains(t, err, "wanted dogu redmine's storage class differs from installed dogu: test-storage-class != <nil>") &&
					assert.ErrorContains(t, err, "wanted dogu postgres's storage class differs from installed dogu: <nil> != test-storage-class") &&
					assert.ErrorContains(t, err, "wanted dogu scm's storage class differs from installed dogu: another-test-storage-class != test-storage-class")
			},
		},
		{
			name: "succeed validation for multiple dogus",
			doguRepositoryFn: func(t *testing.T) DoguInstallationRepository {
				m := NewMockDoguInstallationRepository(t)
				m.EXPECT().GetAll(ctx).Return(map[cescommons.SimpleName]*ecosystem.DoguInstallation{
					"redmine":  {StorageClassName: nil},
					"postgres": {StorageClassName: &storageClassName},
					"scm":      {StorageClassName: &anotherStorageClassName},
				}, nil)
				return m
			},
			effectiveBlueprint: domain.EffectiveBlueprint{
				Dogus: []domain.Dogu{
					{Name: officialRedmine, StorageClassName: nil},
					{Name: officialPostgres, StorageClassName: &storageClassName},
					{Name: officialScm, StorageClassName: &anotherStorageClassName},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useCase := NewValidateStorageClassDomainUseCase(tt.doguRepositoryFn(t))
			tt.wantErr(t, useCase.ValidateDoguStorageClass(ctx, tt.effectiveBlueprint), fmt.Sprintf("ValidateDoguStorageClass(%v, %v)", ctx, tt.effectiveBlueprint))
		})
	}
}
