package reconciler

import (
	"testing"
	"time"

	v2 "github.com/cloudogu/k8s-dogu-lib/v2/api/v2"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

func TestMakeResourcePredicate(t *testing.T) {
	filterFunc := func(o client.Object) bool {
		return o.GetNamespace() == "test-namespace"
	}

	predicate := makeResourcePredicate[*corev1.ConfigMap](filterFunc)

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-cm",
			Namespace: "test-namespace",
		},
	}

	cmWrongNS := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-cm",
			Namespace: "wrong-namespace",
		},
	}

	t.Run("CreateFunc should return false", func(t *testing.T) {
		result := predicate.CreateFunc(event.TypedCreateEvent[*corev1.ConfigMap]{
			Object: cm,
		})
		assert.False(t, result)
	})

	t.Run("UpdateFunc should use filter function", func(t *testing.T) {
		result := predicate.UpdateFunc(event.TypedUpdateEvent[*corev1.ConfigMap]{
			ObjectOld: cm,
			ObjectNew: cm,
		})
		assert.True(t, result)

		result = predicate.UpdateFunc(event.TypedUpdateEvent[*corev1.ConfigMap]{
			ObjectOld: cmWrongNS,
			ObjectNew: cmWrongNS,
		})
		assert.False(t, result)
	})

	t.Run("DeleteFunc should use filter function", func(t *testing.T) {
		result := predicate.DeleteFunc(event.TypedDeleteEvent[*corev1.ConfigMap]{
			Object: cm,
		})
		assert.True(t, result)

		result = predicate.DeleteFunc(event.TypedDeleteEvent[*corev1.ConfigMap]{
			Object: cmWrongNS,
		})
		assert.False(t, result)
	})

	t.Run("GenericFunc should return false", func(t *testing.T) {
		result := predicate.GenericFunc(event.TypedGenericEvent[*corev1.ConfigMap]{
			Object: cm,
		})
		assert.False(t, result)
	})
}

func TestMakeContentPredicate(t *testing.T) {
	debounce := &SingletonDebounce{}
	window := 50 * time.Millisecond

	changeFunc := func(old, new *corev1.ConfigMap) bool {
		return old.Data["key"] != new.Data["key"]
	}

	predicate := makeContentPredicate(debounce, window, changeFunc)

	cmOld := &corev1.ConfigMap{
		Data: map[string]string{"key": "old-value"},
	}

	cmNew := &corev1.ConfigMap{
		Data: map[string]string{"key": "new-value"},
	}

	cmSame := &corev1.ConfigMap{
		Data: map[string]string{"key": "old-value"},
	}

	t.Run("CreateFunc should return false", func(t *testing.T) {
		result := predicate.CreateFunc(event.TypedCreateEvent[*corev1.ConfigMap]{
			Object: cmNew,
		})
		assert.False(t, result)
	})

	t.Run("UpdateFunc should handle content changes with debouncing", func(t *testing.T) {
		// Reset debounce
		debounce = &SingletonDebounce{}
		predicate = makeContentPredicate(debounce, window, changeFunc)

		// First change should be allowed
		result := predicate.UpdateFunc(event.TypedUpdateEvent[*corev1.ConfigMap]{
			ObjectOld: cmOld,
			ObjectNew: cmNew,
		})
		assert.True(t, result)

		// Immediate second change should be debounced
		result = predicate.UpdateFunc(event.TypedUpdateEvent[*corev1.ConfigMap]{
			ObjectOld: cmOld,
			ObjectNew: cmNew,
		})
		assert.False(t, result)
	})

	t.Run("UpdateFunc should return false on no content changes", func(t *testing.T) {
		// Reset debounce
		debounce = &SingletonDebounce{}
		predicate = makeContentPredicate(debounce, window, changeFunc)

		// No change should return false
		result := predicate.UpdateFunc(event.TypedUpdateEvent[*corev1.ConfigMap]{
			ObjectOld: cmOld,
			ObjectNew: cmSame,
		})
		assert.False(t, result)

		// Immediate second change should be allowed
		result = predicate.UpdateFunc(event.TypedUpdateEvent[*corev1.ConfigMap]{
			ObjectOld: cmOld,
			ObjectNew: cmNew,
		})
		assert.True(t, result)
	})

	t.Run("DeleteFunc should trigger debouncing", func(t *testing.T) {
		// Reset debounce
		debounce = &SingletonDebounce{}
		predicate = makeContentPredicate(debounce, window, changeFunc)

		result := predicate.DeleteFunc(event.TypedDeleteEvent[*corev1.ConfigMap]{
			Object: cmOld,
		})
		assert.True(t, result)

		// Immediate second change should be debounced
		result = predicate.DeleteFunc(event.TypedDeleteEvent[*corev1.ConfigMap]{
			Object: cmOld,
		})
		assert.False(t, result)
	})

	t.Run("GenericFunc should return false", func(t *testing.T) {
		result := predicate.GenericFunc(event.TypedGenericEvent[*corev1.ConfigMap]{
			Object: cmOld,
		})
		assert.False(t, result)
	})
}

func TestConfigMapContentChanged(t *testing.T) {
	tests := []struct {
		name     string
		oldCM    *corev1.ConfigMap
		newCM    *corev1.ConfigMap
		expected bool
	}{
		{
			name: "data changed",
			oldCM: &corev1.ConfigMap{
				Data: map[string]string{"key": "old"},
			},
			newCM: &corev1.ConfigMap{
				Data: map[string]string{"key": "new"},
			},
			expected: true,
		},
		{
			name: "binary data changed",
			oldCM: &corev1.ConfigMap{
				BinaryData: map[string][]byte{"key": []byte("old")},
			},
			newCM: &corev1.ConfigMap{
				BinaryData: map[string][]byte{"key": []byte("new")},
			},
			expected: true,
		},
		{
			name: "no change",
			oldCM: &corev1.ConfigMap{
				Data: map[string]string{"key": "same"},
			},
			newCM: &corev1.ConfigMap{
				Data: map[string]string{"key": "same"},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := configMapContentChanged(tt.oldCM, tt.newCM)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSecretContentChanged(t *testing.T) {
	immutableTrue := true
	immutableFalse := false

	tests := []struct {
		name      string
		oldSecret *corev1.Secret
		newSecret *corev1.Secret
		expected  bool
	}{
		{
			name: "data changed",
			oldSecret: &corev1.Secret{
				Data: map[string][]byte{"key": []byte("old")},
			},
			newSecret: &corev1.Secret{
				Data: map[string][]byte{"key": []byte("new")},
			},
			expected: true,
		},
		{
			name: "immutable changed",
			oldSecret: &corev1.Secret{
				Immutable: &immutableTrue,
			},
			newSecret: &corev1.Secret{
				Immutable: &immutableFalse,
			},
			expected: true,
		},
		{
			name: "no change",
			oldSecret: &corev1.Secret{
				Data: map[string][]byte{"key": []byte("same")},
			},
			newSecret: &corev1.Secret{
				Data: map[string][]byte{"key": []byte("same")},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := secretContentChanged(tt.oldSecret, tt.newSecret)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDoguSpecChanged(t *testing.T) {
	tests := []struct {
		name     string
		oldDogu  *v2.Dogu
		newDogu  *v2.Dogu
		expected bool
	}{
		{
			name: "spec changed",
			oldDogu: &v2.Dogu{
				Spec: v2.DoguSpec{
					Name:    "old-name",
					Version: "1.2.3-1",
				},
			},
			newDogu: &v2.Dogu{
				Spec: v2.DoguSpec{
					Name:    "new-name",
					Version: "1.2.3-2",
				},
			},
			expected: true,
		},
		{
			name: "no change",
			oldDogu: &v2.Dogu{
				Spec: v2.DoguSpec{
					Name:    "same-name",
					Version: "1.2.3-1",
				},
			},
			newDogu: &v2.Dogu{
				Spec: v2.DoguSpec{
					Name:    "same-name",
					Version: "1.2.3-1",
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := doguSpecChanged(tt.oldDogu, tt.newDogu)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHasCesLabel(t *testing.T) {
	tests := []struct {
		name     string
		obj      client.Object
		expected bool
	}{
		{
			name: "has ces dogu config labels",
			obj: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":       "ces",
						"dogu.name": "test-dogu",
					},
				},
			},
			expected: true,
		},
		{
			name: "has ces global config labels",
			obj: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":       "ces",
						"dogu.name": "test-dogu",
					},
				},
			},
			expected: true,
		},
		{
			name: "missing app label",
			obj: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"dogu.name": "test-dogu",
					},
				},
			},
			expected: false,
		},
		{
			name: "missing dogu.name label",
			obj: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "ces",
					},
				},
			},
			expected: false,
		},
		{
			name: "wrong app label",
			obj: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":       "other",
						"dogu.name": "test-dogu",
					},
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasCesLabel(tt.obj)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHasNotDoguDescriptorLabel(t *testing.T) {
	tests := []struct {
		name     string
		obj      client.Object
		expected bool
	}{
		{
			name: "no type label",
			obj: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{},
				},
			},
			expected: true,
		},
		{
			name: "different type label",
			obj: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"k8s.cloudogu.com/type": "other-type",
					},
				},
			},
			expected: true,
		},
		{
			name: "has dogu descriptor label",
			obj: &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"k8s.cloudogu.com/type": "local-dogu-registry",
					},
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasNotDoguDescriptorLabel(tt.obj)
			assert.Equal(t, tt.expected, result)
		})
	}
}
