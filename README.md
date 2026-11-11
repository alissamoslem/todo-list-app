# To-Do List App

A full-stack To-Do List web application built with **Go, MySQL, HTML, CSS, and JavaScript**.

Users can register, log in, create and manage their own tasks, set priorities, and track their progress. The application also includes an admin dashboard for managing users and viewing tasks.

The frontend, Go backend, and MySQL database are deployed using **Railway**.

## Live Demo

**Website:** https://moslem-todo-list-app.up.railway.app

**GitHub:** https://github.com/alissamoslem/todo-list-app

## Features

### Authentication and Authorization

- User registration and login
- Password hashing using bcrypt
- JWT-based authentication
- Protected API endpoints
- User-specific task management
- Role-based authorization for administrators
- Logout functionality

### Task Management

- Add individual tasks
- Add multiple tasks at once, one per line
- Mark tasks as completed or incomplete
- Delete tasks
- Assign priorities from 1 (highest) to 3 (lowest)
- Change task priorities
- Sort tasks by priority
- Display live counters for total, remaining, and completed tasks

### Admin Dashboard

- View registered users
- View tasks across users
- Delete users through protected admin endpoints
- Restrict administrative actions to accounts with the `admin` role

### User Interface

- Responsive design for desktop and mobile
- Separate registration and login pages
- Task priority colors
- Admin navigation for authorized users
- Accessibility features, including labels and keyboard focus styles

## Tech Stack

| Component | Technology |
|---|---|
| Frontend | HTML, CSS, JavaScript |
| Backend | Go 1.22+, `net/http` |
| Database | MySQL |
| Database driver | `go-sql-driver/mysql` |
| Authentication | JWT (`golang-jwt/jwt/v5`) |
| Password security | bcrypt |
| Environment variables | `godotenv` |
| Hosting | Railway |
| Version control | Git and GitHub |

## Project Structure

```text
todo-list-app/
├── main.go
├── auth.go
├── middleware.go
├── admin.go
├── cors.go
├── schema.sql
├── database/
│   └── 002_auth.sql
├── static/
│   ├── index.html
│   ├── login.html
│   ├── register.html
│   ├── admin.html
│   ├── script.js
│   ├── login.js
│   ├── register.js
│   ├── admin.js
│   └── style.css
├── .env.example
├── go.mod
├── go.sum
└── README.md
```

### Main Files

| File | Purpose |
|---|---|
| `main.go` | Starts the HTTP server, connects to MySQL, and handles task routes |
| `auth.go` | Implements registration, login, and authentication-related handlers |
| `middleware.go` | Verifies JWTs and checks admin permissions |
| `admin.go` | Implements administrative API functionality |
| `cors.go` | Configures allowed browser origins and CORS preflight requests |
| `schema.sql` | Database initialization |
| `database/002_auth.sql` | Database migration related to authentication |
| `static/` | Frontend pages, JavaScript, and CSS |

## Requirements

To run the application locally, install:

- Go 1.22 or newer
- MySQL Server
- Git
- A web browser

MySQL Workbench is optional but useful for managing the database.

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/alissamoslem/todo-list-app.git
cd todo-list-app
```

# To-Do List App

A full-stack To-Do List web application built with **Go, MySQL, HTML, CSS, and JavaScript**.

Users can register, log in, create and manage their own tasks, set priorities, and track their progress. The application also includes an admin dashboard for managing users and viewing tasks.

The frontend, Go backend, and MySQL database are deployed using **Railway**.

## Live Demo

**Website:** https://moslem-todo-list-app.up.railway.app

**GitHub:** https://github.com/alissamoslem/todo-list-app

## Features

### Authentication and Authorization

- User registration and login
- Password hashing using bcrypt
- JWT-based authentication
- Protected API endpoints
- User-specific task management
- Role-based authorization for administrators
- Logout functionality

### Task Management

- Add individual tasks
- Add multiple tasks at once, one per line
- Mark tasks as completed or incomplete
- Delete tasks
- Assign priorities from 1 (highest) to 3 (lowest)
- Change task priorities
- Sort tasks by priority
- Display live counters for total, remaining, and completed tasks

### Admin Dashboard

- View registered users
- View tasks across users
- Delete users through protected admin endpoints
- Restrict administrative actions to accounts with the `admin` role

### User Interface

- Responsive design for desktop and mobile
- Separate registration and login pages
- Task priority colors
- Admin navigation for authorized users
- Accessibility features, including labels and keyboard focus styles

## Tech Stack

| Component | Technology |
|---|---|
| Frontend | HTML, CSS, JavaScript |
| Backend | Go 1.22+, `net/http` |
| Database | MySQL |
| Database driver | `go-sql-driver/mysql` |
| Authentication | JWT (`golang-jwt/jwt/v5`) |
| Password security | bcrypt |
| Environment variables | `godotenv` |
| Hosting | Railway |
| Version control | Git and GitHub |

## Project Structure

```text
todo-list-app/
├── main.go
├── auth.go
├── middleware.go
├── admin.go
├── cors.go
├── schema.sql
├── database/
│   └── 002_auth.sql
├── static/
│   ├── index.html
│   ├── login.html
│   ├── register.html
│   ├── admin.html
│   ├── script.js
│   ├── login.js
│   ├── register.js
│   ├── admin.js
│   └── style.css
├── .env.example
├── go.mod
├── go.sum
└── README.md
```

### Main Files

| File | Purpose |
|---|---|
| `main.go` | Starts the HTTP server, connects to MySQL, and handles task routes |
| `auth.go` | Implements registration, login, and authentication-related handlers |
| `middleware.go` | Verifies JWTs and checks admin permissions |
| `admin.go` | Implements administrative API functionality |
| `cors.go` | Configures allowed browser origins and CORS preflight requests |
| `schema.sql` | Database initialization |
| `database/002_auth.sql` | Database migration related to authentication |
| `static/` | Frontend pages, JavaScript, and CSS |

## Requirements

To run the application locally, install:

- Go 1.22 or newer
- MySQL Server
- Git
- A web browser

MySQL Workbench is optional but useful for managing the database.

## Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/alissamoslem/todo-list-app.git
cd todo-list-app
```

### 2. Set Up MySQL

Create the local database using the project's SQL scripts.

The database schema includes two main tables:

- `users` — stores user accounts, password hashes, and roles.
- `tasks` — stores tasks and associates them with their owners.

For a fresh database, inspect and execute `schema.sql`, then apply `database/002_auth.sql` if its changes are not already included in the initialized schema.

Use MySQL Workbench or the MySQL command-line client.

**Important:** Check the existing schema and migration state before executing scripts against a database containing data. Do not drop or recreate an existing database without a backup.

### 3. Configure Environment Variables

Copy the example environment file:

**PowerShell:**

```powershell
Copy-Item .env.example .env
```

**Linux/macOS:**

```bash
cp .env.example .env
```

Configure `.env` with your local settings:

```dotenv
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_mysql_password
DB_NAME=todo_app

JWT_SECRET=replace_with_a_long_random_secret

FRONTEND_ORIGIN=http://localhost:8080
PORT=8080
```

| Variable | Description |
|---|---|
| `DB_HOST` | MySQL server hostname |
| `DB_PORT` | MySQL server port |
| `DB_USER` | MySQL username |
| `DB_PASSWORD` | MySQL password |
| `DB_NAME` | Database name |
| `JWT_SECRET` | Secret used to sign and verify JWTs |
| `FRONTEND_ORIGIN` | Browser origin allowed by the CORS middleware |
| `PORT` | HTTP server port, defaulting to `8080` |

Generate a strong random value for `JWT_SECRET`. Never commit your actual `.env` file or production credentials to GitHub.

### 4. Install Dependencies

```bash
go mod download
```

### 5. Run the Application

```bash
go run .
```

Open:

```text
http://localhost:8080
```

The main pages are:

| Page | URL |
|---|---|
| To-Do List | `/` |
| Login | `/login.html` |
| Registration | `/register.html` |
| Admin Dashboard | `/admin.html` |

Stop the server using `Ctrl+C`.

## Database Design

### Users Table

The `users` table stores account information including:

- User ID
- Name
- Email
- Password hash
- Role (`user` or `admin`)
- Creation timestamp

Passwords are stored as hashes rather than plaintext.

### Tasks Table

The `tasks` table includes:

| Column | Description |
|---|---|
| `id` | Unique task ID |
| `title` | Task text |
| `is_completed` | Completion status |
| `priority` | Priority level from 1 to 3 |
| `user_id` | ID of the user who owns the task |
| `created_at` | Task creation timestamp |
| `updated_at` | Last update timestamp |

The `user_id` column links each task to an account.

Regular users can access and modify only their own tasks through the task API.

## API Endpoints

### Authentication

| Method | Endpoint | Description |
|---|---|---|
| POST | `/api/auth/register` | Register a user |
| POST | `/api/auth/login` | Log in and receive a JWT |
| GET | `/api/auth/me` | Get authenticated user information |

### Tasks

These endpoints require authentication.

| Method | Endpoint | Description |
|---|---|---|
| GET | `/api/tasks` | Get the authenticated user's tasks |
| POST | `/api/tasks` | Create a task |
| PATCH | `/api/tasks/{id}` | Update completion or priority |
| DELETE | `/api/tasks/{id}` | Delete an owned task |

Example task:

```json
{
  "id": 1,
  "text": "Learn Go",
  "completed": false,
  "priority": 3
}
```

Create a task:

```http
POST /api/tasks
Content-Type: application/json
Authorization: Bearer YOUR_JWT_TOKEN
```

```json
{
  "text": "Learn Go"
}
```

Update a task:

```http
PATCH /api/tasks/1
Content-Type: application/json
Authorization: Bearer YOUR_JWT_TOKEN
```

```json
{
  "completed": true,
  "priority": 1
}
```

Bulk creation is supported by the frontend, which processes multiple task entries using the task creation API.

### Admin Endpoints

These endpoints require an authenticated account with the `admin` role.

| Method | Endpoint | Description |
|---|---|---|
| GET | `/api/admin/users` | List registered users |
| GET | `/api/admin/tasks` | View tasks across users |
| DELETE | `/api/admin/users/{id}` | Delete a user |

Public registration does not grant administrative privileges. Admin accounts must be provisioned through a trusted administrative process.

### Common HTTP Status Codes

| Code | Meaning |
|---|---|
| `200` | Successful request |
| `201` | Resource created |
| `204` | Successful request with no response body |
| `400` | Invalid input |
| `401` | Missing or invalid authentication |
| `403` | Forbidden action |
| `404` | Resource not found |
| `500` | Internal server error |

## How It Is Deployed

The application is deployed using **Railway**.

### Backend and Frontend

The Go application runs as a Railway service.

It uses `net/http` for routing and serves the HTML, CSS, and JavaScript files from the `static` directory.

The frontend and API are available under the same application domain.

### MySQL Database

A separate Railway MySQL service stores user accounts and tasks.

The Go application connects using Railway environment variables rather than hardcoded database credentials.

### Environment Configuration

The deployed application uses variables such as:

```text
DB_HOST
DB_PORT
DB_USER
DB_PASSWORD
DB_NAME
JWT_SECRET
FRONTEND_ORIGIN
PORT
```

These are configured in the Railway service's Variables settings.

Production secrets are not stored in the source code.

### HTTPS

The application is served through a Railway-provided HTTPS domain:

https://moslem-todo-list-app.up.railway.app

Railway manages HTTPS termination and certificate provisioning for its hosted domains.

### CORS

The backend uses restricted CORS middleware.

The allowed browser origin is configured through `FRONTEND_ORIGIN`.

In production:

```dotenv
FRONTEND_ORIGIN=https://moslem-todo-list-app.up.railway.app
```

The middleware handles CORS preflight requests and returns `Access-Control-Allow-Origin` only for the configured origin.

Unknown origins are not permitted by the CORS policy.

CORS is not a replacement for authentication or authorization; JWT verification and user permissions are enforced separately.

### Deployment Workflow

1. Develop and test features locally.
2. Create a Git feature branch.
3. Commit and push changes to GitHub.
4. Open and review a Pull Request.
5. Merge changes into `main`.
6. Railway builds and deploys the connected branch when automatic deployment is enabled.
7. Test the deployed application.

## Security

The application implements several security measures:

- Password hashing with bcrypt
- JWT authentication for protected endpoints
- Admin role verification
- User-specific task ownership checks
- Parameterized SQL queries to reduce SQL injection risk
- Environment variables for sensitive configuration
- Restricted CORS policy
- HTTPS for deployed traffic
- Validation of task IDs and priority values

The app currently stores authentication tokens in browser local storage. Moving tokens to appropriately configured secure, HttpOnly cookies and adding further browser-side protections are possible future improvements.

## Troubleshooting

| Problem | Possible solution |
|---|---|
| MySQL access denied | Check `DB_USER` and `DB_PASSWORD` |
| Database connection refused | Verify MySQL is running and check `DB_HOST` and `DB_PORT` |
| Unknown database | Check `DB_NAME` and database initialization |
| Missing JWT secret | Set `JWT_SECRET` |
| API returns `401` | Log in again and verify the JWT |
| API returns `403` | Check user permissions or the configured CORS origin |
| Admin Dashboard unavailable | Verify the user's role |
| CORS preflight fails | Check `FRONTEND_ORIGIN` and the deployed middleware |
| Port already in use | Stop the other server or change `PORT` |
| Railway deployment fails | Review Railway deployment logs and environment variables |

## Limitations

- No automated test suite is currently documented
- Task descriptions cannot be edited after creation
- No password reset or email verification workflow
- Additional security hardening is recommended before handling sensitive production data

## Possible Improvements

- Edit task descriptions
- Filter tasks by completed, incomplete, or all
- Search tasks
- Add due dates and reminders
- Add automated Go tests
- Implement password reset
- Improve authentication session management
- Add Docker support
- Add CI checks with GitHub Actions

## Author

Developed by [alissamoslem](https://github.com/alissamoslem) as a Go full-stack development training project.

### 3. Configure Environment Variables

Copy the example environment file:

**PowerShell:**

```powershell
Copy-Item .env.example .env
```

**Linux/macOS:**

```bash
cp .env.example .env
```

Configure `.env` with your local settings:

```dotenv
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_mysql_password
DB_NAME=todo_app

JWT_SECRET=replace_with_a_long_random_secret

FRONTEND_ORIGIN=http://localhost:8080
PORT=8080
```

| Variable | Description |
|---|---|
| `DB_HOST` | MySQL server hostname |
| `DB_PORT` | MySQL server port |
| `DB_USER` | MySQL username |
| `DB_PASSWORD` | MySQL password |
| `DB_NAME` | Database name |
| `JWT_SECRET` | Secret used to sign and verify JWTs |
| `FRONTEND_ORIGIN` | Browser origin allowed by the CORS middleware |
| `PORT` | HTTP server port, defaulting to `8080` |

Generate a strong random value for `JWT_SECRET`. Never commit your actual `.env` file or production credentials to GitHub.

### 4. Install Dependencies

```bash
go mod download
```

### 5. Run the Application

```bash
go run .
```

Open:

```text
http://localhost:8080
```

The main pages are:

| Page | URL |
|---|---|
| To-Do List | `/` |
| Login | `/login.html` |
| Registration | `/register.html` |
| Admin Dashboard | `/admin.html` |

Stop the server using `Ctrl+C`.

## Database Design

### Users Table

The `users` table stores account information including:

- User ID
- Name
- Email
- Password hash
- Role (`user` or `admin`)
- Creation timestamp

Passwords are stored as hashes rather than plaintext.

### Tasks Table

The `tasks` table includes:

| Column | Description |
|---|---|
| `id` | Unique task ID |
| `title` | Task text |
| `is_completed` | Completion status |
| `priority` | Priority level from 1 to 3 |
| `user_id` | ID of the user who owns the task |
| `created_at` | Task creation timestamp |
| `updated_at` | Last update timestamp |

The `user_id` column links each task to an account.

Regular users can access and modify only their own tasks through the task API.

## API Endpoints

### Authentication

| Method | Endpoint | Description |
|---|---|---|
| POST | `/api/auth/register` | Register a user |
| POST | `/api/auth/login` | Log in and receive a JWT |
| GET | `/api/auth/me` | Get authenticated user information |

### Tasks

These endpoints require authentication.

| Method | Endpoint | Description |
|---|---|---|
| GET | `/api/tasks` | Get the authenticated user's tasks |
| POST | `/api/tasks` | Create a task |
| PATCH | `/api/tasks/{id}` | Update completion or priority |
| DELETE | `/api/tasks/{id}` | Delete an owned task |

Example task:

```json
{
  "id": 1,
  "text": "Learn Go",
  "completed": false,
  "priority": 3
}
```

Create a task:

```http
POST /api/tasks
Content-Type: application/json
Authorization: Bearer YOUR_JWT_TOKEN
```

```json
{
  "text": "Learn Go"
}
```

Update a task:

```http
PATCH /api/tasks/1
Content-Type: application/json
Authorization: Bearer YOUR_JWT_TOKEN
```

```json
{
  "completed": true,
  "priority": 1
}
```

Bulk creation is supported by the frontend, which processes multiple task entries using the task creation API.

### Admin Endpoints

These endpoints require an authenticated account with the `admin` role.

| Method | Endpoint | Description |
|---|---|---|
| GET | `/api/admin/users` | List registered users |
| GET | `/api/admin/tasks` | View tasks across users |
| DELETE | `/api/admin/users/{id}` | Delete a user |

Public registration does not grant administrative privileges. Admin accounts must be provisioned through a trusted administrative process.

### Common HTTP Status Codes

| Code | Meaning |
|---|---|
| `200` | Successful request |
| `201` | Resource created |
| `204` | Successful request with no response body |
| `400` | Invalid input |
| `401` | Missing or invalid authentication |
| `403` | Forbidden action |
| `404` | Resource not found |
| `500` | Internal server error |

## How It Is Deployed

The application is deployed using **Railway**.

### Backend and Frontend

The Go application runs as a Railway service.

It uses `net/http` for routing and serves the HTML, CSS, and JavaScript files from the `static` directory.

The frontend and API are available under the same application domain.

### MySQL Database

A separate Railway MySQL service stores user accounts and tasks.

The Go application connects using Railway environment variables rather than hardcoded database credentials.

### Environment Configuration

The deployed application uses variables such as:

```text
DB_HOST
DB_PORT
DB_USER
DB_PASSWORD
DB_NAME
JWT_SECRET
FRONTEND_ORIGIN
PORT
```

These are configured in the Railway service's Variables settings.

Production secrets are not stored in the source code.

### HTTPS

The application is served through a Railway-provided HTTPS domain:

https://moslem-todo-list-app.up.railway.app

Railway manages HTTPS termination and certificate provisioning for its hosted domains.

### CORS

The backend uses restricted CORS middleware.

The allowed browser origin is configured through `FRONTEND_ORIGIN`.

In production:

```dotenv
FRONTEND_ORIGIN=https://moslem-todo-list-app.up.railway.app
```

The middleware handles CORS preflight requests and returns `Access-Control-Allow-Origin` only for the configured origin.

Unknown origins are not permitted by the CORS policy.

CORS is not a replacement for authentication or authorization; JWT verification and user permissions are enforced separately.

### Deployment Workflow

1. Develop and test features locally.
2. Create a Git feature branch.
3. Commit and push changes to GitHub.
4. Open and review a Pull Request.
5. Merge changes into `main`.
6. Railway builds and deploys the connected branch when automatic deployment is enabled.
7. Test the deployed application.

## Security

The application implements several security measures:

- Password hashing with bcrypt
- JWT authentication for protected endpoints
- Admin role verification
- User-specific task ownership checks
- Parameterized SQL queries to reduce SQL injection risk
- Environment variables for sensitive configuration
- Restricted CORS policy
- HTTPS for deployed traffic
- Validation of task IDs and priority values

The app currently stores authentication tokens in browser local storage. Moving tokens to appropriately configured secure, HttpOnly cookies and adding further browser-side protections are possible future improvements.

## Troubleshooting

| Problem | Possible solution |
|---|---|
| MySQL access denied | Check `DB_USER` and `DB_PASSWORD` |
| Database connection refused | Verify MySQL is running and check `DB_HOST` and `DB_PORT` |
| Unknown database | Check `DB_NAME` and database initialization |
| Missing JWT secret | Set `JWT_SECRET` |
| API returns `401` | Log in again and verify the JWT |
| API returns `403` | Check user permissions or the configured CORS origin |
| Admin Dashboard unavailable | Verify the user's role |
| CORS preflight fails | Check `FRONTEND_ORIGIN` and the deployed middleware |
| Port already in use | Stop the other server or change `PORT` |
| Railway deployment fails | Review Railway deployment logs and environment variables |

## Limitations

- No automated test suite is currently documented
- Task descriptions cannot be edited after creation
- No password reset or email verification workflow
- Additional security hardening is recommended before handling sensitive production data

## Possible Improvements

- Edit task descriptions
- Filter tasks by completed, incomplete, or all
- Search tasks
- Add due dates and reminders
- Add automated Go tests
- Implement password reset
- Improve authentication session management
- Add Docker support
- Add CI checks with GitHub Actions

## Author

Developed by [alissamoslem](https://github.com/alissamoslem) as a Go full-stack development training project.
