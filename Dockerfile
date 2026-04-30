FROM golang:1.26-alpine AS builder
WORKDIR /build
COPY go.mod main.go ./
RUN go build -o tinyfs .

FROM alpine:3.20
COPY --from=builder /build/tinyfs /usr/local/bin/tinyfs
EXPOSE 8082
ENTRYPOINT ["tinyfs"]
CMD ["-listen", ":8082", "-dir", "/data"]
