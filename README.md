# Go TODO Rest

## Task

You are tasked with building a RESTful API that allows users to manage a simple TODO list. Users should be able to add new tasks, view all tasks, update tasks, and delete tasks.

## Requirements

- Use the Go programming language to build the API.
- Use the net/http package to handle HTTP requests and responses.
- Use a data storage mechanism of your choice to store the TODO list items.
- The API should follow RESTful principles, including using appropriate HTTP methods (GET, POST, PUT, DELETE) and status codes.
- Write appropriate error handling to handle errors and return informative error messages.
- Include documentation that describes the endpoints, their expected inputs and outputs, and any authentication or authorization requirements.
- Use any additional libraries or frameworks that you deem necessary.

## Endpoints

Your API should implement the following endpoints:

- `GET /api/tasks`

Returns a JSON array of all tasks in the TODO list.

Response Body:

```json
[
  {
    "id": 1,
    "title": "Buy groceries",
    "completed": false
  },
  {
    "id": 2,
    "title": "Do laundry",
    "completed": true
  }
]
```

- `GET /api/tasks/:id`

Returns the JSON representation of a single task specified by the id parameter.

Response Body:

```json
{
  "id": 1,
  "title": "Buy groceries",
  "completed": false
}
```

- `POST /api/tasks`

Adds a new task to the TODO list.

Request Body:

```json
{
  "title": "Clean the house",
  "completed": false
}
```

Response Body:

```json
{
  "id": 3,
  "title": "Clean the house",
  "completed": false
}
```

- `PUT /api/tasks/:id`

Updates an existing task specified by the id parameter.

Request Body:

```json
{
  "title": "Buy groceries",
  "completed": true
}
```

- `DELETE /api/tasks/:id`

Deletes an existing task specified by the id parameter.

Response Body:

```json
{
  "message": "Task deleted successfully."
}
```

## Evaluation

Your solution will be evaluated on the following criteria:

- Does it implement all of the required endpoints?
- Does it follow RESTful principles?
- Is the code well-organized, modular, and readable?
- Is the error handling robust and informative?
- Is the documentation clear and concise?
- Are there any additional features or improvements you made that demonstrate your understanding of Go and RESTful APIs?
