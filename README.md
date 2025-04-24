![Quiz](banner.png)

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

## How to run 
**1. Run the server**
   Use the terminal to run the server:
   
  ```bash
  go run server.go
  ```
The server will start and wait for players to connect via TCP.

**2. Run the clients**
Simply double-click the client application (e.g., kviz.exe) to launch the quiz interface.
Once the client starts, create or join a private room using a code and start playing the quiz in real-time with another player.

## Contributors
 - [@danica-mijajlovic](https://github.com/danica-mijajlovic)
 - [@djordje-mitrovic](https://github.com/djordje-mitrovic)
 - [@dobrivojetrifunovic](https://github.com/dobrivojetrifunovic)
   
