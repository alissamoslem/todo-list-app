
const registerForm = document.querySelector("#register-form");
const errorMessage = document.querySelector("#error-message");

registerForm.addEventListener("submit", async (event) => {
    event.preventDefault();

    const name = document.querySelector("#name").value.trim();
    const email = document.querySelector("#email").value.trim();
    const password = document.querySelector("#password").value;

    errorMessage.textContent = "";

    if (password.length < 8) {
        errorMessage.textContent =
            "Password must be at least 8 characters";
        return;
    }

    try {
        const response = await fetch("/api/auth/register", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                name: name,
                email: email,
                password: password
            })
        });

        if (!response.ok) {
            if (response.status === 409) {
                throw new Error("Email already registered");
            }

            if (response.status === 400) {
                throw new Error("Please check your registration details");
            }

            throw new Error("Registration failed. Please try again.");
        }

        // Registration successful
        window.location.href = "/login.html";

    } catch (error) {
        errorMessage.textContent = error.message;
    }
});
