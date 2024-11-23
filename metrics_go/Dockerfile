FROM golang:1.21-alpine3.20 AS builder
WORKDIR /build/
COPY . .
RUN go build -o main main.go

FROM alpine:3.20
WORKDIR /main/
COPY --from=builder /build/main .
ENTRYPOINT [ "/main/main" ]
