FROM golang:1.25.5-alpine3.23 AS golang
WORKDIR /app

COPY ./go.mod .
COPY ./go.sum .

RUN go mod download 

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/worker_node .

FROM alpine:3.23

WORKDIR /root/
COPY --from=golang /app/worker_node .

EXPOSE 8060

RUN chmod +x ./worker_node

# ENTRYPOINT ["echo", "+cpu +memory +pids", ">", "/sys/fs/cgroup/cgroup.subtree_control"]
CMD ["./worker_node"]


