## Getting Started

Follow these steps to set up and run the project:

1. Clone the repository and navigate to the project directory:

   ```bash
   git clone https://github.com/mdayat/artics-communication.git
   cd artics-communication
   ```

2. Environment variables:

   `.env` files are already included in both the `go` and `react` app directories for ease of testing and review.

3. Start the applications using Docker:

   ```bash
   docker compose up -d
   ```

4. Run database migrations and seed data:

   ```bash
   cd go && make seed
   ```

5. Access the applications:
   - React frontend: [http://localhost:3000](http://localhost:3000)
   - Go backend: [http://localhost:8080](http://localhost:8080)

## Tech Stack

- **Backend**

  - Go
  - go-chi HTTP router
  - Atlas for database migration management

- **Frontend**

  - React with TypeScript for type safety
  - Tailwind CSS for styling
  - shadcn component library for UI components

- **Database**

  - PostgreSQL

- **Infrastructure**
  - Docker for containerization and service orchestration

## Architecture Overview

The application follows a client-server architecture composed of the following components:

- **Frontend (React + TypeScript)**  
  Serves as the user interface layer. Communicates with the backend via HTTP using RESTful APIs.

- **Backend (Go)**  
  Handles incoming HTTP requests, performs business logic, manages authentication, and interacts with the database.

- **Database (PostgreSQL)**  
  Stores and retrieves data.

- **Infrastructure (Docker Compose)**  
  All services (frontend, backend, and database) are containerized and orchestrated using Docker Compose for easy local development.
