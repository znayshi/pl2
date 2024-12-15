# Chat Application

This is a simple chat application written in Go. It supports both a server and client mode, allowing multiple clients to connect to the server and exchange messages.

## Features

- **Server Mode**: The server accepts connections from clients and facilitates messaging between them.
- **Client Mode**: The client connects to the server and allows users to send both public and private messages.
- **Private Messages**: Users can send private messages to specific clients using the `@username` format.
- **Public Messages**: Users can send messages to all connected clients, which will be broadcasted.

## Requirements

- Go 1.16 or later

## Installation

- Clone this repository:

   ```bash
   git clone https://github.com/znayshi/pl2
   cd app
   ```

## Running the Application

To test the chat application with two users, you need to open three terminals:

### Terminal 1: Start the Server

1. Open the first terminal and navigate to the project directory:

   ```bash
   cd app
   ```

2. Run the server:

   ```bash
   go run main.go -server
   ```

   The server will listen on `localhost:9000` by default, but you can change the address with the `-address` flag:

   ```bash
   go run main.go -server -address <address>
   ```

### Terminal 2: Start the First Client

1. Open the second terminal and navigate to the project directory:

   ```bash
   cd app
   ```

2. Run the first client:

   ```bash
   go run main.go
   ```

3. You will be prompted to enter a nickname. For example, enter:

   ```
   Enter your nickname: Matvey
   ```

4. After entering the nickname, you can start sending messages. For example:

   ```
   Hello everyone!
   ```

### Terminal 3: Start the Second Client

1. Open the third terminal and navigate to the project directory:

   ```bash
   cd app
   ```

2. Run the second client:

   ```bash
   go run main.go
   ```

3. You will be prompted to enter a nickname. For example, enter:

   ```
   Enter your nickname: Sasha
   ```

4. Now you can start sending messages as Sasha.

### Testing the Chat

- **Public Message**: Type a regular message to broadcast it to all connected clients. For example:
  
  ```
  Hello everyone!
  ```

- **Private Message**: Use the format `@<username> <message>` to send a private message to a specific user. For example:
  
  ```
  @Matvey Hi, Matvey! How are you?
  ```

  This message will be visible only to Matvey.

## How it Works

- The server listens for incoming client connections.
- Each client can send public and private messages.
- Public messages are broadcasted to all connected clients.
- Private messages are sent only to the specified recipient.