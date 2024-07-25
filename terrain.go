package main

type BLOCK int

const (
    BLOCK_BLANK BLOCK = iota
    BLOCK_PLAYER
)

func get_block_string(block BLOCK) string {
    switch block {
        case BLOCK_BLANK:
            return " " // BLANK SPACE
        case BLOCK_PLAYER:
            return "@"
        default:
            return "" // THIS SHOULD NOT HAPPEN
    }
}
