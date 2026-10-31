const taskForm = document.querySelector("#task-form");
const taskInput = document.querySelector("#task-input");
const bulkForm = document.querySelector("#bulk-form");
const bulkInput = document.querySelector("#bulk-input");
const taskList = document.querySelector("#task-list");
const emptyState = document.querySelector("#empty-state");
const taskCount = document.querySelector("#task-count");
const tasksRemainingCount = document.querySelector("#tasks-remaining-count");
const tasksCompletedCount = document.querySelector("#tasks-completed-count");

// جديد: الكلمات التي تظهر للمستخدم (غيّرها هنا إن أردت، مثل High / Medium / Low)
const PRIORITY_NAMES = { 1: "Priority 1", 2: "Priority 2", 3: "Priority 3" };

async function api(url, options = {}) {
  const response = await fetch(url, {
    headers: { "Content-Type": "application/json" },
    ...options,
  });
  if (!response.ok) throw new Error(`Request failed: ${response.status}`);
  return response.status === 204 ? null : response.json();
}

function updateSummary() {
  const total = taskList.children.length;
  const completed = taskList.querySelectorAll(".task-checkbox:checked").length;
  const remaining = total - completed;
  taskCount.textContent = total === 1 ? "1 task" : `${total} tasks`;
  tasksRemainingCount.textContent = remaining === 1 ? "1 task remaining" : `${remaining} tasks remaining`;
  tasksCompletedCount.textContent = completed === 1 ? "1 task completed" : `${completed} tasks completed`;
  emptyState.hidden = total > 0;
  taskCount.setAttribute(
    "aria-label",
    total > 0 ? `${total} tasks, ${completed} completed` : "No tasks"
  );
}

function renderTask(task) {
  const item = document.createElement("li");
  item.className = "task-item";
  item.classList.toggle("is-complete", task.completed);

  const checkbox = document.createElement("input");
  checkbox.type = "checkbox";
  checkbox.className = "task-checkbox";
  checkbox.checked = task.completed;
  checkbox.setAttribute("aria-label", `Mark “${task.text}” as ${task.completed ? "incomplete" : "complete"}`);

  const label = document.createElement("span");
  label.className = "task-text";
  label.textContent = task.text;

  // جديد: قائمة منسدلة فيها الأولويات الثلاث مكتوبة بالكلمات
  let priority = task.priority || 3;
  const prioritySelect = document.createElement("select");
  prioritySelect.setAttribute("aria-label", `Priority of “${task.text}”`);
  for (const value of [1, 2, 3]) {
    const option = document.createElement("option");
    option.value = value;
    option.textContent = PRIORITY_NAMES[value];
    prioritySelect.append(option);
  }
  prioritySelect.value = priority;

  // جديد: تلوين القائمة حسب الأولوية الحالية
  function paintSelect() {
    prioritySelect.className = `priority-select priority-${priority}`;
  }
  paintSelect();

  const deleteButton = document.createElement("button");
  deleteButton.type = "button";
  deleteButton.className = "delete-button";
  deleteButton.textContent = "×";
  deleteButton.setAttribute("aria-label", `Delete “${task.text}”`);

  checkbox.addEventListener("change", async () => {
    try {
      await api(`/api/tasks/${task.id}`, {
        method: "PATCH",
        body: JSON.stringify({ completed: checkbox.checked }),
      });
      item.classList.toggle("is-complete", checkbox.checked);
      checkbox.setAttribute("aria-label", `Mark “${task.text}” as ${checkbox.checked ? "incomplete" : "complete"}`);
      updateSummary();
    } catch {
      checkbox.checked = !checkbox.checked;
      alert("Could not update the task.");
    }
  });

  // جديد: عند اختيار أولوية من القائمة نحفظها في الخادم
  prioritySelect.addEventListener("change", async () => {
    const chosen = Number(prioritySelect.value);
    try {
      await api(`/api/tasks/${task.id}`, {
        method: "PATCH",
        body: JSON.stringify({ priority: chosen }),
      });
      priority = chosen;
      paintSelect();
    } catch {
      prioritySelect.value = priority; // فشل الحفظ: نعيد القيمة القديمة
      alert("Could not change the priority.");
    }
  });

  deleteButton.addEventListener("click", async () => {
    try {
      await api(`/api/tasks/${task.id}`, { method: "DELETE" });
      item.remove();
      updateSummary();
    } catch {
      alert("Could not delete the task.");
    }
  });

  item.append(checkbox, label, prioritySelect, deleteButton); // جديد: prioritySelect
  taskList.append(item);
}

async function addTask(value) {
  const text = value.trim();
  if (!text) return false;
  try {
    const task = await api("/api/tasks", {
      method: "POST",
      body: JSON.stringify({ text }),
    });
    renderTask(task);
    updateSummary();
    return true;
  } catch {
    alert("Could not save the task.");
    return false;
  }
}

taskForm.addEventListener("submit", async (event) => {
  event.preventDefault();
  if (await addTask(taskInput.value)) taskInput.value = "";
  taskInput.focus();
});

bulkForm.addEventListener("submit", async (event) => {
  event.preventDefault();
  const tasks = bulkInput.value.split(/\r?\n/).map((t) => t.trim()).filter(Boolean);
  for (const text of tasks) await addTask(text); // one at a time keeps the order
  if (tasks.length > 0) bulkInput.value = "";
});

async function loadTasks() {
  try {
    const tasks = await api("/api/tasks");
    tasks.forEach(renderTask);
  } catch {
    alert("Could not load tasks. Is the Go server running?");
  }
  updateSummary();
}

loadTasks();