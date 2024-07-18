package main

import (
	"bufio"
	"crypto/sha512"
	"encoding/base64"
	"net"
	"strings"
	"time"
)

const CLEAR_BLUE = "\033[104m\033[2J\033[H" // Makes blue background, clears the screen and moves cursor to top left
const CLEAR_BLACK = "\033[40m\033[2J\033[H"

const MENU_DASHES = "------------------------" // 24 Dashes
const MENU_SPACES = "                        " // 24 Spaces

const (
	IAC   = 255 // Interpret as Command
	DONT  = 254
	DO    = 253
	WONT  = 252
	WILL  = 251
	SB    = 250
	SE    = 240
	ECHO  = 1
	SUPGO = 3
)

type Client struct {
    conn net.Conn
    ip string
}

type Player struct {
    x int
    y int

    key byte
    answer byte

    health int
}

func main() {
    println("[INFO] AsciiSpace server is setting up")

    ln, err := net.Listen("tcp", ":8080")
    if err != nil {
        println(err)
        return
    }

    println("[INFO] AsciiSpace server is running")

    for {
        conn, err := ln.Accept()
        if err != nil {
            println(err)
            continue
        }

        var client Client

        client.conn = conn
        client.ip = client.conn.RemoteAddr().String()

        // Handle the connection in a new goroutine
        go handleConnection(client)
    }
}

func handleConnection(client Client) {
    defer client.conn.Close()

    var player Player

    client.conn.Write([]byte(CLEAR_BLUE))

    username, password, answer := LoginMenu(client)

    if username == "" || password == "" {
        return
    }

    player.answer = byte(answer)

    hashpassword := sha512.Sum512([]byte(password))

    println(base64.StdEncoding.EncodeToString(hashpassword[:]))

    HandleGame(client, player)
}

func HandleGame(client Client, player Player) {
    reader := bufio.NewReader(client.conn)

    client.conn.Write([]byte(CLEAR_BLACK))

    for {
        key, err := reader.ReadByte()

        player.key = key

        if player.key == 'q' {
            return
        }

        if err != nil {
            println("[INFO]", client.ip, "->", "Something has happend during reading the key")
            return
        }

        println(string(player.key))
    }
}

func LoginMenu(client Client) (string, string, byte) {
    var x int = 1
    reader := bufio.NewReader(client.conn)

    tusername := make([]byte, 15)
    tpassword := make([]byte, 15)
    telnetans := make([]byte, 3)
    var answer byte

    client.conn.Write([]byte("/" + MENU_DASHES + "\\\n"))
    for x < 6 { // 5 Steps
        client.conn.Write([]byte("|" + MENU_SPACES + "|\n"))
        x++
    }
    client.conn.Write([]byte("\\" + MENU_DASHES + "/"))

    client.conn.Write([]byte("\033[2;2HUsername:"))
    tusername, _, err := reader.ReadLine()
    if err != nil {
        println("[INFO]", client.ip, "->", "Something has happend during reading username")
        return "", "", byte(0)
    }
    username := strings.ToLower(strings.TrimSpace(string(tusername)))

    client.conn.Write([]byte("\033[3;2HPassword:"))
    tpassword, _, err = reader.ReadLine()
    if err != nil {
        println("[INFO]", client.ip, "->", "Something has happend during reading password")
        return "", "", byte(0)
    }
    password := strings.ToLower(strings.TrimSpace(string(tpassword)))
    time.Sleep(time.Second * 2)

    client.conn.Write([]byte{IAC, WILL, SUPGO, IAC, WONT, ECHO})

    _, err = client.conn.Read(telnetans) // We are reading the telnet answer to the line buffering

    client.conn.Write([]byte("\033[4;2HAnswer?"))
    answer, err = reader.ReadByte()
    if err != nil {
        println("[INFO]", client.ip, "->", "Something has happend during reading the answer")
    }
    
    return string(username), string(password), answer
}
