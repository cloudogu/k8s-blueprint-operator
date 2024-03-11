package restartcr

import (
	"context"
	"errors"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	v1 "github.com/cloudogu/k8s-dogu-operator/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func Test_doguRestartRepository_RestartAll(t *testing.T) {
	t.Run("no error on restart all", func(t *testing.T) {
		// given
		testContext := context.Background()
		testDoguSimpleName := common.SimpleDoguName("testdogu")
		dogusThatNeedARestart := []common.SimpleDoguName{testDoguSimpleName}
		mockDoguRestartInterface := NewMockDoguRestartInterface(t)

		mockDoguRestartInterface.EXPECT().Create(testContext, mock.Anything, metav1.CreateOptions{}).Return(&v1.DoguRestart{
			ObjectMeta: metav1.ObjectMeta{
				Name: string(testDoguSimpleName),
			}}, nil).Run(func(ctx context.Context, dogu *v1.DoguRestart, opts metav1.CreateOptions) {
			assert.Contains(t, dogu.Name, testDoguSimpleName)
		})

		restartRepository := NewDoguRestartRepository(mockDoguRestartInterface)

		// when
		err := restartRepository.RestartAll(testContext, dogusThatNeedARestart)

		// then
		require.NoError(t, err)
	})

	t.Run("no error on empty restart all", func(t *testing.T) {
		// given
		testContext := context.Background()
		dogusThatNeedARestart := []common.SimpleDoguName{}
		mockDoguRestartInterface := NewMockDoguRestartInterface(t)

		restartRepository := NewDoguRestartRepository(mockDoguRestartInterface)

		// when
		err := restartRepository.RestartAll(testContext, dogusThatNeedARestart)

		// then
		require.NoError(t, err)
	})

	t.Run("fail on error at create", func(t *testing.T) {
		// given
		testContext := context.Background()
		testDoguSimpleName := common.SimpleDoguName("testdogu")
		dogusThatNeedARestart := []common.SimpleDoguName{testDoguSimpleName}
		mockDoguRestartInterface := NewMockDoguRestartInterface(t)

		mockDoguRestartInterface.EXPECT().Create(testContext, mock.Anything, metav1.CreateOptions{}).Return(&v1.DoguRestart{
			ObjectMeta: metav1.ObjectMeta{
				Name: string(testDoguSimpleName),
			}}, errors.New("testerror")).Run(func(ctx context.Context, dogu *v1.DoguRestart, opts metav1.CreateOptions) {
			assert.Contains(t, dogu.Name, testDoguSimpleName)
		})

		restartRepository := NewDoguRestartRepository(mockDoguRestartInterface)

		// when
		err := restartRepository.RestartAll(testContext, dogusThatNeedARestart)

		// then
		require.Error(t, err)
		assert.Equal(t, "testerror", err.Error())
	})
}
