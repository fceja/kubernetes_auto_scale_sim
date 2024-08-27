# Kubernetes Pod

## Docker Container as Kubernetes Pod

To run a Docker container as a `Kubernetes Pod`, you need to create a Kubernetes `Pod` manifest that describes you container and then apply it you your Kubernetes cluster.

### 1. PrepareDocker Image

Ensure that your Docker image is accessible from your Kubernetes cluster. You can either use a public image from Docker Hub or push your image to a private Docker registry.

### 2. Create a Pod Manifest

You need to create a yaml file that defines you Pod. Here's a basic example of a Pod manifest.

- Create pod manifest

  - `my-pod.yaml`

    ```yaml
    apiVersion: v1
    kind: Pod
    metadata:
        name: my-pod
        labels:
        app: my-app
    spec:
        containers:
        - name: my-container
            image: your-docker-image:tag
            ports:
            - containerPort: 80
    ```

- `name`
  - name of Pod.
- `image`
  - Docker image you want to run
- `containerPort`
  - port that your container express

### 3. Apply the Pod Manifest

- Use `kubectl` to create the Pod based on your manifest file:
  - `kubectl apply -f my-pod.yaml`

### 4. Verify the Pod

- Check the status of the pod to ensure it's running:

  - `kubectl get pods`

- You should see your pod listed. If it's not, in a `Running` state, you can check the logs and events for debugging:
  - Check pod Logs:
    - `kubectl logs my-pod`
  - Describe the pod:
    - `kubectl describe pod my-pod`

### 5. Access Your Application

If your container exposes a web service or application, you might need to expose it outside of the cluster. Here's how to create a `Service` to expose the Pod:

- Create a Service Manifest

  - `my-service.yaml`

    ```yaml
    apiVersion: v1
    kind: Pod
    metadata:
    name: my-pod
    labels:
      app: my-app
    spec:
      containers:
      name: my-container
      image: your-docker-image:tag
      ports:
      containerPort: 80
    ```
