package mapper

import (
	"testing"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
)

func TestHealthCheckFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordHealth
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns health check with zero values",
			dto:  &openapi.TgvalidatordHealth{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordHealth {
				tenantId := "tenant-123"
				componentName := "validator"
				componentId := "comp-456"
				group := "core"
				healthCheck := "database"
				status := "HEALTHY"
				lastUpdateDate := time.Now()
				validUntilDate := time.Now().Add(time.Hour)
				reportName := "DB Check"
				reportStatus := "HEALTHY"
				reportMessage := "All OK"
				return &openapi.TgvalidatordHealth{
					TenantId:       &tenantId,
					ComponentName:  &componentName,
					ComponentId:    &componentId,
					Group:          &group,
					HealthCheck:    &healthCheck,
					Status:         &status,
					LastUpdateDate: &lastUpdateDate,
					ValidUntilDate: &validUntilDate,
					Report: &openapi.TgvalidatordHealthReport{
						Name:    &reportName,
						Status:  &reportStatus,
						Message: &reportMessage,
					},
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HealthCheckFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("HealthCheckFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("HealthCheckFromDTO() returned nil for non-nil input")
			}
			// Verify fields if set
			if tt.dto.TenantId != nil && got.TenantID != *tt.dto.TenantId {
				t.Errorf("TenantID = %v, want %v", got.TenantID, *tt.dto.TenantId)
			}
			if tt.dto.ComponentName != nil && got.ComponentName != *tt.dto.ComponentName {
				t.Errorf("ComponentName = %v, want %v", got.ComponentName, *tt.dto.ComponentName)
			}
			if tt.dto.ComponentId != nil && got.ComponentID != *tt.dto.ComponentId {
				t.Errorf("ComponentID = %v, want %v", got.ComponentID, *tt.dto.ComponentId)
			}
			if tt.dto.Group != nil && got.Group != *tt.dto.Group {
				t.Errorf("Group = %v, want %v", got.Group, *tt.dto.Group)
			}
			if tt.dto.HealthCheck != nil && got.Name != *tt.dto.HealthCheck {
				t.Errorf("Name = %v, want %v", got.Name, *tt.dto.HealthCheck)
			}
			if tt.dto.Status != nil && got.Status != *tt.dto.Status {
				t.Errorf("Status = %v, want %v", got.Status, *tt.dto.Status)
			}
			if tt.dto.LastUpdateDate != nil && !got.LastUpdateDate.Equal(*tt.dto.LastUpdateDate) {
				t.Errorf("LastUpdateDate = %v, want %v", got.LastUpdateDate, *tt.dto.LastUpdateDate)
			}
			if tt.dto.ValidUntilDate != nil && !got.ValidUntilDate.Equal(*tt.dto.ValidUntilDate) {
				t.Errorf("ValidUntilDate = %v, want %v", got.ValidUntilDate, *tt.dto.ValidUntilDate)
			}
			// Verify report is mapped if present
			if tt.dto.Report != nil {
				if got.Report == nil {
					t.Error("Report should not be nil when DTO has report")
				}
			}
		})
	}
}

func TestHealthChecksFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordHealth
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1, // Special value to indicate nil check
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordHealth{},
			want: 0,
		},
		{
			name: "converts multiple health checks",
			dtos: func() []openapi.TgvalidatordHealth {
				name1 := "database"
				name2 := "cache"
				return []openapi.TgvalidatordHealth{
					{HealthCheck: &name1},
					{HealthCheck: &name2},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HealthChecksFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("HealthChecksFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("HealthChecksFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestHealthReportFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordHealthReport
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns report with zero values",
			dto:  &openapi.TgvalidatordHealthReport{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordHealthReport {
				name := "DB Check"
				status := "HEALTHY"
				message := "All connections OK"
				duration := "150ms"
				errMsg := ""
				results := map[string]string{"connections": "10"}
				return &openapi.TgvalidatordHealthReport{
					Name:     &name,
					Status:   &status,
					Message:  &message,
					Duration: &duration,
					Error:    &errMsg,
					Results:  &results,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HealthReportFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("HealthReportFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("HealthReportFromDTO() returned nil for non-nil input")
			}
			// Verify fields if set
			if tt.dto.Name != nil && got.Name != *tt.dto.Name {
				t.Errorf("Name = %v, want %v", got.Name, *tt.dto.Name)
			}
			if tt.dto.Status != nil && got.Status != *tt.dto.Status {
				t.Errorf("Status = %v, want %v", got.Status, *tt.dto.Status)
			}
			if tt.dto.Message != nil && got.Message != *tt.dto.Message {
				t.Errorf("Message = %v, want %v", got.Message, *tt.dto.Message)
			}
			if tt.dto.Duration != nil && got.Duration != *tt.dto.Duration {
				t.Errorf("Duration = %v, want %v", got.Duration, *tt.dto.Duration)
			}
			if tt.dto.Error != nil && got.Error != *tt.dto.Error {
				t.Errorf("Error = %v, want %v", got.Error, *tt.dto.Error)
			}
		})
	}
}

func TestClientStatusFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordClientStatus
	}{
		{
			name: "nil input returns zero value",
			dto:  nil,
		},
		{
			name: "empty DTO returns status with zero values",
			dto:  &openapi.TgvalidatordClientStatus{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordClientStatus {
				hostPort := "localhost:8080"
				connected := true
				errMsg := ""
				lastPing := time.Now()
				return &openapi.TgvalidatordClientStatus{
					HostPort:  &hostPort,
					Connected: &connected,
					Error:     &errMsg,
					LastPing:  &lastPing,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ClientStatusFromDTO(tt.dto)
			if tt.dto == nil {
				// For nil, we expect zero value struct
				if got.HostPort != "" || got.Connected || got.Error != "" {
					t.Errorf("ClientStatusFromDTO(nil) = %v, want zero value", got)
				}
				return
			}
			// Verify fields if set
			if tt.dto.HostPort != nil && got.HostPort != *tt.dto.HostPort {
				t.Errorf("HostPort = %v, want %v", got.HostPort, *tt.dto.HostPort)
			}
			if tt.dto.Connected != nil && got.Connected != *tt.dto.Connected {
				t.Errorf("Connected = %v, want %v", got.Connected, *tt.dto.Connected)
			}
			if tt.dto.Error != nil && got.Error != *tt.dto.Error {
				t.Errorf("Error = %v, want %v", got.Error, *tt.dto.Error)
			}
			if tt.dto.LastPing != nil && !got.LastPing.Equal(*tt.dto.LastPing) {
				t.Errorf("LastPing = %v, want %v", got.LastPing, *tt.dto.LastPing)
			}
		})
	}
}

func TestHealthGroupFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordHealthGroup
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns group with nil health checks",
			dto:  &openapi.TgvalidatordHealthGroup{},
		},
		{
			name: "DTO with health checks",
			dto: func() *openapi.TgvalidatordHealthGroup {
				name1 := "db-check"
				return &openapi.TgvalidatordHealthGroup{
					HealthChecks: []openapi.TgvalidatordHealth{
						{HealthCheck: &name1},
					},
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HealthGroupFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("HealthGroupFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("HealthGroupFromDTO() returned nil for non-nil input")
			}
			if tt.dto.HealthChecks != nil {
				if len(got.HealthChecks) != len(tt.dto.HealthChecks) {
					t.Errorf("HealthChecks length = %v, want %v", len(got.HealthChecks), len(tt.dto.HealthChecks))
				}
			}
		})
	}
}

func TestHealthComponentFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordHealthComponent
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns component with nil groups",
			dto:  &openapi.TgvalidatordHealthComponent{},
		},
		{
			name: "DTO with groups",
			dto: func() *openapi.TgvalidatordHealthComponent {
				groups := map[string]openapi.TgvalidatordHealthGroup{
					"core": {},
				}
				return &openapi.TgvalidatordHealthComponent{
					Groups: &groups,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HealthComponentFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("HealthComponentFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("HealthComponentFromDTO() returned nil for non-nil input")
			}
			if tt.dto.Groups != nil {
				if len(got.Groups) != len(*tt.dto.Groups) {
					t.Errorf("Groups length = %v, want %v", len(got.Groups), len(*tt.dto.Groups))
				}
			}
		})
	}
}

func TestHealthComponentsFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos *map[string]openapi.TgvalidatordHealthComponent
		want int
	}{
		{
			name: "nil input returns nil",
			dtos: nil,
			want: -1, // Special value to indicate nil check
		},
		{
			name: "empty map returns empty map",
			dtos: func() *map[string]openapi.TgvalidatordHealthComponent {
				m := make(map[string]openapi.TgvalidatordHealthComponent)
				return &m
			}(),
			want: 0,
		},
		{
			name: "converts multiple components",
			dtos: func() *map[string]openapi.TgvalidatordHealthComponent {
				m := map[string]openapi.TgvalidatordHealthComponent{
					"validator": {},
					"signer":    {},
				}
				return &m
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HealthComponentsFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("HealthComponentsFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("HealthComponentsFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestHealthReportFromDTO_WithVaultdClients(t *testing.T) {
	hostPort := "localhost:8080"
	connected := true
	lastPing := time.Now()

	dto := &openapi.TgvalidatordHealthReport{
		VaultdClients: []openapi.TgvalidatordClientStatus{
			{
				HostPort:  &hostPort,
				Connected: &connected,
				LastPing:  &lastPing,
			},
		},
	}

	got := HealthReportFromDTO(dto)
	if got == nil {
		t.Fatal("HealthReportFromDTO() returned nil for non-nil input")
	}
	if len(got.VaultdClients) != 1 {
		t.Errorf("VaultdClients length = %v, want 1", len(got.VaultdClients))
	}
	if got.VaultdClients[0].HostPort != hostPort {
		t.Errorf("VaultdClients[0].HostPort = %v, want %v", got.VaultdClients[0].HostPort, hostPort)
	}
	if got.VaultdClients[0].Connected != connected {
		t.Errorf("VaultdClients[0].Connected = %v, want %v", got.VaultdClients[0].Connected, connected)
	}
}

func TestHealthReportFromDTO_WithResults(t *testing.T) {
	results := map[string]string{
		"connections":     "10",
		"active_queries":  "5",
		"response_time":   "10ms",
	}

	dto := &openapi.TgvalidatordHealthReport{
		Results: &results,
	}

	got := HealthReportFromDTO(dto)
	if got == nil {
		t.Fatal("HealthReportFromDTO() returned nil for non-nil input")
	}
	if len(got.Results) != len(results) {
		t.Errorf("Results length = %v, want %v", len(got.Results), len(results))
	}
	for key, want := range results {
		if gotVal, ok := got.Results[key]; !ok || gotVal != want {
			t.Errorf("Results[%s] = %v, want %v", key, gotVal, want)
		}
	}
}

func TestHealthCheckFromDTO_NilDates(t *testing.T) {
	status := "HEALTHY"
	dto := &openapi.TgvalidatordHealth{
		Status:         &status,
		LastUpdateDate: nil,
		ValidUntilDate: nil,
	}

	got := HealthCheckFromDTO(dto)
	if got == nil {
		t.Fatal("HealthCheckFromDTO() returned nil for non-nil input")
	}
	// When dates are nil, they should be zero time value
	if !got.LastUpdateDate.IsZero() {
		t.Errorf("LastUpdateDate should be zero time when nil, got %v", got.LastUpdateDate)
	}
	if !got.ValidUntilDate.IsZero() {
		t.Errorf("ValidUntilDate should be zero time when nil, got %v", got.ValidUntilDate)
	}
}

func TestClientStatusFromDTO_NilLastPing(t *testing.T) {
	hostPort := "localhost:8080"
	dto := &openapi.TgvalidatordClientStatus{
		HostPort: &hostPort,
		LastPing: nil,
	}

	got := ClientStatusFromDTO(dto)
	if !got.LastPing.IsZero() {
		t.Errorf("LastPing should be zero time when nil, got %v", got.LastPing)
	}
}
