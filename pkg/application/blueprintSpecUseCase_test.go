package application

import (
	"context"
	"errors"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBlueprintSpecUseCase_ValidateBlueprintSpecStatically_ok(t *testing.T) {
	//given
	repoMock := newMockBlueprintSpecRepository(t)
	registryMock := newMockRemoteDoguRegistry(t)
	ctx := context.Background()
	validateUseCase := domainservice.NewValidateDependenciesDomainUseCase(registryMock)
	useCase := NewBlueprintSpecUseCase(repoMock, validateUseCase, nil)

	repoMock.EXPECT().GetById(ctx, "testBlueprint1").Return(domain.BlueprintSpec{
		Id:     "testBlueprint1",
		Status: domain.StatusPhaseNew,
	}, nil)
	repoMock.EXPECT().Update(ctx, domain.BlueprintSpec{
		Id:     "testBlueprint1",
		Status: domain.StatusPhaseValidated,
		Events: []interface{}{domain.BlueprintSpecValidatedEvent{}},
	}).Return(nil)

	//when
	err := useCase.ValidateBlueprintSpecStatically(ctx, "testBlueprint1")

	//then
	repoMock.Test(t)
	require.NoError(t, err)
}

func TestBlueprintSpecUseCase_ValidateBlueprintSpecStatically_invalid(t *testing.T) {
	//given
	repoMock := newMockBlueprintSpecRepository(t)
	registryMock := newMockRemoteDoguRegistry(t)
	ctx := context.Background()
	validateUseCase := domainservice.NewValidateDependenciesDomainUseCase(registryMock)
	useCase := NewBlueprintSpecUseCase(repoMock, validateUseCase, nil)

	repoMock.EXPECT().GetById(ctx, "testBlueprint1").Return(domain.BlueprintSpec{
		Id:     "",
		Status: domain.StatusPhaseNew,
	}, nil)
	repoMock.EXPECT().Update(ctx, mock.MatchedBy(func(i interface{}) bool {
		spec := i.(domain.BlueprintSpec)
		return spec.Status == domain.StatusPhaseInvalid
	})).Return(nil)

	//when
	err := useCase.ValidateBlueprintSpecStatically(ctx, "testBlueprint1")

	//then
	repoMock.Test(t)
	require.Error(t, err)
	assert.ErrorContains(t, err, "blueprint spec is invalid: blueprint spec don't have an ID")
}

func TestBlueprintSpecUseCase_ValidateBlueprintSpecStatically_repoError(t *testing.T) {

	t.Run("error while loading blueprint spec", func(t *testing.T) {
		//given
		repoMock := newMockBlueprintSpecRepository(t)
		registryMock := newMockRemoteDoguRegistry(t)
		ctx := context.Background()
		validateUseCase := domainservice.NewValidateDependenciesDomainUseCase(registryMock)
		useCase := NewBlueprintSpecUseCase(repoMock, validateUseCase, nil)

		repoMock.EXPECT().GetById(ctx, "testBlueprint1").Return(domain.BlueprintSpec{}, errors.New("test-error"))
		//when
		err := useCase.ValidateBlueprintSpecStatically(ctx, "testBlueprint1")

		//then
		repoMock.Test(t)
		require.Error(t, err)
		assert.ErrorContains(t, err, "cannot load blueprint spec to validate it: test-error")
	})

	t.Run("error while saving blueprint spec", func(t *testing.T) {
		//given
		repoMock := newMockBlueprintSpecRepository(t)
		registryMock := newMockRemoteDoguRegistry(t)
		ctx := context.Background()
		validateUseCase := domainservice.NewValidateDependenciesDomainUseCase(registryMock)
		useCase := NewBlueprintSpecUseCase(repoMock, validateUseCase, nil)

		repoMock.EXPECT().GetById(ctx, "testBlueprint1").Return(domain.BlueprintSpec{
			Id:     "testBlueprint1",
			Status: domain.StatusPhaseNew,
		}, nil)
		repoMock.EXPECT().Update(ctx, mock.Anything).Return(errors.New("test-error"))

		//when
		err := useCase.ValidateBlueprintSpecStatically(ctx, "testBlueprint1")

		//then
		repoMock.Test(t)
		require.Error(t, err)
		assert.ErrorContains(t, err, "cannot update blueprint spec after validation: test-error")
	})

}

func TestBlueprintSpecUseCase_ValidateBlueprintSpecDynamically_ok(t *testing.T) {
	//given
	repoMock := newMockBlueprintSpecRepository(t)
	registryMock := newMockRemoteDoguRegistry(t)
	ctx := context.Background()
	validateUseCase := domainservice.NewValidateDependenciesDomainUseCase(registryMock)
	useCase := NewBlueprintSpecUseCase(repoMock, validateUseCase, nil)

	registryMock.EXPECT().GetDogus([]domainservice.DoguToLoad{}).Return(map[string]*core.Dogu{}, nil)

	repoMock.EXPECT().GetById(ctx, "testBlueprint1").Return(domain.BlueprintSpec{
		Id:     "testBlueprint1",
		Status: domain.StatusPhaseValidated,
	}, nil)
	// at the moment the state won't get persisted
	//repoMock.EXPECT().Update(ctx, domain.BlueprintSpec{
	//	Id:                   "testBlueprint1",
	//	Blueprint:            domain.Blueprint{},
	//	BlueprintMask:        domain.BlueprintMask{},
	//	EffectiveBlueprint:   domain.EffectiveBlueprint{},
	//	StateDiff:            domain.StateDiff{},
	//	BlueprintUpgradePlan: domain.BlueprintUpgradePlan{},
	//	Status:               domain.StatusPhaseValidated,
	//	Events:               []interface{}{},
	//}).Return(nil)

	//when
	err := useCase.ValidateBlueprintSpecDynamically(ctx, "testBlueprint1")

	//then
	repoMock.Test(t)
	registryMock.Test(t)
	require.NoError(t, err)
}

func TestBlueprintSpecUseCase_ValidateBlueprintSpecDynamically_invalid(t *testing.T) {
	//given
	repoMock := newMockBlueprintSpecRepository(t)
	registryMock := newMockRemoteDoguRegistry(t)
	ctx := context.Background()
	validateUseCase := domainservice.NewValidateDependenciesDomainUseCase(registryMock)
	useCase := NewBlueprintSpecUseCase(repoMock, validateUseCase, nil)

	repoMock.EXPECT().GetById(ctx, "testBlueprint1").Return(domain.BlueprintSpec{
		Id: "testBlueprint1",
		EffectiveBlueprint: domain.EffectiveBlueprint{Dogus: []domain.TargetDogu{{
			Namespace:   "official",
			Name:        "redmine",
			Version:     "1.0.0-1",
			TargetState: domain.TargetStatePresent,
		}}},
		Status: domain.StatusPhaseValidated,
	}, nil)

	registryMock.EXPECT().GetDogus([]domainservice.DoguToLoad{
		{
			QualifiedDoguName: "official/redmine",
			Version:           "1.0.0-1",
		},
	}).Return(nil, errors.New("dogu not found for testing"))

	//when
	err := useCase.ValidateBlueprintSpecDynamically(ctx, "testBlueprint1")

	//then
	repoMock.Test(t)
	registryMock.Test(t)
	require.Error(t, err)
	assert.ErrorContains(t, err, "blueprint spec is invalid: cannot load or cannot find dogu specifications from remote registry for dogu dependency validation: dogu not found for testing")
}

func TestBlueprintSpecUseCase_ValidateBlueprintSpecDynamically_repoError(t *testing.T) {
	//given
	repoMock := newMockBlueprintSpecRepository(t)
	registryMock := newMockRemoteDoguRegistry(t)
	ctx := context.Background()
	validateUseCase := domainservice.NewValidateDependenciesDomainUseCase(registryMock)
	useCase := NewBlueprintSpecUseCase(repoMock, validateUseCase, nil)

	repoMock.EXPECT().GetById(ctx, "testBlueprint1").Return(domain.BlueprintSpec{}, errors.New("test-error"))

	//when
	err := useCase.ValidateBlueprintSpecDynamically(ctx, "testBlueprint1")

	//then
	repoMock.Test(t)
	registryMock.Test(t)
	require.Error(t, err)
	assert.ErrorContains(t, err, "cannot load blueprint spec to validate it: test-error")
}
