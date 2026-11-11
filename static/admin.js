const deleteModal = document.querySelector("#delete-modal");
const deleteWarning = document.querySelector("#delete-warning");
const deleteConfirmation = document.querySelector("#delete-confirmation");
const deleteError = document.querySelector("#delete-error");
const confirmDeleteButton = document.querySelector("#confirm-delete");
const cancelDeleteButton = document.querySelector("#cancel-delete");

let selectedUser = null;
const token = localStorage.getItem("token");

if (!token) {
    window.location.replace("/login.html");
} else {
    loadAdminData();
}

async function adminAPI(url, options = {}) {
    const response = await fetch(url, {
        ...options,
        headers: {
            "Authorization": `Bearer ${token}`,
            "Content-Type": "application/json",
            ...options.headers
        }
    });

    if (response.status === 401) {
        window.location.replace("/login.html");
        throw new Error("Please log in");
    }

    if (response.status === 403) {
        throw new Error("Access forbidden");
    }

    if (!response.ok) {
        throw new Error(`Request failed: ${response.status}`);
    }

    if (response.status === 204) {
        return null;
    }

    return response.json();
}

function showMessage(message) {
    document.querySelector("#admin-message").textContent = message;
}

async function loadAdminData() {
    try {
        const [currentUser, users, tasks] = await Promise.all([
            adminAPI("/api/auth/me"),
            adminAPI("/api/admin/users"),
            adminAPI("/api/admin/tasks")
        ]);

        displayUsers(users, currentUser.id);
        displayTasks(tasks);
    } catch (error) {
        showMessage(error.message);
    }
}

function displayUsers(users, currentAdminID) {
    const container = document.querySelector("#users-list");
    container.replaceChildren();

    users.forEach(user => {
        const row = document.createElement("div");
        row.className = "admin-user-row";

        const details = document.createElement("span");
        details.textContent =
            `${user.name} (${user.email}) — ${user.role}`;

        row.appendChild(details);

        // Admins cannot delete their own accounts.
        if (user.id !== currentAdminID) {
            const deleteButton = document.createElement("button");

            deleteButton.type = "button";
            deleteButton.textContent = "Delete User";
            deleteButton.className = "delete-user-btn";

            deleteButton.addEventListener("click", async () => {
                await deleteUser(user.id, user.name);
            });

            row.appendChild(deleteButton);
        }

        container.appendChild(row);
    });
}

function displayTasks(tasks) {
    const container = document.querySelector("#admin-tasks-list");
    container.replaceChildren();

    tasks.forEach(task => {
        const row = document.createElement("p");

        row.textContent =
            `${task.text} — Owner: ${task.owner_name || "Unassigned"}`;

        container.appendChild(row);
    });
}


function deleteUser(userId, userName) {
    // Remember which user the admin selected
    selectedUser = {
        id: userId,
        name: userName
    };

    // Reset the confirmation dialog
    deleteConfirmation.value = "";
    deleteError.textContent = "";
    confirmDeleteButton.disabled = true;

    // Show the selected user's name
    deleteWarning.textContent =
        `Are you sure you want to delete ${userName}?`;

    // Display the modal
    deleteModal.hidden = false;
    deleteConfirmation.focus();
}

// Enable deletion only when DELETE is typed exactly
deleteConfirmation.addEventListener("input", () => {
    confirmDeleteButton.disabled =
        deleteConfirmation.value !== "DELETE";
});

// Cancel deletion
cancelDeleteButton.addEventListener("click", () => {
    closeDeleteModal();
});

function closeDeleteModal() {
    deleteModal.hidden = true;
    selectedUser = null;
    deleteConfirmation.value = "";
    deleteError.textContent = "";
}

// Confirm deletion
confirmDeleteButton.addEventListener("click", async () => {
    if (!selectedUser || deleteConfirmation.value !== "DELETE") {
        return;
    }

    const userId = selectedUser.id;
    const userName = selectedUser.name;

    confirmDeleteButton.disabled = true;
    cancelDeleteButton.disabled = true;

    try {
        const response = await fetch(
            `/api/admin/users/${userId}`,
            {
                method: "DELETE",
                headers: {
                    "Authorization": `Bearer ${token}`,
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({
                    confirmation: "DELETE"
                })
            }
        );

        if (!response.ok) {
            const message = await response.text();
            throw new Error(`${response.status}: ${message}`);
        }

        closeDeleteModal();

        showMessage(`${userName} deleted successfully.`);
        await loadAdminData();

    } catch (error) {
        deleteError.textContent =
            `Deletion failed: ${error.message}`;
    } finally {
        cancelDeleteButton.disabled = false;
        confirmDeleteButton.disabled =
            deleteConfirmation.value !== "DELETE";
    }
});

