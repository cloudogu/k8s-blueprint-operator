package dogucr

import (
	"context"
	"errors"
	"testing"

	cescommons "github.com/cloudogu/ces-commons-lib/dogu"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/v2/pkg/domainservice"
	v2 "github.com/cloudogu/k8s-dogu-lib/v2/api/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

var version3214, _ = core.ParseVersion("3.2.1-4")
var version1231, _ = core.ParseVersion("1.2.3-1")
var version3213, _ = core.ParseVersion("3.2.1-3")

var crResourceVersion = "abc"
var persistenceContext = map[string]interface{}{
	doguInstallationRepoContextKey: doguInstallationRepoContext{resourceVersion: crResourceVersion},
}

var testCtx = context.Background()

func Test_doguInstallationRepo_GetByName(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		// given
		doguClientMock := NewMockDoguInterface(t)
		repo := NewDoguInstallationRepo(doguClientMock)

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
		dogu, err := repo.GetByName(testCtx, "postgresql")

		// then
		require.NoError(t, err)
		quantity1 := resource.MustParse("1G")
		rewriteTarget := "/"
		additionalConfig := "snippet"
		assert.Equal(t, &ecosystem.DoguInstallation{
			Name:               postgresDoguName,
			Version:            version3214,
			Status:             "installed",
			Health:             ecosystem.AvailableHealthStatus,
			UpgradeConfig:      ecosystem.UpgradeConfig{},
			PersistenceContext: persistenceContext,
			MinVolumeSize:      &quantity2,
			ReverseProxyConfig: ecosystem.ReverseProxyConfig{
				MaxBodySize:      &quantity1,
				RewriteTarget:    &rewriteTarget,
				AdditionalConfig: &additionalConfig,
			},
		}, dogu)
	})

	t.Run("not found error", func(t *testing.T) {
		// given
		doguClientMock := NewMockDoguInterface(t)
		repo := NewDoguInstallationRepo(doguClientMock)
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
		repo := NewDoguInstallationRepo(doguClientMock)
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

		sut := &doguInstallationRepo{doguClient: doguClientMock}

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
		volumeQuantity2 := resource.MustParse("2Gi")
		volumeQuantity3 := resource.MustParse("3Gi")
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
					Resources: v2.DoguResources{
						MinDataVolumeSize: volumeQuantity3,
					},
				},
			},
		}}
		doguClientMock.EXPECT().List(testCtx, metav1.ListOptions{}).Return(doguList, nil)

		sut := &doguInstallationRepo{doguClient: doguClientMock}

		// when
		actual, err := sut.GetAll(testCtx)

		// then
		require.NoError(t, err)
		expectedDoguInstallations := map[cescommons.SimpleName]*ecosystem.DoguInstallation{
			"postgresql": {
				Name:               postgresDoguName,
				Version:            version1231,
				MinVolumeSize:      &volumeQuantity2,
				PersistenceContext: map[string]interface{}{"doguInstallationRepoContext": doguInstallationRepoContext{resourceVersion: ""}},
			},
			"ldap": {
				Name:               ldapDoguName,
				Version:            version3213,
				MinVolumeSize:      &volumeQuantity3,
				PersistenceContext: map[string]interface{}{"doguInstallationRepoContext": doguInstallationRepoContext{resourceVersion: ""}},
			},
		}
		//assert.Equal(t, expectedDoguInstallations, actual)
		assert.Equal(t, expectedDoguInstallations["postgresql"], actual["postgresql"])
		assert.Equal(t, expectedDoguInstallations["ldap"], actual["ldap"])
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
					"app":                          "ces",
					"k8s.cloudogu.com/app":         "ces",
					"dogu.name":                    string(postgresDoguName.SimpleName),
					"k8s.cloudogu.com/dogu.name":   string(postgresDoguName.SimpleName),
					"app.kubernetes.io/name":       string(postgresDoguName.SimpleName),
					"app.kubernetes.io/version":    version3214.Raw,
					"app.kubernetes.io/part-of":    "ces",
					"app.kubernetes.io/managed-by": "k8s-blueprint-operator",
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

		expectedDoguPatch := "{\"spec\":{" +
			"\"name\":\"official/postgresql\"," +
			"\"version\":\"3.2.1-4\"," +
			"\"resources\":{" +
			"\"dataVolumeSize\":\"\"," +
			"\"minDataVolumeSize\":\"0\"}," +
			"\"supportMode\":false," +
			"\"pauseReconciliation\":false," +
			"\"upgradeConfig\":{\"allowNamespaceSwitch\":false,\"forceUpgrade\":false}," +
			"\"additionalIngressAnnotations\":null," +
			"\"additionalMounts\":null}" +
			"}"
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
