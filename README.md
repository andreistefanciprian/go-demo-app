# go-demo-app


## Deploy app to k8s
```
# deploy app to k8s
helm install web -n test --create-namespace infra/go-demo-app

# test app is working
kubectl -n test port-forward svc/web-go-demo-app 8080:80
while true; do curl http://localhost:8080/; echo; sleep 0.5; done
kubectl -n test logs -l app.kubernetes.io/name=go-demo-app -f
```