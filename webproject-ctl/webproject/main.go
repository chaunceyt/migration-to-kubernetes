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
	CacheEngine              string
	CacheEngineImage         string
	DatabaseEngine           string
	DatabaseEngineImage      string
	DeleteProject            bool
	DeploymentName           string
	FilesMountPoint          string
	IngressDomainName        string
	Namespace                string
	PrimaryContainerName     string
	PrimaryContainerImageTag string
	PrimaryContainerPort     int
	Replicas                 int32
	SearchEngine             string
	SearchEngineImage        string
	SidecarContainerName     string
	SidecarContainerImageTag string
	SidecarContainerPort     int
	WebrootMountPoint        string
}

var usage = `Usage: webproject-ctl [options...]

Options:
  -cache-engine                 Cache Engine [memcached, redis].
  -cache-engine-image           Cache Engine Image [memcached:1.5.20, redis:5.0.7-alpine].
  -database-engine              DatabaseEngine [mysql or mariadb].
  -database-engine-image        Database image name and tag. i.e. mysql:5.7
  -deployment-name              The Deployment name. (required)
                                If using GitLab consider using the RELEASE_NAME
  -domain-name                  Project ingress domain name. (required)
  -files-mount-point            Site files mount point. Default: /var/www/html/sites/default/files
  -namespace                    Project namespace. (required)
  -primary-container-name       Primary Container name. Default: "web"
  -prinary-container-image-tag  Project Container image and tag. (required)
  -primary-container-port       Primary container port. Default: 8080.
  -replicas                     Number of replicas. Default 1.
  -sidecar-container-name       Sidecar Container name. Default: "cli"
  -sidecar-container-image-tag  Sidecar Container image and tag. Default: docksal/cli:2.10-php7.3 
  -sidecar-container-port       Sidecar container port.
  -webroot-mount-point          Webroot mount point. Default: /var/www/html
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

	deploymentName := flag.String("deployment-name", "", "The Deployment name. (required)")

	// Delete a project
	deleteProject := flag.Bool("delete-web-project", false, "Delete WebProject")

	// Create a project
	cacheEngine := flag.String("cache-engine", "", "")
	cacheEngineImage := flag.String("cache-engine-image", "", "")
	databaseEngine := flag.String("database-engine", "", "")
	databaseEngineImage := flag.String("database-engine-image", "", "")

	// Path to files folder. i.e. /var/www/html/sites/default/files
	filesMountPoint := flag.String("files-mount-point", "/var/www/html/sites/default/files", "")

	ingressDomainName := flag.String("domain-name", "", "")

	primaryContainerName := flag.String("primary-container-name", "web", "")
	primaryContainerImageTag := flag.String("primary-container-image-tag", "", "")
	primaryContainerPort := flag.Int("primary-container-port", 8080, "")
	projectNamespace := flag.String("namespace", "", "")

	replicas := flag.Int("replicas", 1, "")

	// SearchEngine [solr, elasticsearch]
	searchEngine := flag.String("search-engine", "", "")
	searchEngineImage := flag.String("search-engine-image", "", "")

	// Sidercar support
	sidecarContainerName := flag.String("sidecar-container-name", "cli", "")
	sidecarContainerImageTag := flag.String("sidecar-container-image-tag", "docksal/cli:2.10-php7.3", "")
	sidecarContainerPort := flag.Int("sidecar-container-port", 9000, "")

	// Path to DOCROOT i.e. /var/www/html
	webrootMountPoint := flag.String("webroot-mount-point", "/var/www/html", "")

	flag.Parse()

	var requiredParameters int

	if *deleteProject {
		requiredParameters = 2
		if flag.NFlag() < requiredParameters {
			fmt.Println("Deleting a project requires the namespace and deploymentname")
		}
	} else {
		// Ensure we get the required parameters.
		requiredParameters := 3
		if flag.NFlag() < requiredParameters {
			usageAndExit("")
		}

	}

	// TODO: If database-engine isset the database-engine-image can not be empty

	deploymentInput := WebProjectInput{
		DeleteProject:            *deleteProject,
		DeploymentName:           *deploymentName,
		FilesMountPoint:          *filesMountPoint,
		PrimaryContainerName:     *primaryContainerName,
		PrimaryContainerImageTag: *primaryContainerImageTag,
		PrimaryContainerPort:     *primaryContainerPort,
		Replicas:                 int32(*replicas),
		IngressDomainName:        *ingressDomainName,
		Namespace:                *projectNamespace,
		CacheEngine:              *cacheEngine,
		CacheEngineImage:         *cacheEngineImage,
		DatabaseEngine:           *databaseEngine,
		DatabaseEngineImage:      *databaseEngineImage,
		SearchEngine:             *searchEngine,
		SearchEngineImage:        *searchEngineImage,
		SidecarContainerName:     *sidecarContainerName,
		SidecarContainerImageTag: *sidecarContainerImageTag,
		SidecarContainerPort:     *sidecarContainerPort,
		WebrootMountPoint:        *webrootMountPoint,
	}

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		errAndExit("Error getting kubeconfig")
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		errAndExit("Error creating client from config provided")
	}

	if *deleteProject {
		deleteWebprojectWorkloadHandler(client, deploymentInput)
	} else {
		createWebProject(client, deploymentInput)
	}
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

func errAndExit(msg string) {
	fmt.Fprintf(os.Stderr, msg)
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}
