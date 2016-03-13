#!/usr/bin/env bash

set -e

# Build docker containers for each service
docker_build(){
    echo "Docker Build"
    docker build -t shortener/base .
    
    cd urlservice && docker build -t shortener/urlservice . && cd .. 
    cd stats && docker build -t shortener/stats . && cd .. 
    cd frontend && docker build -t shortener/frontend . && cd .. 
}

# Build go executables for each service
go_build(){
    echo "Build URL Service"
    cd urlservice && go get && go build && cd .. 
    echo "Build Stats Service"
    cd stats && go get && go build && cd .. 
    echo "Build FrontEnd"
    cd frontend && go get && go build && cd .. 
}


case "$1" in
    docker)
        docker_build
        ;;
        
    go)
        go_build
        ;; 
    *)
        echo "Usage: $0 {docker|go}"
esac