package ui

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
	"github.com/k8smon/models"
	"github.com/pkg/errors"
)

var bottom = false
var dView *gocui.View

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	// Check to make sure data exists in the next line,
	// otherwise disallow scroll down.
	if v != nil && lineBelow(g, v) == true {
		v.MoveCursor(0, 1, false)
		_, cy := v.Cursor()
		//debuga(g, fmt.Sprintf("DOWN: %d", cy))
		n, _ := v.Line(cy)
		if v.Name() == nodesTitulo {
			nv, _ := g.View(podsTitulo)
			if n != "" {
				noSelecionado = *models.GetNodeByName(nos, n) // models.GetProject(n)
			}
			//debuga(g, fmt.Sprintf("cursorDown node: %s", noSelecionado.Nome))
			redrawPods(g, nv)
		} else if v.Name() == podsTitulo {
			nv, _ := g.View(containersTitulo)
			if n != "" {
				podSelecionado = *models.GetPodByName(pods, n)
			}
			//debuga(g, fmt.Sprintf("cursorDown POD: %s", podSelecionado.Nome))
			redrawContainerss(g, nv)
		} else if v.Name() == containersTitulo {
			nv, _ := g.View(logsTitulo)
			if n != "" {
				containerSelecionado = n
			}
			//debuga(g, fmt.Sprintf("cursorDown POD: %s", podSelecionado.Nome))
			redrawLogs(g, nv)
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		v.MoveCursor(0, -1, false)
		_, cy := v.Cursor()
		//debuga(g, fmt.Sprintf("UP: %d", cy))
		n, _ := v.Line(cy)
		if v.Name() == nodesTitulo {
			nv, _ := g.View(podsTitulo)
			if n != "" {
				noSelecionado = *models.GetNodeByName(nos, n)
			} else {
				noSelecionado = models.Node{}
			}
			//debuga(g, fmt.Sprintf("cursorUp node: %s", noSelecionado.Nome))
			redrawPods(g, nv)
		} else if v.Name() == podsTitulo {
			nv, _ := g.View(containersTitulo)
			if n != "" {
				podSelecionado = *models.GetPodByName(pods, n)
			}
			//debuga(g, fmt.Sprintf("cursorDown POD: %s", podSelecionado.Nome))
			redrawContainerss(g, nv)
		} else if v.Name() == containersTitulo {
			nv, _ := g.View(logsTitulo)
			if n != "" {
				containerSelecionado = n
			}
			//debuga(g, fmt.Sprintf("cursorDown POD: %s", podSelecionado.Nome))
			redrawLogs(g, nv)
		}
	}
	return nil
}

// Returns true if there is a non-empty string in cursor position y+1.
// Otherwise returns false.
func lineBelow(g *gocui.Gui, v *gocui.View) bool {
	_, cy := v.Cursor()
	if l, _ := v.Line(cy + 1); l != "" {
		return true
	}
	return false
}

// Copy the input view (iv) and handle it.
// Used to add project or task.
// TODO: Fix error that pops up when adding two entries in one minute
/*func copyInput(g *gocui.Gui, iv *gocui.View) error {
	var err error
	// We want to read the view’s buffer from the beginning.
	iv.Rewind()
	// Get the output view via its name.
	var ov *gocui.View
	// If there is text input then add the item,
	// else go back to the input view.
	switch iv.Name() {
	case "addProject":
		ov, _ = g.View(NODES_TITULO)
		if iv.Buffer() != "" {
			models.AddProject(iv.Buffer())
		} else {
			inputView(g, ov)
			return nil
		}
	case "addTask":
		ov, _ = g.View(PODS_TITULO)
		if iv.Buffer() != "" {
			models.AddTask(iv.Buffer(), models.CurrentProject)
		} else {
			inputView(g, ov)
			return nil
		}
	}
	// Clear the input view
	iv.Clear()
	// No input, no cursor.
	g.Cursor = false
	// !!!
	// Must delete keybindings before the view, or fatal error !!!
	// !!!
	g.DeleteKeybindings(iv.Name())
	if err = g.DeleteView(iv.Name()); err != nil {
		return err
	}
	// Set the view back.
	if _, err = g.SetCurrentView(ov.Name()); err != nil {
		return err
	}
	switch ov.Name() {
	case NODES_TITULO:
		redrawNos(g, ov)
	case PODS_TITULO:
		redrawPods(g, ov)
	}
	return err
}*/

// Add item to the current view (cv) using the text from the input view (iv).
/*func inputView(g *gocui.Gui, cv *gocui.View) error {
	maxX, maxY := g.Size()
	var title string
	var name string
	switch cv.Name() {
	case NODES_TITULO:
		title = "Name of new project"
		name = "addProject"
	case PODS_TITULO:
		title = "Name of new task"
		name = "addTask"
	}
	if iv, err := g.SetView(name, maxX/2-12, maxY/2, maxX/2+12, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		iv.Title = title
		iv.Editable = true
		g.Cursor = true
		if _, err := g.SetCurrentView(name); err != nil {
			return err
		}
		if err := g.SetKeybinding(name, gocui.KeyEnter, gocui.ModNone, copyInput); err != nil {
			return err
		}
	}
	return nil
}*/

// Get the current view (cv) and transfer cursor to the new view (nv).
// Disallow if there is no string at current cursor.
func selectItem(g *gocui.Gui, cv *gocui.View) error {
	var err error
	_, cy := cv.Cursor()
	n, _ := cv.Line(cy)
	// If line at cursor is not empty (item is selected) then continue.
	if n != "" {
		var nv *gocui.View
		switch cv.Name() {
		case nodesTitulo:
			if nv, err = g.SetCurrentView(podsTitulo); err != nil {
				return err
			}
			// log.Println("selectItem project view CurrentEntry:", models.CurrentEntry.ID)
			nv.SetCursor(0, 0)
			cursorUp(g, nv)
		case podsTitulo:
			if nv, err = g.SetCurrentView(containersTitulo); err != nil {
				return err
			}
			// log.Println("selectItem task view CurrentEntry:", models.CurrentEntry.ID)
			nv.SetCursor(0, 0)
			cursorUp(g, nv)
		}
		// Turn on highlight and set cursor to 0,0 of the new view.
		nv.Highlight = true
		if err = nv.SetCursor(0, 0); err != nil {
			return err
		}
	}
	return nil
}

// Get the name of the item at the cursor and delete it.
// Disallow if there is no string at current cursor.
// Hmm... I'm thinking that if two tasks have the same name it would randomly pick one.
// Maybe not, but I forget how it sets the models.Current[Item]. Something to look out for.
/*func deleteItem(g *gocui.Gui, v *gocui.View) error {
	var err error
	_, cy := v.Cursor()
	n, _ := v.Line(cy)
	// If line at cursor is not empty (item is selected) then continue.
	if n != "" {
		switch v.Name() {
		case NODES_TITULO:
			models.CurrentProject.Delete()
			models.CurrentProject = models.Project{}
			redrawNos(g, v)
			v.SetCursor(0, 0)
			cursorUp(g, v)
		case PODS_TITULO:
			models.CurrentTask.Delete()
			models.CurrentTask = models.Task{}
			redrawTasks(g, v)
			v.SetCursor(0, 0)
			cursorUp(g, v)

		case logsTitulo:
			models.CurrentEntry.Delete()
			models.CurrentEntry = models.Entry{}
			redrawEntries(g, v)
			v.SetCursor(0, 0)
			cursorUp(g, v)
		}
	}
	return err
}*/

// Get the current view (cv) and transfer cursor to the new view (nv).
// Basically the opposite of selectItem.
func goBack(g *gocui.Gui, cv *gocui.View) error {
	var err error
	var nv *gocui.View
	switch cv.Name() {
	// Move from pods to nós view.
	case podsTitulo:
		if nv, err = g.SetCurrentView(nodesTitulo); err != nil {
			return err
		}
		containersView, _ := g.View(containersTitulo)
		redrawContainerss(g, containersView)
	case containersTitulo:
		if nv, err = g.SetCurrentView(podsTitulo); err != nil {
			return err
		}
		logsView, _ := g.View(logsTitulo)
		redrawLogs(g, logsView)
	}

	// Turn off highlight of current view and make sure it's on for the new view.
	cv.Highlight = false
	// Probably redundant.
	nv.Highlight = true
	return nil
}

// Get the view and redraw it with current database info.
// It's important to note that this function will call
// redrawTasks, which will call
// redrawEntries, which will call
// redrawOutput. Make it fucking rain.
func redrawNos(g *gocui.Gui, v *gocui.View) {
	// Clear the view of content and redraw it with a fresh database query.
	v.Clear()
	// Loop through projects to add their names to the view.
	for _, i := range nos {
		// We can simply Fprint to a view.
		_, err := fmt.Fprintln(v, i.Nome)
		if err != nil {
			log.Println("Error writing to the nodes view:", err)
		}
	}
	// While the text may shift lines on insert the cursor does not,
	// so we need to refresh the tasks view with the currently highlighted project.
	_, cy := v.Cursor()
	l, _ := v.Line(cy)
	if l != "" {
		noSelecionado = *models.GetNodeByName(nos, l)
	}
	// log.Println("redrawProjects project.id:", models.CurrentProject.ID)
	podsView, _ := g.View(podsTitulo)
	// Projects is only redrawn if in the projects view, so it's
	// safe to zero the current task and entry.
	podSelecionado = models.Pod{}
	redrawPods(g, podsView)
	podsView.Highlight = false
}

// Get the view and redraw it with current database info.
func redrawPods(g *gocui.Gui, v *gocui.View) {
	// Clear the view of content and redraw it with a fresh database query.
	v.Clear()
	// Loop through tasks to add their names to the view.
	for _, pod := range pods {
		// We can simply Fprint to a view.
		if pod.NomeNode == noSelecionado.Nome {
			_, err := fmt.Fprintln(v, pod.Nome)

			if err != nil {
				log.Println("Erro no desenho da view de PODs:", err)
			}
		}
	}
	if len(pods) != 0 {
		_, cy := v.Cursor()
		l, _ := v.Line(cy)
		p := models.GetPodByName(pods, l)
		if p != nil {
			podSelecionado = *p
		}
	}
	containerView, _ := g.View(containersTitulo)
	containerSelecionado = ""
	redrawContainerss(g, containerView)
	containerView.Highlight = false
}

// Get the view and redraw it with current database info.
func redrawContainerss(g *gocui.Gui, v *gocui.View) {
	// Clear the view of content and redraw it with a fresh database query.
	v.Clear()
	// Loop through tasks to add their names to the view.
	for _, container := range podSelecionado.NomeConteineres {
		// We can simply Fprint to a view.
		_, err := fmt.Fprintln(v, container)
		if err != nil {
			log.Println("Erro no desenho da view de containers:", err)
		}
	}
	if len(podSelecionado.NomeConteineres) != 0 {
		_, cy := v.Cursor()
		l, _ := v.Line(cy)
		if l != "" {
			containerSelecionado = l
		}
	}
	logView, _ := g.View(logsTitulo)
	redrawLogs(g, logView)
	logView.Highlight = false
}

// Get the view and redraw it with current database info.
// The output view should not need to be redrawn while it is itself selected,
// but we'll see...
// v is always the output view.
func redrawLogs(g *gocui.Gui, v *gocui.View) {
	// Clear the view of content and redraw it with a fresh database query.
	v.Clear()
	if cv := g.CurrentView(); cv != nil {
		// Projects
		if cv.Name() == nodesTitulo {
			if _, err := fmt.Fprintf(v, "%d PODs\n%d Containers\n\n",
				noSelecionado.QtPods, noSelecionado.QtContainers); err != nil {
				log.Println("Erro na renderização dos logs do nó:", err)
			}
		}
		if cv.Name() == podsTitulo {
			if _, err := fmt.Fprintf(v, "%d Containeres\n\n",
				len(podSelecionado.NomeConteineres)); err != nil {
				log.Println("Erro na renderização dos logs do pod:", err)
			}
		}
	}
}

// The layout handler calculates all sizes depending on the current terminal size.
func layout(g *gocui.Gui) error {
	// Get the current terminal size.
	tw, th := g.Size()
	// Update the views according to the new terminal size.
	// Nós.
	_, err := g.SetView(nodesTitulo, 0, 0, nodesRight, th-alturaLog-margemInferior)
	if err != nil {
		return errors.Wrap(err, "Cannot update nodes view")
	}
	// Pods
	_, err = g.SetView(podsTitulo, nodesRight+1, 0, podsRight, th-alturaLog-margemInferior)
	if err != nil {
		return errors.Wrap(err, "Cannot update pods view")
	}
	// Containers
	_, err = g.SetView(containersTitulo, podsRight+1, 0, tw-1, th-alturaLog-margemInferior)
	if err != nil {
		return errors.Wrap(err, "Cannot update pods view")
	}
	// Output
	_, err = g.SetView(logsTitulo, 0, th-alturaLog-margemInferior+1, tw-1, th-margemInferior)
	if err != nil {
		return errors.Wrap(err, "Cannot update logs view")
	}
	// Status
	// Not used right now. If uncommented set all above SetView() y1 values to 'th-4'.
	// _, err = g.SetView("status", 0, th-sheight, tw-1, th-1)
	// if err != nil {
	// 	return errors.Wrap(err, "Cannot update input view.")
	// }
	return nil
}

// quit is a handler that gets bound to Ctrl-gocui. It signals the main loop to exit.
func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

// Get the view and redraw it with current database info.
// It's important to note that this function will call
// redrawTasks, which will call
// redrawEntries, which will call
// redrawOutput. Make it fucking rain.
/*func debuga(g *gocui.Gui, s string) {
	v, _ := g.View("debug")
	// Clear the view of content and redraw it with a fresh database query.
	//v.Clear()
	_, err := fmt.Fprintln(v, s)
	if err != nil {
		log.Println("Erro no debug:", err)
	}
}

func limpaDebug(g *gocui.Gui, v *gocui.View) error {
	vd, _ := g.View("debug")
	// Clear the view of content and redraw it with a fresh database query.
	vd.Clear()
	return nil
}*/
