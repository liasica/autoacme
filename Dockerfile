FROM golang:1-alpine AS builder

RUN apk --no-cache --no-progress add make git

WORKDIR /go/autoacme
COPY . .
RUN go mod download && make build

FROM alpine:3
RUN apk update \
    && apk add --no-cache ca-certificates tzdata \
    && ln -s /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && update-ca-certificates \
    && rm -rf /var/cache/apk/* \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /go/autoacme/build/autoacme /usr/bin/autoacme

ENTRYPOINT [ "/usr/bin/autoacme" ]
