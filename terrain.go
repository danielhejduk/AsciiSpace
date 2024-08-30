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
    BLOCK_KILL
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
	case BLOCK_KILL: // KILL TESTING BLOCK
            return "\033[31m&\033[0m"
        default:
            return "" // THIS SHOULD NOT HAPPEN
    }
}

func (terrain *TERRAIN) generate_map() {
    terrain.planetmap[4][4] = BLOCK_ROCK
    terrain.planetmap[5][5] = BLOCK_KILL
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

func (terrain *TERRAIN) is_solid(x int, y int) bool {
	switch (terrain.planetmap[x][y]) {
		case BLOCK_ROCK:
			return true
	}
	return false
}

func (terrain *TERRAIN) collision_handling(client Client, player Player) {
    if terrain.planetmap[player.x][player.y] == BLOCK_KILL {
        client.conn.Close()
    }
}
