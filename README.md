# go-demo-app

This repo has github actions workflows for pushing container image and helm chart to Google Artifact Registry (GAR).
The terraform code for building the GAR is [here](https://github.com/andreistefanciprian/terraform-kubernetes-gke-cluster).

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
PROJECT_ID="YOUR_PROJECT_ID"
gcloud auth print-access-token | helm registry login -u oauth2accesstoken --password-stdin https://australia-southeast2-docker.pkg.dev
helm package infra/go-demo-app
helm push go-demo-app-0.1.0.tgz oci://australia-southeast2-docker.pkg.dev/${PROJECT_ID}/cmek-helm-charts
helm template go-demo-app --namespace test --create-namespace oci://australia-southeast2-docker.pkg.dev/${PROJECT_ID}/cmek-helm-charts/go-demo-app --version 0.1.0
```

## Setup github actions workflow env vars

```
PROJECT_ID="YOUR_PROJECT_ID"
SA_NAME="github-runner"
REGION="australia-southeast2"
WI_POOL_PROVIDER_ID=$(gcloud iam workload-identity-pools providers describe go-demo-app-prvdr --workload-identity-pool=go-demo-app --location global --format='get(name)')

gh secret set PROJECT_ID -b"${PROJECT_ID}"
gh secret set HELM_REPO_ID -b"cmek-helm-charts"
gh secret set IMAGE_REPO_ID -b"cmek-container-images"
gh secret set ARTIFACT_REGISTRY_HOST_NAME -b"${REGION}-docker.pkg.dev"
gh secret set PACKAGER_GSA_ID -b"${SA_NAME}@${PROJECT_ID}.iam.gserviceaccount.com"
gh secret set WI_POOL_PROVIDER_ID -b"${WI_POOL_PROVIDER_ID}"
```