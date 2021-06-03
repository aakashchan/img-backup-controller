package controllers

import (
	registry "github.com/lostbrain.io/img-backup-controller/pkg"
	"k8c.io/utils/pkg/util"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"strings"
)

func processImage(image string) (string, error){

	var (
		newImage string
		err 	 error
	)

	newImage = image
	if !strings.Contains(image, registry.DockerOrg) {
		newImage, err = registry.PullAndUploadImage(image)
		if err != nil {
			return "",err
		}
	}

	return newImage, nil
}

func processContainers(containers []v1.Container) (bool,error) {
	var updated bool
	for i, container := range containers{
		newImg, err := processImage(container.Image)
		if err != nil {
			return false,err
		}
		if newImg != container.Image {
			containers[i].Image = newImg
			updated = true
		}
	}

	return updated,nil
}

func ignoreSystemNamespace() predicate.Predicate {
	return util.PredicateFn(func(obj runtime.Object) bool {
		// we are only interested in non kubesystem objects
		meta, ok := obj.(metav1.Object)
		if !ok {
			return false
		}
		return meta.GetNamespace() != kubeSystem
	})
}