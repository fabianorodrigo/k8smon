package ui

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
	"github.com/k8smon/models"
)

var containerLogs []models.ContainerLog
var nos []models.Node
var noSelecionado models.Node
var pods []models.Pod
var podSelecionado models.Pod

const (
	//Nós
	nodesTitulo = "Nodes"
	nodesRight  = 40
	//Pods
	podsTitulo = "PODs"
	podsRight  = 80
	//Logs
	logsTitulo = "Logs"
)

func AdicionaContainerLog(containerLog models.ContainerLog) {
	if models.NodeIsIn(nos, containerLog.PodNode) == false {
		no := models.Node{Nome: containerLog.PodNode, QtContainers: 1, QtPods: 1}
		nos = append(nos, no)
	}
	if models.PodIsIn(pods, containerLog.PodName) == false {
		pods = append(pods, models.Pod{Nome: containerLog.PodName, QtContainers: 1})
	}
	containerLogs = append(containerLogs, containerLog)
}

func Inicializa(canalLogs chan models.ContainerLog) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicf("Erro ao inicializar interface: " + err.Error())
	}
	defer g.Close()
	// Highlight active view.
	g.Highlight = true
	g.SelFgColor = gocui.ColorBlue
	g.BgColor = gocui.ColorBlack
	g.FgColor = gocui.ColorWhite
	g.Cursor = true

	// The GUI object wants to know how to manage the layout.
	// Unlike termui, gocui does not use a grid layout.
	// Instead, it relies on a custom layout handler function to manage the layout.
	//
	// Here we set the layout manager to a function named layout that is defined further down.
	g.SetManagerFunc(layout)

	// Bind the quit handler function (also defined further down) to Ctrl-C,
	// so that we can leave the application at any time.
	if err := keybindings(g); err != nil {
		panic("Não foi possível fazer o 'key binding': " + err.Error())
	}

	// View definitions *******************************************************************
	// Largura e altura do terminal necessários para cálculos de layout
	terminalWidth, terminalHeight := g.Size()
	// view Nós.
	nodesView, err := g.SetView(nodesTitulo, 0, 0, nodesRight, terminalHeight-4)
	// ErrUnknownView is not a real error condition.
	// It just says that the view did not exist before and needs initialization.
	if err != nil && err != gocui.ErrUnknownView {
		log.Println("Falha a criar visão de Nós", err)
		return
	}
	nodesView.Title = nodesTitulo
	nodesView.FgColor = gocui.ColorCyan
	nodesView.Highlight = true
	nodesView.SelBgColor = gocui.ColorBlue
	nodesView.SelFgColor = gocui.ColorBlack

	// view PODs.
	podsView, err := g.SetView(podsTitulo, nodesRight+1, 0, podsRight, terminalHeight-4)
	// ErrUnknownView is not a real error condition.
	// It just says that the view did not exist before and needs initialization.
	if err != nil && err != gocui.ErrUnknownView {
		log.Println("Falha a criar visão de PODs:", err)
		return
	}
	podsView.Title = podsTitulo
	podsView.FgColor = gocui.ColorCyan
	podsView.SelBgColor = gocui.ColorBlue
	podsView.SelFgColor = gocui.ColorBlack

	// view Logs
	logsView, err := g.SetView("output", podsRight+1, 0, terminalWidth-1, terminalHeight-4)
	if err != nil && err != gocui.ErrUnknownView {
		log.Println("Falha a criar visão de Logs:", err)
		return
	}
	logsView.Title = logsTitulo
	logsView.FgColor = gocui.ColorWhite
	// Let the view scroll if the output exceeds the visible area.
	logsView.Autoscroll = true
	logsView.Wrap = true

	//seleciona o nó do primeiro container
	if len(nos) > 0 {
		noSelecionado = nos[0]
	}

	for cl := range canalLogs {
		AdicionaContainerLog(cl)
	}

	// Must set initial view here, right before program start!!!
	g.SetCurrentView(nodesTitulo)

	if err := g.MainLoop(); err != nil {
		g.Close()
		panic(err.Error())
		//os.Exit(1)
	}
	log.Println("Main loop finalizado:", err)
}

/*func layout(g *gocui.Gui) error {
	if err := setSideLayout(g); err != nil {
		return err
	}
	/*if err := setSummaryLayout(g); err != nil {
		return err
	}
	if err := setDetailLayout(g); err != nil {
		return err
	}
	//return setChangelogLayout(g)
	return nil
}*/

func setSideLayout(g *gocui.Gui) error {
	_, maxY := g.Size()
	if v, err := g.SetView("side", -1, -1, 40, int(float64(maxY)*0.2)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true

		for _, result := range containerLogs {
			fmt.Fprintln(v, result.ContainerName)
		}
		/*if len(containerLogs) == 0 {
			return xerrors.New("Containers inexistentes")
		}*/
		//logs = containerLogs[0].Logs
		if _, err := g.SetCurrentView("side"); err != nil {
			return err
		}
	}
	return nil
}

func keybindings(g *gocui.Gui) (err error) {

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		panic(err.Error())
	}

	/*errs := []error{}

	// Move beetween views
	errs = append(errs, g.SetKeybinding("side", gocui.KeyTab, gocui.ModNone, nextView))
	//  errs = append(errs, g.SetKeybinding("side", gocui.KeyCtrlH, gocui.ModNone, previousView))
	//  errs = append(errs, g.SetKeybinding("side", gocui.KeyCtrlL, gocui.ModNone, nextView))
	//  errs = append(errs, g.SetKeybinding("side", gocui.KeyArrowRight, gocui.ModAlt, nextView))
	errs = append(errs, g.SetKeybinding("side", gocui.KeyArrowDown, gocui.ModNone, cursorDown))
	errs = append(errs, g.SetKeybinding("side", gocui.KeyCtrlJ, gocui.ModNone, cursorDown))
	errs = append(errs, g.SetKeybinding("side", gocui.KeyArrowUp, gocui.ModNone, cursorUp))
	errs = append(errs, g.SetKeybinding("side", gocui.KeyCtrlK, gocui.ModNone, cursorUp))
	errs = append(errs, g.SetKeybinding("side", gocui.KeyCtrlD, gocui.ModNone, cursorPageDown))
	errs = append(errs, g.SetKeybinding("side", gocui.KeyCtrlU, gocui.ModNone, cursorPageUp))
	errs = append(errs, g.SetKeybinding("side", gocui.KeySpace, gocui.ModNone, cursorPageDown))
	errs = append(errs, g.SetKeybinding("side", gocui.KeyBackspace, gocui.ModNone, cursorPageUp))
	errs = append(errs, g.SetKeybinding("side", gocui.KeyBackspace2, gocui.ModNone, cursorPageUp))
	errs = append(errs, g.SetKeybinding("side", gocui.KeyCtrlN, gocui.ModNone, cursorDown))
	errs = append(errs, g.SetKeybinding("side", gocui.KeyCtrlP, gocui.ModNone, cursorUp))
	errs = append(errs, g.SetKeybinding("side", gocui.KeyEnter, gocui.ModNone, nextView))

	//  errs = append(errs, g.SetKeybinding("msg", gocui.KeyEnter, gocui.ModNone, delMsg))
	//  errs = append(errs, g.SetKeybinding("side", gocui.KeyEnter, gocui.ModNone, showMsg))

	// summary
	errs = append(errs, g.SetKeybinding("summary", gocui.KeyTab, gocui.ModNone, nextView))
	errs = append(errs, g.SetKeybinding("summary", gocui.KeyCtrlQ, gocui.ModNone, previousView))
	errs = append(errs, g.SetKeybinding("summary", gocui.KeyCtrlH, gocui.ModNone, previousView))
	//  errs = append(errs, g.SetKeybinding("summary", gocui.KeyCtrlL, gocui.ModNone, nextView))
	//  errs = append(errs, g.SetKeybinding("summary", gocui.KeyArrowLeft, gocui.ModAlt, previousView))
	//  errs = append(errs, g.SetKeybinding("summary", gocui.KeyArrowDown, gocui.ModAlt, nextView))
	errs = append(errs, g.SetKeybinding("summary", gocui.KeyArrowDown, gocui.ModNone, cursorDown))
	errs = append(errs, g.SetKeybinding("summary", gocui.KeyArrowUp, gocui.ModNone, cursorUp))
	errs = append(errs, g.SetKeybinding("summary", gocui.KeyCtrlJ, gocui.ModNone, cursorDown))
	errs = append(errs, g.SetKeybinding("summary", gocui.KeyCtrlK, gocui.ModNone, cursorUp))
	errs = append(errs, g.SetKeybinding("summary", gocui.KeyCtrlD, gocui.ModNone, cursorPageDown))
	errs = append(errs, g.SetKeybinding("summary", gocui.KeyCtrlU, gocui.ModNone, cursorPageUp))
	errs = append(errs, g.SetKeybinding("summary", gocui.KeySpace, gocui.ModNone, cursorPageDown))
	errs = append(errs, g.SetKeybinding("summary", gocui.KeyBackspace, gocui.ModNone, cursorPageUp))
	errs = append(errs, g.SetKeybinding("summary", gocui.KeyBackspace2, gocui.ModNone, cursorPageUp))
	errs = append(errs, g.SetKeybinding("summary", gocui.KeyEnter, gocui.ModNone, nextView))
	errs = append(errs, g.SetKeybinding("summary", gocui.KeyCtrlN, gocui.ModNone, nextSummary))
	errs = append(errs, g.SetKeybinding("summary", gocui.KeyCtrlP, gocui.ModNone, previousSummary))

	// detail
	errs = append(errs, g.SetKeybinding("detail", gocui.KeyTab, gocui.ModNone, nextView))
	errs = append(errs, g.SetKeybinding("detail", gocui.KeyCtrlQ, gocui.ModNone, previousView))
	errs = append(errs, g.SetKeybinding("detail", gocui.KeyCtrlH, gocui.ModNone, nextView))
	//  errs = append(errs, g.SetKeybinding("detail", gocui.KeyCtrlL, gocui.ModNone, nextView))
	//  errs = append(errs, g.SetKeybinding("detail", gocui.KeyArrowUp, gocui.ModAlt, previousView))
	//  errs = append(errs, g.SetKeybinding("detail", gocui.KeyArrowLeft, gocui.ModAlt, nextView))
	errs = append(errs, g.SetKeybinding("detail", gocui.KeyArrowDown, gocui.ModNone, cursorDown))
	errs = append(errs, g.SetKeybinding("detail", gocui.KeyArrowUp, gocui.ModNone, cursorUp))
	errs = append(errs, g.SetKeybinding("detail", gocui.KeyCtrlJ, gocui.ModNone, cursorDown))
	errs = append(errs, g.SetKeybinding("detail", gocui.KeyCtrlK, gocui.ModNone, cursorUp))
	errs = append(errs, g.SetKeybinding("detail", gocui.KeyCtrlD, gocui.ModNone, cursorPageDown))
	errs = append(errs, g.SetKeybinding("detail", gocui.KeyCtrlU, gocui.ModNone, cursorPageUp))
	errs = append(errs, g.SetKeybinding("detail", gocui.KeySpace, gocui.ModNone, cursorPageDown))
	errs = append(errs, g.SetKeybinding("detail", gocui.KeyBackspace, gocui.ModNone, cursorPageUp))
	errs = append(errs, g.SetKeybinding("detail", gocui.KeyBackspace2, gocui.ModNone, cursorPageUp))
	//  errs = append(errs, g.SetKeybinding("detail", gocui.KeyCtrlM, gocui.ModNone, cursorMoveMiddle))
	errs = append(errs, g.SetKeybinding("detail", gocui.KeyCtrlN, gocui.ModNone, nextSummary))
	errs = append(errs, g.SetKeybinding("detail", gocui.KeyCtrlP, gocui.ModNone, previousSummary))
	errs = append(errs, g.SetKeybinding("detail", gocui.KeyEnter, gocui.ModNone, nextView))

	// changelog
	errs = append(errs, g.SetKeybinding("changelog", gocui.KeyTab, gocui.ModNone, nextView))
	errs = append(errs, g.SetKeybinding("changelog", gocui.KeyCtrlQ, gocui.ModNone, previousView))
	errs = append(errs, g.SetKeybinding("changelog", gocui.KeyCtrlH, gocui.ModNone, nextView))
	//  errs = append(errs, g.SetKeybinding("changelog", gocui.KeyCtrlL, gocui.ModNone, nextView))
	//  errs = append(errs, g.SetKeybinding("changelog", gocui.KeyArrowUp, gocui.ModAlt, previousView))
	//  errs = append(errs, g.SetKeybinding("changelog", gocui.KeyArrowLeft, gocui.ModAlt, nextView))
	errs = append(errs, g.SetKeybinding("changelog", gocui.KeyArrowDown, gocui.ModNone, cursorDown))
	errs = append(errs, g.SetKeybinding("changelog", gocui.KeyArrowUp, gocui.ModNone, cursorUp))
	errs = append(errs, g.SetKeybinding("changelog", gocui.KeyCtrlJ, gocui.ModNone, cursorDown))
	errs = append(errs, g.SetKeybinding("changelog", gocui.KeyCtrlK, gocui.ModNone, cursorUp))
	errs = append(errs, g.SetKeybinding("changelog", gocui.KeyCtrlD, gocui.ModNone, cursorPageDown))
	errs = append(errs, g.SetKeybinding("changelog", gocui.KeyCtrlU, gocui.ModNone, cursorPageUp))
	errs = append(errs, g.SetKeybinding("changelog", gocui.KeySpace, gocui.ModNone, cursorPageDown))
	errs = append(errs, g.SetKeybinding("changelog", gocui.KeyBackspace, gocui.ModNone, cursorPageUp))
	errs = append(errs, g.SetKeybinding("changelog", gocui.KeyBackspace2, gocui.ModNone, cursorPageUp))
	//  errs = append(errs, g.SetKeybinding("changelog", gocui.KeyCtrlM, gocui.ModNone, cursorMoveMiddle))
	errs = append(errs, g.SetKeybinding("changelog", gocui.KeyCtrlN, gocui.ModNone, nextSummary))
	errs = append(errs, g.SetKeybinding("changelog", gocui.KeyCtrlP, gocui.ModNone, previousSummary))
	errs = append(errs, g.SetKeybinding("changelog", gocui.KeyEnter, gocui.ModNone, nextView))

	//  errs = append(errs, g.SetKeybinding("msg", gocui.KeyEnter, gocui.ModNone, delMsg))
	//  errs = append(errs, g.SetKeybinding("detail", gocui.KeyEnter, gocui.ModNone, showMsg))

	errs = append(errs, g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit))
	//  errs = append(errs, g.SetKeybinding("side", gocui.KeyEnter, gocui.ModNone, getLine))
	//  errs = append(errs, g.SetKeybinding("msg", gocui.KeyEnter, gocui.ModNone, delMsg))

	for _, e := range errs {
		if e != nil {
			return e
		}
	}*/
	return nil
}
