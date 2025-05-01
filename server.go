package main

import(
	"log"
	"net"
	"sync"
	"time"
	"io"
	"fmt"
	"unicode"
	"sort"
)

const (
	MsgRegister = iota
	MsgList
	MsgCheckMessages
	MsgTell
	MsgSay
	MsgQuit
	MsgShutdown
)

var mutex sync.Mutex
var messages map[string][]string
var shutdown chan struct{}

func server(listenAddress string) {
	shutdown = make(chan struct{})
	messages = make(map[string][]string)

	// set up network listen and accept loop here
	// to receive RPC requests and dispatch each
	// in its own goroutine
	//////
	ln, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	log.Printf("Server listening on %s", listenAddress)

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Printf("Failed to accept connection: %v", err)
				continue
			}
			go dispatch(conn)
		}
	}()
	// wait for a shutdown request
	<-shutdown
	time.Sleep(100 * time.Millisecond)
}

func dispatch(conn net.Conn) {
	defer conn.Close()

	rawLen := make([]byte, 2)
	if _, err := io.ReadFull(conn, rawLen); err != nil {
		log.Printf("Dispatch occured an error while reading the message length: ", err)
		return
	}

	//fmt.Println("Sending request:", rawLen)
	requestLen, _, err := ReadUint16(rawLen)

	if err != nil {
		panic("The ReadUint16 occured an error in dispatch failing to read 2-byte buffer")
	}

	request := make([]byte, requestLen)
	if _, err := io.ReadFull(conn, request); err != nil {
		log.Printf("Dispatch occured an error while reading the message body: ", err)
		return
	}

	msgType, request, err := ReadUint16(request)

	if err != nil {
		log.Printf("The Dispatch occured an error reading the request type: ", err)
		return
	}

	var response []byte
	switch msgType {
	case MsgRegister:
		response = RegisterServerStub(request)
	case MsgList:
		response = ListServerStub(request)
	case MsgCheckMessages:
		response = CheckMessagesServerStub(request)
	case MsgTell:
		response = TellServerStub(request)
	case MsgSay:
		response = SayServerStub(request)
	case MsgQuit:
		response = QuitServerStub(request)
	case MsgShutdown:
		response = ShutdownServerStub(request)
	}

	responseLen := uint16(len(response))
	lengthBytes := WriteUint16(nil, responseLen)

	if _, err := conn.Write(lengthBytes); err != nil {
		log.Printf("Dispatch encountered an error while writing response length: %v", err)
		return
	}

	if _, err := conn.Write(response); err != nil {
		log.Printf("Dispatch encountered an error while writing response body: %v", err)
		return
	}
}

func Register(user string) error {
	if len(user) < 1 || len(user) > 20 {
			return fmt.Errorf("Register: user must be between 1 and 20 letters")
	}
	for _, r := range user {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return fmt.Errorf("Register: user must only contain letters and digits")
		}
	}
	mutex.Lock()
	defer mutex.Unlock()

	msg := fmt.Sprintf("*** %s has logged in", user)
	log.Printf(msg)
	for target, queue := range messages {
		messages[target] = append(queue, msg)
	}
	messages[user] = nil

	return nil
}

func List() []string {
	mutex.Lock()
	defer mutex.Unlock()

	var users []string
	for target := range messages {
		users = append(users, target)
	}
	sort.Strings(users)

	return users
}

func CheckMessages(user string) []string {
	mutex.Lock()
	defer mutex.Unlock()

	if queue, present := messages[user]; present {
		messages[user] = nil
		return queue
	} else {
		return []string{"*** You are not logged in, " + user}
	}
}

func Tell(user, target, message string) {
	mutex.Lock()
	defer mutex.Unlock()

	msg := fmt.Sprintf("%s tells you %s", user, message)
	if queue, present := messages[target]; present {
		messages[target] = append(queue, msg)
	} else if queue, present := messages[user]; present {
		messages[user] = append(queue, "*** No such user: "+target)
	}
}

func Say(user, message string) {
	mutex.Lock()
	defer mutex.Unlock()

	msg := fmt.Sprintf("%s says %s", user, message)
	for target, queue := range messages {
		messages[target] = append(queue, msg)
	}
}

func Quit(user string) {
	mutex.Lock()
	defer mutex.Unlock()

	msg := fmt.Sprintf("*** %s has logged out", user)
	log.Print(msg)
	for target, queue := range messages {
			messages[target] = append(queue, msg)
	}
	delete(messages, user)
}

func Shutdown() {
	shutdown <- struct{}{}
}