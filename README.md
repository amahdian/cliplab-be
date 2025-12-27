# cliplab Backend

AI tool for saving and managing social media links

## 📋 Prerequisites

Before you begin, ensure you have the following installed:

*   **Go 1.24.2+**: [Download Go](https://golang.org/dl/)
*   **PostgreSQL 10.3+**: [Download PostgreSQL](https://www.postgresql.org/download/)
*   **Docker & Docker Compose** (recommended): [Download Docker](https://www.docker.com/products/docker-desktop) 
*   **Make**: For using the provided Makefile commands.
*   **golang-migrate**: For database migrations.

### Installing golang-migrate

You can install `golang-migrate` using one of the following methods:

*   **Go Install (Recommended)**:
    ```bash
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    ```
*   **Homebrew (macOS)**:
    ```bash
    brew install golang-migrate
    ```

## 🛠️ Setup

1.  **Clone the Repository**:
    ```bash
    git clone https://github.com/amahdian/cliplib-be.git
    cd cliplab-be
    ```

2.  **Environment Configuration**:
    Create an environment file by copying the example:
    ```bash
    cp .env .env.dev
    ```
    Update the `.env.dev` file with your configuration, especially the database connection string and JWT secret.

3.  **Database Setup**:
    *   **Using Docker (Recommended)**:
        ```bash
        docker-compose up -d postgres
        make create-db
        make migrate-up
        ```
    *   **Using Local PostgreSQL**:
        ```bash
        # Create the database manually
        createdb cliplab

        # Run migrations
        make migrate-up
        ```

4.  **Install Dependencies**:
    ```bash
    make vendor
    ```

## 🚀 Running the Application

*   **Development Mode**:
    ```bash
    make dev
    ```
*   **Production Mode**:
    ```bash
    make build
    ./build/app-bin
    ```

## 📖 API Documentation

Once the application is running, you can access the Swagger UI for API documentation at:

[http://localhost:8080/swagger/index.html](http://localhost:8090/swagger/index.html)

## 🏗️ Project Structure

The project follows a standard Go project layout:

```
cliplab-be/
├── assets/         # Static assets and migrations
├── docs/           # Swagger documentation
├── domain/         # Domain models and contracts
├── global/         # Global configurations and utilities
├── pkg/            # Reusable packages
├── server/         # HTTP server components
├── storage/        # Data storage layer
├── svc/            # Business logic services
├── main.go         # Application entry point
└── Makefile        # Build and development commands
```