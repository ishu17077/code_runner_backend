FROM golang:1.25.4-alpine3.22 AS golang
WORKDIR /app

COPY ./go.mod .
COPY ./go.sum .

RUN go mod download 

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/worker_node .

FROM alpine:3.22.2

WORKDIR /root/

RUN apk add curl build-base openjdk17 python3 py3-pip cgroup-tools htop iptables

RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
ENV PATH="/root/.cargo/bin:${PATH}"

RUN rustc --version
RUN java -version
RUN gcc --version

RUN mkdir /temp
USER root
RUN addgroup executorgrp --gid 7070 && adduser executor --uid 6969 executorgrp -D -S

RUN chown executor:root /temp
RUN chmod 775 /temp
# RUN iptables -A OUTPUT -m owner --uid-owner 6969 -j DROP

# RUN mkdir /cgroup
# RUN chmod 770 /cgroup

RUN echo "root:1923934edfdfKLJHDKJkwfjkf" | chpasswd 

COPY --from=golang /app/worker_node .

EXPOSE 8060



# ENTRYPOINT ["echo", "+cpu +memory +pids", ">", "/sys/fs/cgroup/cgroup.subtree_control"]

CMD [ "./worker_node" ]

