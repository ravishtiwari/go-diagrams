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
		diagram.Filename("gcp_gemma_on_gke"),
		diagram.Direction("TB"), // Top to Bottom flow often works well for this
	)
	if err != nil {
		log.Fatal(err)
	}

	// --- Define Nodes ---

	// External User
	user := apps.Client.User().Label("External User")

	// GCP Components
	gcpLB := gcp.Network.LoadBalancing().Label("GCP Global\nLoad Balancer")
	
	// Kubernetes Components for GKE
	gkeNode := gcp.Compute.KubernetesEngine().Label("GKE Control Plane") // Represents the GKE service itself
	k8sSvcNode := k8s.Network.Svc().Label("K8s Service\n(LoadBalancer/NodePort)") // Corrected: Svc()
	gemmaDeploymentNode := k8s.Compute.Deploy().Label("Gemma 8b\nDeployment")   // Corrected: Deploy()

	// --- Define Groups ---

	// Top-level group for the diagram title
	diagGroup := diagram.NewGroup("").Label("Gemma 8b on GKE Architecture")

	// GKE Cluster Group
	gkeClusterGroup := diagram.NewGroup("GKECluster").Label("GKE Cluster")
	gkeClusterGroup.Add(gkeNode) // Add GKE icon to represent the cluster service

		// Node Pool Group (within GKE Cluster)
		nodePoolGroup := diagram.NewGroup("NodePool").Label("GKE Node Pool")
		nodePoolGroup.Add(
			k8sSvcNode,          
			gemmaDeploymentNode,
		)
	
	gkeClusterGroup.Group(nodePoolGroup) 
	diagGroup.Group(gkeClusterGroup)     

	// --- Add Nodes and Groups to Diagram ---
	d.Add(user, gcpLB)
	d.Group(diagGroup) 

	// --- Define Connections ---
	d.Connect(user, gcpLB)
	d.Connect(gcpLB, k8sSvcNode, diagram.EdgeLabel("Routes Traffic To")) 
	d.Connect(k8sSvcNode, gemmaDeploymentNode, diagram.EdgeLabel("Exposes"))  

	// --- Render Diagram ---
	if err := d.Render(); err != nil {
		log.Fatal(err)
	}
	log.Println("Diagram generated. Check the 'go-diagrams' directory for the gcp_gemma_on_gke.dot file.")
}
