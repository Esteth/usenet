FROM golang:1.12 AS builder

ENV GO111MODULE on

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o usenet ./cmd

FROM scratch
COPY --from=builder /app/usenet /app/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT [ "/app/usenet" ]