kubectl create sa gitops-operator
kubectl create clusterrolebinding gitops-operator-admin --clusterrole=cluster-admin  --serviceaccount=default:gitops-operator
kubectl apply -f deployment.yaml
kubectl expose deployment gitops --type=LoadBalancer --port=8080
