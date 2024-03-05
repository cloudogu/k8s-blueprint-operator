package dogucr

import (
	"context"
	"errors"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/common"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	v1 "github.com/cloudogu/k8s-dogu-operator/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"testing"
)

var version3_2_1_4, _ = core.ParseVersion("3.2.1-4")

var crResourceVersion = "abc"
var persistenceContext = map[string]interface{}{
	doguInstallationRepoContextKey: doguInstallationRepoContext{resourceVersion: crResourceVersion},
}

var testCtx = context.Background()

func Test_doguInstallationRepo_GetByName(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		// given
		doguClientMock := NewMockDoguInterface(t)
		pvcClientMock := NewMockPvcInterface(t)
		repo := NewDoguInstallationRepo(doguClientMock, pvcClientMock)
		// when
		doguClientMock.EXPECT().Get(testCtx, "postgresql", metav1.GetOptions{}).Return(
			&v1.Dogu{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name:            "postgresql",
					ResourceVersion: crResourceVersion,
				},
				Spec: v1.DoguSpec{
					Name:      "official/postgresql",
					Version:   version3_2_1_4.Raw,
					Resources: v1.DoguResources{},
					UpgradeConfig: v1.UpgradeConfig{
						AllowNamespaceSwitch: false,
					},
				},
				Status: v1.DoguStatus{
					Status: v1.DoguStatusInstalled,
					Health: v1.AvailableHealthStatus,
				},
			}, nil)

		dogu, err := repo.GetByName(testCtx, "postgresql")

		// then
		require.NoError(t, err)
		assert.Equal(t, &ecosystem.DoguInstallation{
			Name:               postgresDoguName,
			Version:            version3_2_1_4,
			Status:             ecosystem.DoguStatusInstalled,
			Health:             ecosystem.AvailableHealthStatus,
			UpgradeConfig:      ecosystem.UpgradeConfig{},
			PersistenceContext: persistenceContext,
		}, dogu)
	})

	t.Run("not found error", func(t *testing.T) {
		// given
		doguClientMock := NewMockDoguInterface(t)
		pvcClientMock := NewMockPvcInterface(t)
		repo := NewDoguInstallationRepo(doguClientMock, pvcClientMock)
		// when
		doguClientMock.EXPECT().Get(testCtx, "postgresql", metav1.GetOptions{}).Return(
			nil,
			k8sErrors.NewNotFound(schema.GroupResource{}, "postgresql"),
		)

		_, err := repo.GetByName(testCtx, "postgresql")

		// then
		require.Error(t, err)
		var expectedError *domainservice.NotFoundError
		assert.ErrorAs(t, err, &expectedError)
	})

	t.Run("internal error", func(t *testing.T) {
		// given
		doguClientMock := NewMockDoguInterface(t)
		pvcClientMock := NewMockPvcInterface(t)
		repo := NewDoguInstallationRepo(doguClientMock, pvcClientMock)
		// when
		doguClientMock.EXPECT().Get(testCtx, "postgresql", metav1.GetOptions{}).Return(
			nil,
			k8sErrors.NewInternalError(errors.New("test-error")),
		)

		_, err := repo.GetByName(testCtx, "postgresql")

		// then
		require.Error(t, err)
		var expectedError *domainservice.InternalError
		assert.ErrorAs(t, err, &expectedError)
	})
}

func Test_doguInstallationRepo_GetAll(t *testing.T) {
	t.Run("should fail to list dogus", func(t *testing.T) {
		// given
		doguClientMock := NewMockDoguInterface(t)
		doguClientMock.EXPECT().List(testCtx, metav1.ListOptions{}).Return(nil, assert.AnError)

		sut := &doguInstallationRepo{doguClient: doguClientMock}

		// when
		_, err := sut.GetAll(testCtx)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		expectedError := &domainservice.InternalError{}
		assert.ErrorAs(t, err, &expectedError)
		assert.ErrorContains(t, err, "error while listing dogu CRs")
	})
	t.Run("should fail for multiple dogus", func(t *testing.T) {
		// given
		doguClientMock := NewMockDoguInterface(t)
		doguList := &v1.DoguList{Items: []v1.Dogu{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "postgresql",
				},
				Spec: v1.DoguSpec{
					Name:    "official/postgresql",
					Version: "invalid",
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "ldap",
				},
				Spec: v1.DoguSpec{
					Name:    "official/ldap",
					Version: "invalid",
				},
			},
		}}
		doguClientMock.EXPECT().List(testCtx, metav1.ListOptions{}).Return(doguList, nil)

		volumeQuantity2 := resource.MustParse("2Gi")
		volumeQuantity3 := resource.MustParse("3Gi")
		pvcClientMock := NewMockPvcInterface(t)
		postgresqlPvc := corev1.PersistentVolumeClaim{Status: corev1.PersistentVolumeClaimStatus{Capacity: map[corev1.ResourceName]resource.Quantity{corev1.ResourceStorage: volumeQuantity2}}}
		ldapPvc := corev1.PersistentVolumeClaim{Status: corev1.PersistentVolumeClaimStatus{Capacity: map[corev1.ResourceName]resource.Quantity{corev1.ResourceStorage: volumeQuantity3}}}
		list := &corev1.PersistentVolumeClaimList{Items: []corev1.PersistentVolumeClaim{postgresqlPvc, ldapPvc}}
		pvcClientMock.EXPECT().List(testCtx, metav1.ListOptions{LabelSelector: "app=ces"}).Return(list, nil)

		sut := &doguInstallationRepo{doguClient: doguClientMock, pvcClient: pvcClientMock}

		// when
		_, err := sut.GetAll(testCtx)

		// then
		require.Error(t, err)
		expectedError := &domainservice.InternalError{}
		assert.ErrorAs(t, err, &expectedError)
		assert.ErrorContains(t, err, "failed to parse some dogu CRs")
	})
	t.Run("should succeed for multiple dogus", func(t *testing.T) {
		// given
		doguClientMock := NewMockDoguInterface(t)
		doguList := &v1.DoguList{Items: []v1.Dogu{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "postgresql",
				},
				Spec: v1.DoguSpec{
					Name:    "official/postgresql",
					Version: "1.2.3-1",
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "ldap",
				},
				Spec: v1.DoguSpec{
					Name:    "official/ldap",
					Version: "3.2.1-3",
				},
			},
		}}
		doguClientMock.EXPECT().List(testCtx, metav1.ListOptions{}).Return(doguList, nil)
		volumeQuantity2 := resource.MustParse("2Gi")
		volumeQuantity3 := resource.MustParse("3Gi")
		pvcClientMock := NewMockPvcInterface(t)
		postgresqlPvc := corev1.PersistentVolumeClaim{ObjectMeta: metav1.ObjectMeta{Name: "postgresql"}, Status: corev1.PersistentVolumeClaimStatus{Capacity: map[corev1.ResourceName]resource.Quantity{corev1.ResourceStorage: volumeQuantity2}}}
		ldapPvc := corev1.PersistentVolumeClaim{ObjectMeta: metav1.ObjectMeta{Name: "ldap"}, Status: corev1.PersistentVolumeClaimStatus{Capacity: map[corev1.ResourceName]resource.Quantity{corev1.ResourceStorage: volumeQuantity3}}}
		list := &corev1.PersistentVolumeClaimList{Items: []corev1.PersistentVolumeClaim{postgresqlPvc, ldapPvc}}
		pvcClientMock.EXPECT().List(testCtx, metav1.ListOptions{LabelSelector: "app=ces"}).Return(list, nil)

		sut := &doguInstallationRepo{doguClient: doguClientMock, pvcClient: pvcClientMock}

		// when
		actual, err := sut.GetAll(testCtx)

		// then
		require.NoError(t, err)
		expectedDoguInstallations := map[common.SimpleDoguName]*ecosystem.DoguInstallation{
			"postgresql": {
				Name:               postgresDoguName,
				Version:            core.Version{Raw: "1.2.3-1", Major: 1, Minor: 2, Patch: 3, Nano: 0, Extra: 1},
				MinVolumeSize:      volumeQuantity2,
				PersistenceContext: map[string]interface{}{"doguInstallationRepoContext": doguInstallationRepoContext{resourceVersion: ""}},
			},
			"ldap": {
				Name: common.QualifiedDoguName{
					Namespace:  "official",
					SimpleName: "ldap",
				},
				Version:            core.Version{Raw: "3.2.1-3", Major: 3, Minor: 2, Patch: 1, Nano: 0, Extra: 3},
				MinVolumeSize:      volumeQuantity3,
				PersistenceContext: map[string]interface{}{"doguInstallationRepoContext": doguInstallationRepoContext{resourceVersion: ""}},
			},
		}
		assert.Equal(t, expectedDoguInstallations, actual)
	})
}

func Test_doguInstallationRepo_appendVolumeSize(t *testing.T) {
	t.Run("should set volume size from pvc in dogu cr", func(t *testing.T) {
		// given
		sut := doguInstallationRepo{}
		cr := &v1.Dogu{}
		size := resource.MustParse("2Gi")
		pvcList := corev1.PersistentVolumeClaimList{
			Items: []corev1.PersistentVolumeClaim{{Status: corev1.PersistentVolumeClaimStatus{Capacity: map[corev1.ResourceName]resource.Quantity{corev1.ResourceStorage: size}}}},
		}

		// when
		err := sut.appendVolumeSize(cr, &pvcList)

		// then
		require.NoError(t, err)
		assert.Equal(t, size.String(), cr.Spec.Resources.DataVolumeSize)
	})
}
