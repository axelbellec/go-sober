# Sōber 🧃

A robust Go-based API for tracking alcohol consumption and calculating Blood Alcohol Content (BAC) with high precision.

## 📋 Table of Contents

- [Features](#features)
- [Getting Started](#getting-started)
- [API Documentation](#api-documentation)
- [Architecture](#architecture)
- [Development](#development)
- [Testing](#testing)
- [Deployment](#deployment)

## ✨ Features

### 🔐 Authentication

- **JWT-based** authentication system
- Secure email/password registration and login
- Password hashing with bcrypt
- Protected routes via middleware
- Token refresh mechanism
- Session management

### 🥃 Drink Management

#### Drink Options Catalog

- Comprehensive database of pre-defined drinks
- Standardized measurements for accuracy
- Each drink includes:
  - Name and type
  - Precise size measurements (ml/oz)
  - ABV (Alcohol By Volume) percentage
  - Standard drink equivalent

#### Smart Drink Logging

- Real-time and historical drink logging
- Accurate timestamp tracking
- User-specific drink history
- Detailed consumption patterns
- Quick-log favorite drinks

### 📊 BAC Analytics

#### Advanced BAC Calculation

- Widmark formula implementation
- Personalized calculations based on:
  - Body weight
  - Biological sex
  - Consumption timeline
  - Drink specifications

#### Comprehensive Analytics

- Real-time BAC monitoring
- Customizable time series data
- Status indicators:
  - 🟢 Sober (0.00-0.02%)
  - 🟡 Minimal (0.02-0.05%)
  - 🟠 Light (0.05-0.08%)
  - 🔴 Over Limit (>0.08%)
- Detailed statistics:
  - Peak BAC levels
  - Time until sober
  - Legal limit warnings
  - Consumption patterns

## 🚀 Getting Started

### Prerequisites

- Go 1.23 or higher
- SQLite 3
- [Task](https://taskfile.dev/) (task runner)
- [Bruno](https://www.usebruno.com/) (API testing)
- [Groq](https://groq.com/) (for natural language processing)
- [Swag](https://github.com/swaggo/swag) (for generating Swagger documentation)
- [Air](https://github.com/air-verse/air) (for hot reloading)

### Installation

1. Clone the repository:

```bash
git clone https://github.com/axelbellec/sober.git
cd sober
```

2. Install dependencies:

```bash
go mod download
```

3. Set up environment:

```bash
cp .env.example .env
# Edit .env with your configurations
```

4. Initialize database:

```bash
task db:migrate
```

5. Generate Swagger documentation:

```bash
task docs
```

6. Run the application:

```bash
task run
```

## 📚 API Documentation

The API is documented using Swagger/OpenAPI. You can access the interactive documentation at:

```
http://localhost:8080/swagger
```

### Bruno API Collection

The project includes comprehensive API tests using [Bruno](https://docs.usebruno.com/).

#### Environment Support

- 🔧 Local (`localhost:8080`)
- 🧪 Test (`test.api.go-sober.com`) 🚧
- 🌐 Production (`api.go-sober.com`) 🚧

#### Running API Tests

```bash
# Install Bruno CLI
npm install -g @usebruno/cli

# Run tests
task test:bruno
```

## 🏗 Architecture

### Clean Architecture

```
├── main.go           # API entrypoint
├── internal/         # Private application code
│   ├── auth/         # Authentication logic
│   ├── drinks/       # Drink management
│   ├── analytics/    # BAC calculations
│   └── models/       # Domain models
└── platform/         # Platform-specific code
```

### Technology Stack

- **Framework**: net/http
- **Database**: SQLite
- **Authentication**: JWT
- **Password Security**: Bcrypt
- **Testing**: Go testing + Bruno
- **Task Runner**: Taskfile

## 💻 Development

### Using Taskfile

```bash
# Start development server
task run

# Run production binary
task run:prod

# Run tests
task test

# Run tests with coverage
task test:coverage

# Run Bruno tests
task test:bruno

# Build development binary
task build

# Build production binary
task build:prod

# Database operations
task db:migrate
task db:rollback

# Docker operations
task docker:build
task docker:run
```

### Task Descriptions

- `run`: Starts the development server using `go run`
- `run:prod`: Runs the compiled production binary
- `build`: Builds the binary for development
- `build:prod`: Creates an optimized production binary in `bin/go-sober`
- `test`: Runs all Go tests with verbose output
- `test:coverage`: Executes tests and generates coverage reports
- `test:bruno`: Runs API tests using Bruno (includes database migration)
- `db:migrate`: Applies all pending database migrations
- `db:rollback`: Rolls back the most recent database migration
- `docker:build`: Builds the Docker image
- `docker:run`: Runs the Docker container with environment variables

### Code Style

- Follow Go standard guidelines
- Use `gofmt` for formatting
- Run `golangci-lint` before commits

## 🧪 Testing

### Running Tests

```bash
# Run all tests
task test

# Run with coverage
task test:coverage
```

## 📦 Deployment

### Building for Production

```bash
# Build optimized binary
task build:prod

# Run production build
./bin/go-sober
```

### Docker Support

```bash
# Build container
docker build -t sober .

# Run container
docker run -p 8080:8080 sober
```
