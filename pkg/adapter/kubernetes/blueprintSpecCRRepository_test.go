package kubernetes

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/serializer/blueprintMaskV1"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/adapter/serializer/blueprintV2"
	v1 "github.com/cloudogu/k8s-blueprint-operator/pkg/api/v1"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domain"
	"github.com/cloudogu/k8s-blueprint-operator/pkg/domainservice"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"testing"
)

var ctx = context.Background()

func Test_blueprintSpecRepo_GetById(t *testing.T) {
	blueprintId := "MyBlueprint"

	t.Run("all ok with empty blueprint", func(t *testing.T) {
		//given
		restClientMock := NewMockBlueprintInterface(t)
		repo := NewBlueprintSpecRepository(restClientMock, blueprintV2.Serializer{}, blueprintMaskV1.Serializer{})

		cr := &v1.Blueprint{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{},
			Spec: v1.BlueprintSpec{
				Blueprint:     `{"blueprintApi": "v2"}`,
				BlueprintMask: `{"blueprintMaskAPI": "v1"}`,
			},
			Status: v1.BlueprintStatus{},
		}
		restClientMock.EXPECT().Get(ctx, blueprintId, metav1.GetOptions{}).Return(cr, nil)

		//when
		spec, err := repo.GetById(ctx, blueprintId)

		//then
		require.NoError(t, err)
		assert.Equal(t, domain.BlueprintSpec{Id: blueprintId}, spec)
	})

	t.Run("invalid blueprint and mask", func(t *testing.T) {
		//given
		restClientMock := NewMockBlueprintInterface(t)
		repo := NewBlueprintSpecRepository(restClientMock, blueprintV2.Serializer{}, blueprintMaskV1.Serializer{})

		cr := &v1.Blueprint{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{},
			Spec: v1.BlueprintSpec{
				Blueprint:     `{}`,
				BlueprintMask: `{}`,
			},
			Status: v1.BlueprintStatus{},
		}
		restClientMock.EXPECT().Get(ctx, blueprintId, metav1.GetOptions{}).Return(cr, nil)

		//when
		_, err := repo.GetById(ctx, blueprintId)

		//then
		require.Error(t, err)
		var expectedErrorType *domain.InvalidBlueprintError
		assert.ErrorAs(t, err, &expectedErrorType)
		assert.ErrorContains(t, err, fmt.Sprintf("could not deserialize Blueprint CR %s: ", blueprintId))
		assert.ErrorContains(t, err, "cannot deserialize blueprint")
		assert.ErrorContains(t, err, "cannot deserialize blueprint mask")
	})

	t.Run("internal error while loading", func(t *testing.T) {
		//given
		restClientMock := NewMockBlueprintInterface(t)
		repo := NewBlueprintSpecRepository(restClientMock, blueprintV2.Serializer{}, blueprintMaskV1.Serializer{})

		restClientMock.EXPECT().Get(ctx, blueprintId, metav1.GetOptions{}).Return(nil, k8sErrors.NewInternalError(errors.New("test-error")))

		//when
		_, err := repo.GetById(ctx, blueprintId)

		//then
		require.Error(t, err)
		var expectedErrorType *domainservice.InternalError
		assert.ErrorAs(t, err, &expectedErrorType)
		assert.ErrorContains(t, err, fmt.Sprintf("error while loading blueprint CR '%s':", blueprintId))
		assert.ErrorContains(t, err, "test-error")
	})

	t.Run("not found error while loading", func(t *testing.T) {
		//given
		restClientMock := NewMockBlueprintInterface(t)
		repo := NewBlueprintSpecRepository(restClientMock, blueprintV2.Serializer{}, blueprintMaskV1.Serializer{})

		restClientMock.EXPECT().
			Get(ctx, blueprintId, metav1.GetOptions{}).
			Return(nil, k8sErrors.NewNotFound(schema.GroupResource{}, blueprintId))

		//when
		_, err := repo.GetById(ctx, blueprintId)

		//then
		require.Error(t, err)
		var expectedErrorType *domainservice.NotFoundError
		assert.ErrorAs(t, err, &expectedErrorType)
		assert.ErrorContains(t, err, fmt.Sprintf("cannot load Blueprint CR '%s' as it does not exist:", blueprintId))
	})
}

func Test_blueprintSpecRepo_Update(t *testing.T) {

}
