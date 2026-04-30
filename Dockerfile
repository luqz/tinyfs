FROM golang:1.26-alpine AS builder
WORKDIR /build
COPY go.mod main.go ./
RUN go build -o fileserver .

FROM alpine:3.20
COPY --from=builder /build/fileserver /usr/local/bin/fileserver
EXPOSE 8082
ENTRYPOINT ["fileserver"]
CMD ["-listen", ":8082", "-dir", "/data"]
