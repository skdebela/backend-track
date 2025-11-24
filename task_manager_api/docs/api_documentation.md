# Task Manager API Documentation

## Base URL
`http://localhost:8080`

## Endpoints

### GET /tasks
List all tasks

### GET /tasks/:id
Get task by ID

### POST /tasks
Create new task
```json
{
  "title": "Learn Go",
  "description": "Study concurrency",
  "due_date": "2025-12-01T10:00:00Z",
  "status": "Pending"
}
...

## Environment variables
- `MONGO_URI` - MongoDB connection URI (default: mongodb://localhost:27017)
- `MONGO_DB` - Mongo database name (default: task_manager_db)