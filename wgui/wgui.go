package wgui

import (
	"fmt"
	"path"

	"Ataxx/ataxx"
	. "Ataxx/utils"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

var (
	turnIndicator *widgets.QPushButton
	final         *widgets.QPushButton
	lineEdit      *widgets.QLineEdit
	centralLayout *widgets.QGridLayout

	assetPath string

	wallIcon    *gui.QIcon
	bluIcon     *gui.QIcon
	redIcon     *gui.QIcon
	emptyIcon   *gui.QIcon
	wshBluIcon  *gui.QIcon
	wshRedIcon  *gui.QIcon
	redHaloIcon *gui.QIcon
	bluHaloIcon *gui.QIcon

	buttons [8][8]Button

	active      int = -1
	activeColor int = -1
)

type Button struct {
	button *widgets.QPushButton
	row    int
	col    int
	state  int
}

func loadAssets() {
	bluIcon = gui.NewQIcon5(path.Join(assetPath, "blu.png"))
	redIcon = gui.NewQIcon5(path.Join(assetPath, "red.png"))
	wallIcon = gui.NewQIcon5(path.Join(assetPath, "wall.png"))
	emptyIcon = gui.NewQIcon5(path.Join(assetPath, "empty.png"))
	wshBluIcon = gui.NewQIcon5(path.Join(assetPath, "wshd_blu.png"))
	wshRedIcon = gui.NewQIcon5(path.Join(assetPath, "wshd_red.png"))
	bluHaloIcon = gui.NewQIcon5(path.Join(assetPath, "blu_halo.png"))
	redHaloIcon = gui.NewQIcon5(path.Join(assetPath, "red_halo.png"))
}

func setupPath(args []string) {
	if len(args) > 1 {
		assetPath = path.Join(args[1], "wgui/assets")
	}
}

//TODO: Put this in DisplayBoard
func updateTurnIcon() {
	turn := ataxx.B().Turn
	var icon *gui.QIcon
	if turn == RED {
		icon = redIcon
	} else {
		icon = bluIcon
	}

	turnIndicator.SetIcon(icon)
}

func setIcon(x, y int, icon *gui.QIcon) {
	buttons[x][y].button.SetIcon(icon)
}

func gameEnd() {
	var str string
	winner := ataxx.Winner()
	if winner == BLU {
		str = "BLU is the winner!"
	} else if winner == RED {
		str = "RED is the winner!"
	} else {
		str = "It's a draw!"
	}

	//Show the 'final' button to reset the game
	final = widgets.NewQPushButton2(str, nil)
	final.SetMinimumSize2(100, 50)
	final.ConnectClicked(reset)
	centralLayout.AddWidget2(final, 1, 0, core.Qt__AlignCenter)
}

func reset(bool) {

	final.Hide()
	ataxx.InitAtaxx()
	DisplayBoard()
}

func updateBoard(bc Button) {
	DisplayBoard()
	b := ataxx.B()
	state := b.B[bc.row][bc.col]
	drawPossibleMoves := false

	if drawPossibleMoves = state == BLU && b.Turn == BLU; drawPossibleMoves {
		setIcon(bc.row, bc.col, bluHaloIcon)
		activeColor = BLU
	} else if drawPossibleMoves = state == RED && b.Turn == RED; drawPossibleMoves {
		setIcon(bc.row, bc.col, redHaloIcon)
		activeColor = RED
	} else if state == NO && active >= 0 {

		//A move has been made
		i, j := PosToCoords(active)
		if DistInf(bc.row, bc.col, i, j) <= 2 {

			ataxx.MakeMove(active, CoordsToPos(bc.row, bc.col), activeColor)
			fen := ataxx.GenFen()
			Assert(ataxx.IsValidFen(fen), "The updated fen isn't valid")
			lineEdit.SetText(fen)
			DisplayBoard()
			if ataxx.Finished() {
				gameEnd()
			}
		}
	}

	if drawPossibleMoves {

		active = CoordsToPos(bc.row, bc.col)
		sqrs := ataxx.PosMoves(CoordsToPos(bc.row, bc.col))
		for _, sqr := range sqrs {
			i, j := PosToCoords(sqr)
			if state == BLU {
				setIcon(i, j, wshBluIcon)
			} else {
				setIcon(i, j, wshRedIcon)
			}
		}
	} else {
		active = -1
		activeColor = -1
	}
}

func getIcon(state int) *gui.QIcon {
	//This could be made with an array and using an offset, like utils.ToByte
	switch state {
	case NO:
		return emptyIcon
	case BLU:
		return bluIcon
	case RED:
		return redIcon
	case WALL:
		return wallIcon
	default:
		return emptyIcon
	}
}

func DisplayBoard() {
	b := ataxx.B()
	for i := 0; i < SIZE; i++ {
		for j := 0; j < SIZE; j++ {
			setIcon(i, j, getIcon(b.B[i][j]))
		}
	}

	updateTurnIcon()
}

func setupGrid() {

	//Grid Layout
	gridWidget := widgets.NewQWidget(nil, 0)

	grid := widgets.NewQGridLayout(gridWidget)
	for i := 0; i < SIZE; i++ {
		for j := 0; j < SIZE; j++ {
			grid.AddWidget2(buttons[i][j].button, i, j, core.Qt__AlignCenter)
		}
	}
	addWidget(gridWidget)
}

func initButtons() {
	for i := 0; i < SIZE; i++ {
		for j := 0; j < SIZE; j++ {
			button := widgets.NewQPushButton3(emptyIcon, "", nil)
			button.SetIconSize(core.NewQSize2(60, 60))
			button.SetMinimumSize2(60, 60)

			bc := Button{button: button, row: i, col: j, state: NO}
			button.ConnectClicked(func(bool) {
				updateBoard(bc)
			})
			buttons[i][j] = bc
		}
	}
}

func Start(num int, s []string) {

	setupPath(s)
	widgets.NewQApplication(num, s)

	loadAssets()

	mainWindow := widgets.NewQMainWindow(nil, 0)

	//Initialize and populate the buttons based on the board
	initButtons()

	scrollWidget := widgets.NewQScrollArea(nil)
	centralWidget := widgets.NewQWidget(nil, 0)
	centralLayout = widgets.NewQGridLayout(centralWidget)

	//Load the array of buttons
	setupGrid()

	//Clocks label
	label1 := widgets.NewQLabel2("12:60", nil, 0)
	label2 := widgets.NewQLabel2("31:21", nil, 0)
	label1.SetFont(gui.NewQFont2("consolas", 21, 1, false))
	label2.SetFont(gui.NewQFont2("consolas", 21, 1, false))
	centralLayout.AddWidget2(label1, 0, 0, core.Qt__AlignLeft)
	centralLayout.AddWidget2(label2, 0, 0, core.Qt__AlignRight)

	//Add the FEN input box
	lineEdit = widgets.NewQLineEdit(nil)
	lineEdit.SetText("FEN")
	lineEdit.SetMinimumSize2(200, 30)
	lineEdit.ConnectTextEdited(func(fen string) {

		fmt.Println(fen)
		if fen == "default" {
			ataxx.InitAtaxx()
			DisplayBoard()
		} else if ataxx.IsValidFen(fen) {
			ataxx.ParseFen(fen)
			DisplayBoard()
		}
	})
	centralLayout.AddWidget2(lineEdit, 2, 0, core.Qt__AlignCenter)

	//Turn indicator button
	turnIndicator = widgets.NewQPushButton3(wallIcon, "", nil)
	turnIndicator.SetMinimumSize2(40, 40)
	turnIndicator.SetIconSize(core.NewQSize2(40, 40))
	centralLayout.AddWidget2(turnIndicator, 2, 0, core.Qt__AlignLeft)

	DisplayBoard()

	//Setup the main window
	scrollWidget.SetWidget(centralWidget)
	mainWindow.SetCentralWidget(scrollWidget)
	mainWindow.SetMinimumSize2(625, 720)
	mainWindow.SetWindowTitle("Attax")
	mainWindow.Show()
	//mainWindow.ShowMaximized()

	widgets.QApplication_Exec()
}

func addWidget(widget widgets.QWidget_ITF) {

	wrappedWidget := widgets.NewQGroupBox2(widget.QWidget_PTR().WindowTitle(), nil)
	wrappedWidgetLayout := widgets.NewQVBoxLayout2(wrappedWidget)
	wrappedWidgetLayout.AddWidget(widget, 0, core.Qt__AlignCenter)
	wrappedWidget.SetFixedSize2(600, 600)

	centralLayout.AddWidget2(wrappedWidget, 1, 0, core.Qt__AlignCenter)
}
