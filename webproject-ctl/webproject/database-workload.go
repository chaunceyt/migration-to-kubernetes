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

func createDatabaseWorkload(client *kubernetes.Clientset, deploymentInput WebProjectInput) {
	databaseDeployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: appsv1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentInput.DeploymentName + "-db",
			Namespace: deploymentInput.Namespace,
			Labels: map[string]string{
				"app":     deploymentInput.DeploymentName + "-db",
				"release": deploymentInput.DeploymentName,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":     deploymentInput.DeploymentName + "-db",
					"release": deploymentInput.DeploymentName,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: deploymentInput.DeploymentName + "-db",
					Labels: map[string]string{
						"app":     deploymentInput.DeploymentName + "-db",
						"release": deploymentInput.DeploymentName,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            deploymentInput.DatabaseEngine,
							Image:           deploymentInput.DatabaseEngineImage,
							ImagePullPolicy: corev1.PullIfNotPresent,
							Env: []v1.EnvVar{
								{Name: "MYSQL_ROOT_PASSWORD", Value: "admin"},
							},
							VolumeMounts: []corev1.VolumeMount{
								createVolumeMount("database-volume", "/var/lib/mysql"),
							},
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 3306,
									Protocol:      corev1.ProtocolTCP,
								},
							},
						},
					},
					RestartPolicy: corev1.RestartPolicyAlways,
					Volumes: []corev1.Volume{
						attachVolumeFromClaim("database-volume", "db", deploymentInput),
					},
				},
			},
		},
	}

	// Create  Database Deployment
	_, foundErr := client.AppsV1().Deployments(deploymentInput.Namespace).Get(deploymentInput.DeploymentName, metav1.GetOptions{})
	if foundErr != nil {
		log.Println("Creating database deployment...")
		resultDatabase, errDatabase := client.AppsV1().Deployments(deploymentInput.Namespace).Create(databaseDeployment)
		if errDatabase != nil {
			panic(errDatabase)
		}
		log.Printf("Created Database Deployment - Name: %q, UID: %q\n", resultDatabase.GetObjectMeta().GetName(), resultDatabase.GetObjectMeta().GetUID())

	}

	databaseServiceName := deploymentInput.DeploymentName + "-db-svc"
	databaseLabels := map[string]string{
		"app": deploymentInput.DeploymentName + "-db",
	}
	databaseService := &v1.Service{
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
				Port:       3306,
				TargetPort: intstr.FromInt(3306),
			}},
		},
	}

	_, foundServiceErr := client.CoreV1().Services(deploymentInput.Namespace).Get(deploymentInput.DeploymentName+"-db-svc", metav1.GetOptions{})
	if foundServiceErr != nil {
		databaseService, errDatabaseService := client.CoreV1().Services(deploymentInput.Namespace).Create(databaseService)
		if errDatabaseService != nil {
			panic(errDatabaseService)
		}
		log.Printf("Created Database Service - Name: %q, UID: %q\n", databaseService.GetObjectMeta().GetName(), databaseService.GetObjectMeta().GetUID())
	}
}
