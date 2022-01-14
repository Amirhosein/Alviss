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

## Routes
### `localhost:8080/`:
Send a GET request and get a warm welcome :)
```JSON
{
  "message": "Welcome to Alviss! Your mythical URL shortener."
}
```


### `localhost:8080/shorten`:
POST a JSON object like below and in return, get the generated short link:
```JSON
{
  "LongURL": "https://gist.github.com/joshbuchea/6f47e86d2510bce28f8e7f42ae84c716",
  "ExpTime": "2d"
}
```
Exp date valid format: `2d` for 2 days, `2h` for 2 hours, `2m` for 2 minutes and `2s` for 2 seconds.

#### Result:
```JSON
{
  "message": "Short url created successfully",
  "ShortURL": "http://localhost:8080/ZLgJHJB2"
}
```

### `localhost:8080/url/{YOUR-SHORTENED-URL}`
Send a GET request and get details of your URL, such as `UsedCount` or `ExpDate`:
```JSON
{
  "ExpDate": "2021-12-20T15:38:26.48860767Z",
  "OriginalURL": "https://gist.github.com/joshbuchea/6f47e86d2510bce28f8e7f42ae84c716",
  "ShortURL": "http://localhost:8080/ZLgJHJB2",
  "UsedCount": 3
}
```
