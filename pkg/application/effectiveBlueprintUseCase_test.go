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
	blueprint := &domain.BlueprintSpec{
		Id:     "testBlueprint1",
		Status: domain.StatusPhaseValidated,
	}

	repoMock := newMockBlueprintSpecRepository(t)
	ctx := context.Background()
	useCase := NewEffectiveBlueprintUseCase(repoMock)

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
	err := useCase.CalculateEffectiveBlueprint(ctx, blueprint)

	// then
	require.NoError(t, err)
}

func TestBlueprintSpecUseCase_CalculateEffectiveBlueprint_repoError(t *testing.T) {
	t.Run("cannot save", func(t *testing.T) {
		//given
		blueprint := &domain.BlueprintSpec{
			Id:     "testBlueprint1",
			Status: domain.StatusPhaseValidated,
		}

		repoMock := newMockBlueprintSpecRepository(t)
		ctx := context.Background()
		useCase := NewEffectiveBlueprintUseCase(repoMock)

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
		err := useCase.CalculateEffectiveBlueprint(ctx, blueprint)

		//then
		require.Error(t, err)
		var errorToCheck *domainservice.InternalError
		assert.ErrorAs(t, err, &errorToCheck)
		assert.ErrorContains(t, err, "cannot save blueprint spec after calculating the effective blueprint: test-error")
	})
}
