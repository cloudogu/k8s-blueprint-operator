package restartcr

import (
	"context"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	v2 "github.com/cloudogu/k8s-dogu-operator/v2/api/v2"
	"github.com/stretchr/testify/assert"
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
		expectedDoguRestartToCreate := &v2.DoguRestart{ObjectMeta: metav1.ObjectMeta{GenerateName: "testdogu-"}, Spec: v2.DoguRestartSpec{DoguName: "testdogu"}}

		mockDoguRestartInterface.EXPECT().Create(testContext, expectedDoguRestartToCreate, metav1.CreateOptions{}).Return(nil, nil)

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
		expectedDoguRestartToCreate := &v2.DoguRestart{ObjectMeta: metav1.ObjectMeta{GenerateName: "testdogu-"}, Spec: v2.DoguRestartSpec{DoguName: "testdogu"}}

		mockDoguRestartInterface.EXPECT().Create(testContext, expectedDoguRestartToCreate, metav1.CreateOptions{}).Return(nil, assert.AnError)

		restartRepository := NewDoguRestartRepository(mockDoguRestartInterface)

		// when
		err := restartRepository.RestartAll(testContext, dogusThatNeedARestart)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
	})
}
