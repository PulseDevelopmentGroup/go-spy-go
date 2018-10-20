# Readme Stuff

This is a game, this is also a readme that needs a good updating

# JSON Info

These are just for development. Some reference material for how web socket messages are to be sent and what responce(s) should be expected

## Create new game without code
### Message from client
```json
{
    "type" : "create-game",
    "data" : "{\"code\":\"\", \"username\":\"USERNAME\"}"
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
    "data" : "ERROR MESSAGE"
}
```

## Create new game with code
### Message from client
```json
{
    "type":"create-game",
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
    "type" : "join-game",
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

