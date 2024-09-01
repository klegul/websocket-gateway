FROM golang:1.23-alpine as build
WORKDIR /opt
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./websocket-gateway

FROM alpine:latest as run
WORKDIR /opt
COPY --from=build /opt/websocket-gateway .
CMD ["/opt/websocket-gateway"]