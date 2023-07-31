FROM golang:1.19-alpine as builder

RUN apk add --no-cache libpcap-dev libc-dev gcc linux-headers
ADD . /app
WORKDIR /app
RUN go build -ldflags "-w -s -linkmode external -extldflags '-static'" -o /sx

FROM alpine:3.17

COPY --from=builder /sx /sx
ENTRYPOINT ["/sx"]
