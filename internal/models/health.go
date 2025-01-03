package models

type HealthStatus string

const (
	HealthStatusOK    HealthStatus = "OK"
	HealthStatusError HealthStatus = "Error"
)

type Health struct {
	Status HealthStatus `json:"status"`
}
