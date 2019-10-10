package wgui

import (
	"fmt"
	"path"
	"time"

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

	labelRED	*widgets.QLabel
	labelBLU	*widgets.QLabel

	buttons [8][8]Button

	active      int = -1
	activeColor int = -1

	remTimeBLU 	time.Duration
	remTimeRED 	time.Duration
	turnChangeTime time.Time

	turn 		int
	stopTimer 	chan bool
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

func getMS(t float64) (int, int) {
	return int(t) / 60, int(t) % 60
}

func LaunchTimer() {
	turnChangeTime = time.Now()
	remTimeRED = 3 * 60 * time.Second
	remTimeBLU = 3 * 60 * time.Second

	stopTimer = make(chan bool, 2)

	setTimers()
	updateTimers()
}

func setTimers() {
	mB, sB := getMS(remTimeBLU.Seconds())
	mR, sR := getMS(remTimeRED.Seconds())
	labelBLU.SetText(fmt.Sprintf("BLU: %d:%02d", mB, sB))
	labelRED.SetText(fmt.Sprintf("RED: %d:%02d", mR, sR))
}


func updateTimers() {

	var timer *time.Timer
	var timeElaped time.Duration

	for remTimeBLU > 0 && remTimeRED > 0 {
		timeElaped = time.Now().Sub(turnChangeTime)

		if turn == RED {
			remTimeRED = remTimeRED - timeElaped
			mins, secs := getMS(remTimeRED.Seconds())
			labelRED.SetText(fmt.Sprintf("RED: %d:%02d", mins, secs))
		} else {
			remTimeBLU = remTimeBLU - timeElaped
			mins, secs := getMS(remTimeBLU.Seconds())
			labelBLU.SetText(fmt.Sprintf("BLU: %d:%02d", mins, secs))
		}
		turnChangeTime = time.Now()
		timer = time.NewTimer(3 * time.Second / 8)

		for channelAct := false; !channelAct; {
			select {
			case <-stopTimer:
				return
			case <-timer.C:
				channelAct = true
				//break //This doesnt work, go is just retarded
			}
		}
	}

	if remTimeBLU <= 0 {
		gameEnd(RED)
	} else if remTimeRED <= 0 {
		gameEnd(BLU)
	}
}

func updateTurnIcon() {
	newTurn := ataxx.B().Turn
	var icon *gui.QIcon
	if newTurn == RED {
		turn = RED
		icon = redIcon
	} else {
		turn = BLU
		icon = bluIcon
	}

	turnIndicator.SetIcon(icon)
}

func setIcon(x, y int, icon *gui.QIcon) {
	buttons[x][y].button.SetIcon(icon)
}

func gameEnd(winner int) {
	var str string
	if winner == BLU {
		str = "BLU is the winner!"
	} else if winner == RED {
		str = "RED is the winner!"
	} else {
		str = "It's a draw!"
	}

	//Show the 'final' button to reset the game
	final.SetVisible(true)
	final.SetText(str)
}

func reset(bool) {
	stopTimer <- true //To avoid the blocking of the channel here stopTimer is buffered (cap = 2)
	final.Hide()
	ataxx.InitAtaxx()
	DisplayBoard()
	close(stopTimer)
	go LaunchTimer()
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
			turnChangeTime = time.Now()
			if ataxx.Finished() {
				gameEnd(ataxx.Winner())
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
	labelRED = widgets.NewQLabel2("R--: 04:20", nil, 0)
	labelBLU = widgets.NewQLabel2("B--: 27:18", nil, 0)
	labelRED.SetFont(gui.NewQFont2("consolas", 21, 1, false))
	labelBLU.SetFont(gui.NewQFont2("consolas", 21, 1, false))
	centralLayout.AddWidget2(labelRED, 0, 0, core.Qt__AlignLeft)
	centralLayout.AddWidget2(labelBLU, 0, 0, core.Qt__AlignRight)

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

	final = widgets.NewQPushButton2("", nil)
	final.SetMinimumSize2(100, 50)
	final.ConnectClicked(reset)
	centralLayout.AddWidget2(final, 1, 0, core.Qt__AlignCenter)
	final.Hide()

	DisplayBoard()

	//Setup the main window
	scrollWidget.SetWidget(centralWidget)
	mainWindow.SetCentralWidget(scrollWidget)
	mainWindow.SetMinimumSize2(625, 720)
	mainWindow.SetWindowTitle("Attax")
	mainWindow.Show()
	//mainWindow.ShowMaximized()

	go LaunchTimer()

	widgets.QApplication_Exec()
}

func addWidget(widget widgets.QWidget_ITF) {

	wrappedWidget := widgets.NewQGroupBox2(widget.QWidget_PTR().WindowTitle(), nil)
	wrappedWidgetLayout := widgets.NewQVBoxLayout2(wrappedWidget)
	wrappedWidgetLayout.AddWidget(widget, 0, core.Qt__AlignCenter)
	wrappedWidget.SetFixedSize2(600, 600)

	centralLayout.AddWidget2(wrappedWidget, 1, 0, core.Qt__AlignCenter)
}
