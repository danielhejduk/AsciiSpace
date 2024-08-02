package main

import (
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
    var x,y int = 1,1

    client.conn.Write([]byte(CLEAR_BLACK))

    for y < 21 {
        for x < 21 {
            var block string = get_block_string(terrain.planetmap[x-1][y-1])
            client.conn.Write([]byte(block))
            x++
        }
        client.conn.Write([]byte("\n"))
        y++
    }
}
