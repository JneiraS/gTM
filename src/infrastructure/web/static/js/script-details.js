let listOfSubtask = document.querySelectorAll("[id^='subtask-']");
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

        fetch('/subtask/status-change/' + subtaskId, {
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
            .catch(error => {
                console.error('Error:', error);
                // Revert the checkbox state if the request failed
                this.checked = !isCompleted;
                this.parentNode.style.textDecoration = isCompleted ? "none" : "line-through";
                this.parentNode.style.opacity = isCompleted ? "1" : "0.25";
            });
    });
});
