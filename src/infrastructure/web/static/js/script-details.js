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
