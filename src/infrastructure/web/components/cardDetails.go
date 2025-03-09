package components

import (
	"fmt"

	"github.com/JneiraS/AMS/src/domain/services"
	"github.com/JneiraS/AMS/src/infrastructure/persistence"
	gom "maragu.dev/gomponents"
	gomh "maragu.dev/gomponents/html"
)

func CardDetails(task persistence.Task) gom.Node {
	return gomh.Section(gomh.Class("detail-task-container"),
		CardTitle(task.Title),
		Container("task-details", task),
	)
}

func CardComments(task persistence.Task, listOfComment []persistence.Comment) gom.Node {
	return gomh.Section(gomh.Class("detail-task-comment"),
		gomh.H1(gom.Text("Comments")),

		FormComment(task.ID),
		ListOfComments(listOfComment),
	)
}

func CardTitle(title string) gom.Node {
	return gomh.H1(gom.Text(title))
}
func Container(class string, task persistence.Task) gom.Node {
	return gomh.Div(gomh.Class(class),
		Paragraph("description_dt", task.Description),
		DetailContainer("task-info", task),
		ProgressBar(task),
		TaskTime(task),
	)
}

func DetailContainer(class string, task persistence.Task) gom.Node {
	return gomh.Div(gomh.Class(class),
		gom.If(!task.DueDate.IsZero(), KeyValue("Due Date", task.DueDate.Format("2006-01-02"))),
		KeyValue("Status", task.Status),
		KeyValue("Priority", services.FormatPryority(task.Priority)),
		KeyValue("Assignee", task.Assignee),
		KeyValue("Creator", task.Creator),
		KeyValue("Project", task.Project),
	)
}

func Paragraph(class, text string) gom.Node {
	return gomh.P(gomh.Class(class),
		gom.Text(text))
}

func TaskTime(task persistence.Task) gom.Node {
	return gomh.Div(gomh.Class("task-time"),
		KeyValue("Estimated Time", fmt.Sprintf("%d hours", task.EstimatedTime)),
		KeyValue("Time Spent", services.FormatTimeSpent(task.TimeSpent)),
	)
}

func KeyValue(key, value string) gom.Node {
	return gomh.P(gomh.Strong(gom.Text(key+": ")),
		gom.Text(value))
}

func ProgressBar(task persistence.Task) gom.Node {
	return gomh.Div(gomh.Class("task-progress"),
		gomh.Div(gomh.Class("progress-bar"),
			gomh.Div(gomh.Class("progress"),
				gomh.Style("width: "+fmt.Sprintf("%d%%", task.Progress))),
		),
	)
}

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

func ListOfComments(listOfComment []persistence.Comment) gom.Node {
	return gomh.Div(gomh.Class("comment-list"),
		gom.Map(listOfComment, func(comment persistence.Comment) gom.Node {
			return ListComments(comment)
		}),
	)
}
func ListComments(listOfComment persistence.Comment) gom.Node {
	return gomh.Div(
		gomh.Class("comment"),

		gomh.Div(gomh.Class("comment-time"),
			gom.Text(listOfComment.CreatedAt.Format("02-01-2006 at 15:04")),
		),

		gomh.Div(gomh.Class("comment-container-inner"),

			gomh.Div(gomh.Class("comment-author"),
				gom.Text(listOfComment.Author),
			),

			gomh.Div(gomh.Class("comment-text"),
				gom.Text(listOfComment.Text),
			),
		),
	)

}
