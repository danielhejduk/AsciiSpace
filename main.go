package main

import (
	"bufio"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

const CLEAR_BLUE = "\033[104m\033[2J\033[H" // Makes blue background, clears the screen and moves cursor to top left
const CLEAR_BLACK = "\033[40m\033[2J\033[H"

const MENU_DASHES = "----------------------------------------" // 40 Dashes
const MENU_SPACES = "                                        " // 40 Spaces

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

var game Game
var gamejolt Gamejolt

func get_gamejolt_credentials() (string, int, bool) {
    err := godotenv.Load(".env")
    if err != nil{
        println("Error loading .env file:", err)
        return "", 0, false
    }

    gamejolt_priv := os.Getenv("GAMEJOLT_KEY")
    idstr := os.Getenv("GAMEJOLT_ID")

    gamejolt_gameid, err := strconv.Atoi(idstr)
    if err != nil {
        println("[ERROR] Something has happend:", err)
        return "", 0, false
    }

    return gamejolt_priv, gamejolt_gameid, true
}


func main() {
    println("[INFO] AsciiSpace server is setting up")

    ln, err := net.Listen("tcp", ":8080")
    if err != nil {
        println("[ERROR]", err)
        return
    }

    priv_key, gameid, success := get_gamejolt_credentials()

    if !success {
        println("[ERROR] Something has happend getting Gamejolt credentials")
        return
    }

    gamejolt.privatekey = priv_key
    gamejolt.gameid = gameid

    println("[INFO] AsciiSpace server is running")

    for {
        conn, err := ln.Accept()
        if err != nil {
            println("[ERROR]", err)
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

    username, password := LoginMenu(client)

    if username == "" || password == "" {
        return
    }

    success := gamejolt.LoginPlayer(username, password)

    if !success {
        println("[INFO]", client.ip + "->", "User entered wrong gamejolt credintals")
        return
    }

    gamejolt.AddTrophy(username, password, TROPHY_LOGIN)

    HandleGame(client, player)
}

func HandleGame(client Client, player Player) {
    reader := bufio.NewReader(client.conn)

    var terrain TERRAIN

    terrain.generate_map()

    for {
        terrain.render_map(client)
        terrain.collision_handling(client, player)
        terrain.print_player(client, player)

        key, err := reader.ReadByte()
        
        if err != nil {
            println("[ERROR]", client.ip, "->", "Something has happend during reading the key")
            return
        }
        
        player.key = key
        quit_signal := player.handle_controls(terrain)
        if quit_signal {
            return
        }
    }
}

func LoginMenu(client Client) (string, string) {
    var x int = 1
    reader := bufio.NewReader(client.conn)

    tusername := make([]byte, 15)
    tpassword := make([]byte, 15)
    telnetans := make([]byte, 3)

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
        return "", ""
    }
    username := strings.ToLower(strings.TrimSpace(string(tusername)))

    client.conn.Write([]byte("\033[3;2HGame Token:"))
    tpassword, _, err = reader.ReadLine()
    if err != nil {
        println("[INFO]", client.ip, "->", "Something has happend during reading password")
        return "", ""
    }
    password := strings.ToLower(strings.TrimSpace(string(tpassword)))

    client.conn.Write([]byte{IAC, WILL, SUPGO, IAC, WONT, ECHO})

    _, err = client.conn.Read(telnetans) // We are reading the telnet answer to the line buffering

    client.conn.Write([]byte("\033[4;2HPress any key to continue"))
    _, err = reader.ReadByte()
    if err != nil {
        println("[INFO]", client.ip, "->", "Something has happend during reading the answer")
    }
    
    return string(username), string(password)
}
