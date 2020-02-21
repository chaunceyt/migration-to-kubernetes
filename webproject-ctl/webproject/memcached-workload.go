// Copyright © 2020 Chauncey Thorn <chaunceyt@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package main

import (
	"log"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

// createMemcachedWorkload - create a memcached deployment and service.
func createMemcachedWorkload(client *kubernetes.Clientset, deploymentInput WebProjectInput) {
	memcachedDeployment := &appsv1.Deployment{
		TypeMeta: genTypeMeta("Deployment"),
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentInput.DeploymentName + "-memcached",
			Namespace: deploymentInput.Namespace,
			Labels: map[string]string{
				"app":     deploymentInput.DeploymentName + "-memcached",
				"release": deploymentInput.DeploymentName,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":     deploymentInput.DeploymentName + "-memcached",
					"release": deploymentInput.DeploymentName,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: deploymentInput.DeploymentName + "-memcached",
					Labels: map[string]string{
						"app":     deploymentInput.DeploymentName + "-memcached",
						"release": deploymentInput.DeploymentName,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            deploymentInput.CacheEngine,
							Image:           deploymentInput.CacheEngineImage,
							ImagePullPolicy: corev1.PullIfNotPresent,
							Ports: []corev1.ContainerPort{
								createContainerPort(11211),
							},
						},
					},
					RestartPolicy: corev1.RestartPolicyAlways,
				},
			},
		},
	}

	_, foundErr := client.AppsV1().Deployments(deploymentInput.Namespace).Get(deploymentInput.DeploymentName+"-memcached", metav1.GetOptions{})
	if foundErr != nil {
		// Create  Memcached Deployment
		result, err := client.AppsV1().Deployments(deploymentInput.Namespace).Create(memcachedDeployment)
		if err != nil {
			panic(err)
		}
		log.Printf("Created Memcached Deployment - Name: %q, UID: %q\n", result.GetObjectMeta().GetName(), result.GetObjectMeta().GetUID())

	}

	serviceName := deploymentInput.DeploymentName + "-memcached-svc"
	memcachedServiceLabels := map[string]string{
		"app": deploymentInput.DeploymentName + "-memcached",
	}
	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   serviceName,
			Labels: genDefaultLabels(deploymentInput),
		},
		Spec: v1.ServiceSpec{
			Selector: memcachedServiceLabels,
			Ports: []v1.ServicePort{{
				Port:       11211,
				TargetPort: intstr.FromInt(11211),
			}},
		},
	}

	_, foundServiceErr := client.CoreV1().Services(deploymentInput.Namespace).Get(deploymentInput.DeploymentName+"-memcached-svc", metav1.GetOptions{})
	if foundServiceErr != nil {
		service, errRedisService := client.CoreV1().Services(deploymentInput.Namespace).Create(service)
		if errRedisService != nil {
			panic(errRedisService)
		}
		log.Printf("Created Memcached Service - Name: %q, UID: %q\n", service.GetObjectMeta().GetName(), service.GetObjectMeta().GetUID())

	}
}
