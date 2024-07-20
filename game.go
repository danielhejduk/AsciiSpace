package main

import (
    "fmt"
)

type Game struct {
    test int
}

func (game *Game) printPlayer(client Client, player Player) {
    var buf string

    buf = fmt.Sprintf("\033[%d;%dH@", player.y+1, player.x+1)

    client.conn.Write([]byte(buf))
}
