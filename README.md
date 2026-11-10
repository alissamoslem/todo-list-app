# To-Do List App

A full-stack to-do list web app. The page is plain HTML, CSS and JavaScript, the backend is written in Go, and the tasks are stored in a MySQL database, so they survive page refreshes and server restarts.

## Features

- Add one task at a time, or several at once (one per line; blank lines are skipped)
- Mark tasks as complete or incomplete
- Delete tasks
- Priority for each task, from 1 (highest) to 3 (lowest), which can be changed after the task is created
- Live counters for total, remaining and completed tasks
- Responsive layout for desktop and mobile
- Accessibility support: screen-reader labels, ARIA live regions, keyboard focus styles and reduced-motion support

## Tech stack

| Part | Technology |
|------|------------|
| Frontend | HTML, CSS, JavaScript (no frameworks) |
| Backend | Go 1.22 or newer, standard library `net/http` |
| Database | MySQL 8.4 or newer |
| Database driver | `github.com/go-sql-driver/mysql` |

## Requirements

- [Go](https://go.dev/dl/) 1.22 or newer
- [MySQL Server](https://dev.mysql.com/downloads/mysql/) 8.4 or newer, running on port 3306
- Optional: MySQL Workbench, to run the SQL script and look inside the database

## Project structure

| File | Purpose |
|------|---------|
| `main.go` | Go server: REST API and database access |
| `schema.sql` | Creates the `todo_app` database and the `tasks` table |
| `go.mod`, `go.sum` | Go module and dependency files |
| `static/index.html` | Page structure and forms |
| `static/style.css` | Layout, colours, animations and responsive rules |
| `static/script.js` | Browser logic: sends requests to the API and updates the page |
| `.env.example` | Template for the settings file. Copy it to `.env` and fill in your values |

## Database

The database is called `todo_app` and has one table, `tasks`:

| Column | Type | Description |
|--------|------|-------------|
| `id` | `INT AUTO_INCREMENT PRIMARY KEY` | Unique task number, assigned by MySQL |
| `title` | `VARCHAR(255) NOT NULL` | The task text (255 characters at most) |
| `is_completed` | `BOOLEAN NOT NULL DEFAULT FALSE` | Whether the task is done |
| `priority` | `TINYINT NOT NULL DEFAULT 3` | 1 = High, 2 = Medium, 3 = Low (checked by the database) |
| `created_at` | `TIMESTAMP` | Set automatically when the task is created |
| `updated_at` | `TIMESTAMP` | Refreshed automatically by MySQL whenever the task changes |

## Getting started

### 1. Clone the repository

```bash
git clone https://github.com/alissamoslem/todo-list-app.git
cd My_to_do_list
```

### 2. Create the database

Make sure the MySQL server is running, then run `schema.sql`.

- **MySQL Workbench:** open a connection, open `schema.sql` from the File menu, and run the whole script.
- **Terminal:**
  ```bash
  mysql -u root -p -e "source schema.sql"
  ```

Check that it worked:

```sql
USE todo_app;
DESCRIBE tasks;
```

The script uses `CREATE DATABASE`, so running it a second time fails because `todo_app` already exists. To start again, run `DROP DATABASE todo_app;` first (this deletes all tasks).

### 3. Configure the database login

Copy the example settings file and put your MySQL password in it:

```bash
cp .env.example .env
```

On PowerShell use `Copy-Item .env.example .env`. Then open `.env` and set `DB_PASSWORD`. The `.env` file is ignored by Git, so your password is never committed.

**PowerShell:**

```powershell
$env:DB_USER = "root"
$env:DB_PASSWORD = "your-mysql-password"
```

**Git Bash, Linux or macOS:**

```bash
export DB_USER=root
export DB_PASSWORD='your-mysql-password'
```

These variables only exist in the terminal window where you set them. Set them again in every new terminal.

### 4. Install the dependencies and run

```bash
go mod tidy
go run .
```

Open http://localhost:8080 in your browser. Press `Ctrl + C` in the terminal to stop the server.

## API

All requests and responses use JSON. A task looks like this:

```json
{ "id": 1, "text": "Buy milk", "completed": false, "priority": 3 }
```

| Method | Route | Body | Description |
|--------|-------|------|-------------|
| GET | `/api/tasks` | none | List all tasks |
| POST | `/api/tasks` | `{"text": "Buy milk"}` | Create a task (priority starts at 3) |
| PATCH | `/api/tasks/{id}` | `{"completed": true}` and/or `{"priority": 1}` | Update a task's completed state or priority |
| DELETE | `/api/tasks/{id}` | none | Delete a task |

Status codes: `200` OK, `201` created, `204` success with no content, `400` invalid input (for example an empty task or a priority outside 1 to 3), `404` task not found, `500` database error.

## Security notes

- All SQL queries use parameterised placeholders (`?`), never string concatenation, which prevents SQL injection.
- The database password comes from an environment variable and is not stored in the repository.
- Task ids are validated as numbers, and priority is checked in the Go code and again by the database.
- The app has no login, so it is meant for local use and is not ready for the public internet.

## Troubleshooting

| Message | Cause and fix |
|---------|---------------|
| `Access denied for user 'root'@'localhost' (using password: NO)` | `.env` is missing or `DB_PASSWORD` is empty. Copy `.env.example` to `.env` and fill it in. |terminal. Set them again and rerun. |
| `Access denied for user 'root'@'localhost' (using password: YES)` | The password is wrong. Use the root password you chose when installing MySQL. |
| `connection refused` or error `10061` | The MySQL server is not running. Check the MySQL service in Windows Services and start it. |
| `Unknown database 'todo_app'` | `schema.sql` has not been run yet. |
| `address already in use` | Another program is using port 8080. Stop the old server or close the program using it. |

## Limitations

GitHub Pages cannot host this project, because it needs the Go server and a MySQL database. Tasks cannot be edited after they are created, there are no user accounts, and there are no automated tests.

## Possible improvements

- Edit existing tasks
- Filter by all, active and completed
- User accounts and login
- Automated tests
- Docker setup for one-command start-up
