# API gateway Service

This is the [API gateway](https://microservices.io/patterns/apigateway.html) of Delinkcious. It is the only externally visible service.
It communicates with the other microservices using [API composition](https://microservices.io/patterns/data/api-composition.html):
- link service
- user service
- social graph service