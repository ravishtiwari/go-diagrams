package main

import (
	"log"

	"github.com/emarais-godaddy/go-diagrams/diagram"
	"github.com/emarais-godaddy/go-diagrams/nodes/apps"
	"github.com/emarais-godaddy/go-diagrams/nodes/aws"
)

func main() {
	d, err := diagram.New(
		// diagram.Label("AWSTieredApp"), // Removed Label to use default directory
		diagram.Filename("aws_generative_app"), 
		diagram.Direction("LR"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// --- Define Nodes ---
	user := apps.Client.User().Label("User")
	cloudfront := aws.Network.Cloudfront().Label("CloudFront CDN")
	s3Static := aws.Storage.SimpleStorageServiceS3().Label("S3 Static Content")
	apiGateway := aws.Network.ApiGateway().Label("API Gateway")
	cognito := aws.Security.Cognito().Label("Cognito")
	lambdaLogic := aws.Compute.Lambda().Label("Lambda Logic")
	bedrock := aws.Ml.MachineLearning().Label("BedRock Service")
	s3Data := aws.Storage.SimpleStorageServiceS3().Label("S3 Data Storage")

	// --- Define Tiers (Groups) ---
	// Add a top-level group to give the diagram a title, since Label was removed
	diagGroup := diagram.NewGroup("").Label("AWS 3-Tier Generative Application") // Using a group for overall title

	presentationTier := diagram.NewGroup("PresentationTier").Label("Presentation Tier").Add(
		cloudfront,
		s3Static,
		apiGateway,
	)
	applicationTier := diagram.NewGroup("ApplicationTier").Label("Application Tier").Add(
		cognito,
		lambdaLogic,
		bedrock,
	)
	dataTier := diagram.NewGroup("DataTier").Label("Data Tier").Add(s3Data)

	diagGroup.Group(presentationTier)
	diagGroup.Group(applicationTier)
	diagGroup.Group(dataTier)

	// --- Add Nodes and Groups to Diagram ---
	d.Add(user) // User is outside the main application tiers but connects to it
	d.Group(diagGroup)


	// --- Define Connections ---
	d.Connect(user, cloudfront)
	d.Connect(user, apiGateway, diagram.EdgeLabel("API Requests"))
	d.Connect(cloudfront, s3Static, diagram.EdgeLabel("Serves Static Content"))
	d.Connect(cloudfront, apiGateway, diagram.EdgeLabel("Routes API Calls"))
	d.Connect(apiGateway, cognito, diagram.EdgeLabel("Auth"))
	d.Connect(apiGateway, lambdaLogic, diagram.EdgeLabel("Invoke"))
	d.Connect(cognito, lambdaLogic, diagram.EdgeLabel("Authorizes")) 
	d.Connect(lambdaLogic, bedrock, diagram.EdgeLabel("Generative AI Tasks"))
	d.Connect(lambdaLogic, s3Data, diagram.EdgeLabel("Store/Retrieve Data"))

	// --- Render Diagram ---
	if err := d.Render(); err != nil {
		log.Fatal(err)
	}
	log.Println("Diagram generated. Check the 'go-diagrams' directory for the aws_generative_app.dot file.")
}
