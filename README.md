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

| Method | Route | Description |
|--------|-------|-------------|
| GET | `/api/tasks` | List all tasks |
| POST | `/api/tasks` | Create a task |
| PATCH | `/api/tasks/{id}` | Mark a task complete or incomplete |
| DELETE | `/api/tasks/{id}` | Delete a task |

## Possible Improvements

- Edit existing tasks
- Filter by all / active / completed
- Task due dates
