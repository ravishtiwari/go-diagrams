package main

import (
	"log"

	"github.com/emarais-godaddy/go-diagrams/diagram"
	"github.com/emarais-godaddy/go-diagrams/nodes/apps"
	"github.com/emarais-godaddy/go-diagrams/nodes/gcp"
	"github.com/emarais-godaddy/go-diagrams/nodes/k8s"
)

func main() {
	d, err := diagram.New(
		diagram.Filename("gcp_go_backend_gke_autopilot"),
		diagram.Direction("LR"), // Left to Right is common for request flow
	)
	if err != nil {
		log.Fatal(err)
	}

	// --- Define Nodes ---

	// External User
	user := apps.Client.User().Label("External User")

	// GCP Network Components
	gcpLB := gcp.Network.LoadBalancing().Label("GCP External HTTP/S\nLoad Balancer")

	// GKE Autopilot Cluster and its internal components
	gkeAutopilotNode := gcp.Compute.KubernetesEngine().Label("GKE Autopilot") // Icon representing the cluster
	k8sIngress := k8s.Network.Ing().Label("K8s Ingress")
	k8sService := k8s.Network.Svc().Label("K8s Service (ClusterIP)")
	goBackendDeployment := k8s.Compute.Deploy().Label("GoLang Backend\nDeployment")

	// GCP Database
	cloudSQL := gcp.Database.Sql().Label("Cloud SQL (PostgreSQL)")

	// --- Define Groups ---

	// Top-level group for the diagram title
	diagGroup := diagram.NewGroup("").Label("Go Backend on GKE Autopilot with Cloud SQL")

	// GKE Autopilot Cluster Group
	// This group will contain all K8s resources. The gkeAutopilotNode can be part of it or outside.
	// For Autopilot, the cluster itself is managed, so we group K8s resources within its conceptual boundary.
	gkeClusterGroup := diagram.NewGroup("GKEAutopilotCluster").Label("GKE Autopilot Cluster")
	gkeClusterGroup.Add(
		gkeAutopilotNode, // Placing the GKE icon within the group boundary as well
		k8sIngress,
		k8sService,
		goBackendDeployment,
	)
	
	diagGroup.Group(gkeClusterGroup)
	// The Cloud SQL instance is outside the GKE cluster group but within the main diagram group
	diagGroup.Add(cloudSQL)


	// --- Add Nodes and Groups to Diagram ---
	// User and LB are outside the main GKE cluster group.
	d.Add(user, gcpLB) 
	d.Group(diagGroup) // This adds the GKE cluster group and Cloud SQL

	// --- Define Connections ---
	d.Connect(user, gcpLB)
	d.Connect(gcpLB, k8sIngress, diagram.EdgeLabel("Routes Traffic To Ingress"))
	d.Connect(k8sIngress, k8sService, diagram.EdgeLabel("Routes To Service"))
	d.Connect(k8sService, goBackendDeployment, diagram.EdgeLabel("Selects Pods"))
	d.Connect(goBackendDeployment, cloudSQL, diagram.EdgeLabel("DB Connection (e.g., Cloud SQL Proxy or Private IP)"))

	// --- Render Diagram ---
	if err := d.Render(); err != nil {
		log.Fatal(err)
	}
	log.Println("Diagram generated. Check the 'go-diagrams' directory for the gcp_go_backend_gke_autopilot.dot file.")
}
