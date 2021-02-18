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
var containerSelecionado string

const (
	//Nós
	nodesTitulo = "Nodes"
	nodesRight  = 50
	//Pods
	podsTitulo = "PODs"
	podsRight  = 130
	//Containers
	containersTitulo = "Containers"

	//Logs
	logsTitulo     = "Logs"
	alturaLog      = 10
	margemInferior = 4
)

func AdicionaPod(pod models.Pod) {
	if models.NodeIsIn(nos, pod.NomeNode) == false {
		no := models.Node{Nome: pod.NomeNode, QtContainers: len(pod.NomeConteineres), QtPods: 1}
		nos = append(nos, no)
	} else {
		no := models.GetNodeByName(nos, pod.NomeNode)
		no.QtPods++
		no.QtContainers = no.QtContainers + len(pod.NomeConteineres)
	}
	if models.PodIsIn(pods, pod.Nome) == false {
		pods = append(pods, pod)
	}
}

//Inicializa - inicializa a UI
func Inicializa(canalPods chan models.Pod, canalLogs chan models.ContainerLog) {
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
	g.Cursor = false

	// The GUI object wants to know how to manage the layout.
	// Unlike termui, gocui does not use a grid layout.
	// Instead, it relies on a custom layout handler function to manage the layout.
	//
	// Here we set the layout manager to a function named layout that is defined further down.
	g.SetManagerFunc(layout)

	// View definitions *******************************************************************
	// Largura e altura do terminal necessários para cálculos de layout
	terminalWidth, terminalHeight := g.Size()
	// view Nós.
	nodesView, err := g.SetView(nodesTitulo, 0, 0, nodesRight, terminalHeight-14)
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
	podsView, err := g.SetView(podsTitulo, nodesRight+1, 0, podsRight, terminalHeight-14)
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

	// view Containers.
	containerView, err := g.SetView(containersTitulo, podsRight+1, 0, terminalWidth-1, terminalHeight-alturaLog-margemInferior)
	// ErrUnknownView is not a real error condition.
	// It just says that the view did not exist before and needs initialization.
	if err != nil && err != gocui.ErrUnknownView {
		log.Println("Falha a criar visão de container:", err)
		return
	}
	containerView.Title = containersTitulo
	containerView.FgColor = gocui.ColorCyan
	containerView.SelBgColor = gocui.ColorBlue
	containerView.SelFgColor = gocui.ColorBlack

	// view Logs
	logsView, err := g.SetView(logsTitulo, 0, terminalHeight-4-alturaLog+1, terminalWidth-1, terminalHeight-4)
	if err != nil && err != gocui.ErrUnknownView {
		log.Println("Falha a criar visão de Logs:", err)
		return
	}
	logsView.Title = logsTitulo
	logsView.FgColor = gocui.ColorWhite
	// Let the view scroll if the output exceeds the visible area.
	logsView.Autoscroll = true
	logsView.Wrap = true

	debugView, err := g.SetView("debug", 0, terminalHeight-4-alturaLog+1, terminalWidth-1, terminalHeight-1)
	if err != nil && err != gocui.ErrUnknownView {
		log.Println("Falha a criar visão de debug:", err)
		return
	}
	debugView.Autoscroll = true
	debugView.Wrap = true

	//Carrega a lista de pods do Canal
	for pod := range canalPods {
		AdicionaPod(pod)
	}

	//seleciona o nó do primeiro container
	if len(nos) > 0 {
		noSelecionado = nos[0]
	}
	//Apply keybindings to program.
	if err = keybindings(g); err != nil {
		log.Panicln(err)
	}
	// Must set initial view here, right before program start!!!
	v, _ := g.SetCurrentView(nodesTitulo)
	redrawNos(g, v)
	// Move the cursor to update the output view with the description.
	// (workaround)
	cursorUp(g, v)

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
