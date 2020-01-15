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
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

func createWebprojectWorkload(client *kubernetes.Clientset, deploymentInput WebProjectInput) {
	// WebProject Deployment.
	deployment := &appsv1.Deployment{
		TypeMeta: genTypeMeta("Deployment"),
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentInput.DeploymentName,
			Namespace: deploymentInput.Namespace,
			Labels:    genDefaultLabels(deploymentInput),
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32ptr(deploymentInput.Replicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: genDefaultLabels(deploymentInput),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   deploymentInput.DeploymentName,
					Labels: genDefaultLabels(deploymentInput),
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            deploymentInput.PrimaryContainerName,
							Image:           deploymentInput.PrimaryContainerImageTag,
							ImagePullPolicy: corev1.PullIfNotPresent,
							VolumeMounts: []corev1.VolumeMount{
								createVolumeMount("webroot", "/var/www/webroot"),
								createVolumeMount("files", "/var/www/html/sites/default/files"),
							},
							Ports: []corev1.ContainerPort{
								createContainerPort(int32(deploymentInput.PrimaryContainerPort)),
							},
						},
					},
					RestartPolicy: corev1.RestartPolicyAlways,
					Volumes: []corev1.Volume{
						createEmptyDirVolume("webroot"),
						attachVolumeFromClaim("files", "webfiles", deploymentInput),
					},
				},
			},
		},
	}

	// Create Web Project Deployment
	foundWebProject, foundErr := client.AppsV1().Deployments(deploymentInput.Namespace).Get(deploymentInput.DeploymentName, metav1.GetOptions{})
	if foundErr != nil {
		fmt.Println("Creating webproject deployment...")
		resultWebProject, errWebProject := client.AppsV1().Deployments(deploymentInput.Namespace).Create(deployment)
		if errWebProject != nil {
			panic(errWebProject)
		}
		fmt.Printf("Created Deployment - Name: %q, UID: %q\n", resultWebProject.GetObjectMeta().GetName(), resultWebProject.GetObjectMeta().GetUID())
	} else {
		fmt.Println("Updating webproject deployment...")
		foundWebProject.Spec.Replicas = int32ptr(deploymentInput.Replicas)
		foundWebProject.Spec.Template.Spec.Containers[0].Image = deploymentInput.PrimaryContainerImageTag
		foundWebProjectResult, errFoundWebProject := client.AppsV1().Deployments(deploymentInput.Namespace).Update(foundWebProject)
		if errFoundWebProject != nil {
			panic(errFoundWebProject)
		}
		fmt.Printf("Updated Deployment - Name: %q, UID: %q\n", foundWebProjectResult.GetObjectMeta().GetName(), foundWebProjectResult.GetObjectMeta().GetUID())
	}

	fmt.Println("Creating service for WebProject.")
	serviceName := deploymentInput.DeploymentName + "-svc"

	webprojectLabels := map[string]string{
		"app": deploymentInput.DeploymentName,
	}
	webprojectService := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   serviceName,
			Labels: genDefaultLabels(deploymentInput),
		},
		Spec: v1.ServiceSpec{
			Selector: webprojectLabels,
			Ports: []v1.ServicePort{{
				Port:       8080,
				Protocol:   "TCP",
				TargetPort: intstr.FromInt(deploymentInput.PrimaryContainerPort),
			}},
		},
	}

	_, foundWebprojectServiceErr := client.CoreV1().Services(deploymentInput.Namespace).Get(deploymentInput.DeploymentName+"-svc", metav1.GetOptions{})
	if foundWebprojectServiceErr != nil {
		webprojectServiceResp, errWebprojectService := client.CoreV1().Services(deploymentInput.Namespace).Create(webprojectService)
		if errWebprojectService != nil {
			panic(errWebprojectService)
		}
		fmt.Printf("Created Webproject Service - Name: %q, UID: %q\n", webprojectServiceResp.GetObjectMeta().GetName(), webprojectServiceResp.GetObjectMeta().GetUID())

	}
}

func createContainerPort(portNumber int32) corev1.ContainerPort {
	return corev1.ContainerPort{
		ContainerPort: portNumber,
		Protocol:      corev1.ProtocolTCP,
	}
}
