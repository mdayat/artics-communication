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

- Go
- React with TypeScript and Tailwind CSS
- PostgreSQL with atlas for database migration
- Docker for containerization and service orchestration

## Architecture Overview

The application follows a client-server architecture composed of the following components:

- **Frontend (React + TypeScript)**  
  Serves as the user interface layer. Communicates with the backend via RESTful APIs.

- **Backend (Go)**  
  Handles incoming HTTP requests, performs business logic, and interacts with the database.

- **Database (PostgreSQL)**  
  Stores and retrieves data.

- **Infrastructure (Docker Compose)**  
  All services (frontend, backend, and database) are containerized and orchestrated using Docker Compose for easy local development.

## Notes for Reviewers

### Login Credentials

Please refer to `go/cmd/seed/main.go` file to find the pre-created account you can use for login.

### Run Tests

To run tests for the Go app:

```bash
cd go
make test
```

### Run Code Quality Checks

To evaluate the code with linting, vulnerability scan, and static analysis:

```bash
cd go
make vet && make staticcheck && make govulncheck && make revive
```

## Project Checklist

Below is the list of requirements from the assessment along with their completion status:

### Functional Specifications

- [x] **User Roles**

  - [x] User can register and log in
  - [x] Admin can view all bookings

- [x] **Main Features**

  - [x] User can view available meeting rooms
  - [x] User can book a room by selecting date and time (time slot)
  - [x] User can view their booking history
  - [x] Admin can view all bookings
  - [x] Admin can cancel any booking

- [x] **Validations**
  - [x] Prevent double booking for the same room and time
  - [x] User can only cancel their own bookings (unless Admin)

### Technical Requirements

- [x] **Frontend**

  - [x] Built using React
  - [x] Pages: Login/Register, User and Admin Dashboard, Booking Form (User)

- [x] **Backend**

  - [x] Built using Go
  - [x] REST API

- [x] **Database**

  - [x] Using PostgreSQL
  - [x] Schema described in SQL file located in `go/schema.sql`

- [x] **Authentication**: JWT-based authentication stored in cookie

### Bonus (Optional)

- [x] **Responsive Design** (works well on mobile)
- [x] **Integration Tests** for HTTP handlers or REST API
- [x] **Deployment**
  - [x] Frontend via Vercel
  - [x] Backend via VPS + Coolify

### Extra

- [x] **CI/CD Integration**: Runs code quality checks and integration tests on push/PR to `main` branch using GitHub Actions
- [x] **API Specs**: https://app.swaggerhub.com/apis-docs/MUHNURDAYAT/artics-communication-test/1.0.0-oas3.1
