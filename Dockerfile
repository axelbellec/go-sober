FROM golang:1.22 AS base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify
# Install Taskfile for task runner
RUN go install github.com/go-task/task/v3/cmd/task@latest

COPY . .
RUN go build -o /app/bin .

FROM gcr.io/distroless/base-debian12
COPY --from=base /app/bin /app/sober


ENV ENVIRONMENT=dev

USER nonroot:nonroot

EXPOSE 8080

ENTRYPOINT ["/app/sober/go-sober"]

