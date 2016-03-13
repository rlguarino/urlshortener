
# Example Kubernetes Deployment

An example deployment of the URL Shortener using Kubernetes. Specifically tested using [Google Container Engine](https://cloud.google.com/container-engine/).
This README assumes that you are using GCE for the deployment. Replace <project-id> with the id of your project or if you are using standalone Kubernetes 
you can replace project-id with anything.

## Dependencies
 * Access to a working Kubernetes cluster (Probably the simplest is to use Google Container Engine)
 * Docker

## Redis Cluster Setup

You can deploy your own Redis Sentinel cluster or use the [helm.sh](http://helm.sh/) redis-cluster.  


## Mongo Setup

To setup the Mongo cluster run `kubectl create -f kubetnetes-example/mongo/mongo-controller.yaml` and `kubectl create -f kubetnetes-example/mongo/mongo-service.yaml`


## Base Docker Image
Build the base docker image with `docker build -t shortener/base .` in the root directory of the app.


## URL Service Setup

Build the URL Service container with `docker build -t <project-id>/urlservice`
Tag the image for your Docker Registry `docker tag <project-id>/urlservice gcr.io/<project-id>/urlservice`
Push your image to the registry `gcloud docker push gcr.io/<project-id>/urlservice`

Setup your replication controler `kubectl create -f kubetnetes-example/urlservice-controller.yaml` and service `kubectl create -f kubetnetes-example/urlservice-service.yaml`

## Stats Service Setup

Build the Stats Service container with `docker build -t <project-id>/stats`
Tag the image for your Docker Registry `docker tag <project-id>/stats gcr.io/<project-id>/stats`
Push your image to the registry `gcloud docker push gcr.io/<project-id>/stats`

Setup your replication controler `kubectl create -f kubetnetes-example/stats-controller.yaml` and service `kubectl create -f kubetnetes-example/stats-service.yaml`


## Frontend Service Setup

Build the Frontend Service container with `docker build -t <project-id>/frontend`
Tag the image for your Docker Registry `docker tag <project-id>/frontend gcr.io/<project-id>/frontend`
Push your image to the registry `gcloud docker push gcr.io/<project-id>/frontend`

Setup your replication controler `kubectl create -f kubetnetes-example/frontend-controller.yaml` and service `kubectl create -f kubetnetes-example/frontend-service.yaml`


## Connect

Find the external IP of the Frontend service using `kubectl get services` and use the external ip to connec to the url shortener.