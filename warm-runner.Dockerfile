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

RUN apk update && apk add --no-cache curl build-base openjdk21 python3 py3-pip htop rustup go
#* dotnet9-sdk-aot 

RUN addgroup executorgrp --gid 7070 && adduser executor --uid 6969 executorgrp -D -S

#? Rust
ENV RUST_HOME="/opt/Rust"
ENV RUSTUP_HOME="$RUST_HOME/.rustup"
ENV CARGO_HOME="$RUST_HOME/.cargo"
ENV PATH="/opt/Rust/.cargo/bin:${PATH}"
RUN rustup-init -y && chown executor:executorgrp /opt/Rust -R

#? Dotnet
# RUN dotnet tool install -g dotnet-script
# ENV PATH="$PATH:/root/.dotnet/tools"
# RUN export DOTNET_NOLOGO=true

#? Sanity Checks
RUN rustc --version
RUN java -version
RUN gcc --version
RUN python --version
# RUN dotnet --list-sdks


RUN mkdir /temp
RUN chmod -R 755 /temp

RUN echo "root:1923934edfdfKLJHDKJkwfjkf" | chpasswd 

COPY --from=golang --chown=root:root /app/runner/runner .
COPY --from=golang --chown=root:root /app/java_output/JavaExecutor.jar /opt/JavaExecutor.jar

RUN chmod 700 ./runner
RUN chmod 755 /opt/JavaExecutor.jar

CMD ["sleep", "infinity"]