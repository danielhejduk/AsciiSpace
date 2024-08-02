package main

import (
    "fmt"
)

type BLOCK int

type TERRAIN struct {
    planetmap [20][20]BLOCK
}

const (
    BLOCK_BLANK BLOCK = iota
    BLOCK_PLAYER
    BLOCK_ROCK
)

const (
    MOVE_TO_LEFT = "\033[H"
)

func get_block_string(block BLOCK) string {
    switch block {
        case BLOCK_BLANK:
            return " " // BLANK SPACE
        case BLOCK_PLAYER:
            return "@"
        case BLOCK_ROCK: // ROCK BLOCK (SOLID BLOCK)
            return "="
        default:
            return "" // THIS SHOULD NOT HAPPEN
    }
}

func (terrain *TERRAIN) generate_map() {
    terrain.planetmap[4][4] = BLOCK_ROCK
}

func (terrain *TERRAIN) render_map(client Client) {
    var x,y int = 0,0

    client.conn.Write([]byte(CLEAR_BLACK))

    for y < 20 {
        x = 0
        for x < 20 {
            var block string = get_block_string(terrain.planetmap[x][y])
            client.conn.Write([]byte(block))
            x++
        }
        client.conn.Write([]byte("\n"))
        y++
    }
}

func (terrain *TERRAIN) print_player(client Client, player Player) {
    var buf string

    buf = fmt.Sprintf("\033[%d;%dH@", player.y+1, player.x+1)

    client.conn.Write([]byte(buf))
    client.conn.Write([]byte(MOVE_TO_LEFT))    
}
