package wgui

import (
	"fmt"
	"os"
	"path"
	"time"
	"strconv"
	"strings"

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

	stopTimer 	chan bool

	gd GameData
)


type Params struct {
	Path string
	Tc TimeControl
}

type Button struct {
	button *widgets.QPushButton
	row    int
	col    int
	state  int
}

type TimeControl struct {
	initialTime time.Duration
	increment time.Duration
	infinite bool
}

type GameData struct {
	tc TimeControl
	remTime [2]time.Duration
	turnChangeTime time.Time
	turn int
}

func NewGameData(tc TimeControl) GameData {
	remTime := [2]time.Duration{tc.initialTime, tc.initialTime}
	return GameData {tc, remTime, time.Now(), BLU}
}

func ParseTC(s string) (tc TimeControl) {

	if strings.Compare(s, "inf") == 0 {
		tc = TimeControl{0, 0, true}
	} else {
		initialTime := 0
		increment := 0
		if strings.Contains(s, "+") {
			res := strings.Split(s, "+")
			initialTime, _ = strconv.Atoi(res[0])
			increment, _ = strconv.Atoi(res[1])
		} else {
			initialTime, _ = strconv.Atoi(s)
		}

		tc = TimeControl{time.Duration(initialTime) * time.Second, time.Duration(increment) * time.Second, false}
	}

	return
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


func getMS(t float64) (int, int) {
	return int(t) / 60, int(t) % 60
}

func incrementTime(side int) {
	if gd.remTime[side] > 0 {
		gd.remTime[side] += gd.tc.increment
	}
}

func LaunchTimer() {
	gd.turnChangeTime = time.Now()

	stopTimer = make(chan bool)

	setTimers()
	updateTimers()
}

func setTimers() {
	mB, sB := getMS(gd.remTime[BLU].Seconds())
	mR, sR := getMS(gd.remTime[BLU].Seconds())
	labelBLU.SetText(fmt.Sprintf("BLU: %d:%02d", mB, sB))
	labelRED.SetText(fmt.Sprintf("RED: %d:%02d", mR, sR))
}


func updateTimers() {

	var timer *time.Timer
	var timeElaped time.Duration

	if gd.tc.infinite {
		labelRED.SetText("RED: -:--")
		labelBLU.SetText("BLU: -:--")
		return
	}
	for gd.remTime[BLU] > 0 && gd.remTime[RED] > 0 {
		timeElaped = time.Now().Sub(gd.turnChangeTime)

		if gd.turn == RED {
			gd.remTime[RED] = gd.remTime[RED] - timeElaped
		} else {
			gd.remTime[BLU] = gd.remTime[BLU] - timeElaped
		}

		mins, secs := getMS(gd.remTime[RED].Seconds())
		labelRED.SetText(fmt.Sprintf("RED: %d:%02d", mins, secs))

		mins, secs = getMS(gd.remTime[BLU].Seconds())
		labelBLU.SetText(fmt.Sprintf("BLU: %d:%02d", mins, secs))

		gd.turnChangeTime = time.Now()
		timer = time.NewTimer(time.Second / 4)

		for channelAct := false; !channelAct; {
			select {
			case <-stopTimer:
				return
			case <-timer.C:
				channelAct = true
			}
		}
	}

	if gd.remTime[BLU] <= 0 {
		gameEndMsg(RED, WIN_ON_TIME)
	} else if gd.remTime[RED] <= 0 {
		gameEndMsg(BLU, WIN_ON_TIME)
	}
}

func updateTurnIcon() {
	newTurn := ataxx.B().Turn
	var icon *gui.QIcon
	if newTurn == RED {
		gd.turn = RED
		icon = redIcon
	} else {
		gd.turn = BLU
		icon = bluIcon
	}

	turnIndicator.SetIcon(icon)
}

func setIcon(x, y int, icon *gui.QIcon) {
	buttons[x][y].button.SetIcon(icon)
}

func gameEndMsg(winner int, condition int) {

	var str string
	if winner == BLU {
		if condition == WIN_ON_TIME{
			str = "BLU wins on time"
		} else {
			str = "BLU wins"
		}
	} else if winner == RED {
		if condition == WIN_ON_TIME{
			str = "RED wins on time"
		} else {
			str = "RED wins"
		}
	} else {
		str = "It's a draw!"
	}

	//Show the 'final' button to reset the game
	final.SetVisible(true)
	final.SetText(str)
}

//This generates the function that is called when asked to reset the board
func reset(bool) {
	gd = NewGameData(gd.tc)

	final.Hide()
	ataxx.InitAtaxx()
	DisplayBoard()

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

			incrementTime(activeColor)

			ataxx.MakeMove(active, CoordsToPos(bc.row, bc.col), activeColor)
			fen := ataxx.GenFen()
			Assert(ataxx.IsValidFen(fen), "The updated fen isn't valid")
			lineEdit.SetText(fen)
			DisplayBoard()
			gd.turnChangeTime = time.Now()
			if ataxx.Finished() {
				//If tc is infinite there is nobody listening
				if !gd.tc.infinite {
					//We do this to stop the timers
					stopTimer <- true
					close(stopTimer)
				}
				gameEndMsg(ataxx.Winner(), WIN)
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
	case BLU:
		return bluIcon
	case RED:
		return redIcon
	case WALL:
		return wallIcon
	case NO:
		fallthrough //In case I change this in the future
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

func Start(num int, s []string, prms Params) {

	assetPath = path.Join(prms.Path, "wgui", "assets")
	_, err := os.Stat(assetPath)
	if os.IsNotExist(err) {
		fmt.Println(fmt.Sprintf("%s is not a valid directory, pass the argument -h for help", assetPath))
		return
	}

	widgets.NewQApplication(num, s)

	loadAssets()

	mainWindow := widgets.NewQMainWindow(nil, 0)

	//Initialize and populate the buttons based on the board
	gd = NewGameData(prms.Tc)

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
