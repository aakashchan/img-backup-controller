/*
Copyright 2021 lostbrain101.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// DaemonSetReconciler reconciles a DaemonSet object
type DaemonSetReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=apps,resources=daemonsets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=daemonsets/status,verbs=get;update;patch

func (r *DaemonSetReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	var err error
	log := r.Log.WithValues("daemonset", req.NamespacedName)

	daemonset := &appsv1.DaemonSet{}
	if err = r.Get(context.TODO(), req.NamespacedName, daemonset); err != nil {
		if apierrs.IsNotFound(err) {
			log.Info("Daemenset not found ", "name:", req.Name)
			log.Info("reconcile completed")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get Deployment, reconcile failed")
		return ctrl.Result{}, err
	}

	daemonsetSpecUpdate,err := processContainers(daemonset.Spec.Template.Spec.Containers)
	if err != nil {
		log.Error(err, "Error while processing containers")
		return ctrl.Result{RequeueAfter: requeuePeriod}, err
	}

	if daemonsetSpecUpdate {
		if err := r.Update(context.TODO(), daemonset); err != nil {
			log.Error(err, "Deployment could not be updated.")
			return ctrl.Result{RequeueAfter: requeuePeriod}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *DaemonSetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.DaemonSet{}).
		WithEventFilter(ignoreSystemNamespace()).
		Complete(r)
}
