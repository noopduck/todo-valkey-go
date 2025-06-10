## Project Title

Todo Application

## Description

This project is a Todo application that utilizes Valkey for backend storage and Nginx as a reverse proxy. It is built using Go and Docker to provide a simple and efficient way to manage todos.

## Technologies Used

- **Go**: The programming language used to develop the application.
- **Valkey**: A lightweight key-value store for backend storage.
- **Nginx**: Acts as a reverse proxy to handle incoming requests.
- **Docker**: For containerization of the application and its dependencies.

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/todo-valkey-go.git
   cd todo-valkey-go
   ```

2. Ensure you have Docker and Docker Compose installed.
3. Start the application using Docker Compose:

   ```bash
   docker-compose up
   ```

## Usage

- Access the application at `http://localhost:8080`
- Use the following endpoints:
  - `GET /listKeys`: Retrieve all keys in the Valkey database.
  - `POST /addItem`: Add a new todo item (JSON body required).

    ```bash
    curl -X POST localhost:3000/addItem -H "Content-Type: application/json" -d '{ "Description": "Need to prettify error handling and add authentication eventually", "Status": 2 }'
    ```

  ```
  - `GET /getItem?uuid=<UUID>`: Retrieve a specific todo item by its UUID.

## Contributing

Contributions are welcome! Please fork this repository and submit a pull request.

- Ensure your code follows the existing style and conventions.
- Add tests for any new functionality.
- Update the documentation as needed.

## Installation

Instructions on how to install and set up the project.

## Usage

Examples of how to use the project.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
