FROM golang:1.23-alpine3.20 AS builder
WORKDIR /app
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o /go/bin/previewer ./cmd/main.go

FROM alpine:latest
WORKDIR /app

COPY --from=builder /go/bin/previewer /app/previewer
COPY --from=builder /app/configs/config-deploy.yaml .

CMD ["/app/previewer","--config","config-deploy.yaml","run"]