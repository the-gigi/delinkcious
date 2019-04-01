FROM golang:1.11 AS builder
ADD ./main.go main.go
# Fetch dependencies
RUN go get -d -v

# Build image as a truly static Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /link_checker -a -tags netgo -ldflags '-s -w' .

FROM scratch
MAINTAINER Gigi Sayfan <the.gigi@gmail.com>
COPY --from=builder /link_checker /app/link_checker
EXPOSE 7070
ENTRYPOINT ["/app/link_checker"]
