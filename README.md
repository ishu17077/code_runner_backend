# Code Runner Backend

## Ensure you have Docker installed and run it in Docker

**Note:** This project must run inside Docker because it uses specific users and permissions.

> Clone the repo

```bash
git clone htps://github.com/ishu17077/code_runner_backend
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
docker rmi code_runner
```

> API Calls

[Postman collection](https://m1racles.postman.co/workspace/SCCSE~18de115a-fe99-43e4-8681-bfccf985ac14/collection/40284511-abc6ba02-0e7d-48cf-9fb2-a3b113eb6fa1?action=share&creator=40284511)
