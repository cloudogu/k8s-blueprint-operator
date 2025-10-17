package application

import (
	"testing"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDogusUpToDateUseCase_CheckDogus(t *testing.T) {
	t.Run("all dogus up to date", func(t *testing.T) {
		blueprint := &domain.BlueprintSpec{
			Conditions: []domain.Condition{},
		}

		repoMock := newMockBlueprintSpecRepository(t)
		doguInstallUseCaseMock := newMockDoguInstallationUseCase(t)
		dogusNotUpToDate := []cescommons.SimpleName{}
		doguInstallUseCaseMock.EXPECT().CheckDogusUpToDate(testCtx).Return(dogusNotUpToDate, nil)
		useCase := NewDogusUpToDateUseCase(repoMock, doguInstallUseCaseMock)

		err := useCase.CheckDogus(testCtx, blueprint)

		require.NoError(t, err)
		assert.Empty(t, blueprint.Events)
	})
	t.Run("multiple dogus not up to date", func(t *testing.T) {
		blueprint := &domain.BlueprintSpec{
			Conditions: []domain.Condition{},
		}

		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, blueprint).Return(nil)
		doguInstallUseCaseMock := newMockDoguInstallationUseCase(t)
		dogusNotUpToDate := []cescommons.SimpleName{ldap, postfix}
		doguInstallUseCaseMock.EXPECT().CheckDogusUpToDate(testCtx).Return(dogusNotUpToDate, nil)
		useCase := NewDogusUpToDateUseCase(repoMock, doguInstallUseCaseMock)

		err := useCase.CheckDogus(testCtx, blueprint)

		require.Error(t, err)
		var expectedErrorType *domain.DogusNotUpToDateError
		assert.ErrorAs(t, err, &expectedErrorType)
		assert.ErrorContains(t, err, "following dogus are not up to date yet:")
		assert.ErrorContains(t, err, ldap.String())
		assert.ErrorContains(t, err, postfix.String())
		require.Equal(t, 1, len(blueprint.Events))
		assert.Equal(t, domain.DogusNotUpToDateEvent{DogusNotUpToDate: dogusNotUpToDate}, blueprint.Events[0])
		require.Empty(t, blueprint.Conditions)
	})

	t.Run("no update without not up to date dogus", func(t *testing.T) {
		blueprint := &domain.BlueprintSpec{
			Conditions: []domain.Condition{},
		}

		repoMock := newMockBlueprintSpecRepository(t)
		doguInstallUseCaseMock := newMockDoguInstallationUseCase(t)
		doguInstallUseCaseMock.EXPECT().CheckDogusUpToDate(testCtx).Return([]cescommons.SimpleName{}, nil)
		useCase := NewDogusUpToDateUseCase(repoMock, doguInstallUseCaseMock)

		err := useCase.CheckDogus(testCtx, blueprint)
		require.NoError(t, err)
		require.Empty(t, blueprint.Events)
		require.Empty(t, blueprint.Conditions)
	})

	t.Run("fail to check dogus", func(t *testing.T) {
		blueprint := &domain.BlueprintSpec{
			Conditions: []domain.Condition{},
		}

		repoMock := newMockBlueprintSpecRepository(t)
		doguInstallUseCaseMock := newMockDoguInstallationUseCase(t)
		doguInstallUseCaseMock.EXPECT().CheckDogusUpToDate(testCtx).Return(nil, assert.AnError)
		useCase := NewDogusUpToDateUseCase(repoMock, doguInstallUseCaseMock)

		err := useCase.CheckDogus(testCtx, blueprint)

		require.ErrorIs(t, err, assert.AnError)
		require.Equal(t, 0, len(blueprint.Events))
	})

	t.Run("fail to update blueprint", func(t *testing.T) {
		blueprint := &domain.BlueprintSpec{
			Conditions: []domain.Condition{},
			StateDiff: domain.StateDiff{
				DoguDiffs: domain.DoguDiffs{
					{NeededActions: []domain.Action{domain.ActionInstall}},
				},
			},
		}

		repoMock := newMockBlueprintSpecRepository(t)
		repoMock.EXPECT().Update(testCtx, blueprint).Return(assert.AnError)
		doguInstallUseCaseMock := newMockDoguInstallationUseCase(t)
		dogusNotUpToDate := []cescommons.SimpleName{"ldap"}
		doguInstallUseCaseMock.EXPECT().CheckDogusUpToDate(testCtx).Return(dogusNotUpToDate, nil)
		useCase := NewDogusUpToDateUseCase(repoMock, doguInstallUseCaseMock)

		err := useCase.CheckDogus(testCtx, blueprint)

		require.ErrorIs(t, err, assert.AnError)
		require.Equal(t, 1, len(blueprint.Events))
		assert.Equal(t, domain.DogusNotUpToDateEvent{DogusNotUpToDate: dogusNotUpToDate}, blueprint.Events[0])
	})
}
