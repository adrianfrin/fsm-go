FROM golang:1.26-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /out/fsm-demo ./cmd/fsm-demo

FROM alpine:3.22

WORKDIR /app
COPY --from=build /out/fsm-demo /app/fsm-demo
COPY configs /app/configs
EXPOSE 8080
ENTRYPOINT ["/app/fsm-demo"]
