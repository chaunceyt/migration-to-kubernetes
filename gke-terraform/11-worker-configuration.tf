locals {
  worker_pool_max_nodes = 10
  worker_pool_node_role = "worker"
  worker_pool_disk_type = "pd-standard"
}

# Separately Managed Node Pool
resource "google_container_node_pool" "worker_pool" {
  project    = data.google_project.project.project_id
  provider   = google
  name       = "${google_container_cluster.cluster.name}-${local.worker_pool_node_role}-pool"
  location   = google_container_cluster.cluster.location
  cluster    = google_container_cluster.cluster.name

  // Start with a single node
  initial_node_count = 1

  // node_count = 5
  // Autoscale the cluster as needed.
  autoscaling {
    min_node_count = 1
    max_node_count = local.worker_pool_max_nodes
  }

  node_config {
    oauth_scopes = [
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
    ]

    labels = {
      env = data.google_project.project.project_id
      purpose = "${local.worker_pool_node_role}-node"
    }

    machine_type = "n1-standard-1"
    disk_size_gb       = 200
    disk_type           = "${local.worker_pool_disk_type}"
    image_type         = "COS"
    tags         = ["gke-node", "${data.google_project.project.project_id}-${local.worker_pool_node_role}-node"]
    metadata = {
      disable-legacy-endpoints = "true"
    }
  }

  // Create the new one before destroying the old one.
  lifecycle {
    create_before_destroy = true
  }  
}
