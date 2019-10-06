package ataxx

import (
	"strings"

	. "Ataxx/utils"
)

func InitAtaxx() {
	initBoard()
}

func isValidFenChar(c rune) bool {
	if c == 'o' || c == 'x' || c == '-' || c == ' ' || c == '/' {
		return true
	}
	return int(c-'0') > 0 && int(c-'0') <= SIZE
}

func IsValidFen(fen string) bool {

	if len(fen) < 15 {
		return false
	}
	last := fen[len(fen)-1]
	if (last != 'x' && last != 'o') || fen[len(fen)-2] != ' ' {
		return false
	}
	for _, c := range fen {
		if !isValidFenChar(c) {
			return false
		}
	}

	return true
}

func GenFen() string {

	var sb strings.Builder
	sb.Grow(SIZE * SIZE)

	for i := 0; i < SIZE; i++ {
		counter := 0
		for j := 0; j < SIZE; j++ {
			val := board.B[i][j]

			if val == NO {
				counter++
			} else {
				if counter > 0 {
					sb.WriteByte(byte('0' + counter))
				}
				sb.WriteByte(ToByte(val))
				counter = 0
			}
		}

		if counter > 0 {
			sb.WriteByte(byte('0' + counter))
		}
		sb.WriteByte('/')
	}

	sb.WriteByte(' ')
	sb.WriteByte(ToByte(board.Turn))
	return sb.String()
}

func ParseFen(fen string) {

	Assert(len(fen) > 14, "The fen isn't long enough")
	currPos := 0
	cleanBoard()
	for _, c := range fen {
		x, y := PosToCoords(currPos)
		if c == 'o' {
			board.B[x][y] = RED
			currPos += 1
		} else if c == 'x' {
			board.B[x][y] = BLU
			currPos += 1
		} else if c == '-' {
			board.B[x][y] = WALL
			currPos += 1
		} else if c == ' ' {
			break
		} else if c != '/' && int(c-'0') < 10 {
			currPos += int(c - '0')
		}
	}

	if fen[len(fen)-1] == 'x' {
		board.Turn = BLU
	} else if fen[len(fen)-1] == 'o' {
		board.Turn = RED
	}
}

func BoardToString() string {

	var sb strings.Builder
	sb.Grow(2*SIZE*SIZE + SIZE)
	for i := 0; i < SIZE; i++ {
		for j := 0; j < SIZE; j++ {
			switch board.B[i][j] {
			case WALL:
				sb.WriteByte('#')
			case NO:
				sb.WriteByte('-')
			case BLU:
				sb.WriteByte('x')
			case RED:
				sb.WriteByte('o')
			}
			sb.WriteByte(' ')
		}
		sb.WriteByte('\n')
	}

	return sb.String()
}
