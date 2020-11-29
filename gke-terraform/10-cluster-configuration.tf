locals {
  cluster_name            = "developer-cthorn"
  cluster_location        = "us-central1-c"
  cluster_master_version  = "latest"
  cluster_release_channel = "REGULAR"
}

# GKE cluster
resource "google_container_cluster" "cluster" {
  name     = local.cluster_name
  location = local.cluster_location

  // Disable Stackdriver Kubernetes Monitoring
  logging_service    = "none"
  monitoring_service = "none"

  provider = google-beta
  project = data.google_project.project.project_id

  min_master_version = local.cluster_master_version

  release_channel {
    channel = local.cluster_release_channel
  }

  // Start with a single node, since we're going to delete the default pool
  initial_node_count       = 1
  
  // Remove the default node pool
  remove_default_node_pool = true

  // Network config
  network = "default"
  networking_mode = "VPC_NATIVE"
  ip_allocation_policy {}

  master_auth {
    username = ""
    password = ""

    client_certificate_config {
      issue_client_certificate = false
    }
  }

  // Enable google-groups for RBAC
  authenticator_groups_config {
    security_group = "gke-security-groups@cthorn.com"
  }

    // Enable GKE Network Policy
  network_policy {
    enabled  = true
    provider = "CALICO"
  }

  // Configure cluster addons
  addons_config {
    horizontal_pod_autoscaling {
      disabled = true
    }
    http_load_balancing {
      disabled = true
    }
    network_policy_config {
      disabled = false
    }

    // We want to use VolumeSnapshots enabled CSI driver
    gce_persistent_disk_csi_driver_config {
      enabled = true
    }
  }

  // PodSecurityPolicy enforcement
  pod_security_policy_config {
    enabled = false
  }

  // VPA 
  vertical_pod_autoscaling {
    enabled = false
  }

  // GKE clusters are critical objects and should not be destroyed
  // IMPORTANT: should be false on test clusters
  lifecycle {
    prevent_destroy = false
  }

  // Set maintenance time
  maintenance_policy {
    daily_maintenance_window {
      start_time = "11:00" // (in UTC), 03:00 PST
    }
  }

}