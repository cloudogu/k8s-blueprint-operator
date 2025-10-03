package application

import (
	"context"
	"testing"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
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
		debugModeRepoMock := newMockDebugModeRepository(t)
		ctx := context.Background()
		useCase := NewEffectiveBlueprintUseCase(repoMock, debugModeRepoMock)

		debugModeRepoMock.EXPECT().GetSingleton(ctx).Return(nil, nil)
		repoMock.EXPECT().Update(ctx, blueprint).Return(nil)

		// when
		err := useCase.CalculateEffectiveBlueprint(ctx, blueprint)

		// then
		require.NoError(t, err)
		assert.Equal(t, 0, len(blueprint.Events))
		assert.Equal(t, blueprint.EffectiveBlueprint, domain.EffectiveBlueprint{})
	})

	t.Run("should ignore error on debug mode repo not found error", func(t *testing.T) {
		//given
		blueprint := &domain.BlueprintSpec{
			Id: "testBlueprint1",
		}

		repoMock := newMockBlueprintSpecRepository(t)
		debugModeRepoMock := newMockDebugModeRepository(t)
		ctx := context.Background()
		useCase := NewEffectiveBlueprintUseCase(repoMock, debugModeRepoMock)

		debugModeRepoMock.EXPECT().GetSingleton(ctx).Return(nil, domainservice.NewNotFoundError(assert.AnError, "test-error"))
		repoMock.EXPECT().Update(ctx, blueprint).Return(nil)

		// when
		err := useCase.CalculateEffectiveBlueprint(ctx, blueprint)

		// then
		require.NoError(t, err)
		assert.Equal(t, 0, len(blueprint.Events))
		assert.Equal(t, blueprint.EffectiveBlueprint, domain.EffectiveBlueprint{})
	})

	t.Run("should ignore loglevel configs on active debug mode", func(t *testing.T) {
		//given
		blueprint := &domain.BlueprintSpec{
			Id: "testBlueprint1",
			Blueprint: domain.Blueprint{
				Dogus: []domain.Dogu{{Name: cescommons.QualifiedName{SimpleName: "testDogu", Namespace: "test"}}},
				Config: domain.Config{
					Dogus: domain.DoguConfig{
						"testDogu": {
							{
								Key:    "logging/root",
								Absent: true,
							},
							{
								Key:    "someKey",
								Absent: true,
							},
						},
					},
				},
			},
		}

		repoMock := newMockBlueprintSpecRepository(t)
		debugModeRepoMock := newMockDebugModeRepository(t)
		ctx := context.Background()
		useCase := NewEffectiveBlueprintUseCase(repoMock, debugModeRepoMock)

		debugMode := ecosystem.DebugMode{Phase: ecosystem.DebugModeStatusSet}

		debugModeRepoMock.EXPECT().GetSingleton(ctx).Return(&debugMode, nil)
		repoMock.EXPECT().Update(ctx, blueprint).Return(nil)

		// when
		err := useCase.CalculateEffectiveBlueprint(ctx, blueprint)

		// then
		require.NoError(t, err)
		assert.Equal(t, 0, len(blueprint.Events))
		assert.Equal(t, blueprint.EffectiveBlueprint, domain.EffectiveBlueprint{
			Dogus: []domain.Dogu{{Name: cescommons.QualifiedName{SimpleName: "testDogu", Namespace: "test"}}},
			Config: domain.Config{
				Dogus: domain.DoguConfig{
					"testDogu": {
						{
							Key:    "someKey",
							Absent: true,
						},
					},
				},
			},
		})
	})

	t.Run("should not ignore loglevel configs on deactivated debug mode", func(t *testing.T) {
		//given
		blueprint := &domain.BlueprintSpec{
			Id: "testBlueprint1",
			Blueprint: domain.Blueprint{
				Dogus: []domain.Dogu{{Name: cescommons.QualifiedName{SimpleName: "testDogu", Namespace: "test"}}},
				Config: domain.Config{
					Dogus: domain.DoguConfig{
						"testDogu": {
							{
								Key:    "logging/root",
								Absent: true,
							},
							{
								Key:    "someKey",
								Absent: true,
							},
						},
					},
				},
			},
		}

		repoMock := newMockBlueprintSpecRepository(t)
		debugModeRepoMock := newMockDebugModeRepository(t)
		ctx := context.Background()
		useCase := NewEffectiveBlueprintUseCase(repoMock, debugModeRepoMock)

		debugMode := ecosystem.DebugMode{Phase: "completed"}

		debugModeRepoMock.EXPECT().GetSingleton(ctx).Return(&debugMode, nil)
		repoMock.EXPECT().Update(ctx, blueprint).Return(nil)

		// when
		err := useCase.CalculateEffectiveBlueprint(ctx, blueprint)

		// then
		require.NoError(t, err)
		assert.Equal(t, 0, len(blueprint.Events))
		assert.Equal(t, blueprint.EffectiveBlueprint, domain.EffectiveBlueprint{
			Dogus: []domain.Dogu{{Name: cescommons.QualifiedName{SimpleName: "testDogu", Namespace: "test"}}},
			Config: domain.Config{
				Dogus: domain.DoguConfig{
					"testDogu": {
						{
							Key:    "logging/root",
							Absent: true,
						},
						{
							Key:    "someKey",
							Absent: true,
						},
					},
				},
			},
		})
	})

	t.Run("should throw error on debug mode repo error", func(t *testing.T) {
		//given
		blueprint := &domain.BlueprintSpec{
			Id: "testBlueprint1",
		}

		repoMock := newMockBlueprintSpecRepository(t)
		debugModeRepoMock := newMockDebugModeRepository(t)
		ctx := context.Background()
		useCase := NewEffectiveBlueprintUseCase(repoMock, debugModeRepoMock)

		debugModeRepoMock.EXPECT().GetSingleton(ctx).Return(nil, &domainservice.InternalError{Message: "test-error"})

		//when
		err := useCase.CalculateEffectiveBlueprint(ctx, blueprint)

		//then
		require.Error(t, err)
		var errorToCheck *domainservice.InternalError
		assert.ErrorAs(t, err, &errorToCheck)
		assert.ErrorContains(t, err, "cannot calculate effective blueprint due to an error when loading the debug mode cr")
	})

	t.Run("should throw error on update error", func(t *testing.T) {
		//given
		blueprint := &domain.BlueprintSpec{
			Id: "testBlueprint1",
		}

		repoMock := newMockBlueprintSpecRepository(t)
		debugModeRepoMock := newMockDebugModeRepository(t)
		ctx := context.Background()
		useCase := NewEffectiveBlueprintUseCase(repoMock, debugModeRepoMock)

		debugModeRepoMock.EXPECT().GetSingleton(ctx).Return(nil, nil)
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
