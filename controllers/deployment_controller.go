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
	"github.com/lostbrain101/img-backup-controller/pkg/registry"
	appsv1 "k8s.io/api/apps/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

const (
	requeuePeriod = 10 * time.Second
	kubeSystem    = "kube-system"
)

// DeploymentReconciler reconciles a Deployment object
type DeploymentReconciler struct {
	client.Client
	Log          logr.Logger
	Scheme       *runtime.Scheme
	recorder     record.EventRecorder
	ClusterCache cache.Cache
	Registry     *registry.RegistryOptions
}

// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get;update;patch

func (r *DeploymentReconciler) Reconcile(_ context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("deployment", req.NamespacedName)

	var deployment = &appsv1.Deployment{}
	if err := r.Get(context.TODO(), req.NamespacedName, deployment); err != nil {
		if apierrs.IsNotFound(err) {
			log.Info("Deployment not found ", "name:", req.Name)
			log.Info("reconcile completed")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get Deployment, reconcile failed")
		return ctrl.Result{}, err
	}

	deploymentSpecUpdate, err := processContainers(deployment.Spec.Template.Spec.Containers, r.Registry)
	if err != nil {
		log.Error(err, "Error while processing containers")
		return ctrl.Result{RequeueAfter: requeuePeriod}, err
	}

	if deploymentSpecUpdate {
		if err := r.Update(context.TODO(), deployment); err != nil {
			log.Error(err, "Deployment could not be updated.")
			return ctrl.Result{RequeueAfter: requeuePeriod}, err
		}
	}

	log.Info("Deployment reconciled")
	return ctrl.Result{}, nil
}

func (r *DeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.Deployment{}).
		WithEventFilter(ignoreSystemNamespace()).
		Complete(r)
}
