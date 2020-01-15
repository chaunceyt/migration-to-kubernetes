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

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// createPersistentVolumeClaim - creates persistent volume claim
func createPersistentVolumeClaim(pvcType string, client *kubernetes.Clientset, deploymentInput WebProjectInput) {
	pvc := &corev1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolumeClaim",
			APIVersion: corev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   deploymentInput.DeploymentName + "-" + pvcType + "-pvc",
			Labels: genDefaultLabels(deploymentInput),
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse("1Gi"),
				},
			},
		},
	}

	// Check to see if persistent volume claim exists already.
	_, foundErr := client.CoreV1().PersistentVolumeClaims(deploymentInput.Namespace).Get(deploymentInput.DeploymentName+"-"+pvcType+"-pvc", metav1.GetOptions{})
	if foundErr != nil {
		log.Println("Creating pvc...")
		result, err := client.CoreV1().PersistentVolumeClaims(deploymentInput.Namespace).Create(pvc)
		if err != nil {
			panic(err)
		}
		log.Printf("Created PVC - Name: %q, UID: %q\n", result.GetObjectMeta().GetName(), result.GetObjectMeta().GetUID())
	}
}

// createEmptyDirVolume
func createEmptyDirVolume(name string) corev1.Volume {
	return corev1.Volume{
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
		Name: name,
	}
}

// createVolumeMount - should be used to attach a volume to a deplopyment.
func createVolumeMount(name string, mountPath string) corev1.VolumeMount {
	return corev1.VolumeMount{
		Name:      name,
		MountPath: mountPath,
	}
}

// attachVolumeClaim attaches the Persistent volume to container.
func attachVolumeFromClaim(name string, pvType string, deploymentInput WebProjectInput) corev1.Volume {
	return corev1.Volume{
		VolumeSource: corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: deploymentInput.DeploymentName + "-" + pvType + "-pvc",
			},
		},
		Name: name,
	}
}
