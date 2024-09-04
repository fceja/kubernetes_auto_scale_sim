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

`docker-compose.yml` file isn't directly compatible with Kubernetes, but it can be converted.
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

- 7. A - Local. Create `deployment` and `service` files for Kubernetes

  - Make sure local image is `tagged`

    - Build

      - `kind get clusters`
      - `docker build -t zookeeper-3.9.2:latest .`

      ```bash
      docker build -t zookeeper-3.9.2:latest ./zookeeper \
      && docker build -t kafka-3.4.1-server:latest ./kafka_server \
      && docker build -t kafka-3.4.1-producer:latest ./kafka_producer \
      && docker build -t kafka-3.4.1-consumer:latest ./kafka_consumer \
      && docker build -t kafdrop-4.0.3-snapshot:latest ./kafdrop_server
      ```

      - `kind load docker-image zookeeper-3.9.2:latest --name kafka-cluster`

    - Load into kind

      - `kind load docker-image zookeeper-3.9.2 --name kafka-cluster`

      ```bash
        kind load docker-image zookeeper-3.9.2 --name kafka-cluster \
        && kind load docker-image kafka-3.4.1-server --name kafka-cluster \
        && kind load docker-image kafka-3.4.1-producer --name kafka-cluster \
        && kind load docker-image kafka-3.4.1-consumer --name kafka-cluster \
        && kind load docker-image kafdrop-4.0.3-snapshot --name kafka-cluster
      ```

      ```bash
        kind load docker-image zookeeper-3.9.2 --name kafka-cluster \
        && kind load docker-image kafka-3.4.1-server.kubernetes --name kafka-cluster \
        && kind load docker-image kafka-3.4.1-producer --name kafka-cluster \
        && kind load docker-image kafka-3.4.1-consumer --name kafka-cluster \
        && kind load docker-image kafdrop-4.0.3-snapshot --name kafka-cluster
      ```

  - Verify list all available images to the kind-cluster

    - `docker exec -it kafka-cluster-control-plane crictl images`

    - to remove

      - `docker exec -it kafka-cluster-control-plane crictl rmi docker.io/ciscoceja/zookeeper-3.9.2`

    - copy full image name and paste into `docker-compose.local.yaml`

  - (In not already installed) - Install run `kompose`

    - run `kompose convert`
      - `kompose convert -f docker-compose.local-image.yaml`

  - place generated file into `/kompose/local`

  - manually add

    - `imagePullPolicy: IfNotPresent`

  - apply to kubernetes, navigate to`/kompose/local`
    - `kubectl apply -f .`
    - OR
    - `kubectl apply -f zookeeper-deployment.yaml -f zookeeper-service.yaml`
      - standard to apply deployment first, however does not matter since handled by kubernetes

- 7. B - DockerHub - Create `deployment` and `service` files for Kubernetes

  - (In not already installed) - Install run `kompose`

    - run `kompose convert`
      - `kompose convert -f docker-compose.yaml`
      - `kompose convert -f docker-compose.kubernetes.yaml`

  - place generated file into `/kompose/dockerhub` and navigate to it

  - manually add

    - `imagePullPolicy: IfNotPresent`

  - `kubectl apply -f .`
  - OR
  - `kubectl apply -f zookeeper-deployment.yaml -f zookeeper-service.yaml`
    - standard to apply deployment first, however does not matter since handled by kubernetes

- 8. Adding secrets

  - `kubectl create secret generic my-secret --from-env-file=.env`

  ```bash
    kubectl create secret generic kafka-server-env --from-env-file=./kafka_server/.env \
    && kubectl create secret generic kafka-producer-env --from-env-file=./kafka_producer/.env \
    && kubectl create secret generic kafka-consumer-env --from-env-file=./kafka_consumer/.env \
    && kubectl create secret generic kafdrop-server-env --from-env-file=./kafdrop_server/.env
  ```

  - `kubectl delete secret my-secret`
  - `kubectl get secrets`
  - `kubectl get secrets --namespace=my-namespace`
  - `kubectl describe secret my-secret`

- 9. View Kafdrop UI

  - kubectl get svc
    - get name of kafdrop service , `kafdrop-server`
  - forward port
    - `kubectl port-forward svc/kafdrop-server 9000:9000`

- **DEBUG**

- Note: In Kubernetes, when you delete a pod, if it is a part of a deployment, replica set, or stateful set, a new pod will automatically be created to maintain desired state

  - To remove pod entirely and prevent from being recreated, you need to delete the higher-level resource that manages the pod.
    - Delete Deployment, StatefulSet or ReplicaSet
      - `kubectl delete deployment <deployment-name>`
        - `kubectl delete deployment --all`
      - `kubectl delete statefulset <statefulset-name>`
      - `kubectl delete replicaset <replicaset-name>`
    - If Pod managed by Job or CronJob
      - `kubectl delete job <job-name>`
      - `kubectl delete cronjob <cronjob-name>`

- `docker exec -it kafka-cluster-control-plane crictl images`

- `kind get clusters`
- `kind create cluster --name CLUSTER_NAME`
  - `kind create cluster --name kafka-cluster`
- `kind delete cluster --name CLUSTER_NAME`

- `kubectl get pods`
- `kubectl delete pod kafdrop-server-ffc788cb5-xp8mr`
- `kubectl delete pods --all`
