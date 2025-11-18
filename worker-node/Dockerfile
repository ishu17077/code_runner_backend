FROM golang:1.25.4-alpine3.22 AS golang
WORKDIR /app

COPY ./go.mod .
COPY ./go.sum .

RUN go mod download 

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/workerNode .

FROM alpine:latest

WORKDIR /root/

RUN apk add curl build-base openjdk17 python3 py3-pip

RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
ENV PATH="/root/.cargo/bin:${PATH}"

RUN rustc --version
RUN java -version
RUN gcc --version
RUN mkdir temp

COPY --from=golang /app/workerNode .

EXPOSE 8060

CMD [ "./workerNode" ]

