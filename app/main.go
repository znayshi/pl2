package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

var (
	isServer = flag.Bool("server", false, "Run as server")
	address  = flag.String("address", "localhost:9000", "Address to bind or connect to")
)

type Client struct {
	nick string
	conn net.Conn
}

var (
	clients     = make(map[string]*Client)
	clientsLock = sync.Mutex{}
)

func main() {
	flag.Parse()
	if *isServer {
		runServer(*address)
	} else {
		runClient(*address)
	}
}

func runServer(addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server started on", addr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	var nick string
	for {
		fmt.Fprintln(conn, "Enter your nickname:")
		n, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading nickname:", err)
			return
		}
		nick = strings.TrimSpace(n)
		if nick != "" {
			clientsLock.Lock()
			if _, exists := clients[nick]; exists {
				clientsLock.Unlock()
				fmt.Fprintln(conn, "Nickname already taken. Please choose another.")
			} else {
				clients[nick] = &Client{nick: nick, conn: conn}
				clientsLock.Unlock()
				break
			}
		} else {
			fmt.Fprintln(conn, "Nickname cannot be empty.")
		}
	}

	fmt.Println("New client connected:", nick)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Client disconnected:", nick)
			break
		}
		handleMessage(nick, strings.TrimSpace(msg))
	}

	// Cleanup
	clientsLock.Lock()
	delete(clients, nick)
	clientsLock.Unlock()
}

func handleMessage(senderNick, msg string) {
	clientsLock.Lock()
	defer clientsLock.Unlock()

	if strings.HasPrefix(msg, "@") {
		// Personal message
		parts := strings.SplitN(msg, " ", 2)
		if len(parts) < 2 {
			sender := clients[senderNick]
			fmt.Fprintf(sender.conn, "Invalid message format. Use @username <message>\n")
			return
		}
		recipientNick := strings.TrimPrefix(parts[0], "@")
		message := parts[1]

		recipient, ok := clients[recipientNick]
		if ok {
			fmt.Fprintf(recipient.conn, "[Private] %s -> %s: %s\n", senderNick, recipientNick, message)
			sender := clients[senderNick]
			fmt.Fprintf(sender.conn, "[Private] %s -> %s: %s\n", senderNick, recipientNick, message)
		} else {
			sender := clients[senderNick]
			fmt.Fprintf(sender.conn, "User '%s' not found.\n", recipientNick)
		}
	} else {
		for nick, client := range clients {
			if nick != senderNick {
				fmt.Fprintf(client.conn, "[Public] %s: %s\n", senderNick, msg)
			} else {
				fmt.Fprintf(client.conn, "[You] %s: %s\n", senderNick, msg)
			}
		}
	}
}

func runClient(addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)
	serverReader := bufio.NewReader(conn)

	var nick string
	for {
		fmt.Print("Enter your nickname: ")
		n, _ := reader.ReadString('\n')
		nick = strings.TrimSpace(n)
		if nick != "" {
			break
		}
		fmt.Println("Nickname cannot be empty.")
	}
	fmt.Fprintln(conn, nick)

	confirmation, _ := serverReader.ReadString('\n')
	fmt.Print(confirmation)

	go func() {
		for {
			msg, err := serverReader.ReadString('\n')
			if err != nil {
				fmt.Println("Disconnected from server.")
				os.Exit(0)
			}
			fmt.Print(msg)
		}
	}()

	fmt.Println("You can start chatting. Use @username to message a specific user.")
	for {
		msg, _ := reader.ReadString('\n')
		msg = strings.TrimSpace(msg)
		if msg != "" {
			fmt.Fprintln(conn, msg)
		}
	}
}
