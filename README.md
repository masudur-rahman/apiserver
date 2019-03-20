# API Server

An API Server containing the worker profile. Source code was written in golang. For routing purpose `gorilla-mux` was used. To give it a CLI (command line interface) `cobra` was used...  


## Run apiserver - from SourceCode

At first we need to get something to run the api server.

#### Main commands
`$ go install` - to install the api api

`$ apiserver --help` - to get basic commands of the api

`$ apiserver start` - to start the server

`$ apiserver version` - to get api version

`$ apiserver start --help` - to get to know about flags associated with start

`$ apiserver start --bypass true` - to get a bypass authorization

`$ apiserver start --port 8080 --stopTime 5` - to assign a port to run and to set time to stop the server

 

## Run apiserver - from Dockerfile

`$ docker run -p 4000:8080 <user-name>/<image-name> <additional-argument-if-any` - to directly run from hub.docker.com

`$ docker build -t <new-image-name> .`

`$ docker run -p 8000:8080 <image-name> <additional-argument>` - example : `$ docker run -p 8000:8080 api start --bypass true`

`$ docker run -d --name <new-name> -p <new-port>:<existing-port> <image-name>` - `-d` is used to run as daemon (in background), `--name <name>` to give the image a friendly-name

## Upload to docker hub

`docker push masudjuly02/apiserver:v1alpha1`

i.e. : `docker push <user-name>/<image-name>:<tag>` - if tag is not provided :<tag> can be omitted...

