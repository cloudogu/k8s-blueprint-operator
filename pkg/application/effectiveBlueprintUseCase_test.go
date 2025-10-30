package application

import (
	"context"
	"testing"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
)

func TestBlueprintSpecUseCase_CalculateEffectiveBlueprint_ok(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		// given
		blueprint := &domain.BlueprintSpec{
			Id: "testBlueprint1",
		}

		repoMock := newMockBlueprintSpecRepository(t)
		ctx := context.Background()
		useCase := NewEffectiveBlueprintUseCase(repoMock)

		repoMock.EXPECT().Update(ctx, blueprint).Return(nil)

		// when
		err := useCase.CalculateEffectiveBlueprint(ctx, blueprint)

		// then
		require.NoError(t, err)
		assert.Equal(t, 0, len(blueprint.Events))
		assert.Equal(t, blueprint.EffectiveBlueprint, domain.EffectiveBlueprint{})
	})

	t.Run("should throw error on update error", func(t *testing.T) {
		//given
		blueprint := &domain.BlueprintSpec{
			Id: "testBlueprint1",
		}

		repoMock := newMockBlueprintSpecRepository(t)
		ctx := context.Background()
		useCase := NewEffectiveBlueprintUseCase(repoMock)

		repoMock.EXPECT().Update(ctx, blueprint).Return(&domainservice.InternalError{Message: "test-error"})

		//when
		err := useCase.CalculateEffectiveBlueprint(ctx, blueprint)

		//then
		require.Error(t, err)
		var errorToCheck *domainservice.InternalError
		assert.ErrorAs(t, err, &errorToCheck)
		assert.ErrorContains(t, err, "cannot save blueprint spec after calculating the effective blueprint: test-error")
	})
}
