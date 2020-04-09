package utils

const (
	WALL int = -2
	NO   int = -1
	RED  int = 0
	BLU  int = 1

	SIZE int = 7

	WIN int = 9
	WIN_ON_TIME int = 10
	WIN_ON_RESIGNATION int = 11
)

var Bytes = [4]byte{'-', 'M', 'o', 'x'} //You shouldn't get M
//The program can be improved using bitboards
type Board struct {
	B    [SIZE][SIZE]int //0 -> RED pieces, 1 -> BLU pieces
	Turn int
}

func Max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func Abs(a int) int {
	if a < 0 {
		return -a
	} else {
		return a
	}
}

func DistInf(aX, aY, bX, bY int) int {
	return Max(Abs(aX-bX), Abs(aY-bY))
}

func CoordsToPos(i, j int) int {
	return SIZE*i + j
}

func PosToCoords(pos int) (int, int) {
	return pos / SIZE, pos % SIZE
}

func ValidXY(x, y int) bool {
	return 0 <= x && x < SIZE && 0 <= y && y < SIZE
}

func ToByte(i int) byte {
	return Bytes[2+i]
	/*
		switch i{
		case BLU:
			return 'x'
		case RED:
			return 'o'
		case WALL:
			return '-'
		default:
			return 0
		}
	*/
}

func Assert(b bool, msg string) {
	if !b {
		panic(msg)
	}
}
