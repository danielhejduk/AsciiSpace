package main

type Player struct { // x,y starts with 0,0 but ANSI starts with 1,1
	x int
	y int

	key byte

	health int
}

func (player *Player) handle_controls(terrain TERRAIN) bool /*Quit Signal*/ {
	// MOVEMENT CONTROLS
	if player.key == 'w' {
		if player.y == 0 || terrain.is_solid(player.x, player.y-1) {
			return false
		}
		player.y--
	} else if player.key == 'a' {
		if player.x == 0 || terrain.is_solid(player.x-1, player.y) {
			return false
		}
		player.x--
	} else if player.key == 's' {
		if terrain.is_solid(player.x, player.y+1) {
			return false
		}
		player.y++
	} else if player.key == 'd' {
		if terrain.is_solid(player.x+1, player.y) {
			return false
		}
		player.x++
	} else if player.key == 'q' { // QUIT COMMAND
		return true
	}
	return false
}
