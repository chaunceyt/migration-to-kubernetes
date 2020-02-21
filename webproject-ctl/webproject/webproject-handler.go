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
	"k8s.io/client-go/kubernetes"
)

// createWebProject - create all of the desire components.
func createWebProject(client *kubernetes.Clientset, deploymentInput WebProjectInput) {

	// Create Persistent Volume claims first.
	createPersistentVolumeClaim("webfiles", client, deploymentInput)

	// Determine if we are needing to deploy a database.
	var useDatabase bool

	if deploymentInput.DatabaseEngine == "" || deploymentInput.DatabaseEngineImage == "" {
		useDatabase = false
	} else {
		useDatabase = true
	}

	// Create database workload.
	if useDatabase == true {
		createPersistentVolumeClaim("db", client, deploymentInput)
		createDatabaseWorkload(client, deploymentInput)
	}

	// Create cacheEngine.
	switch deploymentInput.CacheEngine {
	case "redis":
		createRedisWorkload(client, deploymentInput)
	case "memcached":
		createMemcachedWorkload(client, deploymentInput)
	default:
		//fmt.Println("Unsupported CacheEngine selected or not defined")
	}

	// Create searchEngine.
	switch deploymentInput.SearchEngine {
	case "solr":
		createSolrWorkload(client, deploymentInput)
	default:
		//fmt.Println("Unsupported SearchEngine selected or not defined")
	}
	// Create project's primary workload.
	createWebprojectWorkload(client, deploymentInput)

	// Setup domain(s) for Webproject
	createIngress(client, deploymentInput)

}
