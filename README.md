![Quiz](banner.png)
![Language | Go](https://img.shields.io/badge/Language-Go-brightgreen)
![OS | Linux](https://img.shields.io/badge/OS-Linux-brightgreen)


# Quiz

This project is two-player quiz game written in go programming language. In this project TCP was used for networking and SQLite was used for question storage.

## Features
- **Two player game**: Players can connect using TCP and play against each other in real time.
- **SQLite database**: All quiz questions and answers are stored in SQLite database.
- **Real-time score updates**: The current score for both players is updated during the game. In the end of the game player can see number of points he achieved and also if he won, lost or it was draw.
- **Private rooms**: Players can create or join private rooms using random generated codes.

## Technologies
- **Go**: Handles backend logic and the TCP server.
- **Fyne**: GUI library used to create the client interface.
- **SQLite**: Lightweight database used to store quiz content.
- **TCP**: Enables real-time network communication for multiplayer gameplay.

## Installation
Copy this in your Linux terminal to get necessary libraries:

```bash
sudo apt update
sudo apt install sqlite3
go get github.com/mattn/go-sqlite3
go get fyne.io/fyne/v2
go install fyne.io/fyne/v2/cmd/fyne-cross@latest
```
After that:

1. Download all the files from this repository (or clone it using `git clone`).
2. Open a terminal inside the project folder.
3. Run:

```bash
make
```

You can now find and launch Kviz from your application menu.

## How to run 

**1. Run the server**
   Use the terminal to run the server:
   
  ```bash
  go run server.go
  ```
The server will start and wait for players to connect via TCP. 

**2. Run the clients**

Before running the application, you need to specify the server IP address in the _config.txt_ file (e.g. 192.168.111.111:8082).
Important: Do not change the port number :8082. Only replace the IP address with the correct one for your server.

Simply double-click the client application (e.g., kviz.exe) to launch the quiz interface.
Once the client starts, create or join a private room using a code and start playing the quiz in real-time with another player.

## Contributors
 - [@danica-mijajlovic](https://github.com/danica-mijajlovic)
 - [@djordje-mitrovic](https://github.com/djordje-mitrovic)
 - [@dobrivojetrifunovic](https://github.com/dobrivojetrifunovic)
   
