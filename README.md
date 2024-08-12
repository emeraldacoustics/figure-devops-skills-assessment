# Welcome 

Welcome to Figure's DevOps skills assessment! 

The goal of this assessment is to get an idea of how you work and your ability to speak in depth about the details in your work. Generally, this assessment should not take you longer than 30 minutes to complete. 

Your answers will be reviewed with you in a subsequent interview.

## Instructions

1. Click on the green "Use This Template" button in the upper-right corner and create a copy of this repository in your own GitHub account.
2. Name your respository and ensure that it's public, as you will need to share it with us for review.
3. When you have completed the questions, please send the URL to the recruiter.

## Assessments

### Kubernetes

1. Fix the issues with this Kubernetes manifest to ensure it is ready for deployment. 
2. Add the following limits and requests to the manifest:
- CPU limit of 0.5 CPU cores
- Memory limit of 256 Mebibytes
- CPU request of 0.2 CPU cores
- Memory request of 128 Mebibytes 

```yaml
apiVersion: apps/v1
kind: Deploy
metadata:
  name: nginx-deploy
  labels:
    app: nginx
spec:
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:current
        ports:
        - containerPort: 80
---
apiVersion: v1
kind: Service
metadata:
  name: nginx-service
spec:
  selector:
    app: nginx
  ports:
    - protocol: TCP
      port: 80
  type: ClusterIP
  ```

### Go

Write a script in Go that redeploys all pods in a Kubernetes cluster that have the word `database` in the name.

Requirements:
- Assume local credentials in your kube config have full access. There is no need to connect via a service account, etc.
- You must use the [client-go](https://github.com/kubernetes/client-go) library.
- Your script must perform a graceful restart, similar to kubectl rollout restart. Do not just delete pods.
- You must use Go modules (no vendor directory).
