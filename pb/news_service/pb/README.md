# News service gRPC contract

This directory contains the news.proto file that defines the gRPC contract of the news service.
Whenever you change it run the following commands:

```
protoc news.proto --go_out=plugins=grpc:.
python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. news.proto
```

The first command generates `news.pb.go` that the [news service](https://github.com/the-gigi/delinkcious/tree/master/svc/news_service) imports.

The second command generates two Python modules for Python consumers of the service.
