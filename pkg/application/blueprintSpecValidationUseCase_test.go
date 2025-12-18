package application

import (
	"context"
	"errors"
	"testing"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/api/meta"
)

var redmineQualifiedDoguName = cescommons.QualifiedName{
	Namespace:  "official",
	SimpleName: "redmine",
}

func TestBlueprintSpecUseCase_ValidateBlueprintSpecStatically_ok(t *testing.T) {
	//given
	blueprint := &domain.BlueprintSpec{
		Id: "testBlueprint1",
	}

	repoMock := newMockBlueprintSpecRepository(t)
	ctx := context.Background()
	dependencyUseCase := newMockValidateDependenciesDomainUseCase(t)
	mountsUseCase := newMockValidateAdditionalMountsDomainUseCase(t)
	storageClassUseCase := newMockValidateDoguStorageClassDomainUseCase(t)
	useCase := NewBlueprintSpecValidationUseCase(repoMock, dependencyUseCase, mountsUseCase, storageClassUseCase)

	repoMock.EXPECT().Update(ctx, &domain.BlueprintSpec{
		Id: "testBlueprint1",
	}).Return(nil)

	//when
	err := useCase.ValidateBlueprintSpecStatically(ctx, blueprint)

	//then
	require.NoError(t, err)
	assert.Nil(t, blueprint.Conditions, "should not set conditions")
}

func TestBlueprintSpecUseCase_ValidateBlueprintSpecStatically_invalid(t *testing.T) {
	//given
	blueprint := &domain.BlueprintSpec{
		//missing ID
	}

	repoMock := newMockBlueprintSpecRepository(t)
	ctx := context.Background()
	dependencyUseCase := newMockValidateDependenciesDomainUseCase(t)
	mountsUseCase := newMockValidateAdditionalMountsDomainUseCase(t)
	storageClassUseCase := newMockValidateDoguStorageClassDomainUseCase(t)
	useCase := NewBlueprintSpecValidationUseCase(repoMock, dependencyUseCase, mountsUseCase, storageClassUseCase)

	repoMock.EXPECT().
		Update(ctx, blueprint).
		Return(nil)

	//when
	err := useCase.ValidateBlueprintSpecStatically(ctx, blueprint)

	//then
	require.Error(t, err)
	var invalidError *domain.InvalidBlueprintError
	assert.ErrorAs(t, err, &invalidError, "error should be an InvalidBlueprintError")
	assert.ErrorContains(t, err, "blueprint spec is invalid: blueprint spec doesn't have an ID")
}

func TestBlueprintSpecUseCase_ValidateBlueprintSpecStatically_repoError(t *testing.T) {
	t.Run("error while saving blueprint spec", func(t *testing.T) {
		//given
		blueprint := &domain.BlueprintSpec{
			Id: "testBlueprint1",
		}
		repoMock := newMockBlueprintSpecRepository(t)
		ctx := context.Background()
		dependencyUseCase := newMockValidateDependenciesDomainUseCase(t)
		mountsUseCase := newMockValidateAdditionalMountsDomainUseCase(t)
		storageClassUseCase := newMockValidateDoguStorageClassDomainUseCase(t)
		useCase := NewBlueprintSpecValidationUseCase(repoMock, dependencyUseCase, mountsUseCase, storageClassUseCase)

		repoMock.EXPECT().Update(ctx, mock.Anything).Return(&domainservice.InternalError{Message: "test-error"})

		//when
		err := useCase.ValidateBlueprintSpecStatically(ctx, blueprint)

		//then
		require.Error(t, err)
		var invalidError *domainservice.InternalError
		assert.ErrorAs(t, err, &invalidError)
		assert.ErrorContains(t, err, "cannot update blueprint spec after static validation: test-error")
	})

}

func TestBlueprintSpecUseCase_ValidateBlueprintSpecDynamically_ok(t *testing.T) {
	// given
	blueprint := &domain.BlueprintSpec{
		Id:         "testBlueprint1",
		Conditions: []domain.Condition{},
	}
	repoMock := newMockBlueprintSpecRepository(t)
	ctx := context.Background()
	dependencyUseCase := newMockValidateDependenciesDomainUseCase(t)
	mountsUseCase := newMockValidateAdditionalMountsDomainUseCase(t)
	storageClassUseCase := newMockValidateDoguStorageClassDomainUseCase(t)
	useCase := NewBlueprintSpecValidationUseCase(repoMock, dependencyUseCase, mountsUseCase, storageClassUseCase)

	dependencyUseCase.EXPECT().ValidateDependenciesForAllDogus(ctx, mock.Anything).Return(nil)
	mountsUseCase.EXPECT().ValidateAdditionalMounts(ctx, mock.Anything).Return(nil)
	storageClassUseCase.EXPECT().ValidateDoguStorageClass(ctx, mock.Anything).Return(nil)

	repoMock.EXPECT().Update(ctx, blueprint).Return(nil)

	// when
	err := useCase.ValidateBlueprintSpecDynamically(ctx, blueprint)

	// then
	require.NoError(t, err)
	assert.True(t, meta.IsStatusConditionTrue(blueprint.Conditions, domain.ConditionValid))
}

func TestBlueprintSpecUseCase_ValidateBlueprintSpecDynamically_invalid(t *testing.T) {
	// given
	repoMock := newMockBlueprintSpecRepository(t)
	ctx := context.Background()
	dependencyUseCase := newMockValidateDependenciesDomainUseCase(t)
	mountsUseCase := newMockValidateAdditionalMountsDomainUseCase(t)
	storageClassUseCase := newMockValidateDoguStorageClassDomainUseCase(t)
	useCase := NewBlueprintSpecValidationUseCase(repoMock, dependencyUseCase, mountsUseCase, storageClassUseCase)

	version, _ := core.ParseVersion("1.0.0-1")
	blueprint := &domain.BlueprintSpec{
		Id: "testBlueprint1",
		EffectiveBlueprint: domain.EffectiveBlueprint{Dogus: []domain.Dogu{{
			Name:    redmineQualifiedDoguName,
			Version: &version,
			Absent:  false,
		}}},
		Conditions: []domain.Condition{},
	}
	invalidDependencyError := errors.New("invalid dependencies")
	invalidMountsError := errors.New("invalid mounts")
	invalidStorageClassError := errors.New("invalid storage class")
	dependencyUseCase.EXPECT().ValidateDependenciesForAllDogus(ctx, mock.Anything).Return(invalidDependencyError)
	mountsUseCase.EXPECT().ValidateAdditionalMounts(ctx, mock.Anything).Return(invalidMountsError)
	storageClassUseCase.EXPECT().ValidateDoguStorageClass(ctx, mock.Anything).Return(invalidStorageClassError)
	repoMock.EXPECT().Update(ctx, blueprint).Return(nil)

	// when
	err := useCase.ValidateBlueprintSpecDynamically(ctx, blueprint)

	// then
	assert.True(t, meta.IsStatusConditionFalse(blueprint.Conditions, domain.ConditionValid))

	require.Error(t, err)
	var invalidError *domain.InvalidBlueprintError
	assert.ErrorAs(t, err, &invalidError)
	assert.ErrorIs(t, err, invalidDependencyError)
	assert.ErrorIs(t, err, invalidMountsError)
	assert.ErrorContains(t, err, "blueprint spec is invalid")

	assert.Equal(t, "testBlueprint1", blueprint.Id)
	require.Equal(t, 1, len(blueprint.Events))
	assert.IsType(t, domain.BlueprintSpecInvalidEvent{}, blueprint.Events[0])
	assert.ErrorContains(t, blueprint.Events[0].(domain.BlueprintSpecInvalidEvent).ValidationError, "blueprint spec is invalid: ")
}
