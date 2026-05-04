FROM golang:1.25.5-alpine3.23 AS golang
WORKDIR /app

COPY ./go.mod .
COPY ./go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./runner ./runner

FROM alpine:3.23

USER root

WORKDIR /root/

RUN apk update && apk add --no-cache curl build-base openjdk21 python3 py3-pip htop rustup go perl dotnet10-sdk-aot

RUN addgroup executorgrp --gid 7070 && adduser executor --uid 6969 executorgrp -D -S

#? Rust
ENV RUST_HOME="/opt/Rust"
ENV RUSTUP_HOME="$RUST_HOME/.rustup"
ENV CARGO_HOME="$RUST_HOME/.cargo"
ENV PATH="/opt/Rust/.cargo/bin:${PATH}"
RUN rustup-init -y && chown executor:executorgrp /opt/Rust -R

#? Dotnet Warmup
# RUN mkdir -p /opt/dotnet/tools
# RUN chmod 755 /opt/dotnet/ -R
# ENV PATH="$PATH:/opt/dotnet/tools"
ENV DOTNET_NOLOGO=true
RUN mkdir -p /opt/dotnet-project
WORKDIR /opt/dotnet-project
RUN dotnet new console && \
    dotnet restore
RUN mkdir -p /home/executor/.nuget && chown -R executor:executorgrp /home/executor
RUN echo 'using System;public class Program{public static void Main(string[] args){Console.Write("DOTNET Warmup successful");}}' > /opt/dotnet-project/Program.cs && \
    cd /opt/dotnet-project && \
    dotnet build -c Release && \
    chown executor:executorgrp /opt/dotnet-project -R
RUN  ./bin/Release/net10.0/dotnet-project

WORKDIR /


#? Sanity Checks
RUN rustc --version
RUN java -version
RUN gcc --version
RUN perl --version
RUN python --version
RUN dotnet --list-sdks

#? Go Warmup
ENV GOOS=linux
ENV GOCACHE=/opt/go-cache
RUN mkdir -p /opt/go-cache
RUN echo 'package main; import "fmt"; func main(){ fmt.Println("Go Cache Warmput Complete....")}' > /tmp/warmup.go && \
    go build -o /tmp/warmup /tmp/warmup.go && \
    rm -rf /tmp/warmup /tmp/warmup.go  && \
    chown -R executor:executorgrp /opt/go-cache
USER root




RUN mkdir /temp
RUN chmod -R 755 /temp

RUN echo "root:1923934edfdfKLJHDKJkwfjkf" | chpasswd 

COPY --from=golang --chown=root:root /app/runner/runner .
COPY --from=golang --chown=root:root /app/java_output/JavaExecutor.jar /opt/JavaExecutor.jar

RUN chmod 700 ./runner
RUN chmod 755 /opt/JavaExecutor.jar



CMD ["sleep", "infinity"]