# Code Runner Backend

## Ensure you have Docker installed and run it in Docker

**Note:** This project must run inside Docker because it uses specific users and permissions.

### Clone the repo

```bash
git clone https://github.com/ishu17077/code_runner_backend
cd code_runner_backend
```

> API Calls

[Postman collection](https://m1racles.postman.co/workspace/SCCSE~18de115a-fe99-43e4-8681-bfccf985ac14/collection/40284511-abc6ba02-0e7d-48cf-9fb2-a3b113eb6fa1?action=share&creator=40284511)

## The Kubernetes Method

### To use with more production grade microk8s instead of minikube(New New Method)

> Install microk8s

```bash
sudo apt install snapd
sudo snap install microk8s --classic
```

**Note:** For Fedora use sudo dnf install snapd instead of first command

> Build the image and output it using docker
>> 1

```bash
docker build -f ./code-runner.Dockerfile . -t code-runner;
docker build -f ./warm-runner.Dockerfile . -t warm-runner;
```

>> 2

```bash
docker save -o ./code-runner.tar code-runner;
docker save -o ./warm-runner.tar warm-runner;
```

### Secret Config

**Note:** For new method, use this:

>> First create 1-code-runner-secret.yaml from 1-code-runner-secret.yaml.sample provided. (The values must be base64 encoded in the params)

```bash
cp ./1-code-runner-secret.yaml.sample ./1-code-runner-secret.yaml
```

Change the values of 1-code-runner-secret.yaml accordingly

> Import image in microk8s

```bash
microk8s images import ./code-runner.tar
microk8s images import ./warm-runner.tar
```

> Check if image is present

```bash
microk8s ctr images ls | grep code-runner
microk8s ctr images ls | grep code-runner
```

If present all good, else revert to building the image again and importing it.

**Note:** You must create these files from template mentioned [here](#secret-config)</a>

> Now applying kubernetes config and secrets(Optional)

```bash
microk8s kubectl apply -f ./1-code-runner-secret.yaml
microk8s kubectl apply -f ./2-code-runner.config.yaml
```

> Deploy Everything(One shot)

```bash
microk8s kubectl apply -f .
```

> Check the deployment

```bash
microk8s kubectl get pods
```

Choose any pod and describe

```bash
microk8s kubectl get code-runner-*
```

**Note:** * is any pod id you found

#### Enable dashboard(optional)

```bash
microk8s enable dashboard
microk8s dashboard-proxy
```

**Note:** You can use the token provided by dashboard-proxy in the web address that will be automatically opened

## To call the API

## Converting the programs to base64

You need to convert your entire code to base64 encoding before passing down to api

[Base64 encode/decode here](https://www.base64encode.org/)

Below is an example of base64 encoding

```base64
I2luY2x1ZGUgPHN0ZGlvLmg+CmludCBtYWluKCkgewogICAgaW50IHJlczsKICAgIHNjYW5mKCIlZCIsICZyZXMpOwogICAgaWYgKHJlcyAlIDIgPT0gMCkgewogICAgICAgIHByaW50ZigiWWVzXG5ZZXMiKTsKICAgIH0gZWxzZSB7CiAgICAgICAgcHJpbnRmKCJOb1xuTm8iKTsKICAgIH0KfQ==
```

Make sure you select encode option at the top for encoding

> [!NOTE]
> Please make sure you use class Solution for java programs, and don't forget to add 'package main' at the top of golang code.

> JSON Payload to call the api

**Note:** You can either use postman or curl

Request Type: POST

>> {host_url}/submission/test/private

>> {host_url}/submission/test/public

Replace {host_url} with actual url whether that be localhost or somewhere else.

**Note:** /submission/test/public endpoint cannot have more than 3 tests defined

```json
{
    "problem_id": "69",
    "language": "C",
    "code": "I2luY2x1ZGUgPHN0ZGlvLmg+CmludCBtYWluKCkgewogICAgaW50IHJlczsKICAgIHNjYW5mKCIlZCIsICZyZXMpOwogICAgaWYgKHJlcyAlIDIgPT0gMCkgewogICAgICAgIHByaW50ZigiWWVzXG5ZZXMiKTsKICAgIH0gZWxzZSB7CiAgICAgICAgcHJpbnRmKCJOb1xuTm8iKTsKICAgIH0KfQ==",
      "tests": [
    {
  "problem_id":     "69",
  "is_public":      true,
  "stdin":          "12\n",
  "expected_output": "Yes\nYes",
  "test_id":        "1"
 }
    ]
}
```

### Postman

> Install Postman from [here](https://www.postman.com/downloads/)

### Curl

```bash
curl -v -X POST "127.0.0.1:300080/submission/test/private" -H "Content-Type: application/json" -H "Connection: close" --data @./examples/stress_payload.json
```

**Note:** You can pipe the above command through jq to pretty print json, in simple terms add " | jq", at the end of the above code

> [!NOTE]
> This api will be accessible through 30080 port, as it is configured to use NodePort as load balancer
> For more paraller processing go to file 3-warm-runner-deployment.yaml and change 'replicas' in accordance with your computer specification.

### The minikube method(Old Method)

> Install Kubernetes (kubectl to create deployments preferrably through docker desktop)

Install kubernetes and kubectl & docker desktop and click on kubernetes and create kubernetes cluster

> Install Containerd Runtime
You can either use kind or minikube, minikube is preferrable because it has been tested.
>> Install Minikube
Install it on your system using [Minikube Tutorial](https://minikube.sigs.k8s.io/docs/start/).
>> Install KIND(Kunbernetes In Docker)
Install it on your system using [KIND Tutorial](https://kind.sigs.k8s.io/docs/user/quick-start/).

> [!WARNING]
> With KIND, you are on your own, it has not been tested

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

> Now apply the kubernetes conf

```bash
kubectl apply -f .
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

> [!TIP]
> Use the microk8s method to avoid incompabilities.

> [!NOTE]
> The postman collection sample is attached with this program [here](code-runner.postman-collection.json.sample).
