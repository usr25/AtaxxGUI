package ataxx

import (
	. "Ataxx/utils"
)

var board Board

//This is used to return a duplicate to avoid messing up the internal board
func B() Board {
	return board
}

func SetCoord(x, y, state int) {
	board.B[x][y] = state
}
func ChangeTurn() {
	board.Turn = 1 - board.Turn
	if !playerCanMove() {
		board.Turn = 1 - board.Turn
	}
}

func distInfPos(a, b int) int {
	aX, aY := PosToCoords(a)
	bX, bY := PosToCoords(b)
	return DistInf(aX, aY, bX, bY)
}

//Checks if any piece can move, this could be better if the board uses bitboards
func playerCanMove() bool {

	for i := 0; i < SIZE; i++ {
		for j := 0; j < SIZE; j++ {
			if board.B[i][j] == board.Turn {
				pm := PosMoves(CoordsToPos(i, j))
				if len(pm) > 0 {
					return true
				}
			}
		}
	}

	return false
}

//Pos moves for a given position, it is color independent
func PosMoves(pos int) []int {
	x, y := PosToCoords(pos)
	var res [25]int
	counter := 0

	//TODO: This can be improved
	for i := -2; i < 3; i++ {
		for j := -2; j < 3; j++ {
			nX := x + i
			nY := y + j
			if ValidXY(nX, nY) && board.B[nX][nY] == NO {
				res[counter] = CoordsToPos(nX, nY)
				counter += 1
			}
		}
	}

	return res[:counter]
}

//Determines who is the winner, NO is for draw
func Winner() int {
	blu := 0
	red := 0
	for i := 0; i < SIZE; i++ {
		for j := 0; j < SIZE; j++ {
			if board.B[i][j] == BLU {
				blu++
			} else if board.B[i][j] == RED {
				red++
			}
		}
	}

	if red > blu {
		return RED
	} else if red < blu {
		return BLU
	} else {
		return NO
	}
}

//Determines if the match has been finished
func Finished() bool {
	red := false
	blu := false
	empty := false
	for i := 0; i < SIZE; i++ {
		for j := 0; j < SIZE; j++ {
			if board.B[i][j] == NO {
				empty = true
			} else if board.B[i][j] == RED {
				red = true
			} else if board.B[i][j] == BLU {
				blu = true
			}

			if empty && red && blu {
				return false
			}
		}
	}

	return true
}

//Returns all the enemies within inf-dist == 1
func closeEnemies(pos, col int) []int {
	x, y := PosToCoords(pos)
	var res [8]int
	counter := 0
	opp := 1 ^ col

	for i := -1; i < 2; i++ {
		for j := -1; j < 2; j++ {
			if i == 0 && j == 0 {
				continue
			}

			nX := x + i
			nY := y + j
			if ValidXY(nX, nY) && board.B[nX][nY] == opp {
				res[counter] = CoordsToPos(nX, nY)
				counter += 1
			}
		}
	}

	return res[:counter]
}

func MakeMove(from, to, col int) {
	SetTile(to, col)
	if distInfPos(from, to) == 2 {
		SetTile(from, NO)
	}
	enemies := closeEnemies(to, col)
	for _, p := range enemies {
		SetTile(p, col)
	}
	ChangeTurn()
}

func SetTile(pos, col int) {
	x, y := PosToCoords(pos)
	board.B[x][y] = col
}

func cleanBoard() {
	board.Turn = BLU
	for i := 0; i < SIZE; i++ {
		for j := 0; j < SIZE; j++ {
			board.B[i][j] = NO
		}
	}
}

func initBoard() {
	cleanBoard()
	board.B[0][0] = RED
	board.B[SIZE-1][SIZE-1] = RED
	board.B[SIZE-1][0] = BLU
	board.B[0][SIZE-1] = BLU
}
