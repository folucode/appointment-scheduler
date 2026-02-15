# Project Name
Appointment Scheduling App

---

## Prerequisites

* **Go**: Version 1.21 or higher ([Download](https://go.dev/dl/))
* **Docker & Docker Compose**: To manage the database container ([Download](https://www.docker.com/products/docker-desktop/))
* **Git**: For version control

---

## Getting Started

### 1. Database Setup

1. Create a new file named `docker-compose.yml`.
2. Copy the contents from `docker-compose.example.yml` into your new file.
3. Update the database credentials (username, password, etc.).
4. Start the database container:
   ```bash
        docker compose up -d
   ```

### 2. Environment Configuration

1. Create a `.env` file in the root directory.
2. Open `.env` and update the connection string to match your Docker DB settings.

### 3. Running the Go Service

1. Download dependencies: `go mod download`.
2. Start the application: `go run server/main.go`.

The Go service should now be running locally!

### 4. Install Frontend Dependencies

1. In the frontend directory, run `npm i` to install the dependencies.
2. Start the application by running `npm run dev`.

The frontend service should now be running locally!