package model

import "time"

// HealthCheck represents a single health check entry.
type HealthCheck struct {
	// TenantID is the tenant ID for the health check.
	TenantID string `json:"tenant_id,omitempty"`
	// ComponentName is the name of the component being checked.
	ComponentName string `json:"component_name,omitempty"`
	// ComponentID is the unique identifier for the component.
	ComponentID string `json:"component_id,omitempty"`
	// Group is the group this health check belongs to.
	Group string `json:"group,omitempty"`
	// Name is the name of the health check.
	Name string `json:"name,omitempty"`
	// Status is the health status (e.g., "HEALTHY", "UNHEALTHY").
	Status string `json:"status,omitempty"`
	// Report contains detailed health check report information.
	Report *HealthReport `json:"report,omitempty"`
	// LastUpdateDate is when the health check was last updated.
	LastUpdateDate time.Time `json:"last_update_date,omitempty"`
	// ValidUntilDate is when the health check result expires.
	ValidUntilDate time.Time `json:"valid_until_date,omitempty"`
}

// HealthReport contains detailed information about a health check.
type HealthReport struct {
	// Name is the name of the health check report.
	Name string `json:"name,omitempty"`
	// Status is the status of the report.
	Status string `json:"status,omitempty"`
	// Message contains additional status message.
	Message string `json:"message,omitempty"`
	// Duration is how long the check took to execute.
	Duration string `json:"duration,omitempty"`
	// Error contains error message if the check failed.
	Error string `json:"error,omitempty"`
	// Results contains additional result data as key-value pairs.
	Results map[string]string `json:"results,omitempty"`
	// VaultdClients contains status information about vault clients.
	VaultdClients []ClientStatus `json:"vaultd_clients,omitempty"`
}

// ClientStatus represents the status of a client connection.
type ClientStatus struct {
	// HostPort is the host:port of the client.
	HostPort string `json:"host_port,omitempty"`
	// Connected indicates whether the client is connected.
	Connected bool `json:"connected"`
	// Error contains error message if connection failed.
	Error string `json:"error,omitempty"`
	// LastPing is when the client was last pinged.
	LastPing time.Time `json:"last_ping,omitempty"`
}

// HealthGroup represents a group of health checks.
type HealthGroup struct {
	// HealthChecks is the list of health checks in this group.
	HealthChecks []*HealthCheck `json:"health_checks,omitempty"`
}

// HealthComponent represents a component with health check groups.
type HealthComponent struct {
	// Groups is a map of group name to health group.
	Groups map[string]*HealthGroup `json:"groups,omitempty"`
}

// GetAllHealthChecksOptions contains options for getting all health checks.
type GetAllHealthChecksOptions struct {
	// TenantID filters health checks by tenant ID.
	TenantID string
	// FailIfUnhealthy causes the API to return an error if any checks are unhealthy.
	FailIfUnhealthy bool
}

// GetAllHealthChecksResult contains the result of getting all health checks.
type GetAllHealthChecksResult struct {
	// Components is a map of component name to health component.
	Components map[string]*HealthComponent `json:"components,omitempty"`
}
