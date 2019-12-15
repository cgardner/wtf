package task

import (
	"fmt"
	"sort"

	"github.com/cgardner/go-taskwarrior"
	"github.com/rivo/tview"
	"github.com/wtfutil/wtf/utils"
	"github.com/wtfutil/wtf/view"
)

const (
	publishedDateLayout = "Mon, 02 2006 15:04:05"
)

// Widget is the container for RSS and Atom data
type Widget struct {
	view.KeyboardWidget
	view.ScrollableWidget

	settings *Settings
	err      error
	twClient *taskwarrior.TaskWarrior
	tasks    []taskwarrior.Task
}

// NewWidget creates a new instance of a widget
func NewWidget(app *tview.Application, pages *tview.Pages, settings *Settings) *Widget {
	twClient, _ := taskwarrior.NewTaskWarrior("~/.taskrc")

	widget := &Widget{
		KeyboardWidget:   view.NewKeyboardWidget(app, pages, settings.common),
		ScrollableWidget: view.NewScrollableWidget(app, settings.common),

		settings: settings,
		twClient: twClient,
	}

	widget.SetRenderFunction(widget.Render)
	widget.initializeKeyboardControls()
	widget.View.SetInputCapture(widget.InputCapture)

	widget.KeyboardWidget.SetView(widget.View)

	return widget
}

/* -------------------- Exported Functions -------------------- */

func (widget *Widget) FetchTasks() ([]taskwarrior.Task, error) {
	widget.twClient.FetchAllTasks()
	tasks := []taskwarrior.Task{}

	for _, task := range widget.twClient.Tasks {
		if task.Status != "pending" {
			continue
		}
		tasks = append(tasks, task)
	}

	tasks = widget.SortTasks(tasks)
	return tasks, nil
}

// Refresh updates the data in the widget
func (widget *Widget) Refresh() {
	tasks, err := widget.FetchTasks()
	if err != nil {
		widget.err = err
		widget.tasks = nil
		widget.SetItemCount(0)
	} else {
		widget.err = nil
		widget.tasks = tasks
		widget.SetItemCount(len(tasks))
	}

	widget.Render()
}

// Render sets up the widget data for redrawing to the screen
func (widget *Widget) Render() {
	widget.Redraw(widget.content)
}

/* -------------------- Unexported Functions -------------------- */

func (widget *Widget) content() (string, string, bool) {
	title := widget.CommonSettings().Title
	if widget.err != nil {
		return title, widget.err.Error(), true
	}
	data := widget.tasks
	if data == nil || len(data) == 0 {
		return title, "No data", false
	}
	var str string
	descriptionLength, projectLength := widget.getLongestColumnLengths(data)
	for idx, task := range data {
		rowColor := widget.RowColor(idx)

		displayDescription :=
			tview.Escape(trimToMaxLength(task.Description, widget.settings.maxDescriptionLength))

		row := fmt.Sprintf(
			`[%s]%2d %-*s %-*s %.2f[%s]`,
			rowColor,
			task.Id,
			descriptionLength+1,
			displayDescription,
			projectLength+1,
			tview.Escape(trimToMaxLength(task.Project, widget.settings.maxProjectLength)),
			task.Urgency,
			rowColor,
		)

		str += utils.HighlightableHelper(widget.View, row, idx, len(displayDescription))
	}

	return title, str, false
}

func (widget *Widget) SortTasks(tasks []taskwarrior.Task) []taskwarrior.Task {
	sort.Slice(tasks, func(i, j int) bool { return tasks[i].Urgency > tasks[j].Urgency })

	return tasks
}

func (widget *Widget) openTask() {
	sel := widget.GetSelected()
	if sel <= 0 || widget.tasks == nil || sel > len(widget.tasks) {
		return
	}
}

func (widget *Widget) getLongestColumnLengths(tasks []taskwarrior.Task) (int, int) {
	longestDescriptionLength := 0
	longestProjectLength := 0

	for _, task := range tasks {
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

	return longestDescriptionLength, longestProjectLength
}

func trimToMaxLength(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	}
	return text[:maxLength]
}
