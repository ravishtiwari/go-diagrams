# Go-Diagrams

## Fast and easy application diagrams

Go-Diagrams is a loose port of [diagrams](https://github.com/mingrammer/diagrams).

## Contents

- [Features](#features)
- [Usage](#usage)

## Features

Turn this:

```go
d, err := diagram.New(diagram.Filename("app"), diagram.Label("App"), diagram.Direction("LR"))
if err != nil {
    log.Fatal(err)
}

dns := gcp.Network.Dns(diagram.NodeLabel("DNS"))
lb := gcp.Network.LoadBalancing(diagram.NodeLabel("NLB"))
cache := gcp.Database.Memorystore(diagram.NodeLabel("Cache"))
db := gcp.Database.Sql(diagram.NodeLabel("Database"))

dc := diagram.NewGroup("GCP")
dc.NewGroup("services").
    Label("Service Layer").
    Add(
        gcp.Compute.ComputeEngine(diagram.NodeLabel("Server 1")),
        gcp.Compute.ComputeEngine(diagram.NodeLabel("Server 2")),
        gcp.Compute.ComputeEngine(diagram.NodeLabel("Server 3")),
    ).
    ConnectAllFrom(lb.ID(), diagram.Forward()).
    ConnectAllTo(cache.ID(), diagram.Forward())

dc.NewGroup("data").Label("Data Layer").Add(cache, db).Connect(cache, db)

d.Connect(dns, lb, diagram.Forward()).Group(dc)

if err := d.Render(); err != nil {
    log.Fatal(err)
}
```

Into this:

![app-diagram](images/app-diagram.png)

## Usage

```sh
go get github.com/emarais-godaddy/go-diagrams
```

Create a diagram:

```go
d, err := diagram.New(diagram.Label("my-diagram"), diagram.Filename("diagram"))
if err != nil {
    log.Fatal(err)
}

fw := generic.Network.Firewall().Label("fw")
sw := generic.Network.Switch().Label("sw")

d.Connect(fw, sw)
```

Render the output:

```go
if err := d.Render(); err != nil {
    log.Fatal(err)
}
```

Go-Diagrams will create a folder in the current working directory with the graphviz DOT file and any image assets.

Create an ouput image with any graphviz compatible renderer:

```sh
dot -Tpng diagram.dot > diagram.png
```

## Development

This project uses a Makefile to streamline common development tasks. Here are the available targets:

*   `make all`: (Default) Builds the `go-diagrams` CLI tool.
*   `make compile`: Compiles the Go code for all packages.
*   `make validate`: Validates the Go code using `go vet`.
*   `make build`: Builds the `go-diagrams` CLI tool (output: `go-diagrams`).
*   `make test`: Runs all Go tests.
*   `make package`: Creates a gzipped tar archive of the source code (`go-diagrams-src.tar.gz`).
*   `make run-examples`: Runs all examples in the `_examples` directory. For examples with a `taskfile.yml`, it also runs `task -t taskfile.yml all`.
*   `make clean`: Removes build artifacts, the packaged archive, and example output directories.
*   `make help`: Displays a list of available Makefile targets.

## Examples

The `_examples` directory contains various examples demonstrating how to use Go-Diagrams to create diagrams for different architectures.

### 3-Tier Generative Application on AWS

*   **File**: `_examples/aws/generative_app.go`
*   **Description**: This example illustrates a common 3-tier architecture for a generative AI application on AWS.
*   **Architecture**:
    *   **Presentation Tier**: Amazon CloudFront (CDN), Amazon S3 (for static content), Amazon API Gateway.
    *   **Application Tier**: AWS Lambda (for business logic), Amazon BedRock (for generative AI capabilities), Amazon Cognito (for user authentication).
    *   **Data Tier**: Amazon S3 (for data storage).
    *   A User node interacts with CloudFront and API Gateway.

### 3-Tier Generative Application on GCP

*   **File**: `_examples/gcp/generative_app.go`
*   **Description**: This example illustrates a common 3-tier architecture for a generative AI application on Google Cloud Platform.
*   **Architecture**:
    *   **Presentation Tier**: Google Cloud CDN, Google Cloud Storage (for static content), Apigee API Gateway (represented by a Load Balancer icon).
    *   **Application Tier**: Google Cloud Run (for business logic), Google Vertex AI (for generative AI capabilities), Google OAuth (represented by an IAM icon for user authentication).
    *   **Data Tier**: Google Cloud Storage (for data storage).
    *   A User node interacts with Cloud CDN and Apigee API Gateway.

### GCP: Hosting Gemma 8b on GKE

*   **File**: `_examples/gcp/gemma_on_gke.go`
*   **Description**: This example demonstrates hosting the Gemma 8b large language model on Google Kubernetes Engine (GKE).
*   **Architecture**:
    *   An External User connects to a GCP Global Load Balancer.
    *   The Load Balancer routes traffic to a Kubernetes Service within a GKE Cluster.
    *   The GKE Cluster contains a Node Pool where the Gemma 8b model is deployed using a Kubernetes Deployment, exposed by the Kubernetes Service.

### GCP: GoLang Backend on GKE Autopilot with Cloud SQL

*   **File**: `_examples/gcp/go_backend_gke_autopilot.go`
*   **Description**: This example showcases a GoLang backend application running on GKE Autopilot, connected to a Google Cloud SQL (PostgreSQL) database.
*   **Architecture**:
    *   An External User connects to a GCP External HTTP/S Load Balancer.
    *   The Load Balancer directs traffic to a Kubernetes Ingress resource within a GKE Autopilot cluster.
    *   The Ingress routes requests to a Kubernetes Service (ClusterIP type).
    *   The Service selects Pods from a GoLang Backend Deployment.
    *   The GoLang backend application connects to a Google Cloud SQL (PostgreSQL) instance for data persistence.
