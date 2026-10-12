# Code Runner Backend

Backend service for compiling and executing submitted code inside isolated containers.

## Demo

- Frontend project demo: https://github.com/Subhojit534/Code_Arena/
- Demo video: https://github.com/user-attachments/assets/e212aa2a-4bcb-4693-8fb3-34c3595380f4

## Prerequisites

- Docker (required)
- Kubernetes runtime (MicroK8s recommended; Minikube/KIND supported as legacy flow)

> [!IMPORTANT]
> This project is designed to run in Docker/Kubernetes because it depends on specific users and permissions.

## Clone the repository

```bash
git clone https://github.com/ishu17077/code_runner_backend
cd code_runner_backend
```

## API collection

- Postman collection: https://m1racles.postman.co/workspace/SCCSE~18de115a-fe99-43e4-8681-bfccf985ac14/collection/40284511-abc6ba02-0e7d-48cf-9fb2-a3b113eb6fa1?action=share&creator=40284511
- Sample collection file: [code-runner.postman-collection.json.sample](./code-runner.postman-collection.json.sample)

---

## Kubernetes Setup (Recommended: MicroK8s)

### 1) Install MicroK8s

```bash
sudo apt install snapd
sudo snap install microk8s --classic
```

For Fedora, install `snapd` with:

```bash
sudo dnf install snapd
```

### 2) Build Docker images

```bash
docker build -f ./code-runner.Dockerfile . -t code-runner
docker build -f ./warm-runner.Dockerfile . -t warm-runner
```

### 3) Export images

```bash
docker save -o ./code-runner.tar code-runner
docker save -o ./warm-runner.tar warm-runner
```

### 4) Configure secrets

Create a secret file from the sample:

```bash
cp ./1-code-runner-secret.yaml.sample ./1-code-runner-secret.yaml
```

Then update values in `1-code-runner-secret.yaml` (values must be base64 encoded).

### 5) Import images into MicroK8s

```bash
microk8s images import ./code-runner.tar
microk8s images import ./warm-runner.tar
```

Verify import:

```bash
microk8s ctr images ls | grep code-runner
microk8s ctr images ls | grep warm-runner
```

### 6) Apply configuration

Optional: apply only secret/config first:

```bash
microk8s kubectl apply -f ./1-code-runner-secret.yaml
microk8s kubectl apply -f ./2-code-runner-config.yaml
```

Apply everything:

```bash
microk8s kubectl apply -f .
```

Check deployment:

```bash
microk8s kubectl get pods
```

### 7) (Optional) Dashboard

```bash
microk8s enable dashboard
microk8s dashboard-proxy
```

Use the token printed by `dashboard-proxy` to sign in.

---

## Calling the API

### Encode source code to base64

The `code` field must contain base64-encoded source code.

Online tool: https://www.base64encode.org/

Example encoded payload:

```base64
I2luY2x1ZGUgPHN0ZGlvLmg+CmludCBtYWluKCkgewogICAgaW50IHJlczsKICAgIHNjYW5mKCIlZCIsICZyZXMpOwogICAgaWYgKHJlcyAlIDIgPT0gMCkgewogICAgICAgIHByaW50ZigiWWVzXG5ZZXMiKTsKICAgIH0gZWxzZSB7CiAgICAgICAgcHJpbnRmKCJOb1xuTm8iKTsKICAgIH0KfQ==
```

> [!NOTE]
> - For Java, use `class Solution`.
> - For Go, include `package main` at the top.

### Endpoints

- `POST {host_url}/submission/test/private`
- `POST {host_url}/submission/test/public`

Replace `{host_url}` with your service address.

> [!NOTE]
> `/submission/test/public` supports up to 3 tests in one request.

### Request body example

```json
{
  "problem_id": "69",
  "language": "C",
  "code": "I2luY2x1ZGUgPHN0ZGlvLmg+CmludCBtYWluKCkgewogICAgaW50IHJlczsKICAgIHNjYW5mKCIlZCIsICZyZXMpOwogICAgaWYgKHJlcyAlIDIgPT0gMCkgewogICAgICAgIHByaW50ZigiWWVzXG5ZZXMiKTsKICAgIH0gZWxzZSB7CiAgICAgICAgcHJpbnRmKCJOb1xuTm8iKTsKICAgIH0KfQ==",
  "tests": [
    {
      "problem_id": "69",
      "is_public": true,
      "stdin": "12\\n",
      "expected_output": "Yes\\nYes",
      "test_id": "1"
    }
  ]
}
```

### Curl example

```bash
curl -v -X POST "127.0.0.1:30080/submission/test/private" \
  -H "Content-Type: application/json" \
  -H "Connection: close" \
  --data @./examples/stress_payload.json
```

You can append `| jq` to pretty-print the response.

> [!NOTE]
> Service is exposed via NodePort `30080` in the provided Kubernetes manifests.

### Scaling warm runners

For higher parallel processing, increase `replicas` in:

- `3-warm-runner-deployment.yaml`

---

## Legacy Setup (Minikube / KIND)

Use this only if you are not using the MicroK8s flow.

### Install and start cluster

- Minikube guide: https://minikube.sigs.k8s.io/docs/start/
- KIND guide: https://kind.sigs.k8s.io/docs/user/quick-start/

Create cluster with KIND:

```bash
kind create cluster
```

Or start Minikube:

```bash
minikube start
```

### Build and load image

```bash
docker build -t code_runner . -o code_runner.tar
minikube load ./code_runner.tar
```

### Deploy and access service

```bash
kubectl apply -f .
minikube service code-runner
```

### Minikube dashboard

```bash
minikube dashboard
```

> [!TIP]
> Prefer the MicroK8s method to avoid compatibility issues.
