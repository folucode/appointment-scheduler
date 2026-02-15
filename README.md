# Project Name
Appointment Scheduling App

---

## Prerequisites

* **Go**: Version 1.21 or higher ([Download](https://go.dev/dl/))
* **Docker & Docker Compose**: To manage the database container ([Download](https://www.docker.com/products/docker-desktop/))
* **Git**: For version control

---

## Getting Started

### 1. Environment Configuration

1. Create a `.env` file in the root directory.
2. Open `.env` and copy the variables in the `.env.example` file and adjust the credentials to match your prosgres credentials.

### 2. Running the Services

1. Run this command to build and start all services: `docker compose up --build`

* [Frontend](http://localhost:5173)
* [Backend](http://localhost:8080)
