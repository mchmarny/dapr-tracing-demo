# dapr-tracing-demo

Dapr tracing demo integrating multiple self-container microservices (no dependencies) to illustrate end-to-end tracing where the tracing headers are being propagated across:

* service to service invocation 
* input binding
* state operation 
* pubsub publish
* pubsub subscribe 
* output binding  

![alt text](img/overview.png "Overview")

## Prerequisites

* To run this demo locally, you will have to install [Dapr](https://github.com/dapr/docs/blob/master/getting-started/environment-setup.md).
* Additionally to view the resulting traces you will need Zipkin. Instructions on how to setup Zipkin for Dapr are [here](https://github.com/dapr/docs/blob/master/howto/diagnose-with-tracing/zipkin.md)

## Setup

Assuming you have all the prerequisites mentioned above, first, start by cloning this repo:

```shell
git clone https://github.com/mchmarny/dapr-tracing-demo.git
```

navigate into the `dapr-tracing-demo` directory

```shell
cd dapr-tracing-demo
```

and build the executables for your OS

```shell
bin/build
```

You should see now 4 files in the [dist](dist) directory: `producer`, `formatter`, `subscriber`, `sender`

## Run

Inside of the `dapr-tracing-demo` directory, start each one of the service individually:

### Formatter

```shell
dapr run dist/formatter --app-id formatter --app-port 808 --protocol http
```

### Subscriber

```shell
dapr run dist/subscriber --app-id subscriber --app-port 8083 --protocol http
```

### Producer

```shell
dapr run dist/producer --app-id producer --app-port 8081 --protocol http --port 3500
```

### Sender

> Sender is not a Dapr service, it will serve as a target for output binding 

```shell
dist/sender
```

## Observability 


![](../img/trace.png)

http://localhost:9411/zipkin/

> Note, if your Zipkin isn't deployed in the `default` namespace you will have to edit the `exporterAddress` in [deployment/tracing/zipkin.yaml](deployment/tracing/zipkin.yaml)


Then just restart all the deployments 

```shell
kubectl rollout restart deployment processor sentimenter  viewer
```

## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.

## License
This software is released under the [Apache v2 License](./LICENSE)


