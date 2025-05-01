package main

import(
	"strings"
	"net"
	"os"
    "log"
)

func main() {
    log.SetFlags(log.Ltime)

    var listenAddress string
    var serverAddress string
    var username string

    switch len(os.Args) {
    case 2:
        listenAddress = net.JoinHostPort("", os.Args[1])
    case 3:
        serverAddress = os.Args[1]
        if strings.HasPrefix(serverAddress, ":") {
            serverAddress = "localhost" + serverAddress
        }
        username = strings.TrimSpace(os.Args[2])
        if username == "" {
            log.Fatal("empty user name")
        }
    default:
        log.Fatalf("Usage: %s <port>   OR   %s <server> <user>",
            os.Args[0], os.Args[0])
    }   

    if len(listenAddress) > 0 { 
        server(listenAddress)
    } else {
        client(serverAddress, username)
    }   
}