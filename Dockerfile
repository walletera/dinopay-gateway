FROM golang:1.22.8-alpine3.20 as builder

WORKDIR /root

COPY . ./

RUN go build -o bin/dinopay-gateway cmd/main.go

FROM alpine:3.20

# Allow customization of user ID and group ID (it's useful when you use Docker bind mounts)
ARG UID=1000
ARG GID=1000

RUN addgroup -g ${GID} -S app && adduser -u ${UID} -S -G app app

WORKDIR /home/app

COPY --from=builder /root/bin/dinopay-gateway ./


RUN chown app:app ./dinopay-gateway
RUN chmod +x ./dinopay-gateway

USER app

CMD ["./dinopay-gateway"]
