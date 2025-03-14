let listOfSubtask = document.querySelectorAll("[id^='subtask-']");
const taskId = window.location.pathname.split('/').pop();



// console.log(listOfSubtask)

// Barrer les sous-tâches cochées
listOfSubtask.forEach(subtask => {
    if (subtask.checked) {
        subtask.parentNode.style.textDecoration = "line-through";
        subtask.parentNode.style.opacity = "0.25";
    }

    subtask.addEventListener("click", function () {
        if (this.checked) {
            this.parentNode.style.textDecoration = "line-through";
            this.parentNode.style.opacity = "0.25";
        }
        else {
            this.parentNode.style.textDecoration = "none";
            this.parentNode.style.opacity = "1";
        }
    });
});



listOfSubtask.forEach(subtask => {
    subtask.addEventListener("change", function () {
        const subtaskId = this.id.split('-')[1];
        const isCompleted = this.checked;

        fetch('/subtask/' + taskId + '/status-change/' + subtaskId, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                completed: isCompleted
            })
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                return response.json();
            })
            .then(() => {
                if (isCompleted) {
                    this.parentNode.style.textDecoration = "line-through";
                    this.parentNode.style.opacity = "0.25";
                } else {
                    this.parentNode.style.textDecoration = "none";
                    this.parentNode.style.opacity = "1";
                }
            })

        StatusBarEvo();

    });
});

function StatusBarEvo() {
    let totalSubtasks = document.querySelectorAll("[id^='subtask-']").length;
    let checkedSubtasks = document.querySelectorAll("[id^='subtask-']:checked").length;
    let porcentOfCompletedSubtasks = (checkedSubtasks / totalSubtasks) * 100;
    document.querySelector(".progress").style.width = `${porcentOfCompletedSubtasks}%`;
}
