# Kubernetes Installation

## Installation

To simplify installation, homebrew package manager is used. You can install dependencies separately or with your preferred package manager.

\*\*NOTE - You can view helpful de-bugging commands at the bottom of this page.

- Install `kind`

  - Tool that allows you to use Kubernetes clusters inside Docker containers.
  - [Kind-website](https://kind.sigs.k8s.io/)

  - ```bash
    brew install kind
    ```

- Install `kompose`

  - Tool that helps convert Docker Compose files into Kubernetes manifests.
  - [Kompose-website](https://kompose.io/)

  - ```bash
    brew install kompose
    ```

### Kind setup

- 1. Create kind cluster

  - ```bash
    kind create cluster --name kafka-cluster
    ```

- 2. Set `kubectl` context to kind cluster

  - ```bash
      kubectl config use-context kind-kafka-cluster
    ```

- 3. Load docker images into kind cluster

  - Note - you can skip loading image there were docker image was not built.

  - ```bash
      kind load docker-image zookeeper-3.9.2 --name kafka-cluster \
      && kind load docker-image kafka-3.4.1-server.kubernetes:latest --name kafka-cluster \
      && kind load docker-image kafka-3.4.1-producer:latest --name kafka-cluster \
      && kind load docker-image kafka-3.4.1-consumer:latest --name kafka-cluster \
      && kind load docker-image kafdrop-4.0.3-snapshot:latest --name kafka-cluster \
      && kind load docker-image prometheus-2.54.1:latest --name kafka-cluster \
      && kind load docker-image grafana-11.2.0:latest --name kafka-cluster
    ```

### Kompose setup

- 1. Convert `docker-compose` into Kubernetes manifests

  - ```bash
      kompose convert -f docker-compose.kubernetes.yaml
    ```

- 2. Place generated file into `~/project_root/kompose/local`

- 3. Navigate to `~/project_root/kompose/local`

- 4. For each `*-deployment.yaml` file that was created from `kompose`, you must append the following in the .yaml:

  - spec > template > spec > containers > env > `imagePullPolicy: IfNotPresent`
  - This tells Kubernetes to pull from local docker image and not from docker hub.

    - Example:

      - ```yaml
            ...
            spec:
            containers:
                - env:
                    - name: KAFKA_BROKERCONNECT
                    valueFrom:
                        configMapKeyRef:
                        key: KAFKA_BROKERCONNECT
                        name: kafdrop-server-env
                image: kafdrop-4.0.3-snapshot:latest
                imagePullPolicy: IfNotPresent # HERE
                name: kafdrop-server
            ...
        ```

### Kubernetes setup

- 1. Apply manifest to Kubernetes

  - `kubectl apply -f .`

- 2. Adding secrets to Kubernetes

  - navigate to `~/project_root`

  - ```bash
      kubectl create secret generic kafka-server-env --from-env-file=./kafka/server/.env \
      && kubectl create secret generic kafka-producer-env --from-env-file=./kafka/producer/.env \
      && kubectl create secret generic kafka-consumer-env --from-env-file=./kafka/consumer/.env \
      && kubectl create secret generic kafdrop-server-env --from-env-file=./kafdrop_server/.env
    ```

- 3. Install Kubernetes dashboard

  - Install dependency
    - `kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.7.0/aio/deploy/recommended.yaml`
  - Apply Kubernetes user settings
    - Navigate to `~/project_root/kompose/local` and run
      - `kubectl apply -f dashboard-adminuser.yaml`
      - `kubectl apply -f dashboard-adminuser-binding.yaml`
  - Obtain token to copy
    - `kubectl -n kubernetes-dashboard create token admin-user`
    - Save this, will need for next step.

- 5. Start proxy for Kubernetes dashboard

  - ```bash
    kubectl proxy
    ```

  - Dashboard running on:
    - `http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/`
    - Apply copied token for password

### **DE-BUGGING** - Helpful commands

- ```bash
    kind get clusters
  ```

- ```bash
    kind create cluster --name <cluster_name>
  ```

- ```bash
    kind delete cluster --name <cluster_name>
  ```

- ```bash
    kubectl config get-contexts
  ```

- ```bash
    kubectl config current-context
  ```

- ```bash
    kubectl config use-context kind-<cluster_name>
  ```

- ```bash
    docker exec -it kafka-cluster-control-plane crictl images
  ```

- ```bash
    docker exec -it kafka-cluster-control-plane crictl rmi <image_name>
  ```
