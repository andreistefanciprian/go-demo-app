# go-demo-app

A simple Go web application that demonstrates HTTP status code monitoring with Prometheus metrics. The app provides endpoints to return custom HTTP status codes and exposes metrics for monitoring request patterns.

## Features
- **Hello World endpoint** (`/`) - Returns a greeting with hostname
- **Custom Status endpoint** (`/status?code=XXX`) - Returns any HTTP status code (100-599)
- **Prometheus metrics** (`/metrics`) - Tracks HTTP requests by status code
- **Kubernetes ready** - Includes Helm chart with ServiceMonitor for Prometheus scraping

This repo has github actions workflows for pushing container image and helm chart to Google Artifact Registry (GAR).
* Terraform code for building the GAR is [here](https://github.com/andreistefanciprian/terraform-kubernetes-gke-cluster).
* K8s deployment code via flux [here](https://github.com/andreistefanciprian/flux-demo/blob/main/clusters/home/go-demo-app.yaml)

## Run Locally with Docker Compose

```bash
# Build and run the application
docker-compose up --build

# Test the endpoints
curl http://localhost:8080/                    # Hello World (200)
curl http://localhost:8080/status?code=500     # Error response (500)
curl http://localhost:8080/metrics             # Prometheus metrics

# Generate varied test traffic (100 iterations of different status codes)
for i in {1..100}; do
  for code in 200 400 404 500 502 503; do
    curl -s "http://localhost:8080/status?code=$code" > /dev/null
  done
  echo "Completed iteration $i/100 (6 requests per iteration)"
done

# Stop the application
docker-compose down
```

## Deploy app to k8s
```
# deploy app to k8s
helm install web -n test --create-namespace infra/go-demo-app

# test app is working
kubectl -n test port-forward svc/web-go-demo-app 8080:80
while true; do curl http://localhost:8080/; echo; sleep 0.5; done
kubectl -n test logs -l app.kubernetes.io/name=go-demo-app -f
```


## Manual helm push in GCP Artifact Registry

```
GCP_PROJECT="YOUR-GCP-PROJECT"
gcloud auth print-access-token | helm registry login -u oauth2accesstoken --password-stdin https://australia-southeast2-docker.pkg.dev
helm package infra/go-demo-app
helm push go-demo-app-0.1.0.tgz oci://australia-southeast2-docker.pkg.dev/${GCP_PROJECT}/cmek-helm-charts
helm template go-demo-app --namespace test --create-namespace oci://australia-southeast2-docker.pkg.dev/${GCP_PROJECT}/cmek-helm-charts/go-demo-app --version 0.1.0
```

## Setup github actions workflow env vars

```
GCP_PROJECT="YOUR-GCP-PROJECT"
SA_NAME="github-runner"
REGION="australia-southeast2"
WI_POOL_PROVIDER_ID=$(gcloud iam workload-identity-pools providers describe go-demo-app --workload-identity-pool=go-demo-app --location global --format='get(name)')
echo $WI_POOL_PROVIDER_ID

gh secret set PROJECT_ID -b"${GCP_PROJECT}"
gh secret set HELM_REPO_ID -b"cmek-helm-charts"
gh secret set IMAGE_REPO_ID -b"cmek-container-images"
gh secret set ARTIFACT_REGISTRY_HOST_NAME -b"${REGION}-docker.pkg.dev"
gh secret set PACKAGER_GSA_ID -b"${SA_NAME}@${GCP_PROJECT}.iam.gserviceaccount.com"
gh secret set WI_POOL_PROVIDER_ID -b"${WI_POOL_PROVIDER_ID}"
```