package application

import (
	"context"
	"errors"
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

var redmineQualifiedDoguName = cescommons.QualifiedName{
	Namespace:  "official",
	SimpleName: "redmine",
}

func TestBlueprintSpecUseCase_ValidateBlueprintSpecStatically_ok(t *testing.T) {
	//given
	repoMock := newMockBlueprintSpecRepository(t)
	ctx := context.Background()
	DependencyUseCase := newMockValidateDependenciesDomainUseCase(t)
	MountsUseCase := newMockValidateAdditionalMountsDomainUseCase(t)
	useCase := NewBlueprintSpecValidationUseCase(repoMock, DependencyUseCase, MountsUseCase)

	repoMock.EXPECT().GetById(ctx, "testBlueprint1").Return(&domain.BlueprintSpec{
		Id:     "testBlueprint1",
		Status: domain.StatusPhaseNew,
	}, nil)
	repoMock.EXPECT().Update(ctx, &domain.BlueprintSpec{
		Id:     "testBlueprint1",
		Status: domain.StatusPhaseStaticallyValidated,
		Events: []domain.Event{domain.BlueprintSpecStaticallyValidatedEvent{}},
	}).Return(nil)

	//when
	err := useCase.ValidateBlueprintSpecStatically(ctx, "testBlueprint1")

	//then
	require.NoError(t, err)
}

func TestBlueprintSpecUseCase_ValidateBlueprintSpecStatically_invalid(t *testing.T) {
	//given
	repoMock := newMockBlueprintSpecRepository(t)
	ctx := context.Background()
	DependencyUseCase := newMockValidateDependenciesDomainUseCase(t)
	MountsUseCase := newMockValidateAdditionalMountsDomainUseCase(t)
	useCase := NewBlueprintSpecValidationUseCase(repoMock, DependencyUseCase, MountsUseCase)

	repoMock.EXPECT().GetById(ctx, "testBlueprint1").Return(&domain.BlueprintSpec{
		Id:     "",
		Status: domain.StatusPhaseNew,
	}, nil)
	repoMock.EXPECT().Update(ctx, mock.MatchedBy(func(i interface{}) bool {
		spec := i.(*domain.BlueprintSpec)
		return spec.Status == domain.StatusPhaseInvalid
	})).Return(nil)

	//when
	err := useCase.ValidateBlueprintSpecStatically(ctx, "testBlueprint1")

	//then
	require.Error(t, err)
	var invalidError *domain.InvalidBlueprintError
	assert.ErrorAs(t, err, &invalidError, "error should be an InvalidBlueprintError")
	assert.ErrorContains(t, err, "blueprint spec is invalid: blueprint spec doesn't have an ID")
}

func TestBlueprintSpecUseCase_ValidateBlueprintSpecStatically_repoError(t *testing.T) {

	t.Run("blueprint spec not found while loading", func(t *testing.T) {
		//given
		repoMock := newMockBlueprintSpecRepository(t)
		ctx := context.Background()
		DependencyUseCase := newMockValidateDependenciesDomainUseCase(t)
		MountsUseCase := newMockValidateAdditionalMountsDomainUseCase(t)
		useCase := NewBlueprintSpecValidationUseCase(repoMock, DependencyUseCase, MountsUseCase)

		repoMock.EXPECT().GetById(ctx, "testBlueprint1").Return(nil, &domainservice.NotFoundError{Message: "test-error"})
		//when
		err := useCase.ValidateBlueprintSpecStatically(ctx, "testBlueprint1")

		//then
		require.Error(t, err)
		var invalidError *domainservice.NotFoundError
		assert.ErrorAs(t, err, &invalidError)
		assert.ErrorContains(t, err, "cannot load blueprint spec to validate it: test-error")
	})

	t.Run("cannot parse blueprint in repository", func(t *testing.T) {
		//given
		repoMock := newMockBlueprintSpecRepository(t)
		ctx := context.Background()
		DependencyUseCase := newMockValidateDependenciesDomainUseCase(t)
		MountsUseCase := newMockValidateAdditionalMountsDomainUseCase(t)
		useCase := NewBlueprintSpecValidationUseCase(repoMock, DependencyUseCase, MountsUseCase)
		invalidError := domain.InvalidBlueprintError{Message: "test-error"}
		var events []domain.Event
		events = append(events, domain.BlueprintSpecInvalidEvent{ValidationError: &invalidError})
		repoMock.EXPECT().GetById(ctx, "testBlueprint1").Return(&domain.BlueprintSpec{Id: "testBlueprint1"}, &invalidError)
		repoMock.EXPECT().Update(ctx, &domain.BlueprintSpec{Id: "testBlueprint1", Status: domain.StatusPhaseInvalid, Events: events}).Return(nil)
		//when
		err := useCase.ValidateBlueprintSpecStatically(ctx, "testBlueprint1")

		//then
		require.Error(t, err)
		var expectedErrorType *domain.InvalidBlueprintError
		assert.ErrorAs(t, err, &expectedErrorType)
		assert.ErrorContains(t, err, "blueprint spec syntax is invalid: test-error")
	})

	t.Run("internal error while loading blueprint spec", func(t *testing.T) {
		//given
		repoMock := newMockBlueprintSpecRepository(t)
		ctx := context.Background()
		DependencyUseCase := newMockValidateDependenciesDomainUseCase(t)
		MountsUseCase := newMockValidateAdditionalMountsDomainUseCase(t)
		useCase := NewBlueprintSpecValidationUseCase(repoMock, DependencyUseCase, MountsUseCase)

		repoMock.EXPECT().GetById(ctx, "testBlueprint1").Return(nil, &domainservice.InternalError{Message: "test-error"})
		//when
		err := useCase.ValidateBlueprintSpecStatically(ctx, "testBlueprint1")

		//then
		require.Error(t, err)
		var invalidError *domainservice.InternalError
		assert.ErrorAs(t, err, &invalidError)
		assert.ErrorContains(t, err, "cannot load blueprint spec to validate it: test-error")
	})

	t.Run("error while saving blueprint spec", func(t *testing.T) {
		//given
		repoMock := newMockBlueprintSpecRepository(t)
		ctx := context.Background()
		DependencyUseCase := newMockValidateDependenciesDomainUseCase(t)
		MountsUseCase := newMockValidateAdditionalMountsDomainUseCase(t)
		useCase := NewBlueprintSpecValidationUseCase(repoMock, DependencyUseCase, MountsUseCase)

		repoMock.EXPECT().GetById(ctx, "testBlueprint1").Return(&domain.BlueprintSpec{
			Id:     "testBlueprint1",
			Status: domain.StatusPhaseNew,
		}, nil)
		repoMock.EXPECT().Update(ctx, mock.Anything).Return(&domainservice.InternalError{Message: "test-error"})

		//when
		err := useCase.ValidateBlueprintSpecStatically(ctx, "testBlueprint1")

		//then
		require.Error(t, err)
		var invalidError *domainservice.InternalError
		assert.ErrorAs(t, err, &invalidError)
		assert.ErrorContains(t, err, "cannot update blueprint spec after static validation: test-error")
	})

}

func TestBlueprintSpecUseCase_ValidateBlueprintSpecDynamically_ok(t *testing.T) {
	// given
	repoMock := newMockBlueprintSpecRepository(t)
	ctx := context.Background()
	DependencyUseCase := newMockValidateDependenciesDomainUseCase(t)
	MountsUseCase := newMockValidateAdditionalMountsDomainUseCase(t)
	useCase := NewBlueprintSpecValidationUseCase(repoMock, DependencyUseCase, MountsUseCase)

	DependencyUseCase.EXPECT().ValidateDependenciesForAllDogus(ctx, mock.Anything).Return(nil)
	MountsUseCase.EXPECT().ValidateAdditionalMounts(ctx, mock.Anything).Return(nil)

	repoMock.EXPECT().GetById(ctx, "testBlueprint1").Return(&domain.BlueprintSpec{
		Id:     "testBlueprint1",
		Status: domain.StatusPhaseValidated,
	}, nil)
	repoMock.EXPECT().Update(ctx, &domain.BlueprintSpec{
		Id:                 "testBlueprint1",
		Blueprint:          domain.Blueprint{},
		BlueprintMask:      domain.BlueprintMask{},
		EffectiveBlueprint: domain.EffectiveBlueprint{},
		StateDiff:          domain.StateDiff{},
		Status:             domain.StatusPhaseValidated,
		Events:             []domain.Event{domain.BlueprintSpecValidatedEvent{}},
	}).Return(nil)

	// when
	err := useCase.ValidateBlueprintSpecDynamically(ctx, "testBlueprint1")

	// then
	require.NoError(t, err)
}

func TestBlueprintSpecUseCase_ValidateBlueprintSpecDynamically_invalid(t *testing.T) {
	// given
	repoMock := newMockBlueprintSpecRepository(t)
	ctx := context.Background()
	DependencyUseCase := newMockValidateDependenciesDomainUseCase(t)
	MountsUseCase := newMockValidateAdditionalMountsDomainUseCase(t)
	useCase := NewBlueprintSpecValidationUseCase(repoMock, DependencyUseCase, MountsUseCase)

	version, _ := core.ParseVersion("1.0.0-1")
	blueprintSpec := &domain.BlueprintSpec{
		Id: "testBlueprint1",
		EffectiveBlueprint: domain.EffectiveBlueprint{Dogus: []domain.Dogu{{
			Name:        redmineQualifiedDoguName,
			Version:     version,
			TargetState: domain.TargetStatePresent,
		}}},
		Status: domain.StatusPhaseValidated,
	}
	repoMock.EXPECT().GetById(ctx, "testBlueprint1").Return(blueprintSpec, nil)
	invalidDependencyError := errors.New("invalid dependencies")
	invalidMountsError := errors.New("invalid mounts")
	DependencyUseCase.EXPECT().ValidateDependenciesForAllDogus(ctx, mock.Anything).Return(invalidDependencyError)
	MountsUseCase.EXPECT().ValidateAdditionalMounts(ctx, mock.Anything).Return(invalidMountsError)
	repoMock.EXPECT().Update(ctx, mock.Anything).Return(nil)

	// when
	err := useCase.ValidateBlueprintSpecDynamically(ctx, "testBlueprint1")

	// then
	require.Error(t, err)
	var invalidError *domain.InvalidBlueprintError
	assert.ErrorAs(t, err, &invalidError)
	assert.ErrorIs(t, err, invalidDependencyError)
	assert.ErrorIs(t, err, invalidMountsError)
	assert.ErrorContains(t, err, "blueprint spec is invalid")

	assert.Equal(t, "testBlueprint1", blueprintSpec.Id)
	assert.Equal(t, domain.StatusPhaseInvalid, blueprintSpec.Status)
	require.Equal(t, 1, len(blueprintSpec.Events))
	assert.IsType(t, domain.BlueprintSpecInvalidEvent{}, blueprintSpec.Events[0])
	assert.ErrorContains(t, blueprintSpec.Events[0].(domain.BlueprintSpecInvalidEvent).ValidationError, "blueprint spec is invalid: ")
}

func TestBlueprintSpecUseCase_ValidateBlueprintSpecDynamically_repoError(t *testing.T) {
	t.Run("internal error", func(t *testing.T) {
		// given
		repoMock := newMockBlueprintSpecRepository(t)
		ctx := context.Background()
		DependencyUseCase := newMockValidateDependenciesDomainUseCase(t)
		MountsUseCase := newMockValidateAdditionalMountsDomainUseCase(t)
		useCase := NewBlueprintSpecValidationUseCase(repoMock, DependencyUseCase, MountsUseCase)

		repoMock.EXPECT().GetById(ctx, "testBlueprint1").Return(nil, &domainservice.InternalError{Message: "test-error"})

		// when
		err := useCase.ValidateBlueprintSpecDynamically(ctx, "testBlueprint1")

		// then
		require.Error(t, err)
		var invalidError *domainservice.InternalError
		assert.ErrorAs(t, err, &invalidError)
		assert.ErrorContains(t, err, "cannot load blueprint spec to validate it: test-error")
	})
	t.Run("not found error", func(t *testing.T) {
		// given
		repoMock := newMockBlueprintSpecRepository(t)
		ctx := context.Background()
		DependencyUseCase := newMockValidateDependenciesDomainUseCase(t)
		MountsUseCase := newMockValidateAdditionalMountsDomainUseCase(t)
		useCase := NewBlueprintSpecValidationUseCase(repoMock, DependencyUseCase, MountsUseCase)

		repoMock.EXPECT().GetById(ctx, "testBlueprint1").Return(nil, &domainservice.NotFoundError{Message: "test-error"})

		// when
		err := useCase.ValidateBlueprintSpecDynamically(ctx, "testBlueprint1")

		// then
		require.Error(t, err)
		var invalidError *domainservice.NotFoundError
		assert.ErrorAs(t, err, &invalidError)
		assert.ErrorContains(t, err, "cannot load blueprint spec to validate it: test-error")
	})
}
