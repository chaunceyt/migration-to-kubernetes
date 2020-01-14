// Copyright © 2018 Chauncey Thorn <chaunceyt@gmail.com>
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

//@TODO
// - Check to see if object exists and update it. Similar to the kubectl apply -f filename.yaml
// - Sidecar support

import (
	"flag"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// WebProjectInput struct.
type WebProjectInput struct {
	DeploymentName           string
	PrimaryContainerName     string
	PrimaryContainerImageTag string
	PrimaryContainerPort     int
	Replicas                 int32
	Namespace                string
	CacheEngine              string
	DatabaseEngine           string
	DatabaseEngineImage      string
	IngressDomainName        string
}

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	deploymentName := flag.String("deployment-name", "", "Webproject Name")
	primaryContainerName := flag.String("primary-container-name", "", "Primary Container name")
	primaryContainerImageTag := flag.String("prinary-container-image-tag", "", "Primary Container image and tag")
	primaryContainerPort := flag.Int("primary-container-port", 8080, "Primary container port")
	replicas := flag.Int("replicas", 1, "Number of replicas")
	ingressDomainName := flag.String("domain-name", "", "Domainname for workload")
	projectNamespace := flag.String("namespace", "", "Project Namespace")
	cacheEngine := flag.String("cache-engine", "", "CacheEngine [memcached or redis]")
	databaseEngine := flag.String("database-engine", "", "DatabaseEngine [mysql or mariadb]")
	databaseEngineImage := flag.String("database-engine-image", "", "Image name and tag")

	flag.Parse()

	deploymentInput := WebProjectInput{
		DeploymentName:           *deploymentName,
		PrimaryContainerName:     *primaryContainerName,
		PrimaryContainerImageTag: *primaryContainerImageTag,
		PrimaryContainerPort:     *primaryContainerPort,
		Replicas:                 int32(*replicas),
		IngressDomainName:        *ingressDomainName,
		Namespace:                *projectNamespace,
		CacheEngine:              *cacheEngine,
		DatabaseEngine:           *databaseEngine,
		DatabaseEngineImage:      *databaseEngineImage,
	}

	//	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	//client, err := dynamic.NewForConfig(config)
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	createWebProject(client, deploymentInput)
}

func int32ptr(i int32) *int32 {
	return &i
}
