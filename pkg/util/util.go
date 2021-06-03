package util

import (
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

type PredicateFn func(obj runtime.Object) bool

func (p PredicateFn) Create(ev event.CreateEvent) bool {
	return p(ev.Object)
}

func (p PredicateFn) Delete(ev event.DeleteEvent) bool {
	return p(ev.Object)
}

func (p PredicateFn) Update(ev event.UpdateEvent) bool {
	return p(ev.ObjectNew)
}

func (p PredicateFn) Generic(ev event.GenericEvent) bool {
	return p(ev.Object)
}

var _ predicate.Predicate = (PredicateFn)(nil)
