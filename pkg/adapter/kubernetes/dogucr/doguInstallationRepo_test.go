package dogucr

import (
	"context"
	"errors"
	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
	v2 "github.com/cloudogu/k8s-dogu-operator/v3/api/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"testing"
)

var version3214, _ = core.ParseVersion("3.2.1-4")

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
			&v2.Dogu{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name:            "postgresql",
					ResourceVersion: crResourceVersion,
				},
				Spec: v2.DoguSpec{
					Name:      "official/postgresql",
					Version:   version3214.Raw,
					Resources: v2.DoguResources{},
					UpgradeConfig: v2.UpgradeConfig{
						AllowNamespaceSwitch: false,
					},
					AdditionalIngressAnnotations: v2.IngressAnnotations{
						"nginx.ingress.kubernetes.io/proxy-body-size":       "1G",
						"nginx.ingress.kubernetes.io/rewrite-target":        "/",
						"nginx.ingress.kubernetes.io/configuration-snippet": "snippet",
					},
				},
				Status: v2.DoguStatus{
					Status: v2.DoguStatusInstalled,
					Health: v2.AvailableHealthStatus,
				},
			}, nil)
		quantity2 := resource.MustParse("2Gi")
		expectedPVCList := &corev1.PersistentVolumeClaimList{
			Items: []corev1.PersistentVolumeClaim{
				{
					ObjectMeta: metav1.ObjectMeta{Name: "postgresql", Labels: map[string]string{"app": "ces"}},
					Status: corev1.PersistentVolumeClaimStatus{
						Capacity: map[corev1.ResourceName]resource.Quantity{corev1.ResourceStorage: quantity2},
					},
				},
			},
		}
		pvcClientMock.EXPECT().List(testCtx, metav1.ListOptions{LabelSelector: "app=ces"}).Return(expectedPVCList, nil)

		dogu, err := repo.GetByName(testCtx, "postgresql")

		// then
		require.NoError(t, err)
		quantity1 := resource.MustParse("1G")
		assert.Equal(t, &ecosystem.DoguInstallation{
			Name:               postgresDoguName,
			Version:            version3214,
			Status:             ecosystem.DoguStatusInstalled,
			Health:             ecosystem.AvailableHealthStatus,
			UpgradeConfig:      ecosystem.UpgradeConfig{},
			PersistenceContext: persistenceContext,
			MinVolumeSize:      &quantity2,
			ReverseProxyConfig: ecosystem.ReverseProxyConfig{
				MaxBodySize:      &quantity1,
				RewriteTarget:    "/",
				AdditionalConfig: "snippet",
			},
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

	t.Run("should return internal error on error getting pvcs", func(t *testing.T) {
		// given
		doguClientMock := NewMockDoguInterface(t)
		pvcClientMock := NewMockPvcInterface(t)
		repo := NewDoguInstallationRepo(doguClientMock, pvcClientMock)
		// when
		doguClientMock.EXPECT().Get(testCtx, "postgresql", metav1.GetOptions{}).Return(&v2.Dogu{}, nil)
		pvcClientMock.EXPECT().List(testCtx, metav1.ListOptions{LabelSelector: "app=ces"}).Return(nil, assert.AnError)

		_, err := repo.GetByName(testCtx, "postgresql")

		// then
		require.Error(t, err)
		var expectedError *domainservice.InternalError
		assert.ErrorAs(t, err, &expectedError)
		assert.ErrorContains(t, err, "error while listing dogu PVCs")
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

	t.Run("should return error on error getting pvcs", func(t *testing.T) {
		// given
		doguClientMock := NewMockDoguInterface(t)
		pvcClientMock := NewMockPvcInterface(t)
		doguClientMock.EXPECT().List(testCtx, metav1.ListOptions{}).Return(nil, nil)
		pvcClientMock.EXPECT().List(testCtx, metav1.ListOptions{LabelSelector: "app=ces"}).Return(nil, assert.AnError)

		sut := &doguInstallationRepo{doguClient: doguClientMock, pvcClient: pvcClientMock}

		// when
		_, err := sut.GetAll(testCtx)

		// then
		require.Error(t, err)
		assert.ErrorIs(t, err, assert.AnError)
		expectedError := &domainservice.InternalError{}
		assert.ErrorAs(t, err, &expectedError)
		assert.ErrorContains(t, err, "error while listing dogu PVCs")
	})

	t.Run("should fail for multiple dogus", func(t *testing.T) {
		// given
		doguClientMock := NewMockDoguInterface(t)
		doguList := &v2.DoguList{Items: []v2.Dogu{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "postgresql",
				},
				Spec: v2.DoguSpec{
					Name:    "official/postgresql",
					Version: "invalid",
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "ldap",
				},
				Spec: v2.DoguSpec{
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
		doguList := &v2.DoguList{Items: []v2.Dogu{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "postgresql",
				},
				Spec: v2.DoguSpec{
					Name:    "official/postgresql",
					Version: "1.2.3-1",
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name: "ldap",
				},
				Spec: v2.DoguSpec{
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
		expectedDoguInstallations := map[cescommons.SimpleName]*ecosystem.DoguInstallation{
			"postgresql": {
				Name:               postgresDoguName,
				Version:            core.Version{Raw: "1.2.3-1", Major: 1, Minor: 2, Patch: 3, Nano: 0, Extra: 1},
				MinVolumeSize:      &volumeQuantity2,
				PersistenceContext: map[string]interface{}{"doguInstallationRepoContext": doguInstallationRepoContext{resourceVersion: ""}},
			},
			"ldap": {
				Name: cescommons.QualifiedName{
					Namespace:  "official",
					SimpleName: "ldap",
				},
				Version:            core.Version{Raw: "3.2.1-3", Major: 3, Minor: 2, Patch: 1, Nano: 0, Extra: 3},
				MinVolumeSize:      &volumeQuantity3,
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
		cr := &v2.Dogu{}
		size := resource.MustParse("2Gi")
		pvcList := corev1.PersistentVolumeClaimList{
			Items: []corev1.PersistentVolumeClaim{{Status: corev1.PersistentVolumeClaimStatus{Capacity: map[corev1.ResourceName]resource.Quantity{corev1.ResourceStorage: size}}}},
		}

		// when
		sut.appendVolumeSizeIfNotSet(cr, &pvcList)

		// then
		assert.Equal(t, size.String(), cr.Spec.Resources.DataVolumeSize)
	})

	t.Run("should do nothing if the volume size is already defined in dogu cr", func(t *testing.T) {
		// given
		sut := doguInstallationRepo{}
		cr := &v2.Dogu{
			Spec: v2.DoguSpec{
				Resources: v2.DoguResources{
					DataVolumeSize: "2Gi",
				},
			},
		}

		// when
		sut.appendVolumeSizeIfNotSet(cr, nil)

		// then
		assert.Equal(t, "2Gi", cr.Spec.Resources.DataVolumeSize)
	})

	t.Run("should do nothing if volume size is not defined and no volume exists", func(t *testing.T) {
		// given
		sut := doguInstallationRepo{}
		cr := &v2.Dogu{
			Spec: v2.DoguSpec{
				Resources: v2.DoguResources{
					DataVolumeSize: "",
				},
			},
		}
		pvcList := corev1.PersistentVolumeClaimList{
			Items: []corev1.PersistentVolumeClaim{},
		}

		// when
		sut.appendVolumeSizeIfNotSet(cr, &pvcList)

		// then
		assert.Equal(t, "", cr.Spec.Resources.DataVolumeSize)
	})
}

func Test_doguInstallationRepo_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		doguClientMock := NewMockDoguInterface(t)
		repo := &doguInstallationRepo{doguClient: doguClientMock}

		doguClientMock.EXPECT().Delete(testCtx, "postgresql", metav1.DeleteOptions{}).Return(nil)

		// when
		err := repo.Delete(testCtx, "postgresql")

		// then
		require.NoError(t, err)
	})

	t.Run("should return internal error on delete error", func(t *testing.T) {
		// given
		doguClientMock := NewMockDoguInterface(t)
		repo := &doguInstallationRepo{doguClient: doguClientMock}

		doguClientMock.EXPECT().Delete(testCtx, "postgresql", metav1.DeleteOptions{}).Return(assert.AnError)

		// when
		err := repo.Delete(testCtx, "postgresql")

		// then
		require.Error(t, err)
		var internalErr *domainservice.InternalError
		assert.ErrorAs(t, err, &internalErr)
		assert.ErrorContains(t, err, "cannot delete dogu CR for dogu \"postgresql\"")
	})
}

func Test_doguInstallationRepo_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		doguClientMock := NewMockDoguInterface(t)
		repo := &doguInstallationRepo{doguClient: doguClientMock}

		expectedDoguCr := &v2.Dogu{
			ObjectMeta: metav1.ObjectMeta{
				Name: string(postgresDoguName.SimpleName),
				Labels: map[string]string{
					"app":       "ces",
					"dogu.name": string(postgresDoguName.SimpleName),
				},
			},
			Spec: v2.DoguSpec{
				Name:    postgresDoguName.String(),
				Version: version3214.Raw,
			},
		}
		doguClientMock.EXPECT().Create(testCtx, expectedDoguCr, metav1.CreateOptions{}).Return(nil, nil)
		dogu := &ecosystem.DoguInstallation{
			Name:    postgresDoguName,
			Version: version3214,
		}

		// when
		err := repo.Create(testCtx, dogu)

		// then
		require.NoError(t, err)
	})

	t.Run("should return internal error on create error", func(t *testing.T) {
		// given
		doguClientMock := NewMockDoguInterface(t)
		repo := &doguInstallationRepo{doguClient: doguClientMock}

		doguClientMock.EXPECT().Create(testCtx, mock.Anything, metav1.CreateOptions{}).Return(nil, assert.AnError)

		dogu := &ecosystem.DoguInstallation{
			Name: postgresDoguName,
		}

		// when
		err := repo.Create(testCtx, dogu)

		// then
		require.Error(t, err)
		var internalErr *domainservice.InternalError
		assert.ErrorAs(t, err, &internalErr)
		assert.ErrorContains(t, err, "cannot create dogu CR for dogu \"official/postgresql\"")
	})
}

func Test_doguInstallationRepo_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// given
		doguClientMock := NewMockDoguInterface(t)
		repo := &doguInstallationRepo{doguClient: doguClientMock}

		expectedDoguPatch := "{\"spec\":{\"name\":\"official/postgresql\",\"version\":\"3.2.1-4\",\"resources\":{\"dataVolumeSize\":\"\"},\"supportMode\":false,\"upgradeConfig\":{\"allowNamespaceSwitch\":false,\"forceUpgrade\":false},\"additionalIngressAnnotations\":null}}"
		doguClientMock.EXPECT().Patch(testCtx, "postgresql", types.MergePatchType, []byte(expectedDoguPatch), metav1.PatchOptions{}).Return(nil, nil)
		dogu := &ecosystem.DoguInstallation{
			Name:    postgresDoguName,
			Version: version3214,
		}

		// when
		err := repo.Update(testCtx, dogu)

		// then
		require.NoError(t, err)
	})

	t.Run("should return internal error on create error", func(t *testing.T) {
		// given
		doguClientMock := NewMockDoguInterface(t)
		repo := &doguInstallationRepo{doguClient: doguClientMock}

		doguClientMock.EXPECT().Patch(testCtx, "postgresql", mock.Anything, mock.Anything, metav1.PatchOptions{}).Return(nil, assert.AnError)

		dogu := &ecosystem.DoguInstallation{
			Name: postgresDoguName,
		}

		// when
		err := repo.Update(testCtx, dogu)

		// then
		require.Error(t, err)
		var internalErr *domainservice.InternalError
		assert.ErrorAs(t, err, &internalErr)
		assert.ErrorContains(t, err, "cannot patch dogu CR for dogu \"official/postgresql\"")
	})
}
