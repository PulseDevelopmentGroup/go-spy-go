# Readme Stuff

Go-Spy-Go is an implementation of the [Spyfall](http://international.hobbyworld.ru/spyfall) party game, built by [Josiah](https://josnun.github.io/) and [Carson](https://carsonseese.com) in React and Go.

This is mostly a learning-focused project, so the codebase is a mess. It will get better over time... hopefully

# Go Development

To develop in this project with Go, you must set your `$GOPATH` enviorment variable to the `/server` directory in the root of this project. The project can be run by running `go run main.go` from the `/server/src/spyfall` directory or running `go run spyfall` from the `/server` directory.

# JSON Info

These are just for development. Some reference material for how web socket messages are to be sent and what responce(s) should be expected

## Create new game without code

### Message from client

```json
{
    "kind":"CREATE_GAME",
    "data": "{\"game-id\":\"\",\"username\":\"USERNAME\"}"
}
```

### Response from server if OK

(Notice it's the exact same thing as the message from the client)

```json
{
    "kind":"CREATE_GAME",
    "data": "{\"game-id\":\"\",\"username\":\"USERNAME\"}"
}
```

### Response from server if Error

(Notice it's the exact same thing as the message from the client, just with the generated game-id (if applicable) and the error)

```json
{
    "kind":"CREATE_GAME",
    "data":"{\"game-id\":\"GAMEID\",\"username\":\"USERNAME\"}",
    "error":"{\"error\":\"ERRORCODE\",\"description\":\"ERRORDESC\"}"
    }
```

## Create new game with code

### Message from client

```json
{
    "trigger":"create-game",
    "data":"{\"code\":\"GAMECODE\", \"username\":\"USERNAME\"}"
}
```

### Response from server if OK

```json
{
    "response" : "OK",
    "data" : "GAMECODE"
}
```

### Response from server if Error

```json
{
    "response" : "ERROR",
    "data" : "ERROR_MESSAGE"
}
```

## Join game with code

### Message from client

```json
{
    "trigger" : "join-game",
    "data" : "{\"code\":\"GAMECODE\", \"username\":\"USERNAME\"}"
}
```

### Response from server if OK

```json
{
    "response" : "OK",
    "data" : "GAMECODE"
}
```

### Response from server if Error

```json
{
    "response" : "ERROR",
    "data" : "ERROR_MESSAGE"
}
```