
const loginForm = document.querySelector("#login-form");
const errorMessage = document.querySelector("#error-message");

loginForm.addEventListener("submit", async (event) => {
    event.preventDefault();

    const email = document.querySelector("#email").value.trim();
    const password = document.querySelector("#password").value;

    errorMessage.textContent = "";

    try {
        const response = await fetch("/api/auth/login", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                email: email,
                password: password
            })
        });

        if (!response.ok) {
            if (response.status === 401) {
                throw new Error("Invalid email or password");
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
        errorMessage.textContent = error.message;
    }
});
