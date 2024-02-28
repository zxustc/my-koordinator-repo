package workload

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ handler.EventHandler = &EnqueueRequestForWorkload{}

type EnqueueRequestForWorkload struct {
	client.Client
}

func (n *EnqueueRequestForWorkload) Create(e event.CreateEvent, q workqueue.RateLimitingInterface) {
	if pod, ok := e.Object.(*corev1.Pod); !ok {
		return
	} else {
		namespace, name := GetOwnerReferenceNamespaceName(pod)
		q.Add(reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      name,
				Namespace: namespace,
			},
		})
	}
}

func (n *EnqueueRequestForWorkload) Update(e event.UpdateEvent, q workqueue.RateLimitingInterface) {
	newPod, _ := e.ObjectNew.(*corev1.Pod), e.ObjectOld.(*corev1.Pod)
	namespace, name := GetOwnerReferenceNamespaceName(newPod)
	q.Add(reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      name,
			Namespace: namespace,
		},
	})
}

func (n *EnqueueRequestForWorkload) Delete(e event.DeleteEvent, q workqueue.RateLimitingInterface) {
	if pod, ok := e.Object.(*corev1.Pod); !ok {
		return
	} else {
		namespace, name := GetOwnerReferenceNamespaceName(pod)
		q.Add(reconcile.Request{
			NamespacedName: types.NamespacedName{
				Name:      name,
				Namespace: namespace,
			},
		})
	}
}

func (n *EnqueueRequestForWorkload) Generic(e event.GenericEvent, q workqueue.RateLimitingInterface) {
}
