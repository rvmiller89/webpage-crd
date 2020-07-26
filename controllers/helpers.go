package controllers

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	api "github.com/rvmiller89/webpage-crd/api/v1beta1"
)

func (r *WebPageReconciler) desiredConfigMap(webpage api.WebPage) (corev1.ConfigMap, error) {
	cm := corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{APIVersion: corev1.SchemeGroupVersion.String(), Kind: "ConfigMap"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      webpage.Name + "-config",
			Namespace: webpage.Namespace,
		},
		Data: map[string]string{
			"index.html": webpage.Spec.Html,
		},
	}

	// For garbage collector to clean up resource
	if err := ctrl.SetControllerReference(&webpage, &cm, r.Scheme); err != nil {
		return cm, err
	}

	return cm, nil
}

func (r *WebPageReconciler) desiredDeployment(webpage api.WebPage, cm corev1.ConfigMap) (appsv1.Deployment, error) {
	dep := appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{APIVersion: appsv1.SchemeGroupVersion.String(), Kind: "Deployment"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      webpage.Name,
			Namespace: webpage.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"webpage": webpage.Name},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"webpage": webpage.Name},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: "nginx",
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "config-volume",
									MountPath: "/usr/share/nginx/html",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "config-volume",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: cm.Name,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	// For garbage collector to clean up resource
	if err := ctrl.SetControllerReference(&webpage, &dep, r.Scheme); err != nil {
		return dep, err
	}

	return dep, nil
}
