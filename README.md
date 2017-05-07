# mcs-go
Very small server written with Go. Created in a day.

Works with [mcs-go-react](https://github.com/TKasekamp/mcs-gui-react), the front-end for this.

Should be run like a standard Go app. 

## How it works
The idea is that the user will make a POST to `/commands` and is immediately returned the command id and a notification that the command has been accepted by the server.

The server then "works" (actually time.sleep(3000)) and then publishes the command result to all websockets. There are two random results.

`/commands` is for checking the command statuses for development.

## API
* POST /commands to submit a command for processing
* GET /commands for a list of all commands submitted
* GET /ws to connect to commands websocket

## Features
* UUID generation with `github.com/satori/go.uuid`
* API stuff with `github.com/gin-gonic/gin`
* Websocket management with `github.com/gorilla/websocket`