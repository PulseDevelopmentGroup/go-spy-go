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
  "kind": "CREATE_GAME",
  "data": "{\"gameid\":\"\",\"username\":\"USERNAME\"}"
}
```

### Response from server if OK

(Notice it's the exact same thing as the message from the client. There will be a second response wich matches the sucessful join response)

```json
{
  "kind": "CREATE_GAME",
  "data": "{\"gameid\":\"\",\"username\":\"USERNAME\"}"
}
```

### Response from server if Error

(Notice it's the exact same thing as the message from the client, just with the generated gameId (if applicable) and the error)

```json
{
  "kind": "CREATE_GAME",
  "data": "{\"gameid\":\"GAMEID\",\"username\":\"USERNAME\"}",
  "error": "{\"error\":\"ERRORCODE\",\"description\":\"ERRORDESC\"}"
}
```

## Create new game with code

### Message from client

```json
{
  "kind": "CREATE_GAME",
  "data": "{\"gameid\":\"GAMEID\", \"username\":\"USERNAME\"}"
}
```

### Response from server if OK

```json
{
  "kind": "CREATE_GAME",
  "data": "{\"gameid\":\"GAMEID\",\"username\":\"USERNAME\"}"
}
```

### Response from server if Error

```json
{
  "kind": "CREATE_GAME",
  "data": "{\"gameid\":\"GAMEID\",\"username\":\"USERNAME\"}",
  "error": "{\"error\":\"ERRORCODE\",\"description\":\"ERRORDESC\"}"
}
```

## Join game with code

### Message from client

```json
{
  "kind": "JOIN_GAME",
  "data": "{\"gameid\":\"GAMEID\", \"username\":\"USERNAME\"}"
}
```

### Response from server if OK

```json
{
  "kind":"JOIN_GAME",
  "data":"{\"gameid\":\"GAMEID\",\"username\":\"USERNAME\"}"
}

```

### Response from server if Error

```json
{
  "kind":"JOIN_GAME",
  "data":"{\"gameid\":\"GAMEID\",\"username\":\"USERNAME\"}",
  "error":"{\"error\":\"ERRORCODE\",\"description\":\"ERRORDESC\"}"
}

```

## Start Game

### Message from client

Start game is based on the websocket connection (so that it is not possible to start a game you are not a part of) so the message from the client does not require any associated data.

```json
{
  "kind":"START_GAME",
  "data":"{}"
}
```

### Message from server if OK

*//TODO Add stuff here*

### Message from server if Error

*//TODO Add stuff here*

## Leave Game

The same results can be achieved by simply closing the websocket connection

### Message from client

Start game is based on the websocket connection (so that it is not possible to leave a game you are not a part of) so the message from the client does not require any associated data.

```json
{
  "kind":"LEAVE_GAME",
  "data":"{}"
}
```

### Message from server if OK

*//TODO Add stuff here*

### Message from server if Error

*//TODO Add stuff here*

## Error Codes

### Create Game Errors

| Error Code            | Error Description                                         | Response Description                            |
|-----------------------|-----------------------------------------------------------|-------------------------------------------------|
| `GAME_ALREADY_EXISTS` | Game already exists in the database.                      | `Game: \"GAMEID\"  already exists in database.` |
| `UNKNOWN_ERROR`       | Something bad happened, but the server doesn't know what. | `This shouldn't happen, see the server log for details.`                            |

### Join Game Errors

| Error Code            | Error Description                                            | Response Description                                                         |
|-----------------------|--------------------------------------------------------------|------------------------------------------------------------------------------|
| `NO_GAME_CODE`        | The user didn't supply a game code.                          | `Good luck joining a game with no code!`                                     |
| `USER_ALREADY_EXISTS` | That user already exists in the database for that game code. | `A user with the username: \"USERNAME\" already exists in game: \"GAMEID\".` |
| `UNKNOWN_ERROR`       | Something bad happened, but the server doesn't know what.    | `This shouldn't happen, see the server log for details.`                     |