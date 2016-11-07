package main

import (
		"net"
		"net/textproto"
		"strings"
		"bufio"
		"encoding/json"
		"os"
		"fmt"
)

func main() {

	// parse the config file
	type Configuration struct {
		Username string
		Password string
		Channel string
	}

	file, _ := os.Open("conf.json")
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error", err)
	}
	fmt.Println(configuration.Username)

	// connect to the twitch server
	conn, err := net.Dial("tcp", "irc.chat.twitch.tv:6667")
	if err != nil {
		panic(err)
	}

	// token, username, channel
	conn.Write([]byte("PASS " + configuration.Password + "\r\n"))
	conn.Write([]byte("NICK " + configuration.Username + "\r\n"))
	conn.Write([]byte("JOIN " + configuration.Channel + "\r\n"))
	defer conn.Close()

	// handles reading from the connection
	tp := textproto.NewReader(bufio.NewReader(conn))

	// listens/responds to chat messages
	for {
		msg, err := tp.ReadLine()
		if err != nil {
			panic(err)
		}

		// split the msg by spaces
		msgParts := strings.Split(msg, " ")

		// if the msg contains PING you're required to
		// respond with PONG else you get kicked
		if msgParts[0] == "PING" {
			conn.Write([]byte("PONG " + msgParts[1]))
			continue
		}

		// if msg contains PRIVMSG then respond
		if msgParts[1] == "PRIVMSG" {
			// echo back the same message
			conn.Write([]byte("PRIVMSG " + msgParts[2] + " " + msgParts[3] + "\r\n"))
		}
	}
}