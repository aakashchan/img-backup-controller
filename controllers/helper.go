package controllers

import (
	registry2 "github.com/lostbrain101/img-backup-controller/pkg/registry"
	"github.com/lostbrain101/img-backup-controller/pkg/util"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)


func processImage(image string) (string, error) {

	var (
		newImage string
		err      error
	)

	newImage, err = registry2.PullAndUploadImage(image)
	if err != nil {
		return "", err
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
	return util.PredicateFn(func(obj runtime.Object) bool {
		meta, ok := obj.(metav1.Object)
		if !ok {
			return false
		}

		return meta.GetNamespace() != kubeSystem
	})
}