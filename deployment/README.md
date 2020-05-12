# dapr-tracing-demo

This document will overview the deployment of `dapr-tracing-demo` onto Kubernetes. For illustration purposes, all commands in this document will be based on Microsoft Azure. 


## Create namespace 

```shell
kubectl create namespace trace-demo
```

## Deploy components 

```shell
kubectl apply -f components/
```

## Deploying demo 

```shell
kubectl apply -f deployment/
```

Watch pods and make sure the four services are up

```shell
kubectl apply -f ../components/
```


```shell
NAME                           READY   STATUS    RESTARTS   AGE
processor-89666d54b-hkd5t      2/2     Running   0          18s
sentimenter-85cfbf5456-lc85g   2/2     Running   0          18s
viewer-76448d65fb-bm2dc        2/2     Running   0          18s
```

### Exposing viewer UI

To expose the viewer application externally, create Kubernetes `service` using [deployment/service/viewer.yaml](deployment/service/viewer.yaml)

```shell
kubectl apply -f service/viewer.yaml
```

> Note, the provisioning of External IP may take up to 1-2 min 

To view the viewer application by capturing the load balancer public IP and opening it in the browser:

```shell
export VIEWER_IP=$(kubectl get svc viewer --output 'jsonpath={.status.loadBalancer.ingress[0].ip}')
open "http://${VIEWER_IP}/"
```

> To change the Twitter topic query, first edit the [deployment/component/twitter.yaml](deployment/component/twitter.yaml), then apply it (`kubectl apply -f component/twitter.yaml`), and finally, restart the processor (`kubectl rollout restart deployment processor`) to ensure the new configuration is applied. 


## Observability 


> Instructions on how to setup Zipkin for Dapr are [here](https://github.com/dapr/docs/blob/master/howto/diagnose-with-tracing/zipkin.md)

![](../img/trace.png)

http://localhost:9411/zipkin/

> Note, if your Zipkin isn't deployed in the `default` namespace you will have to edit the `exporterAddress` in [deployment/tracing/zipkin.yaml](deployment/tracing/zipkin.yaml)


Then just restart all the deployments 

```shell
kubectl rollout restart deployment processor sentimenter  viewer
```






