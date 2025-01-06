# k8s gitops operator
GitOps operator for Kubernetes implemented in Go.

Uses Istio for traffic management, GitHub as the source of truth, and interacts with k8s API to create necessary resources.
GCP is used as a cloud provider, but it's mostly for a reference and can be easily replaced with any other cloud provider.

For this to work, the repositories managed by this operator need to have a webhook configured to send events to the operator's endpoint.
Their root directory should contain an `ops.yaml` file with the following structure:

```yaml
name: test-service
image:
  repository: nginx
  tag: latest
extraEnv:
  - name: FOO
    value: bar
routes:
  - name: grpc-web
    scope: external
    port: 80
permissions:
  gcp:
    - name: cloudsql
      role: roles/cloudsql.client    
```
