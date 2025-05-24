package main

import (
	"log"

	"github.com/emarais-godaddy/go-diagrams/diagram"
	"github.com/emarais-godaddy/go-diagrams/nodes/apps"
	"github.com/emarais-godaddy/go-diagrams/nodes/gcp"
)

func main() {
	d, err := diagram.New(
		diagram.Filename("gcp_generative_app"),
		diagram.Direction("LR"), // Left to Right direction
	)
	if err != nil {
		log.Fatal(err)
	}

	// --- Define Nodes ---

	// User (generic client)
	user := apps.Client.User().Label("User")

	// Presentation Tier Nodes
	cdn := gcp.Network.Cdn().Label("Cloud CDN")
	gcsStatic := gcp.Storage.Storage().Label("Cloud Storage (Static Content)")
	// Using LoadBalancing as a placeholder for Apigee if a specific icon is not available, common for gateway patterns.
	apiGateway := gcp.Network.LoadBalancing().Label("Apigee API Gateway") 

	// Application Tier Nodes
	// Using Iam as a placeholder for Google OAuth
	oauth := gcp.Security.Iam().Label("Google OAuth")
	cloudRun := gcp.Compute.Run().Label("Cloud Run Logic")
	// Using AiPlatform and labeling as Vertex AI
	vertexAI := gcp.Ml.AiPlatform().Label("Vertex AI Service")

	// Data Tier Nodes
	gcsData := gcp.Storage.Storage().Label("Cloud Storage (Data)")

	// --- Define Tiers (Groups) ---
	// Top-level group for the diagram title
	diagGroup := diagram.NewGroup("").Label("GCP 3-Tier Generative Application")

	presentationTier := diagram.NewGroup("PresentationTier").Label("Presentation Tier").Add(
		cdn,
		gcsStatic,
		apiGateway,
	)

	applicationTier := diagram.NewGroup("ApplicationTier").Label("Application Tier").Add(
		oauth,
		cloudRun,
		vertexAI,
	)
	
	dataTier := diagram.NewGroup("DataTier").Label("Data Tier").Add(gcsData)

	// Nest tiers under the main diagram group
	diagGroup.Group(presentationTier)
	diagGroup.Group(applicationTier)
	diagGroup.Group(dataTier)


	// --- Add Nodes and Groups to Diagram ---
	d.Add(user) // User is outside the main application tiers
	d.Group(diagGroup)


	// --- Define Connections ---

	// User to Presentation Tier
	d.Connect(user, cdn)
	d.Connect(cdn, gcsStatic, diagram.EdgeLabel("Serves Static Content"))
	d.Connect(user, apiGateway, diagram.EdgeLabel("API Requests")) // User can also go direct to APIGW

	// Presentation Tier to Application Tier
	d.Connect(apiGateway, oauth, diagram.EdgeLabel("Auth"))
	d.Connect(apiGateway, cloudRun, diagram.EdgeLabel("Invoke Backend"))
	
	// Application Tier Connections
	d.Connect(oauth, cloudRun, diagram.EdgeLabel("Authorizes")) // OAuth validates token for Cloud Run
	d.Connect(cloudRun, vertexAI, diagram.EdgeLabel("Generative AI Tasks"))

	// Application Tier to Data Tier
	d.Connect(cloudRun, gcsData, diagram.EdgeLabel("Store/Retrieve Data"))


	// --- Render Diagram ---
	if err := d.Render(); err != nil {
		log.Fatal(err)
	}
	log.Println("Diagram generated. Check the 'go-diagrams' directory for the gcp_generative_app.dot file.")
}
