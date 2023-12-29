package dogucr

import (
	"context"
	"errors"
	"github.com/cloudogu/cesapp-lib/core"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain/ecosystem"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	v1 "github.com/cloudogu/k8s-dogu-operator/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"reflect"
	"testing"
)

var version3_2_1_4, _ = core.ParseVersion("3.2.1-4")

var crResourceVersion = "abc"
var persistenceContext = map[string]interface{}{
	doguInstallationRepoContextKey: doguInstallationRepoContext{resourceVersion: crResourceVersion},
}

func Test_parseDoguCR(t *testing.T) {
	type args struct {
		cr *v1.Dogu
	}
	tests := []struct {
		name    string
		args    args
		want    *ecosystem.DoguInstallation
		wantErr bool
	}{
		{
			name:    "nil",
			args:    args{cr: nil},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ok",
			args: args{cr: &v1.Dogu{
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
						AllowNamespaceSwitch: true,
					},
				},
				Status: v1.DoguStatus{
					Status: v1.DoguStatusInstalled,
					Health: v1.AvailableHealthStatus,
				},
			}},
			want: &ecosystem.DoguInstallation{
				Namespace: "official",
				Name:      "postgresql",
				Version:   version3_2_1_4,
				Status:    ecosystem.DoguStatusInstalled,
				Health:    ecosystem.AvailableHealthStatus,
				UpgradeConfig: ecosystem.UpgradeConfig{
					AllowNamespaceSwitch: true,
				},
				PersistenceContext: persistenceContext,
			},
			wantErr: false,
		},
		{
			name: "cannot parse version",
			args: args{cr: &v1.Dogu{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Name:            "postgresql",
					ResourceVersion: "abc",
				},
				Spec: v1.DoguSpec{
					Name:      "official/postgresql",
					Version:   "vxyz",
					Resources: v1.DoguResources{},
					UpgradeConfig: v1.UpgradeConfig{
						AllowNamespaceSwitch: false,
					},
				},
				Status: v1.DoguStatus{
					Status: v1.DoguStatusInstalled,
					Health: v1.AvailableHealthStatus,
				},
			}},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDoguCR(tt.args.cr)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDoguCR() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseDoguCR() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_doguInstallationRepo_GetByName(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		//given
		ctx := context.Background()
		k8sMock := NewMockDoguInterface(t)
		repo := NewDoguInstallationRepo(k8sMock)
		//when
		k8sMock.EXPECT().Get(ctx, "postgresql", metav1.GetOptions{}).Return(
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

		dogu, err := repo.GetByName(ctx, "postgresql")

		//then
		require.NoError(t, err)
		assert.Equal(t, &ecosystem.DoguInstallation{
			Namespace:          "official",
			Name:               "postgresql",
			Version:            version3_2_1_4,
			Status:             ecosystem.DoguStatusInstalled,
			Health:             ecosystem.AvailableHealthStatus,
			UpgradeConfig:      ecosystem.UpgradeConfig{},
			PersistenceContext: persistenceContext,
		}, dogu)
	})

	t.Run("not found error", func(t *testing.T) {
		//given
		ctx := context.Background()
		k8sMock := NewMockDoguInterface(t)
		repo := NewDoguInstallationRepo(k8sMock)
		//when
		k8sMock.EXPECT().Get(ctx, "postgresql", metav1.GetOptions{}).Return(
			nil,
			k8sErrors.NewNotFound(schema.GroupResource{}, "postgresql"),
		)

		_, err := repo.GetByName(ctx, "postgresql")

		//then
		require.Error(t, err)
		var expectedError *domainservice.NotFoundError
		assert.ErrorAs(t, err, &expectedError)
	})

	t.Run("internal error", func(t *testing.T) {
		//given
		ctx := context.Background()
		k8sMock := NewMockDoguInterface(t)
		repo := NewDoguInstallationRepo(k8sMock)
		//when
		k8sMock.EXPECT().Get(ctx, "postgresql", metav1.GetOptions{}).Return(
			nil,
			k8sErrors.NewInternalError(errors.New("test-error")),
		)

		_, err := repo.GetByName(ctx, "postgresql")

		//then
		require.Error(t, err)
		var expectedError *domainservice.InternalError
		assert.ErrorAs(t, err, &expectedError)
	})
}
