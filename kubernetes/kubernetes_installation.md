# Kubernetes Installation

## Kind

Kind allows us to use Kubernetes in a Docker container.
Meant for testing locally.

- Create cluster
  - `kind create cluster --name kafka-cluster`
- Verify cluster is running
  - `kubectl cluster-info --context kind-kafka-cluster`

## Dashboard via Proxy

Kind does not come with dashboard. To view dashboard, must follow steps below.

- Install
  - `kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.7.0/aio/deploy/recommended.yaml`
- Start
  - `kubectl proxy`
- Dashboard running on:
  - [dashboard_link](http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/https:kubernetes-dashboard:/proxy/)

## Access to Dashboard

- Apply admin user yaml configs
  - `kubectl apply -f dashboard-adminuser.yaml`
  - `kubectl apply -f dashboard-adminuser-binding.yaml`
- Obtain token
  - `kubectl -n kubernetes-dashboard create token admin-user`
