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

// createSolrWorkload - create Solr deployment and service.
func createSolrWorkload(client *kubernetes.Clientset, deploymentInput WebProjectInput) {
	solrDeployment := &appsv1.Deployment{
		TypeMeta: genTypeMeta("Deployment"),
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentInput.DeploymentName + "-solr",
			Namespace: deploymentInput.Namespace,
			Labels: map[string]string{
				"app":     deploymentInput.DeploymentName + "-solr",
				"release": deploymentInput.DeploymentName,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":     deploymentInput.DeploymentName + "-solr",
					"release": deploymentInput.DeploymentName,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: deploymentInput.DeploymentName + "-solr",
					Labels: map[string]string{
						"app":     deploymentInput.DeploymentName + "-solr",
						"release": deploymentInput.DeploymentName,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            deploymentInput.SearchEngine,
							Image:           deploymentInput.SearchEngineImage,
							ImagePullPolicy: corev1.PullAlways,
							VolumeMounts: []corev1.VolumeMount{
								createVolumeMount("datadir", "/opt/solr/server/data"),
							},
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 8983,
									Protocol:      corev1.ProtocolTCP,
								},
							},
						},
					},
					RestartPolicy: corev1.RestartPolicyAlways,
					Volumes: []corev1.Volume{
						attachVolumeFromClaim("datadir", "search", deploymentInput),
					},
				},
			},
		},
	}

	// Create Solr Deployment
	_, foundErr := client.AppsV1().Deployments(deploymentInput.Namespace).Get(deploymentInput.DeploymentName+"-solr", metav1.GetOptions{})
	if foundErr != nil {
		resultSolr, errSolr := client.AppsV1().Deployments(deploymentInput.Namespace).Create(solrDeployment)
		if errSolr != nil {
			panic(errSolr)
		}
		log.Printf("Created Solr Deployment - Name: %q, UID: %q\n", resultSolr.GetObjectMeta().GetName(), resultSolr.GetObjectMeta().GetUID())

	}

	databaseServiceName := deploymentInput.DeploymentName + "-solr-svc"
	databaseLabels := map[string]string{
		"app": deploymentInput.DeploymentName + "-solr",
	}

	// Database service.
	solrService := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: databaseServiceName,
			Labels: map[string]string{
				"app":     deploymentInput.DeploymentName,
				"release": deploymentInput.DeploymentName,
			},
		},
		Spec: v1.ServiceSpec{
			Selector: databaseLabels,
			Ports: []v1.ServicePort{{
				Port:       8983,
				TargetPort: intstr.FromInt(8983),
			}},
		},
	}

	_, foundServiceErr := client.CoreV1().Services(deploymentInput.Namespace).Get(deploymentInput.DeploymentName+"-solr-svc", metav1.GetOptions{})
	if foundServiceErr != nil {
		solrService, errSolrService := client.CoreV1().Services(deploymentInput.Namespace).Create(solrService)
		if errSolrService != nil {
			panic(errSolrService)
		}
		log.Printf("Created Solr Service - Name: %q, UID: %q\n", solrService.GetObjectMeta().GetName(), solrService.GetObjectMeta().GetUID())
	}
}
