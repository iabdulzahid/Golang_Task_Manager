# Golang_Task_Manager

**Golang_Task_Manager** is a task management API built using **Go (Golang)** and **SQLite**. This API provides a lightweight solution for managing tasks, supporting key features such as **CRUD operations**, **task prioritization**, **due dates**, **labels**, and **data export** in **JSON** and **CSV** formats. The API also includes **rate limiting** to ensure fair usage and prevent abuse.

## Features
- **Task CRUD Operations**: Create, read, update, and delete tasks.
- **Task Prioritization**: Assign priority levels (e.g., high, medium, low) to tasks.
- **Due Dates**: Set due dates for tasks to ensure timely completion.
- **Labels**: Categorize tasks with custom labels for better organization.
- **Task Export**: Export tasks in **JSON** or **CSV** format for backup, sharing, or integration.
- **Rate Limiting**: Prevents abuse and ensures fair API usage by limiting requests.

## Technology Stack
- **Go (Golang)**: Backend programming language for building a fast and scalable API.
- **SQLite**: Lightweight, serverless database used to store task data.
- **Gin**: Web framework for building RESTful APIs in Go.
- **Swagger**: Interactive documentation for exploring and testing API endpoints.

## Installation & Setup

### Prerequisites
- **Go (Golang)** v1.16 or higher.
- **SQLite** for local database storage.

### Steps to Run the API

1. **Clone the repository**:
    ```bash
    git clone https://github.com/your-username/golang_task_manager.git
    cd golang_task_manager
    ```

2. **Install dependencies**:
    ```bash
    go mod tidy
    ```

3. **Run the API server**:
    ```bash
    go run main.go
    ```

4. **Access the API**:
    The server will be available at `http://localhost:8080`.

5. **Swagger Documentation**:
    Access the interactive Swagger documentation at:
    ```
    http://localhost:8080/swagger
    ```

---

## API Endpoints

### 1. **Create Task**
- **Endpoint**: `POST /tasks`
- **Request Body**:
    ```json
    {
      "title": "Task Title",
      "description": "Task Description",
      "priority": "high",
      "due_date": "2024-12-01T00:00:00Z",
      "labels": ["work", "urgent"]
    }
    ```

### 2. **Get All Tasks**
- **Endpoint**: `GET /tasks`
- **Response**:
    ```json
    [
      {
        "id": 1,
        "title": "Task Title",
        "description": "Task Description",
        "priority": "high",
        "due_date": "2024-12-01T00:00:00Z",
        "labels": ["work", "urgent"]
      }
    ]
    ```

### 3. **Get Task by ID**
- **Endpoint**: `GET /tasks/{id}`

### 4. **Update Task**
- **Endpoint**: `PUT /tasks/{id}`
- **Request Body**:
    ```json
    {
      "title": "Updated Task Title",
      "description": "Updated Description",
      "priority": "medium",
      "due_date": "2024-12-05T00:00:00Z",
      "labels": ["personal", "low"]
    }
    ```

### 5. **Delete Task**
- **Endpoint**: `DELETE /tasks/{id}`
- **Response**:
    ```json
    {
      "message": "Task deleted successfully"
    }
    ```

### 6. **Export Tasks**
- **Endpoint**: `GET /tasks/export`
- **Description**: Exports all tasks in **JSON** or **CSV** format.
- **Query Parameters**: `format=json` or `format=csv`

---

## Rate Limiting

To ensure fair usage of the API, **rate limiting** is implemented. By default, users are limited to **5 requests per minute**. If this limit is exceeded, the API will respond with a `429 Too Many Requests` error.

---

## Contributing

We welcome contributions to **Golang_Task_Manager**! If you find any bugs or want to add features, feel free to submit a pull request. When contributing, please ensure you:
- Follow the existing code style.
- Write tests for any new features.
- Update the documentation if necessary.

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## Contact

For questions or feedback, please open an issue or contact the project maintainers.

---

This **README** provides a comprehensive guide to using, installing, and contributing to **Golang_Task_Manager**. It includes installation steps, detailed information about API endpoints, and additional information about rate limiting and how to contribute to the project.
