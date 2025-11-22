```
cd worker-node
docker build -t code_runner .

docker run -p 8060:8060 -d code_runner
```
