# Code Runner Backend

## Ensure you have Docker installed and run it in Docker

**Note:** This project must run inside Docker because it uses specific users and permissions.

> Clone the repo

```bash
git clone https://github.com/ishu17077/code_runner_backend
cd code_runner_backend
```

> Usage to run

```bash
cd worker-node
docker compose up -d
```

> To stop the container

```bash
docker stop code_runner
```

To run it again use the above two command inside this project

> To remove the container

```bash
docker container stop code_runner
docker container rm code_runner
docker rmi code_runner:latest
```

> API Calls

[Postman collection](https://m1racles.postman.co/workspace/SCCSE~18de115a-fe99-43e4-8681-bfccf985ac14/collection/40284511-abc6ba02-0e7d-48cf-9fb2-a3b113eb6fa1?action=share&creator=40284511)

## The Kubernetes Method

> ./worker-node

```bash
cd worker-node
```

> Install Kubernetes (kubectl to create deployments preferrably through docker desktop)

Install kubernetes and kubectl & docker desktop and click on kubernetes and create kubernetes cluster

> Install Containerd Runtime
You can either use kind or minikube, minikube is preferrable because it has been tested.
>> Install Minikube
Install it on your system using [Minikube Tutorial](https://minikube.sigs.k8s.io/docs/start/).
>> Install KIND(Kunbernetes In Docker)
Install it on your system using [KIND Tutorial](https://kind.sigs.k8s.io/docs/user/quick-start/).

```bash
kind create cluster
```

> Using Minikube

```bash
minikube start
```

> Build Docker image

```bash
docker build -t code_runner . -o code_runner.tar
```

> Load image to Minikube

```bash
minikube load ./code_runner.tar
```

<!-- > Old Method of creating env
>> Create Environment Variable Map in kubernetes(Old Method)

Create .env file in /worker-node directory

```bash
cp ./.env.sample ./.env 
```

Now provide values into .env

>> Create a configmap in kubernetes from .env file (Old Method)

```bash
kubectl create configmap env --from-file ./.env
``` -->

**Note:** For new method, use this:

>> First create code-runner-secret.yaml from code-runner-secret.sample.yaml provided. (The values must be base64 encoded in the params)

```bash
kubectl apply -f ./code-runner-config.yaml && kubectl apply -f ./code-runner-secret.yaml
```

> Now apply the kubernetes conf

```bash
kubectl apply -f ./code-runner-deployment.yaml
kubectl apply -f ./code-runner-service.yaml
```

> Now we need to get a port exposed from minikube

```bash
minikube service code-runner
```

**Note:** The kubernetes is set up with load balancer with different pod, so it will be accessible via 30080 port instead of 8060 with docker-desktop kubernetes, but with minikube it would be running on port provided by above command

> Access minikube dashboard

```bash
minikube dashboard
```
