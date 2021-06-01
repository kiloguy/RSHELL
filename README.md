## RSHELL

A multi-clients reverse shell implementation in Go.

### Build

#### Server

```
go build -o rshell_serv -ldflags="-w -s" rshell_serv.go rshell_serv_help.go
```
(**Note**: The server code should works both on Windows and Linux, but the Ctrl-C signal can't
work properly on my Windows and I don't know why, so I use it on Linux.)

#### Client

```
go build -o rshell.exe -ldflags="-w -s -H=windowsgui" rshell.go
```
##### On Linux

The client code is for Windows, but can be modified to run on Linux by following steps:
1. Comment out import of `syscall`.
2. Comment out this line: `cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}`.
3. Change `cmd := exec.Command("C:\\Windows\\System32\\cmd.exe", "/k", "chcp", "65001")` to `cmd := exec.Command("/bin/bash", "-i")`,
or other shell you want.
4. Build again without the `-H=windowsgui` option.

Then it should be run without problem.

### Guide

#### Server commands
`l`: List currently accepted connections (first field is connection IDs).

`<Number>`: Spawn cmd.exe on client of specified connection ID.

`h`: This help message.

`q`: Quit server (clients won't quit).

`Ctrl-C` (When communicating with client): Close current client's connection.
        
#### Client behavior

Once client program starts, it will keep trying to connect server. Once it got accepted by the server, it waits server inform 
to start cmd.exe. After cmd.exe exit or server disconnect (like server press Ctrl-C), it will close the connection and cmd.exe,
then try to connect server again (server will get a new connection ID). So, **the client runs forever if nothing bad happended
even after quitting server.**

(**Note: This client spawns a cmd.exe process, so even using techniques like
DLL injection to hide the client program itself, cmd.exe will still appear in
task manager.**)