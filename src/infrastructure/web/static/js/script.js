document.querySelectorAll(".description").forEach(function (el) {
    const id = el.id.split("-")[1];
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
                    });
                    this.contentEditable = false;
                    this.removeEventListener("keyup", save);
                }
            };

            this.addEventListener("keyup", save);
        });
    });
});

document.querySelectorAll("[class^='task-card-header-'] li").forEach(function (li) {
    li.addEventListener("click", function () {
        this.contentEditable = true;
        this.focus();

        const save = (e) => {
            if (e.key === "Enter" && (e.ctrlKey || e.metaKey)) {
                const id = this.closest(".description").id.split("-")[1];
                fetch("/update-title-task", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json"
                    },
                    body: JSON.stringify({
                        id: id,
                        description: this.textContent
                    })
                });
                this.contentEditable = false;
                this.removeEventListener("keyup", save);
            }
        };

        this.addEventListener("keyup", save);
    });
});

