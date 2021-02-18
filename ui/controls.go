package ui

import "github.com/jroimartin/gocui"

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	// if err := g.SetKeybinding(nodesTitulo, gocui.KeyEsc, gocui.ModNone, limpaDebug); err != nil {
	// 	return err
	// }
	if err := g.SetKeybinding(nodesTitulo, gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding(nodesTitulo, gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	// if err := g.SetKeybinding(nodesTitulo, gocui.KeyCtrlA, gocui.ModNone, inputView); err != nil {
	// 	return err
	// }
	if err := g.SetKeybinding(nodesTitulo, gocui.KeyArrowRight, gocui.ModNone, selectItem); err != nil {
		return err
	}
	// if err := g.SetKeybinding(nodesTitulo, gocui.KeyCtrlR, gocui.ModNone, deleteItem); err != nil {
	// 	return err
	// }
	// if err := g.SetKeybinding(nodesTitulo, gocui.KeyCtrlD, gocui.ModNone, addDescription); err != nil {
	// 	return err
	// }
	if err := g.SetKeybinding(podsTitulo, gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding(podsTitulo, gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	// if err := g.SetKeybinding(podsTitulo, gocui.KeyCtrlA, gocui.ModNone, inputView); err != nil {
	// 	return err
	// }
	if err := g.SetKeybinding(podsTitulo, gocui.KeyArrowRight, gocui.ModNone, selectItem); err != nil {
		return err
	}
	// if err := g.SetKeybinding(podsTitulo, gocui.KeyCtrlR, gocui.ModNone, deleteItem); err != nil {
	// 	return err
	// }
	if err := g.SetKeybinding(podsTitulo, gocui.KeyArrowLeft, gocui.ModNone, goBack); err != nil {
		return err
	}
	if err := g.SetKeybinding(containersTitulo, gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding(containersTitulo, gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	// if err := g.SetKeybinding(podsTitulo, gocui.KeyCtrlA, gocui.ModNone, inputView); err != nil {
	// 	return err
	// }
	// if err := g.SetKeybinding(containersTitulo, gocui.KeyArrowRight, gocui.ModNone, selectItem); err != nil {
	// 	return err
	// }
	// if err := g.SetKeybinding(podsTitulo, gocui.KeyCtrlR, gocui.ModNone, deleteItem); err != nil {
	// 	return err
	// }
	if err := g.SetKeybinding(containersTitulo, gocui.KeyArrowLeft, gocui.ModNone, goBack); err != nil {
		return err
	}

	// if err := g.SetKeybinding(podsTitulo, gocui.KeyCtrlD, gocui.ModNone, addDescription); err != nil {
	// 	return err
	// }
	// if err := g.SetKeybinding(logsTitulo, gocui.KeyCtrlS, gocui.ModNone, save); err != nil {
	// 	return err
	// }

	return nil
}
