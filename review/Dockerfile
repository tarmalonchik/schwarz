# build stage
# It is better to use mutlistage build, because you will not have source code in docker image.
# Also in case of issue, you can copy any sensetive configs to docker image
# You docker image will be lighter
FROM golang:1.23.2 AS builder

WORKDIR /go/src/coupon

COPY cmd cmd
COPY internal internal
COPY go.mod go.mod
COPY go.sum go.sum

RUN go build -o bin/coupon/main cmd/coupon/main.go

FROM alpine:latest as base
RUN apk add bash nano build-base gcompat
WORKDIR /app
COPY --from=builder /go/src/coupon/bin/coupon/main ./
CMD ["/app/main"]
