FROM golang:1.11 AS builder
ADD ./main.go main.go
ADD ./service service

# Fetch dependencies
RUN go get -d -v

# Build image as a truly static Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /social_graph_service -a -tags netgo -ldflags '-w' .

FROM scratch
MAINTAINER Gigi Sayufan <the.gigi@gmail.com>
COPY --from=builder /social_graph_service /social_graph_service
EXPOSE 9090
ENTRYPOINT ["/social_graph_service"]