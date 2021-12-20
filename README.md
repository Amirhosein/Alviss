# Alviss

## Introduction
Simple URL shortener project, written in Golang.

## Setup and run

You can easily use the following command to simply run the project:
```
docker-compose -f deployments/docker-compose.yml up
```

Otherwise, you can install the project with `go install` in `cmd/alviss/` directory and then use:
```
alviss runserver
```
To run the project on your machine; and there is a `-p` or `--port` flag to specify server port too.

**Note: if you are going to run the project on your local machine, you must have a `redis-server` running on background.
