package main

import (
	"bufio"
	"crypto/sha512"
	"encoding/base64"
	"net"
	"os"
	"strings"
    "strconv"

    "github.com/joho/godotenv"
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

type Player struct { // x,y starts with 0,0 but ANSI starts with 1,1
    x int
    y int

    key byte
    answer byte

    health int
    planetmap [100][100]byte
}

var game Game
var gamejolt Gamejolt

func get_gamejolt_credentials() (string, int) {
    err := godotenv.Load(".env")
    if err != nil{
        println("Error loading .env file:", err)
    }

    gamejolt_priv := os.Getenv("GAMEJOLT_KEY")
    idstr := os.Getenv("GAMEJOLT_ID")

    gamejolt_gameid, err := strconv.Atoi(idstr)
    if err != nil {
        println("[ERROR] Something has happend:", err)
    }

    return gamejolt_priv, gamejolt_gameid
}


func main() {
    println("[INFO] AsciiSpace server is setting up")

    ln, err := net.Listen("tcp", ":8080")
    if err != nil {
        println("[ERROR]", err)
        return
    }

    println("[INFO] AsciiSpace server is running")

    for {
        conn, err := ln.Accept()
        if err != nil {
            println("[ERROR]", err)
            continue
        }

        var client Client

        priv_key, gameid := get_gamejolt_credentials()
        println(gameid)
        gamejolt.privatekey = priv_key
        gamejolt.gameid = gameid

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

    gamejolt.LoginPlayer(username, password)
    gamejolt.AddTrophy(username, password, 239682) // TODO: Add modifiable trophy system

    player.answer = byte(answer)

    hashpassword := sha512.Sum512([]byte(password))

    println(base64.StdEncoding.EncodeToString(hashpassword[:]))

    HandleGame(client, player)
}

func HandleGame(client Client, player Player) {
    reader := bufio.NewReader(client.conn)

    for {
        client.conn.Write([]byte(CLEAR_BLACK))

        game.printPlayer(client, player)

        key, err := reader.ReadByte()

        if err != nil {
            println("[ERROR]", client.ip, "->", "Something has happend during reading the key")
            return
        }
        
        player.key = key

        println("[INFO]", client.ip, "->", string(player.key))

        if player.key == 'q' {
            return
        } else if player.key == 'd' {
            player.x++
        } else if player.key == 's' {
            player.y++
        } else if player.key == 'a' {
            if player.x == 0 {
                continue
            } else {
                player.x--
            }
        } else if player.key == 'w' {
            if player.y == 0 {
                continue
            } else {
                player.y--
            }
        }

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

    client.conn.Write([]byte{IAC, WILL, SUPGO, IAC, WONT, ECHO})

    _, err = client.conn.Read(telnetans) // We are reading the telnet answer to the line buffering

    client.conn.Write([]byte("\033[4;2HIs this Gamejolt Login?"))
    answer, err = reader.ReadByte()
    if err != nil {
        println("[INFO]", client.ip, "->", "Something has happend during reading the answer")
    }
    
    return string(username), string(password), answer
}
