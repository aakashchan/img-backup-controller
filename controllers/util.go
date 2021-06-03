package controllers

import (
	registry "github.com/lostbrain101/img-backup-controller/pkg"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"strings"
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

func processImage(image string) (string, error) {

	var (
		newImage string
		err      error
	)

	newImage = image
	if !strings.Contains(image, registry.DockerOrg) {
		newImage, err = registry.PullAndUploadImage(image)
		if err != nil {
			return "", err
		}
	}

	return newImage, nil
}

func processContainers(containers []v1.Container) (bool, error) {
	var updated bool
	for i, container := range containers {
		newImg, err := processImage(container.Image)
		if err != nil {
			return false, err
		}
		if newImg != container.Image {
			containers[i].Image = newImg
			updated = true
		}
	}

	return updated, nil
}

func ignoreSystemNamespace() predicate.Predicate {
	return PredicateFn(func(obj runtime.Object) bool {
		meta, ok := obj.(metav1.Object)
		if !ok {
			return false
		}

		return meta.GetNamespace() != kubeSystem
	})
}
