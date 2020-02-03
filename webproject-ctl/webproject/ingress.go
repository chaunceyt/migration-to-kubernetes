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

	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

// createIngress - create an ingress object for the web project serivce.
// TODO: add annotations for ingress controllers. nginx + proxy auth
func createIngress(client *kubernetes.Clientset, deploymentInput WebProjectInput) {
	ingress := &v1beta1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentInput.DeploymentName + "-ing",
			Namespace: deploymentInput.Namespace,
			Labels: map[string]string{
				"app":     deploymentInput.DeploymentName,
				"release": deploymentInput.DeploymentName,
			},
		},
		Spec: v1beta1.IngressSpec{
			Rules: []v1beta1.IngressRule{
				{
					Host: deploymentInput.IngressDomainName,
					IngressRuleValue: v1beta1.IngressRuleValue{
						HTTP: &v1beta1.HTTPIngressRuleValue{
							Paths: []v1beta1.HTTPIngressPath{
								{
									Path: "/",
									Backend: v1beta1.IngressBackend{
										ServiceName: deploymentInput.DeploymentName + "-svc",
										ServicePort: intstr.FromInt(deploymentInput.PrimaryContainerPort),
									},
								},
							},
						},
					},
				},
			},
		},
	}

	_, foundErr := client.ExtensionsV1beta1().Ingresses(deploymentInput.Namespace).Get(deploymentInput.DeploymentName+"-ing", metav1.GetOptions{})
	if foundErr != nil {
		result, errIngress := client.ExtensionsV1beta1().Ingresses(deploymentInput.Namespace).Create(ingress)
		if errIngress != nil {
			panic(errIngress)
		}
		log.Printf("Created Memcahed Deployment - Name: %q, UID: %q\n", result.GetObjectMeta().GetName(), result.GetObjectMeta().GetUID())

	}
}
