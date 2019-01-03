FROM golang:1.11 AS builder
ADD ./main.go main.go
ADD ./service service

# Fetch dependencies
RUN go get -d -v

# Build image
RUN go build -o /social_graph_service

FROM scratch
MAINTAINER Gigi Sayufan <the.gigi@gmail.com>
COPY --from=builder /social_graph_service /social_graph_service
EXPOSE 9090
ENTRYPOINT ["/social_graph_service"]