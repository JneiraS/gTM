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
                    .then(response => response.json())
                    .then(data => {
                        if (data.ok) {
                            this.contentEditable = false;
                            this.removeEventListener("keyup", saveDate);
                        } else if (data.error === "Due date must be in the future") {
                            console.error("Error updating task date:", data.error);
                            const errorDiv = document.getElementById('error-message');
                            errorDiv.textContent = data.error;
                            errorDiv.style.display = 'block';
                            errorDiv.className = 'alert alert-danger';
                        } else {
                            console.error("Error updating task date:", data.error);
                        }
                    })
                    .catch(error => {
                        console.error("Error updating task date:", error);
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


document.querySelectorAll("[id^='status-']").forEach(function (li) {
    li.addEventListener("click", function () {
        const statusText = this.textContent.trim();
        const dateId = this.id.split("-")[1];

        fetch(`/update-task-status/`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                id: dateId,
                status: statusText === "Pending" || statusText === "Stopped" ? "In progress" : "Stopped"
            })
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                return response.json();
            })
            .then(data => {
                if (data.status === 200) {
                    this.textContent = statusText === "Pending" || statusText === "Stopped" ? "In progress" : "Stopped";
                    this.classList.toggle('pending');
                    this.classList.toggle('stopped');
                    this.classList.toggle('in-progress');
                    setStatusColors();
                } else {
                    throw new Error('Server returned error status');
                }
            })
            .catch(function (err) {
                console.error("Error updating task status:", err);
                alert('Failed to update task status. Please try again.');
            });
    });
});

document.addEventListener('DOMContentLoaded', function () {
    const doneButtons = document.querySelectorAll('.done-bottom');

    doneButtons.forEach(button => {
        let pressTimer;

        button.addEventListener('mousedown', function (e) {
            pressTimer = setTimeout(() => {
                const href = this.getAttribute('href');
                window.location.href = href;
            }, 2000); // 2 seconds
            e.preventDefault();
        });

        button.addEventListener('click', function (e) {
            e.preventDefault();
        });

        button.addEventListener('mouseup', function () {
            clearTimeout(pressTimer);
        });

        button.addEventListener('mouseleave', function () {
            clearTimeout(pressTimer);
        });
    });
});

document.querySelectorAll(".done-bottom button").forEach((button) => {
    let pressTimer;

    button.addEventListener('mousedown', function (e) {
        pressTimer = setTimeout(() => {
            this.classList.add("filled");
        }, 100); // 1 second
        e.preventDefault();
    });

    button.addEventListener('click', function (e) {
        e.preventDefault();
    });

    button.addEventListener('mouseup', function () {
        clearTimeout(pressTimer);
        this.classList.remove("filled");
    });

    button.addEventListener('mouseleave', function () {
        clearTimeout(pressTimer);
        this.classList.remove("filled");
    });
});

function setStatusColors() {
    document.querySelectorAll(".change-status").forEach((element) => {
        if (element.textContent === "In progress") {
            element.style.backgroundColor = "#61A875";
            element.style.color = "whitesmoke";
            element.style.border = "1px solid rgba(255, 255, 255, 0.25)";
        } else if (element.textContent === "Stopped") {
            element.style.backgroundColor = "#C07972";
            element.style.color = "whitesmoke";
            element.style.border = "1px solid rgba(255, 255, 255, 0.25)";

        } else if (element.textContent === "Pending") {
            element.style.backgroundColor = "#737677";
            element.style.color = "#C6C5B9";
            element.style.border = "1px solid rgba(255, 255, 255, 0.25)";

        }
    });
}

document.querySelectorAll(".description_dt").forEach(function (el) {
    const taskElement = el.closest(".task-details");
    const taskId = taskElement ? taskElement.id.split("-")[1] : null;

    el.addEventListener("click", function () {
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
                        id: taskId,
                        description: this.innerText.trim()
                    })
                })
                    .then(response => {
                        if (response.ok) {
                            this.contentEditable = false;
                            this.removeEventListener("keyup", save);
                        } else {
                            console.error("Échec de la mise à jour de la description");
                        }
                    })
                    .catch(error => {
                        console.error("Erreur lors de la mise à jour :", error);
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


setStatusColors()
