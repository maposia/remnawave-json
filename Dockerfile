FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download


COPY ./cmd/app .

COPY ./internal ./internal
COPY ./remnawawe ./remnawawe
COPY templates ./template

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/app .

FROM alpine:latest

RUN apk update && apk add --no-cache ca-certificates
RUN update-ca-certificates

RUN mkdir -p /app/templates

COPY --from=builder /bin/app /app/app

COPY --from=builder /app/template /app/templates

ENV APP_PORT=4000

EXPOSE ${APP_PORT}

CMD ["/app/app"]

