# To-Do List App

A to-do list web app with a Go backend and a SQLite database, so tasks are saved permanently.

## Features

- Add a single task, or several at once (one per line; blank lines are skipped)
- Mark tasks as complete or incomplete
- Delete tasks
- Live counters for total, remaining and completed tasks
- Tasks are stored in a SQLite database and survive page refreshes and server restarts
- Responsive layout for desktop and mobile
- Accessibility support: screen-reader labels, ARIA live regions, keyboard focus styles and reduced-motion support

## Requirements

- [Go](https://go.dev/dl/) 1.22 or newer

## Project Structure

| File | Purpose |
|------|---------|
| `main.go` | Go server: REST API and database access |
| `go.mod`, `go.sum` | Go module and dependency files |
| `static/index.html` | Page structure and forms |
| `static/style.css` | Layout, colours, animations and responsive rules |
| `static/script.js` | Browser logic: sends requests to the API and updates the page |
| `todos.db` | SQLite database (created automatically, not tracked by Git) |

## Getting Started

```
git clone https://github.com/Muslim03/My_to_do_list.git
cd My_to_do_list
go run .
```

Dependencies are downloaded automatically from `go.mod` and `go.sum` the first time you run the app.

Then open [http://localhost:8080](http://localhost:8080) in your browser. Press `Ctrl + C` in the terminal to stop the server.
## API

| Method | Route             | Description                                       |
|--------|-------------------|---------------------------------------------------|
| GET    | `/api/tasks`      | List all tasks (each with its priority)           |
| POST   | `/api/tasks`      | Create a task (starts with priority 3)            |
| PATCH  | `/api/tasks/{id}` | Update a task's completed state and/or priority   |
| DELETE | `/api/tasks/{id}` | Delete a task                                     |

### Examples

A task as returned by the API:

```json
{ "id": 1, "text": "Plan the week", "completed": false, "priority": 3 }
```

`POST /api/tasks` body:

```json
{ "text": "Plan the week" }
```

`PATCH /api/tasks/{id}` body (send either field or both):

```json
{ "completed": true }
{ "priority": 1 }
```

`priority` must be 1, 2 or 3 (1 is the highest). Any other value returns `400`.
## Possible Improvements

- Edit existing tasks
- Filter by all / active / completed
- Task due dates
