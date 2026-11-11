
const loginForm = document.querySelector("#login-form");
const errorMessage = document.querySelector("#error-message");
const submitButton = loginForm.querySelector('button[type="submit"]');

function showError(message) {
    errorMessage.textContent = message;
    errorMessage.hidden = false;
}

function clearError() {
    errorMessage.textContent = "";
    errorMessage.hidden = true;
}

loginForm.addEventListener("submit", async (event) => {
    event.preventDefault();
    clearError();

    const email = document.querySelector("#email").value.trim();
    const password = document.querySelector("#password").value;

    submitButton.disabled = true;
    submitButton.textContent = "Logging in...";

    try {
        const response = await fetch("/api/auth/login", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                email,
                password
            })
        });

        if (!response.ok) {
            if (response.status === 401) {
                throw new Error("Invalid email or password.");
            }

            throw new Error("Login failed. Please try again.");
        }

        const data = await response.json();

        localStorage.setItem("token", data.token);

        if (data.user.role === "admin") {
            window.location.href = "/admin.html";
        } else {
            window.location.href = "/";
        }

    } catch (error) {
        showError(error.message);
    } finally {
        submitButton.disabled = false;
        submitButton.textContent = "Log In";
    }
});
