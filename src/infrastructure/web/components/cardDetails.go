package components

import (
	"fmt"

	as "github.com/JneiraS/AMS/src/application/services"
	p "github.com/JneiraS/AMS/src/infrastructure/persistence"
	gom "maragu.dev/gomponents"
	gomh "maragu.dev/gomponents/html"
)

func CardDetails(task p.Task) gom.Node {
	return gomh.Section(gomh.Class("detail-task-container"),
		CardTitle(task.Title),
		Container("task-details", task),
	)
}

func CardComments(task p.Task, listOfComment []p.Comment) gom.Node {
	return gomh.Section(gomh.Class("detail-task-comment"),
		gomh.H1(gom.Text("Comments")),
		FormComment(task.ID),
		ListOfComments(listOfComment),
	)
}

func CardTitle(title string) gom.Node {
	return gomh.H1(gom.Text(title))
}

func Container(class string, task p.Task) gom.Node {
	return gomh.Div(gomh.Class(class), gomh.ID("task-"+fmt.Sprintf("%d", task.ID)),
		ParagraphWithClass("description_dt", task.Description),
		DetailContainer("task-info", task),
		ProgressBar(task),
		TaskTime(task),
	)
}

func DetailContainer(class string, task p.Task) gom.Node {
	return gomh.Div(gomh.Class(class),
		gom.If(!task.DueDate.IsZero(), KeyValue("Due Date", task.DueDate.Format("2006-01-02"), "fa-solid fa-calendar-xmark")),
		KeyValue("Status", task.Status, "fa-solid fa-ellipsis-vertical"),
		KeyValue("Priority", as.FormatPryority(task.Priority), "fa-solid fa-circle-exclamation"),
		KeyValue("Assignee", task.Assignee, "fas fa-user"),
		KeyValue("Creator", task.Creator, "fas fa-user"),
		KeyValue("Project", task.Project, "fa-solid fa-diagram-project"),
	)
}

func ParagraphWithClass(class, text string) gom.Node {
	return gomh.P(gomh.Class(class),
		gom.Text(text))
}

func TaskTime(task p.Task) gom.Node {
	return gomh.Div(gomh.Class("task-time"),
		KeyValue("Estimated Time", fmt.Sprintf("%d hours", task.EstimatedTime), "fa-solid fa-stopwatch-20"),
		KeyValue("Time Spent", as.FormatTimeSpent(task.TimeSpent), "fa-solid fa-hourglass-end"),
	)
}

// KeyValue returns a p element containing a key-value pair. The icon parameter
// allows to specify an icon to display before the key-value pair.
func KeyValue(key, value, icon string) gom.Node {
	return gomh.P(gomh.I(gomh.Class(icon)), gomh.Strong(gom.Text(key+": ")),
		gom.Text(value))
}

// ProgressBar returns a div element representing the progress bar of a task.
func ProgressBar(task p.Task) gom.Node {
	return gomh.Div(gomh.Class("task-progress"),
		gomh.Div(gomh.Class("progress-bar"),
			gomh.Div(gomh.Class("progress"),
				gomh.Style("width: "+fmt.Sprintf("%d%%", task.Progress))),
		),
	)
}

// FormComment returns a form element for adding a comment to a task with the given id.
func FormComment(id uint) gom.Node {
	return gomh.Form(
		gomh.Class("comment-form"),
		gomh.Action("/comment/"+fmt.Sprintf("%d", id)),
		gomh.Method("POST"),
		gomh.Textarea(
			gomh.Class("comment-input"),
			gomh.Placeholder("Add a comment..."),
			gomh.Name("comment"),
		),
		gomh.Button(
			gomh.Type("submit"),
			gom.Text("Add Comment"),
		),
	)
}

// ListOfComments returns a list element of comments from a slice of comments.
func ListOfComments(listOfComment []p.Comment) gom.Node {
	return gomh.Div(gomh.Class("comment-list"),
		gom.Map(listOfComment, func(comment p.Comment) gom.Node {
			return Comment(comment)
		}),
	)
}

// Comment returns a comment element for a list of comments.
func Comment(comment p.Comment) gom.Node {
	return gomh.Div(
		gomh.Class("comment"),
		gomh.Div(gomh.Class("comment-time"),
			gom.Text(comment.CreatedAt.Format("02-01-2006 at 15:04")),
		),
		gomh.Div(gomh.Class("comment-container-inner"),
			gomh.Div(gomh.Class("comment-author"),
				gom.Text(comment.Author),
			),
			gomh.Div(gomh.Class("comment-text"),
				gom.Text(comment.Text),
			),
		),
	)

}
