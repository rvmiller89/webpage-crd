/*


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
	"time"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	api "github.com/rvmiller89/webpage-crd/api/v1beta1"
)

// WebPageReconciler reconciles a WebPage object
type WebPageReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=sandbox.rvmiller.com,resources=webpages,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=sandbox.rvmiller.com,resources=webpages/status,verbs=get;update;patch

// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete

func (r *WebPageReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("webpage", req.NamespacedName)

	log.Info("starting reconcile")

	// Get custom resource
	var webpage api.WebPage
	if err := r.Get(ctx, req.NamespacedName, &webpage); err != nil {
		log.Error(err, "unable to fetch WebPage")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Desired ConfigMap
	cm, err := r.desiredConfigMap(webpage)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Desired Deployment
	dep, err := r.desiredDeployment(webpage, cm)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Patch (create/update) both owned resources
	applyOpts := []client.PatchOption{client.ForceOwnership, client.FieldOwner("webpage-controller")}

	err = r.Patch(ctx, &cm, client.Apply, applyOpts...)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = r.Patch(ctx, &dep, client.Apply, applyOpts...)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Set the last update time
	webpage.Status.LastUpdateTime = &metav1.Time{Time: time.Now()}
	if err = r.Status().Update(ctx, &webpage); err != nil {
		log.Error(err, "unable to update status")
	}

	log.Info("finished reconcile")

	return ctrl.Result{}, nil
}

func (r *WebPageReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&api.WebPage{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.ConfigMap{}).
		Complete(r)
}
