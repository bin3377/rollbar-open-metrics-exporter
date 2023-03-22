FROM golang:alpine as builder

WORKDIR /app

RUN apk update && apk upgrade && apk add --no-cache ca-certificates
RUN update-ca-certificates

COPY . .

RUN go get -d -v
RUN go build

FROM scratch

COPY --from=builder /app/rollbar-open-metrics-exporter .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["./rollbar-open-metrics-exporter"]
