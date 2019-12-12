package taskwarrior

/**
 * Settings:
 *  - Filtering (Default: +PENDING)
 *  - Number to display (Default: 10)
 * Keyboard Commands
 *  - Start a task
 *  - Delete a task
 *  - Mark a task complete
 *  - View task details
 * - Can I display a modal with an input form to add tasks
 */

import (
	"fmt"
	"github.com/rivo/tview"
	"github.com/wtfutil/wtf/utils"
	"github.com/wtfutil/wtf/view"
	"sort"
	"strconv"
)

type Widget struct {
	view.KeyboardWidget
	view.MultiSourceWidget
	view.ScrollableWidget

	settings *Settings
	tasks    []Task
}

// NewWidget creates a new instance of a widget
func NewWidget(app *tview.Application, pages *tview.Pages, settings *Settings) *Widget {
	tw, _ := NewTaskWarrior("~/.taskrc")
	tw.FetchAllTasks()
	widget := Widget{
		KeyboardWidget:    view.NewKeyboardWidget(app, pages, settings.common),
		MultiSourceWidget: view.NewMultiSourceWidget(settings.common, "task", "tasks"),
		ScrollableWidget:  view.NewScrollableWidget(app, settings.common),

		settings: settings,
		tasks:    tw.Tasks,
	}

	// Don't use a timer for this widget, watch for filesystem changes instead
	widget.settings.common.RefreshInterval = 0

	widget.initializeKeyboardControls()
	widget.View.SetInputCapture(widget.InputCapture)

	widget.SetDisplayFunction(widget.Refresh)

	widget.KeyboardWidget.SetView(widget.View)

	return &widget
}

/* -------------------- Exported Functions -------------------- */

// Refresh is only called once on start-up. Its job is to display the
// text files that first time. After that, the watcher takes over
func (widget *Widget) Refresh() {
	widget.Redraw(widget.content)
}

func (widget *Widget) HelpText() string {
	return widget.KeyboardWidget.HelpText()
}

/* -------------------- Unexported Functions -------------------- */
func (widget *Widget) content() (string, string, bool) {
	out := ""

	descriptionLength, projectLength, urgencyLength := widget.getLongestColumnLengths(widget.tasks)

	sort.Slice(widget.tasks, func(i, j int) bool { return widget.tasks[i].Urgency > widget.tasks[j].Urgency })

	displayRow := 0
	for _, task := range widget.tasks {
		if task.Status != "pending" {
			continue
		}
		taskId := strconv.Itoa(task.Id)
		taskUrgency := tview.Escape(fmt.Sprintf("%.2f", task.Urgency))
		row := fmt.Sprintf(
			`[%s]%-*s %-*s %-*s %-*s[%s]`,
			widget.RowColor(displayRow),
			len(taskId)+1,
			taskId,
			descriptionLength+1,
			tview.Escape(trimToMaxLength(task.Description, widget.settings.maxDescriptionLength)),
			projectLength+1,
			tview.Escape(trimToMaxLength(task.Project, widget.settings.maxProjectLength)),
			urgencyLength+1,
			taskUrgency,
			widget.RowColor(displayRow),
		)
		out += utils.HighlightableHelper(widget.View, row, displayRow, len(task.Description))
		displayRow++
	}

	return widget.CommonSettings().Title, out, false
}

func trimToMaxLength(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}
	return text[:maxLength]
}

func (widget *Widget) getLongestColumnLengths(tasks []Task) (int, int, int) {
	longestDescriptionLength := 0
	longestProjectLength := 0
	longestUrgencyLength := 0

	for _, task := range tasks {
		urgencyLength := len(fmt.Sprintf("%.2f", task.Urgency))
		if urgencyLength > longestUrgencyLength {
			longestUrgencyLength = urgencyLength
		}

		descriptionLength := len(task.Description)
		if descriptionLength > longestDescriptionLength {
			longestDescriptionLength = descriptionLength
		}

		projectLength := len(task.Project)
		if projectLength > longestProjectLength {
			longestProjectLength = projectLength
		}
	}

	if longestDescriptionLength > widget.settings.maxDescriptionLength {
		longestDescriptionLength = widget.settings.maxDescriptionLength
	}

	if longestProjectLength > widget.settings.maxProjectLength {
		longestProjectLength = widget.settings.maxProjectLength
	}

	return longestDescriptionLength, longestProjectLength, longestUrgencyLength
}
