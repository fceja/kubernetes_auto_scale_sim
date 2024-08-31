# Kubernetes Installation

## Kind

Kind allows us to use Kubernetes in a Docker container.
Meant for testing locally.

- Create cluster
  - `kind create cluster --name kafka-cluster`
  - Then, set kubectl context to newly created kind cluster
    - `kind create cluster --name kafka-cluster`
    - `kubectl config use-context kind-kafka-cluster`
    - Note - when you have multiple Kubernetes clusters (e.g., development, staging, production), setting correct context ensure you're executing commands again intended cluster
- Verify cluster is running
  - `kind create cluster --name kafka-cluster`
  - `kind get clusters`
  - `kubectl config get-contexts`
  - `kubectl config use-context kind-kafka-cluster`
  - `kubectl config current-context`
  - `kubectl cluster-info --context kind-kafka-cluster`
  - `kind delete cluster --name CLUSTER_NAME`
  - `kind load docker-image kafka-3.4.1-producer --name kafka-cluster`

## Dashboard via Proxy

Kind does not come with dashboard. To view dashboard, must follow steps below.

- Install
  - `kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.7.0/aio/deploy/recommended.yaml`

### Pre-install - Access to Dashboard

- Apply admin user yaml configs
  - `kubectl apply -f dashboard-adminuser.yaml`
  - `kubectl apply -f dashboard-adminuser-binding.yaml`
- Obtain token
  - `kubectl -n kubernetes-dashboard create token admin-user`

### Access Dashboard

- Start
  - `kubectl proxy`
- Dashboard running on:
  - `http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/`
  - apply copied token

## Creating Pods

`docker-compose.yml` file isn't directly compatabile with Kubernetes, but it can be converted.
Kubernetes uses its own format, Kubernetes manifest.

- Convert by using `kompose` tool

  - you can convert a `docker-compose.yaml` to Kubernetes manifest files using `kompose`
  - `brew install kompose`
  - Navigate to `docker-compose.yaml`
    - `kompose convert`
  - Apply to Kubernetes cluster
    - `kubectl apply -f .`

- Or, you can manually convert `docker-compose.yaml`

  - Naming convention
    - `pod.yaml`
      - `pod-zookeeper.yaml`
    - `deployment.yaml`
      - `deployment-zookeeper.yaml`
    - `service.yaml`
      - `service-zookeeper.yaml`

- Once `pod.yaml` is crated

  - `kubectl apply -f pod.yaml`

- Verify creation

  - `kubectl get pods`
  - `kubectl get svc`

- View logs - In Kubernetes Dashboard, or:

  - `kubectl logs <pod-name>`
  - pod has multiple containers:
    - `kubectl logs <pod-name> -c <container-name>`

- Connecting to Pod - In Kubernetes Dashboard, or:
  - `kubectl logs <pod-name> -c <container-name>`
  - open shell
    - `kubectl logs zookeeper -c /bin/sh`
  - if container has bash
    - `kubectl exec -it zookeeper -- /bin/bash`
  - run a single command
    - `kubectl exec <pod-name> -- ls /app`

### Load images into Kind

Images are pulled in variety of ways, docker hub, kind, docker daemon, etc.

- To load into kind
  - kind load docker-image.
    - kind load docker-image.

## Full Steps to create Kubernetes cluster with pods

- 1. Delete and Verify

  - `kind get clusters`
  - `kind delete cluster --name CLUSTER_NAME`
  - `kind get clusters`

- 2. Create and verify

  - `kind create cluster --name kafka-cluster`
    - Note - this will create `kind-kafka-cluster`
  - `kind get clusters`
    - should return 'kafka-cluster'

- 3. Verify current context

  - `kubectl config current-context`
  - incase - `kubectl config get-contexts`

- 4. Install dashboard

  - install
    - `kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.7.0/aio/deploy/recommended.yaml`
  - apply
    - `kubectl apply -f dashboard-adminuser.yaml`
    - `kubectl apply -f dashboard-adminuser-binding.yaml`
  - obtain token to copy
    - `kubectl -n kubernetes-dashboard create token admin-user`

- 5. Start proxy for Kubernetes dashboard

  - `kubectl proxy`
    - Dashboard running on:
      - `http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/`
      - apply copied token

- 6. (In not already installed) - Install run `kompose`

  - run `kompose convert`

- 7. Navigate to `/kompose`

  - `kubectl apply -f .`
    - or one by one
      - `kubectl apply -f zookeeper-deployment.yaml -f zookeeper-service.yaml`
      - standard to apply deployment first, however does not matter since handled by kubernetes

- 7. Navigate to `/kompose`

- **DEBUG**

- `kind get clusters`
- `kind create cluster --name CLUSTER_NAME`
- `kind delete cluster --name CLUSTER_NAME`

- `kubectl get pods`
- `kubectl delete pod kafdrop-server-ffc788cb5-xp8mr`
- `kubectl delete pods --all`
