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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	util "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

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

func (r *WebPageReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("webpage", req.NamespacedName)

	// Get custom resource
	var webpage api.WebPage
	if err := r.Get(ctx, req.NamespacedName, &webpage); err != nil {
		log.Error(err, "unable to fetch WebPage")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// CreatOrUpdate ConfigMap
	var cm corev1.ConfigMap
	cm.Name = webpage.Name + "-config"
	cm.Namespace = webpage.Namespace
	_, err := ctrl.CreateOrUpdate(ctx, r, &cm, func() error {
		cm.Data = map[string]string{
			"index.html": webpage.Spec.Html,
		}
		// For garbage collector to clean up resource
		return util.SetControllerReference(&webpage, &cm, r.Scheme)
	})
	if err != nil {
		log.Error(err, "unable to CreateOrUpdate configmap")
		return ctrl.Result{}, err
	}

	// Create nginx pod with mounted ConfigMap volume if it does not exist
	if webpage.Status.LastUpdateTime == nil {
		var pod corev1.Pod
		pod.Name = webpage.Name + "-nginx"
		pod.Namespace = webpage.Namespace
		_, err = ctrl.CreateOrUpdate(ctx, r, &pod, func() error {
			pod.Spec.Containers = []corev1.Container{
				corev1.Container{
					Name:  "nginx",
					Image: "nginx",
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "config-volume",
							MountPath: "/usr/share/nginx/html",
						},
					},
				},
			}
			pod.Spec.Volumes = []corev1.Volume{
				corev1.Volume{
					Name: "config-volume",
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: webpage.Name + "-config",
							},
						},
					},
				},
			}
			// For garbage collector to clean up resource
			return util.SetControllerReference(&webpage, &pod, r.Scheme)
		})
		if err != nil {
			log.Error(err, "unable to CreateOrUpdate pod")
			return ctrl.Result{}, err
		}
	}

	// Set the last update time
	webpage.Status.LastUpdateTime = &metav1.Time{Time: time.Now()}
	if err = r.Status().Update(ctx, &webpage); err != nil {
		log.Error(err, "unable to update status")
	}

	return ctrl.Result{}, nil
}

func (r *WebPageReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&api.WebPage{}).
		Complete(r)
}
