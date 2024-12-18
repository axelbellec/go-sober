# Modern Go Backend Structure

## Tech Stack Updates

- Go 1.23+ (latest stable version)

## Project Structure

Use controller, service, repository, and model folders.

## Architecture Guidelines

- Go 1.23
- net/http for routing
- SQLite for database
- JWT for authentication
- Bcrypt for password hashing

## Bruno API Collection

The project includes a Bruno API collection for testing the endpoints. The collection supports multiple environments:

- Local
- Test
- Production

Install Bruno CLI: https://docs.usebruno.com/bru-cli/overview

```bash
# Go to bruno directory and run tests
cd ./bruno
bru run --env Local
```
