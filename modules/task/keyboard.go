package task

import "github.com/gdamore/tcell"

func (widget *Widget) initializeKeyboardControls() {
	widget.InitializeCommonControls(widget.Refresh)

	widget.SetKeyboardChar("j", widget.Next, "Select next item")
	widget.SetKeyboardChar("k", widget.Prev, "Select previous item")
	widget.SetKeyboardChar("o", widget.openTask, "View task details")

	widget.SetKeyboardKey(tcell.KeyDown, widget.Next, "Select next item")
	widget.SetKeyboardKey(tcell.KeyUp, widget.Prev, "Select previous item")
	widget.SetKeyboardKey(tcell.KeyEnter, widget.openTask, "View task details")
	widget.SetKeyboardKey(tcell.KeyEsc, widget.Unselect, "Clear selection")
}
