document.querySelectorAll(".description").forEach(function (el) {
    const id = el.id ? el.id.split("-")[1] : null;
    el.querySelectorAll("li").forEach(function (li) {
        li.addEventListener("click", function () {
            this.contentEditable = true;
            this.focus();

            const save = (e) => {
                if (e.key === "Enter" && (e.ctrlKey || e.metaKey)) {
                    fetch("/update-task", {
                        method: "POST",
                        headers: {
                            "Content-Type": "application/json"
                        },
                        body: JSON.stringify({
                            id: id,
                            description: this.textContent
                        })
                    })
                        .then(response => {
                            if (response.ok) {
                                this.contentEditable = false;
                                this.removeEventListener("keyup", save);
                            } else {
                                console.error("Failed to update task description");
                            }
                        })
                        .catch(error => {
                            console.error("Error updating task description:", error);
                        });
                }
            };

            this.addEventListener("blur", () => {
                this.contentEditable = false;
                this.removeEventListener("keyup", save);
            });

            this.addEventListener("keyup", save);
        });
    });
});

document.querySelectorAll("[id^='title-']").forEach(function (li) {
    li.addEventListener("click", function () {
        this.contentEditable = true;
        this.focus();

        const saveTitle = function (event) {
            if (event.key === "Enter" && (event.ctrlKey || event.metaKey)) {
                const taskId = this.id.split("-")[1];
                const taskTitle = this.textContent;
                fetch(`/update-task-title/`, {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json"
                    },
                    body: JSON.stringify({
                        id: taskId,
                        title: taskTitle
                    })
                })
                    .then(function (res) {
                        if (res.status === 200) {
                            this.contentEditable = false;
                            this.removeEventListener("keyup", saveTitle);
                        } else {
                            console.log("Error updating task title");
                        }
                    })
                    .catch(function (err) {
                        console.log("Error updating task title", err);
                    });
            }
        };

        this.addEventListener("blur", () => {
            this.contentEditable = false;
            this.removeEventListener("keyup", saveTitle);
        });

        this.addEventListener("keyup", saveTitle);
    });
});

