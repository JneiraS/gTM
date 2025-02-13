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
                            description: this.innerHTML.replace(/<div>/g, '\n').replace(/<\/div>/g, '').replace(/<br\s*\/?>/g, '\n').replace(/<pre>/g, '').replace(/<\/pre>/g, '').replace('&gt;', ">")
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

const formContainer = document.querySelector('.form-container');
function toggleForm() {
    const formElement = document.querySelector('.form-container');
    formElement.classList.toggle('is-shifted-left');
}

const styleElement = document.createElement('style');
styleElement.innerHTML = `
    .is-shifted-left {
        transition: .4s ease; 
    }

    .show-button:active {
        transform: rotate(45deg);
        transition: .4s ease;     }
      
    .show-button::after {
        content: "+";
        font-size: 1.5em;
        font-weight: bold;
    }
`;
document.head.appendChild(styleElement);



document.querySelectorAll("[id^='due-date-']").forEach(function (li) {
    li.addEventListener("click", function () {
        this.contentEditable = true;
        this.focus();

        const saveDate = function (event) {
            if (event.key === "Enter" && (event.ctrlKey || event.metaKey)) {
                const dateId = this.id.split("-")[2];
                const newDate = this.textContent;
                // Convert date from "DD.MM.YYYY à HH:mm" to "YYYY-MM-DD HH:mm"
                const [datePart, timePart] = newDate.split(" à ");
                const [day, month, year] = datePart.trim().split(".").map(part => part.trim());
                const formattedDate = `${year}-${month}-${day} ${timePart}`;

                fetch(`/update-task-due-date/`, {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json"
                    },
                    body: JSON.stringify({
                        id: dateId,
                        date: formattedDate
                    })
                })
                    .then(function (res) {
                        if (res.status === 200) {
                            this.contentEditable = false;
                            this.removeEventListener("keyup", saveDate);
                        } else {
                            console.log("Error updating task date");
                        }
                    })
                    .catch(function (err) {
                        console.log("Error updating task date", err);
                    });
            }
        };

        this.addEventListener("blur", () => {
            this.contentEditable = false;
            this.removeEventListener("keyup", saveDate);
        });

        this.addEventListener("keyup", saveDate);
    });
});