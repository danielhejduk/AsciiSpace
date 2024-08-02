package main

type Player struct { // x,y starts with 0,0 but ANSI starts with 1,1
	x int
	y int

	key byte

	health int
}

func (player *Player) handle_controls() bool /*Quit Signal*/ {
	// MOVEMENT CONTROLS
	if player.key == 'w' {
		if player.y == 0 {
			return false
		}
		player.y--
	} else if player.key == 'a' {
		if player.x == 0 {
			return false
		}
		player.x--
	} else if player.key == 's' {
		player.y++
	} else if player.key == 'd' {
		player.x++
	} else if player.key == 'q' { // QUIT COMMAND
		return true
	}

	return false
}
