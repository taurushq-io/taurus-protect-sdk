package mapper

import (
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// HealthCheckFromDTO converts an OpenAPI Health to a domain HealthCheck.
func HealthCheckFromDTO(dto *openapi.TgvalidatordHealth) *model.HealthCheck {
	if dto == nil {
		return nil
	}

	check := &model.HealthCheck{
		TenantID:      safeString(dto.TenantId),
		ComponentName: safeString(dto.ComponentName),
		ComponentID:   safeString(dto.ComponentId),
		Group:         safeString(dto.Group),
		Name:          safeString(dto.HealthCheck),
		Status:        safeString(dto.Status),
	}

	if dto.Report != nil {
		check.Report = HealthReportFromDTO(dto.Report)
	}

	if dto.LastUpdateDate != nil {
		check.LastUpdateDate = *dto.LastUpdateDate
	}

	if dto.ValidUntilDate != nil {
		check.ValidUntilDate = *dto.ValidUntilDate
	}

	return check
}

// HealthChecksFromDTO converts a slice of OpenAPI Health to domain HealthChecks.
func HealthChecksFromDTO(dtos []openapi.TgvalidatordHealth) []*model.HealthCheck {
	if dtos == nil {
		return nil
	}
	checks := make([]*model.HealthCheck, len(dtos))
	for i := range dtos {
		checks[i] = HealthCheckFromDTO(&dtos[i])
	}
	return checks
}

// HealthReportFromDTO converts an OpenAPI HealthReport to a domain HealthReport.
func HealthReportFromDTO(dto *openapi.TgvalidatordHealthReport) *model.HealthReport {
	if dto == nil {
		return nil
	}

	report := &model.HealthReport{
		Name:     safeString(dto.Name),
		Status:   safeString(dto.Status),
		Message:  safeString(dto.Message),
		Duration: safeString(dto.Duration),
		Error:    safeString(dto.Error),
	}

	if dto.Results != nil {
		report.Results = *dto.Results
	}

	if dto.VaultdClients != nil {
		report.VaultdClients = make([]model.ClientStatus, len(dto.VaultdClients))
		for i, client := range dto.VaultdClients {
			report.VaultdClients[i] = ClientStatusFromDTO(&client)
		}
	}

	return report
}

// ClientStatusFromDTO converts an OpenAPI ClientStatus to a domain ClientStatus.
func ClientStatusFromDTO(dto *openapi.TgvalidatordClientStatus) model.ClientStatus {
	if dto == nil {
		return model.ClientStatus{}
	}

	status := model.ClientStatus{
		HostPort:  safeString(dto.HostPort),
		Connected: safeBool(dto.Connected),
		Error:     safeString(dto.Error),
	}

	if dto.LastPing != nil {
		status.LastPing = *dto.LastPing
	}

	return status
}

// HealthGroupFromDTO converts an OpenAPI HealthGroup to a domain HealthGroup.
func HealthGroupFromDTO(dto *openapi.TgvalidatordHealthGroup) *model.HealthGroup {
	if dto == nil {
		return nil
	}

	return &model.HealthGroup{
		HealthChecks: HealthChecksFromDTO(dto.HealthChecks),
	}
}

// HealthComponentFromDTO converts an OpenAPI HealthComponent to a domain HealthComponent.
func HealthComponentFromDTO(dto *openapi.TgvalidatordHealthComponent) *model.HealthComponent {
	if dto == nil {
		return nil
	}

	component := &model.HealthComponent{}

	if dto.Groups != nil {
		component.Groups = make(map[string]*model.HealthGroup)
		for name, group := range *dto.Groups {
			component.Groups[name] = HealthGroupFromDTO(&group)
		}
	}

	return component
}

// HealthComponentsFromDTO converts a map of OpenAPI HealthComponent to domain HealthComponents.
func HealthComponentsFromDTO(dtos *map[string]openapi.TgvalidatordHealthComponent) map[string]*model.HealthComponent {
	if dtos == nil {
		return nil
	}

	components := make(map[string]*model.HealthComponent)
	for name, dto := range *dtos {
		components[name] = HealthComponentFromDTO(&dto)
	}
	return components
}
