FROM golang:latest AS builder
ENV HOME=/
ENV CGO_ENABLED=0
ENV GOOS=linux
WORKDIR /
COPY . .
RUN go get -d && go mod download && go build -a -ldflags "-s -w" -installsuffix cgo -o tlscdn-controller .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=builder /tlscdn-controller .
ENTRYPOINT ["/tlscdn-controller"]
