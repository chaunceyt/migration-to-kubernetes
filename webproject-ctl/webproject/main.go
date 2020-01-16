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

//@TODO
// - Check to see if object exists and update it. Similar to the kubectl apply -f filename.yaml
// - Sidecar support

import (
	"flag"
	"fmt"
	"os"
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

var usage = `Usage: webproject-ctl [options...]

Options:
  -cache-engine                 Cache Engine [memcached, redis].
  -database-engine              DatabaseEngine [mysql or mariadb].
  -database-engine-image        Database image name and tag. i.e. mysql:5.7
  -deployment-name              The Deployment name. (required)
                                If using GitLab consider using the RELEASE_NAME
  -domain-name                  Project ingress domain name. (required)
  -namespace                    Project namespace. (required)
  -primary-container-name       Primary Container name. Default "web"
  -prinary-container-image-tag  Project Container image and tag. (required)
  -primary-container-port       Primary container port. Default 8080.
  -replicas                     Number of replicas. Default 1.
`

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, fmt.Sprintf(usage))
	}

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	deploymentName := flag.String("deployment-name", "", "")
	primaryContainerName := flag.String("primary-container-name", "web", "")
	primaryContainerImageTag := flag.String("primary-container-image-tag", "", "")
	primaryContainerPort := flag.Int("primary-container-port", 8080, "")
	// Sidercar support
	// sidecarContainerName := flag.String("sidecar-container-name", "cli", "")
	// sidecarContainerImageTag := flag.String("sidecar-container-image-tag", "", "")
	// sidecarContainerPort := flag.Int("sidecar-container-port", 8080, "")
	replicas := flag.Int("replicas", 1, "")
	ingressDomainName := flag.String("domain-name", "", "")
	projectNamespace := flag.String("namespace", "", "")
	cacheEngine := flag.String("cache-engine", "", "")
	// Container image and tag i.e. redis:5.0.7-alpine or memcached:1.5.20-alpine
	// cacheEngineImage := flag.String("cache-engine-image", "", "")
	databaseEngine := flag.String("database-engine", "", "")
	databaseEngineImage := flag.String("database-engine-image", "", "")
	// Path to DOCROOT i.e. /var/www/html
	// webrootMountPoint := flag.String("webroot-mount-point", "", "")

	// Path to files folder. i.e. /var/www/html/sites/default/files
	// filesMountPoint := flag.String("files-mount-point", "", "")

	flag.Parse()
	// We currently have 6 required parameters.
	if flag.NFlag() < 6 {
		usageAndExit("")
	}

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

func usageAndExit(msg string) {
	if msg != "" {
		fmt.Fprintf(os.Stderr, msg)
		fmt.Fprintf(os.Stderr, "\n\n")
	}
	flag.Usage()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}
