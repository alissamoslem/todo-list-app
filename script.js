const taskForm = document.querySelector("#task-form");
const taskInput = document.querySelector("#task-input");
const bulkForm = document.querySelector("#bulk-form");
const bulkInput = document.querySelector("#bulk-input");
const taskList = document.querySelector("#task-list");
const emptyState = document.querySelector("#empty-state");
const taskCount = document.querySelector("#task-count");
const tasksRemainingCount = document.querySelector("#tasks-remaining-count");
const tasksCompletedCount = document.querySelector("#tasks-completed-count");

function updateSummary() {
  const total = taskList.children.length;
  const completed = taskList.querySelectorAll(".task-checkbox:checked").length;
  const remaining = total - completed;
  taskCount.textContent = total === 1 ? "1 task" : `${total} tasks`;
  tasksRemainingCount.textContent = remaining === 1 ? "1 task remaining" : `${remaining} tasks remaining`;
  tasksCompletedCount.textContent = completed === 1 ? "1 task completed" : `${completed} tasks completed`;
  emptyState.hidden = total > 0;
  if (total > 0) {
    taskCount.setAttribute("aria-label", `${total} tasks, ${completed} completed`);
  } else {
    taskCount.setAttribute("aria-label", "No tasks");
  }
}

function addTask(value) {
  const text = value.trim();
  if (!text) return false;

  const item = document.createElement("li");
  item.className = "task-item";

  const checkbox = document.createElement("input");
  checkbox.type = "checkbox";
  checkbox.className = "task-checkbox";
  checkbox.setAttribute("aria-label", `Mark “${text}” as complete`);

  const label = document.createElement("span");
  label.className = "task-text";
  label.textContent = text;

  const deleteButton = document.createElement("button");
  deleteButton.type = "button";
  deleteButton.className = "delete-button";
  deleteButton.textContent = "×";
  deleteButton.setAttribute("aria-label", `Delete “${text}”`);

  checkbox.addEventListener("change", () => {
    item.classList.toggle("is-complete", checkbox.checked);
    checkbox.setAttribute("aria-label", `Mark “${text}” as ${checkbox.checked ? "incomplete" : "complete"}`);
    updateSummary();
  });
  deleteButton.addEventListener("click", () => {
    item.remove();
    updateSummary();
  });

  item.append(checkbox, label, deleteButton);
  taskList.append(item);
  updateSummary();
  return true;
}

taskForm.addEventListener("submit", (event) => {
  event.preventDefault();
  if (addTask(taskInput.value)) {
    taskInput.value = "";
    taskInput.focus();
  } else {
    taskInput.focus();
  }
});

bulkForm.addEventListener("submit", (event) => {
  event.preventDefault();
  const tasks = bulkInput.value.split(/\r?\n/).map((task) => task.trim()).filter(Boolean);
  tasks.forEach(addTask);
  if (tasks.length > 0) bulkInput.value = "";
});

updateSummary();
