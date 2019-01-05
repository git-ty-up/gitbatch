package gui

import (
	"fmt"

	"github.com/isacikgoz/gitbatch/core/command"
	"github.com/jroimartin/gocui"
)

// not staged view
func (gui *Gui) openUnStagedView(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	v, err := g.SetView(unstageViewFeature.Name, maxX/2+1, 5, maxX-6, int(0.75*float32(maxY))-1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = unstageViewFeature.Title
	}

	return refreshUnstagedView(g)
}

func (gui *Gui) addChanges(g *gocui.Gui, v *gocui.View) error {
	r := gui.getSelectedRepository()

	_, cy := v.Cursor()
	_, oy := v.Origin()
	if len(unstagedFiles) <= 0 || len(unstagedFiles) < cy+oy {
		return nil
	}
	if err := command.Add(r, unstagedFiles[cy+oy], &command.AddOptions{}); err != nil {
		return err
	}

	return refreshAllStatusView(g, r, true)
}

func (gui *Gui) addAllChanges(g *gocui.Gui, v *gocui.View) error {
	r := gui.getSelectedRepository()
	if err := command.AddAll(r, &command.AddOptions{}); err != nil {
		return err
	}

	return refreshAllStatusView(g, r, true)
}

// refresh the main view and re-render the repository representations
func refreshUnstagedView(g *gocui.Gui) error {
	stageView, err := g.View(unstageViewFeature.Name)
	if err != nil {
		return err
	}
	stageView.Clear()
	_, cy := stageView.Cursor()
	_, oy := stageView.Origin()
	for i, file := range unstagedFiles {
		var prefix string
		if i == cy+oy {
			prefix = prefix + selectionIndicator
		}
		fmt.Fprintf(stageView, "%s%s%s %s\n", prefix, red.Sprint(string(file.X)), red.Sprint(string(file.Y)), file.Name)
	}
	return nil
}
