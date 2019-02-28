# gRPC contracts and generated code

This directory contains a sub-directory for each gRPC service

The service directory contains another sub-directory called pb as well: `pb/<service>/pb`.
The internal pb directory contains the following files:

- <service>.proto file: gRPC contract)
- <service>.pb.go: Go server and client generated code
- <service>_pb2.py: Python generated code
- <service>_pb2_grpc.py: Python class called <service>Stub used by Python clients

The purpose of this directory is to serve as shared code for both gRPC services and gRPC clients
because the generated code is a mix of both service and client code. Putting this code with the gRPC service itself
means that client code would have to import the service code, which is a violation of separation of concerns.