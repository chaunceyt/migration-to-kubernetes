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
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// deleteWebprojectWorkloadHandler - delete all of the objects related to a specific deployment name.
func deleteWebprojectWorkloadHandler(client *kubernetes.Clientset, deploymentInput WebProjectInput) {

	/*
	 * Delete all of the objects related to the RELEASE_NAME
	 */

	// Delete cache workload.
	deleteRedisWorkload(client, deploymentInput)
	deleteMemcachedWorkload(client, deploymentInput)

	// Delete search workload.
	deleteSolrWorkload(client, deploymentInput)

	// Delete Web project workload.
	deleteWebProjectWorkload(client, deploymentInput)

	// Delete Database workload.
	deleteDatabaseWorkload(client, deploymentInput)

}

// deletePVC - delete persistent volume claim.
func deletePVC(client *kubernetes.Clientset, deploymentInput WebProjectInput, pvcType string) {
	// Delete PVC based on pvcType db and webfiles
	_, foundErr := client.CoreV1().PersistentVolumeClaims(deploymentInput.Namespace).Get(deploymentInput.DeploymentName+"-"+pvcType+"-pvc", metav1.GetOptions{})
	if foundErr != nil {
		// Do nothing.
	} else {
		if err := client.CoreV1().PersistentVolumeClaims(deploymentInput.Namespace).Delete(deploymentInput.DeploymentName+"-"+pvcType+"-pvc", &metav1.DeleteOptions{}); err != nil {
			log.Println(fmt.Errorf("Error while deleting PVC - %v", err))
		}
		log.Println("Deleted PVC " + deploymentInput.DeploymentName + "-" + pvcType + "-pvc")

	}
}

// deleteRedisWorload - delete service and deployment.
func deleteRedisWorkload(client *kubernetes.Clientset, deploymentInput WebProjectInput) {
	// Delete redis serivce
	_, foundServiceErr := client.CoreV1().Services(deploymentInput.Namespace).Get(deploymentInput.DeploymentName+"-redis-svc", metav1.GetOptions{})
	if foundServiceErr != nil {
		// Do nothing.
	} else {
		if serviceErr := client.CoreV1().Services(deploymentInput.Namespace).Delete(deploymentInput.DeploymentName+"-redis-svc", &metav1.DeleteOptions{}); serviceErr != nil {
			log.Println("The Redis Service " + deploymentInput.DeploymentName + "-redis-svc")
		}
	}

	// Delete Redis deployment.
	_, foundErr := client.AppsV1().Deployments(deploymentInput.Namespace).Get(deploymentInput.DeploymentName+"-redis", metav1.GetOptions{})
	if foundErr != nil {
		// Do nothing.
	} else {
		if err := client.AppsV1().Deployments(deploymentInput.Namespace).Delete(deploymentInput.DeploymentName+"-redis", &metav1.DeleteOptions{}); err != nil {
			log.Println(fmt.Errorf("Error while Redis Deployment - %v", err))
		}
		log.Println("Deleted Redis Deployment " + deploymentInput.DeploymentName + "-redis")

	}
}

// deleteMemcachedWorkload - delete service and deployment.
func deleteMemcachedWorkload(client *kubernetes.Clientset, deploymentInput WebProjectInput) {
	// Delete Memcached service.
	_, foundServiceErr := client.CoreV1().Services(deploymentInput.Namespace).Get(deploymentInput.DeploymentName+"-memcached-svc", metav1.GetOptions{})
	if foundServiceErr != nil {
		// Do nothing.
	} else {
		if serviceErr := client.CoreV1().Services(deploymentInput.Namespace).Delete(deploymentInput.DeploymentName+"-memcached-svc", &metav1.DeleteOptions{}); serviceErr != nil {
			log.Println("The Memcached Service " + deploymentInput.DeploymentName + "-memcached-svc")
		}
	}

	// Delete Memcached deployment.
	_, foundErr := client.AppsV1().Deployments(deploymentInput.Namespace).Get(deploymentInput.DeploymentName+"-memcached", metav1.GetOptions{})
	if foundErr != nil {
		// Do nothing.
	} else {
		if err := client.AppsV1().Deployments(deploymentInput.Namespace).Delete(deploymentInput.DeploymentName+"-memcached", &metav1.DeleteOptions{}); err != nil {
			log.Println(fmt.Errorf("Error while Memcache Deployment - %v", err))
		}
		log.Println("Deleted Memcached Deployment " + deploymentInput.DeploymentName + "-memcached")
	}
}

// deleteSolrWorkload - delete service and deployment.
func deleteSolrWorkload(client *kubernetes.Clientset, deploymentInput WebProjectInput) {
	// Delete Solr service.
	_, foundServiceErr := client.CoreV1().Services(deploymentInput.Namespace).Get(deploymentInput.DeploymentName+"-solr-svc", metav1.GetOptions{})
	if foundServiceErr != nil {
		// Do nothing.
	} else {
		if serviceErr := client.CoreV1().Services(deploymentInput.Namespace).Delete(deploymentInput.DeploymentName+"-solr-svc", &metav1.DeleteOptions{}); serviceErr != nil {
			log.Println("The Solr Service " + deploymentInput.DeploymentName + "-solr-svc")
		}
	}
	_, foundErr := client.AppsV1().Deployments(deploymentInput.Namespace).Get(deploymentInput.DeploymentName+"-solr", metav1.GetOptions{})
	if foundErr != nil {
		// Do nothing.
	} else {
		if err := client.AppsV1().Deployments(deploymentInput.Namespace).Delete(deploymentInput.DeploymentName+"-solr", &metav1.DeleteOptions{}); err != nil {
			log.Println(fmt.Errorf("Error while deleting Solr Deployment - %v", err))
		}
		log.Println("Deleted Solr Deployment " + deploymentInput.DeploymentName + "-solr")
	}
}

// deleteWebProjectWorkload - delete ingress, service deployment and pvc.
func deleteWebProjectWorkload(client *kubernetes.Clientset, deploymentInput WebProjectInput) {
	// Delete ingress
	_, foundIngressErr := client.ExtensionsV1beta1().Ingresses(deploymentInput.Namespace).Get(deploymentInput.DeploymentName+"-ing", metav1.GetOptions{})
	if foundIngressErr != nil {
		// Do nothing.
	} else {
		if ingressErr := client.ExtensionsV1beta1().Ingresses(deploymentInput.Namespace).Delete(deploymentInput.DeploymentName+"-ing", &metav1.DeleteOptions{}); ingressErr != nil {
			log.Println(fmt.Errorf("Error while Deleting Ingress - %v", ingressErr))
		}
	}

	// Delete service
	_, foundServiceErr := client.CoreV1().Services(deploymentInput.Namespace).Get(deploymentInput.DeploymentName+"-svc", metav1.GetOptions{})
	if foundServiceErr != nil {
		// Do nothing.
	} else {
		if serviceErr := client.CoreV1().Services(deploymentInput.Namespace).Delete(deploymentInput.DeploymentName+"-svc", &metav1.DeleteOptions{}); serviceErr != nil {
			log.Println("The Webproject Service " + deploymentInput.DeploymentName + "-svc")
		}
	}
	// Delete web deployment
	_, foundErr := client.AppsV1().Deployments(deploymentInput.Namespace).Get(deploymentInput.DeploymentName, metav1.GetOptions{})
	if foundErr != nil {
		// Do nothing.
	} else {
		if err := client.AppsV1().Deployments(deploymentInput.Namespace).Delete(deploymentInput.DeploymentName, &metav1.DeleteOptions{}); err != nil {
			log.Println(fmt.Errorf("Error while deleting WebProject Deployment - %v", err))
		}
		log.Println("Deleted WebProject Deployment " + deploymentInput.DeploymentName)
	}

	// Delete webfiles PVC
	deletePVC(client, deploymentInput, "webfiles")
}

// deleteDatabaseWorkload - delete service, deployment and pvc.
func deleteDatabaseWorkload(client *kubernetes.Clientset, deploymentInput WebProjectInput) {
	// Delete database service
	_, foundServiceErr := client.CoreV1().Services(deploymentInput.Namespace).Get(deploymentInput.DeploymentName+"-db-svc", metav1.GetOptions{})
	if foundServiceErr != nil {
		// Do nothing.
	} else {
		if serviceErr := client.CoreV1().Services(deploymentInput.Namespace).Delete(deploymentInput.DeploymentName+"-db-svc", &metav1.DeleteOptions{}); serviceErr != nil {
			log.Println("The Database Service " + deploymentInput.DeploymentName + "-db-svc")
		}
	}

	// Delete database deployment
	if err := client.AppsV1().Deployments(deploymentInput.Namespace).Delete(deploymentInput.DeploymentName+"-db", &metav1.DeleteOptions{}); err != nil {
		log.Println(fmt.Errorf("Error while deleting PVC - %v", err))
	}
	log.Println("Deleted Database Deployment" + deploymentInput.DeploymentName + "-db")

	// Delete db PVC
	deletePVC(client, deploymentInput, "db")

}
