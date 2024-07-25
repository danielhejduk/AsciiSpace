package main

type BLOCK int

const (
    BLOCK_BLANK BLOCK = iota
)

func get_block_string(block BLOCK) string {
    switch block {
        case BLOCK_BLANK:
            return " " // BLANK SPACE
        default:
            return "" // THIS SHOULD NOT HAPPEN
    }
}
