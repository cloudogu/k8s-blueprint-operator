package reconciler

import (
	"time"

	bpv3 "github.com/cloudogu/k8s-blueprint-lib/v3/api/v3"
	v2 "github.com/cloudogu/k8s-dogu-lib/v2/api/v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

func makeResourcePredicate[T client.Object](resourceFunc func(client.Object) bool) predicate.TypedFuncs[T] {
	return predicate.TypedFuncs[T]{
		CreateFunc:  func(event.TypedCreateEvent[T]) bool { return false }, // ignore creates
		UpdateFunc:  func(e event.TypedUpdateEvent[T]) bool { return resourceFunc(e.ObjectNew) },
		DeleteFunc:  func(e event.TypedDeleteEvent[T]) bool { return resourceFunc(e.Object) },
		GenericFunc: func(event.TypedGenericEvent[T]) bool { return false }, // ignore generics
	}
}

func makeContentPredicate[T client.Object](
	debounce *SingletonDebounce,
	window time.Duration,
	changed func(oldObj, newObj T) bool,
) predicate.TypedFuncs[T] {
	return predicate.TypedFuncs[T]{
		CreateFunc: func(event.TypedCreateEvent[T]) bool { return false },
		UpdateFunc: func(e event.TypedUpdateEvent[T]) bool {
			if !changed(e.ObjectOld, e.ObjectNew) {
				return false
			}
			return debounce.AllowOrMark(window)
		},
		DeleteFunc:  func(event.TypedDeleteEvent[T]) bool { return debounce.AllowOrMark(window) }, // reconcile on delete
		GenericFunc: func(event.TypedGenericEvent[T]) bool { return false },
	}
}

func configMapContentChanged(oldCM, newCM *corev1.ConfigMap) bool {
	return !equality.Semantic.DeepEqual(oldCM.Data, newCM.Data) ||
		!equality.Semantic.DeepEqual(oldCM.BinaryData, newCM.BinaryData)
}

func secretContentChanged(oldS, newS *corev1.Secret) bool {
	return !equality.Semantic.DeepEqual(oldS.Data, newS.Data) ||
		!equality.Semantic.DeepEqual(oldS.Immutable, newS.Immutable)
}

func doguSpecChanged(oldObj, newObj *v2.Dogu) bool {
	return !equality.Semantic.DeepEqual(oldObj.Spec, newObj.Spec)
}

func blueprintMaskSpecChanged(oldObj, newObj *bpv3.BlueprintMask) bool {
	return !equality.Semantic.DeepEqual(oldObj.Spec, newObj.Spec)
}

func hasCesLabel(o client.Object) bool {
	// Consider only CES ConfigMaps that are doguConfig or globalConfig
	return o.GetLabels()["app"] == "ces" && (o.GetLabels()["dogu.name"] != "" || o.GetLabels()["k8s.cloudogu.com/type"] == "global-config")
}

func hasNotDoguDescriptorLabel(o client.Object) bool {
	return o.GetLabels()["k8s.cloudogu.com/type"] != "local-dogu-registry"
}

func (r *BlueprintReconciler) hasOperatorNamespace(o client.Object) bool {
	return r.namespace == "" || o.GetNamespace() == r.namespace
}
