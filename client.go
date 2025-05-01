package main

import(
	"bufio"
	"fmt"
	"strings"
	"net"
	"os"
	"time"
	"log"
	"io"
)

// Register loop for client 
// Client loop to handle user input and server communication

//client(serverAddress, user)
func client(server string, user string) { // CLIENT LOOP
	// Register the user with the server before proceeding
	err := RegisterRPC(server, user)
	if err != nil {
		fmt.Println("Registration failed:", err)
		return
	}
	fmt.Println("Registration successful.")


	reader := bufio.NewScanner(os.Stdin)
	go pollMessages(server, user)

	for {
		fmt.Print("> ")
		if !reader.Scan() {
			break
		}
		input := strings.TrimSpace(reader.Text())
		if input == "" {
			continue
		}
		fields := strings.Fields(input)
		cmd := fields[0]

		switch cmd {
		case "tell":
			if len(fields) < 3 {
				fmt.Println("Usage: tell <user> <message>")
				continue
			}
			targetUser := fields[1]
			message := strings.Join(fields[2:], " ")
			TellRPC(server, user, targetUser, message)

		case "say":
			if len(fields) < 2 {
				fmt.Println("Usage: say <message>")
				continue
			}
			message := strings.Join(fields[1:], " ")
			SayRPC(server, user, message)

		case "list":
			users, err := ListRPC(server)
			if err != nil {
				fmt.Println("Error retrieving user list:", err)
			} else {
				fmt.Println("Online users:")
				for _, user := range users {
					fmt.Println(" -", user)
				}
			}

		case "quit":
			QuitRPC(server, user)
			fmt.Println("Goodbye!")
			return

		case "shutdown":
			ShutdownRPC(server)
			return

		case "help":
			fmt.Println("Commands:")
			fmt.Println("  tell <user> <message>  - Send a private message")
			fmt.Println("  say <message>          - Broadcast a message")
			fmt.Println("  list                   - List online users")
			fmt.Println("  quit                   - Logout and exit")
			fmt.Println("  shutdown               - Shutdown the server")
			fmt.Println("  help                   - Show this message")

		default:
			fmt.Println("Unrecognized command. Type 'help' for a list of commands.")
		}
	}
}

func pollMessages(server string, user string) {
	for {
		messages, err := CheckMessagesRPC(server, user)
		if err != nil{
			log.Fatal("CheckMessagesRPC ", err)
		}

		if err == nil && len(messages) > 1 {
			for _, msg := range messages {
				fmt.Println(msg)
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func SendAndReceive(server string, request []byte) ([]byte, error){
	conn, err := net.Dial("tcp", server)
	if err != nil {
		log.Print("SendAndRecieve had error connecting to server", server, ": ", err)
		return nil, err
	}
	
	defer conn.Close()

	//fmt.Println("Sending request:", request)

	var requestLen []byte
	requestLen = WriteUint16(requestLen, uint16(len(request)))
	if _, err := conn.Write(requestLen); err != nil{
		log.Print("SendAndRecieve had error writing message length: ", err)
		return nil, err
	}

	if _, err := conn.Write(request); err != nil{
		log.Print("SendAndRecieve had error writing body message: ", err)
		return nil, err
	}
	
	//Recieve response
	rawLen := make([]byte, 2)

	//n, err := io.ReadFull(conn, rawLen);
	_, err = io.ReadFull(conn, rawLen);
	if err != nil{
		log.Print("SendAndRecieve had error reading message length: ", err)
		return nil, err
	}

	responseLen, _, err := ReadUint16(rawLen)

	if err != nil{
		panic("SendAndRecieve ReadUint16 failed to read the 2-byte value")
	}
	response := make([]byte, responseLen)
	if _, err := io.ReadFull(conn, response); err != nil{
		log.Print("SendAndRecieve had error reading message body: ", err)
		return nil, err
	}
	
	// Continue on from here
	return response, nil
}


func RegisterRPC(server string, user string) error {
	// CREATE/FORM REQUEST
	var request []byte
	request = WriteUint16(request, MsgRegister)
	request = WriteString(request, user)
	

	// SEND REQUEST AND RECEIVE RESPONSE
	
	response, err := SendAndReceive(server, request)
	// RETURN an error if error != nil 
	if err != nil {
		return err // Return error if communication failed
	}

	// Decode server response
	msg, _, err := ReadString(response)
	if err != nil {
		return err
	}
	if msg != "" {
		return fmt.Errorf("registration failed: %s", msg)
	}

	return nil
}

func RegisterServerStub(request []byte) []byte {
	user, _, err := ReadString(request)
	if err != nil {
		return WriteString(nil, "*** Error reading user")
	}

	if err := Register(user); err != nil {
		return WriteString(nil, err.Error())
	}

	return WriteString(nil, "")
}


func ListRPC(server string) ([]string, error) {
	// CREATE/FORM REQUEST
	var request []byte
	request = WriteUint16(request, MsgList)
	//request = WriteUint16(request, server)


	// SEND REQUEST AND RECEIVE RESPONSE
	response, err := SendAndReceive(server, request)
	// RETURN an error if error != nil 
	if err != nil {
		return nil, err // Return error if communication failed
	}

	// DECODE RESPONSE
	var users []string
	for len(response) > 0 {
		var user string
		user, response, err = ReadString(response)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}


	// RETURN DECODED RESPONSE
	return users, nil
}

func ListServerStub(request []byte) []byte {
	users := List()
	var response []byte
	for _, user := range users {
		response = WriteString(response, user)
	}
	return response
}


//CheckMessagesRPC
func CheckMessagesRPC(server string, user string) ([]string, error) {
	// CREATE/FORM REQUEST
	var request []byte
	request = WriteUint16(request, MsgCheckMessages)
	request = WriteString(request, user) // Encode the user

	// SEND REQUEST AND RECEIVE RESPONSE
	response, err := SendAndReceive(server, request)

	// RETURN an error if error != nil 
	if err != nil {
		return nil, err // Return error if communication failed
	}

	// DECODE RESPONSE
	var messages []string
	for len(response) > 0 {
		var msg string
		var err error
		msg, response, err = ReadString(response)
		
		// CHECK if response has the error
		if err != nil {
			log.Printf("CheckMessagesRPC: error decoding response: %v", err)
			return nil, err
		}
		messages = append(messages, msg)
	}

	// RETURN DECODED MESSAGES
	return messages, nil
}

//CheckMessagesServerStub
func CheckMessagesServerStub(request []byte) []byte {
	user, _, err := ReadString(request)
	if err != nil {
		return WriteString(nil, "*** Error reading user")
	}

	messages := CheckMessages(user)
	var response []byte
	for _, msg := range messages {
		response = WriteString(response, msg)
	}

	//encode the error into the response as bytes using WriteString
	if len(messages) > 0 && strings.HasPrefix(messages[0], "*** You are not logged in") {
		response = WriteString(response, "1") // 1 indicates an error??
	} else {
		response = WriteString(response, "") // 0 indicates success??
	}

	return response
}

func TellRPC(server string, user string, targetUser string, message string) error {
	var request []byte
	request = WriteUint16(request, MsgTell)
	request = WriteString(request, user)
	request = WriteString(request, targetUser)
	request = WriteString(request, message)

	response, err := SendAndReceive(server, request)
	if err != nil {
		return err
	}

	msg, _, err := ReadString(response)
	if err != nil {
		return err
	}
	if msg != "" {
		return fmt.Errorf("%s", msg)
	}

	return nil
}

func TellServerStub(request []byte) []byte {
	user, request, err := ReadString(request)
	if err != nil {
		return WriteString(nil, "*** Error reading sender user")
	}
	target, request, err := ReadString(request)
	if err != nil {
		return WriteString(nil, "*** Error reading target user")
	}
	message, _, err := ReadString(request)
	if err != nil {
		return WriteString(nil, "*** Error reading message")
	}

	Tell(user, target, message)
	return nil
}

func SayRPC(server string, user string, message string) error {
	var request []byte

	request = WriteUint16(request, MsgSay)
	request = WriteString(request, user)
	request = WriteString(request, message)

	response, err := SendAndReceive(server, request)
	if err != nil{
		return err
	}

	msg, _, err := ReadString(response)
	if err != nil{
		return err
	}
	if len(msg) != 0 {
		return fmt.Errorf("%s", msg)
	}

	return nil
}

func SayServerStub(request []byte) []byte {
	user, request, err := ReadString(request)
	if err != nil {
		return WriteString(nil, "*** Error reading user")
	}
	message, _, err := ReadString(request)
	if err != nil {
		return WriteString(nil, "*** Error reading message")
	}

	Say(user, message)
	return nil
}



func QuitRPC(server string, user string) error {
	var request []byte
	request = WriteUint16(request, MsgQuit)
	request = WriteString(request, user)

	response, err := SendAndReceive(server, request)
	if err != nil {
		return err
	}

	msg, _, err := ReadString(response)
	if err != nil {
		return err
	}
	if msg != "" {
		return fmt.Errorf("%s", msg)
	}

	return nil
}

func QuitServerStub(request []byte) []byte {
	user, _, err := ReadString(request)
	if err != nil {
		return WriteString(nil, "*** Error reading user")
	}

	Quit(user)
	return nil
}


//ShutdownRPC
func ShutdownRPC(server string) error {
	var request []byte
	request = WriteUint16(request, MsgShutdown)

	response, err := SendAndReceive(server, request)
	if err != nil {
		return err
	}

	//////
	msg, _, err := ReadString(response)
	if err != nil {
		return err
	}
	if msg != "" {
		return fmt.Errorf("%s", msg)
	}

	return nil
}

//ShutdownServerStub
func ShutdownServerStub(request []byte) []byte {
	///
	Shutdown()
	return nil
}