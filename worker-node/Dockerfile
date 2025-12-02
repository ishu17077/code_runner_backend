FROM golang:1.25.4-alpine3.22 AS golang
WORKDIR /app

COPY ./go.mod .
COPY ./go.sum .

RUN go mod download 

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/worker_node .

FROM alpine:3.22.2

WORKDIR /root/

RUN apk add curl build-base openjdk17 python3 py3-pip cgroup-tools htop iptables dotnet9-sdk-aot rustup

#? Rust
ENV RUST_HOME="/opt/Rust"
ENV RUSTUP_HOME="$RUST_HOME/.rustup"
ENV CARGO_HOME="$RUST_HOME/.cargo"
ENV PATH="/opt/Rust/.cargo/bin:${PATH}"
RUN rustup-init -y

#? Dotnet
RUN dotnet tool install -g dotnet-script
ENV PATH="$PATH:/root/.dotnet/tools"
RUN export DOTNET_NOLOGO=true

#? Sanity Checks
RUN rustc --version
RUN java -version
RUN gcc --version
RUN python --version
RUN dotnet --list-sdks

#? Pre-Code environment set-up 
RUN mkdir /temp
USER root
RUN addgroup executorgrp --gid 7070 && adduser executor --uid 6969 executorgrp -D -S

#! The root will be creating new subdirectories
# RUN chown -R executor:root /temp
RUN chmod -R 775 /temp
# RUN iptables -A OUTPUT -m owner --uid-owner 6969 -j DROP

# RUN mkdir /cgroup
# RUN chmod 770 /cgroup

RUN echo "root:1923934edfdfKLJHDKJkwfjkf" | chpasswd 

COPY --from=golang /app/entrypoint.sh .
COPY --from=golang /app/worker_node .

RUN chmod +x ./entrypoint.sh

EXPOSE 8060

RUN chmod +x ./worker_node

# ENTRYPOINT ["echo", "+cpu +memory +pids", ">", "/sys/fs/cgroup/cgroup.subtree_control"]
ENTRYPOINT ["./entrypoint.sh"]
CMD ["./worker_node"]


