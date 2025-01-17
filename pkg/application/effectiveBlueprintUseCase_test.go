package application

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
)

func TestBlueprintSpecUseCase_CalculateEffectiveBlueprint_ok(t *testing.T) {
	// given
	repoMock := newMockBlueprintSpecRepository(t)
	ctx := context.Background()
	useCase := NewEffectiveBlueprintUseCase(repoMock)

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
		Status:             domain.StatusPhaseEffectiveBlueprintGenerated,
		Events:             []domain.Event{domain.EffectiveBlueprintCalculatedEvent{}},
	}).Return(nil)

	// when
	err := useCase.CalculateEffectiveBlueprint(ctx, "testBlueprint1")

	// then
	require.NoError(t, err)
}

func TestBlueprintSpecUseCase_CalculateEffectiveBlueprint_repoError(t *testing.T) {
	t.Run("blueprint spec not found", func(t *testing.T) {
		//given
		repoMock := newMockBlueprintSpecRepository(t)
		ctx := context.Background()
		useCase := NewEffectiveBlueprintUseCase(repoMock)

		repoMock.EXPECT().GetById(ctx, "testBlueprint1").Return(nil, &domainservice.NotFoundError{Message: "test-error"})

		//when
		err := useCase.CalculateEffectiveBlueprint(ctx, "testBlueprint1")

		//then
		require.Error(t, err)
		var errorToCheck *domainservice.NotFoundError
		assert.ErrorAs(t, err, &errorToCheck)
		assert.ErrorContains(t, err, "cannot load blueprint spec to calculate effective blueprint: test-error")
	})

	t.Run("internal error while loading", func(t *testing.T) {
		//given
		repoMock := newMockBlueprintSpecRepository(t)
		ctx := context.Background()
		useCase := NewEffectiveBlueprintUseCase(repoMock)

		repoMock.EXPECT().GetById(ctx, "testBlueprint1").Return(nil, &domainservice.InternalError{Message: "test-error"})

		//when
		err := useCase.CalculateEffectiveBlueprint(ctx, "testBlueprint1")

		//then
		require.Error(t, err)
		var errorToCheck *domainservice.InternalError
		assert.ErrorAs(t, err, &errorToCheck)
		assert.ErrorContains(t, err, "cannot load blueprint spec to calculate effective blueprint: test-error")
	})

	t.Run("cannot save", func(t *testing.T) {
		//given
		repoMock := newMockBlueprintSpecRepository(t)
		ctx := context.Background()
		useCase := NewEffectiveBlueprintUseCase(repoMock)

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
			Status:             domain.StatusPhaseEffectiveBlueprintGenerated,
			Events:             []domain.Event{domain.EffectiveBlueprintCalculatedEvent{}},
		}).Return(&domainservice.InternalError{Message: "test-error"})

		//when
		err := useCase.CalculateEffectiveBlueprint(ctx, "testBlueprint1")

		//then
		require.Error(t, err)
		var errorToCheck *domainservice.InternalError
		assert.ErrorAs(t, err, &errorToCheck)
		assert.ErrorContains(t, err, "cannot save blueprint spec after calculating the effective blueprint: test-error")
	})
}
