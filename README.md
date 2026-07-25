# Tic-Tac-Toe

A multiplayer Tic-Tac-Toe game written in Go using TCP sockets. The project consists of a TCP server that manages game sessions and a terminal-based client that allows two players to play against each other over a network.

## Features

- Multiplayer gameplay over TCP
- Automatic player matchmaking
- Configurable server address
- Turn-based gameplay
- Client-side input validation
- Server-side game validation
- Win and draw detection
- Graceful handling of player disconnects

## Project Structure

```text
.
├── cmd/
│   ├── client/        # Client application
│   └── server/        # Server application
├── internal/
│   ├── game/          # Game engine and rules
│   └── server/        # TCP server and session management
├── go.mod
├── Makefile
└── README.md
```

## Requirements

- Go 1.25 or later

## Building

Build both the client and server:

```bash
make build
```

The executables are placed in the `bin/` directory.

## Running the Server

Start the server using the default address (`:8080`):

```bash
make run-server
```

or

```bash
go run ./cmd/server
```

To listen on a different address:

```bash
go run ./cmd/server -address :9000
```

## Running the Client

Connect to the default server (`localhost:8080`):

```bash
make run-client
```

or

```bash
go run ./cmd/client
```

Connect to a different server:

```bash
go run ./cmd/client -server localhost:9000
```

Connect from another machine on the local network:

```bash
go run ./cmd/client -server 192.168.1.245:9000
```

## Gameplay

1. Start the server.
2. Connect two clients.
3. Players are automatically assigned as **X** and **O**.
4. Take turns entering the number corresponding to the desired board position.

Board positions:

```text
1 | 2 | 3
---------
4 | 5 | 6
---------
7 | 8 | 9
```

Example session:

```text
Connected to server. Board Size: 3 x 3
Waiting for another player to connect...

You are X

___
___
___

Your turn X
```

## Testing

Run all unit tests:

```bash
make test
```

or

```bash
go test ./...
```

## Available Make Targets

| Command | Description |
|---------|-------------|
| `make build` | Build the client and server |
| `make run-server` | Start the server |
| `make run-client` | Start the client |
| `make test` | Run all tests |
| `make fmt` | Format the source code |
| `make vet` | Run `go vet` |
| `make clean` | Remove generated binaries |

## Future Improvements

Some possible enhancements include:

- Support for configurable board sizes from the command line
- Multiple concurrent game sessions
- Spectator mode
- Better handling of disconnected waiting players
- AI opponent
- Graphical user interface

## License

This project was developed as a learning exercise to explore Go networking, concurrency, and TCP client-server application development.