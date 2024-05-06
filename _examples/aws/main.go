package main

import (
	"log"

	"github.com/emarais-godaddy/go-diagrams/diagram"
	"github.com/emarais-godaddy/go-diagrams/nodes/aws"
)

func main() {
	attrs := map[string]string{
		"ranksep":     "1.5",
		"compound":    "true",
		"concentrate": "true",
		"splines":     "true",
	}

	d, err := diagram.New(diagram.Filename("diagram"), diagram.WithAttributes(attrs), diagram.Label("background process"), diagram.Direction("LR"))
	if err != nil {
		log.Fatal(err)
	}

	// Elements
	eventbridge := aws.Integration.Eventbridge(diagram.NodeLabel("default"))
	sqs := aws.Integration.SimpleQueueServiceSqs(diagram.NodeLabel("sqs"), diagram.Name("sqs"))
	scheduleLambda := aws.Compute.Lambda(diagram.NodeLabel("lambda"))
	processLambda := aws.Compute.Lambda(diagram.NodeLabel("lambda"))
	db := aws.Database.Aurora(diagram.NodeLabel("rds"))

	dbGrp := diagram.NewGroup("database").
		Label("aurora-projects").
		Add(db)

	sGrp := diagram.NewGroup("scheduler").
		Label("scheduler-lambda").
		Add(eventbridge, scheduleLambda).
		Connect(eventbridge, scheduleLambda, diagram.EdgeLabel("5min interval")).
		Connect(scheduleLambda, sqs, diagram.EdgeLabel("Batches of 100"))

	pGrp := diagram.NewGroup("processor").
		Label("processor-lambda").
		Add(sqs, processLambda).
		Connect(sqs, processLambda, diagram.EdgeLabel("invokes"))

	d.Connect(scheduleLambda, db, diagram.SnapToGroup(dbGrp), diagram.Bidirectional())
	d.Connect(processLambda, db, diagram.SnapToGroup(dbGrp), diagram.Bidirectional())

	d.Group(sGrp).Group(pGrp).Group(dbGrp)

	if err := d.Render(); err != nil {
		log.Fatal(err)
	}
}
