package components

import (
	"fmt"

	p "github.com/JneiraS/AMS/src/infrastructure/persistence"
	gom "maragu.dev/gomponents"
	gomh "maragu.dev/gomponents/html"
)

func CardSubtasks(task p.Task, listOfSubtask []p.Subtask) gom.Node {
	return gomh.Section(gomh.Class("detail-task-comment"),
		gomh.H1(gom.Text("SubTasks")),
		FormSubtask(task.ID),
		ListOfSubtasks(listOfSubtask),
	)
}

func FormSubtask(id uint) gom.Node {
	return gomh.Form(
		gomh.Class("comment-form"),
		gomh.Action("/subtask/"+fmt.Sprintf("%d", id)),
		gomh.Method("POST"),
		gomh.Textarea(
			gomh.Class("comment-input"),
			gomh.Placeholder("Tile..."),
			gomh.Name("title"),
		),
		gomh.Button(
			gomh.Type("submit"),
			gom.Text("Add SubTask"),
		),
	)
}

func ListOfSubtasks(listOfSubtask []p.Subtask) gom.Node {
	return gomh.Div(gomh.Class("subtask-list"),
		gom.Map(listOfSubtask, func(subtask p.Subtask) gom.Node {
			return Subtask(subtask, subtask.ID)
		}))
}
func Subtask(subtask p.Subtask, id uint) gom.Node {
	return gomh.Form(
		gomh.Class("subtask-check"),
		gomh.Action("/subtask/status-change/"+fmt.Sprintf("%d", id)),
		gomh.Method("POST"),
		gomh.Input(
			gomh.Type("checkbox"),
			gomh.ID(fmt.Sprintf("subtask-%d", id)),
			gomh.Name("completed"),
			gom.If(subtask.Completed(), gomh.Checked()),
		),
		gomh.Label(
			gomh.For(fmt.Sprintf("subtask-%d", id)),
			gom.Text(subtask.Title),
		),
	)
}
