package mapper

import (
	"testing"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
)

func TestJobFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordJob
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns job with zero values",
			dto:  &openapi.TgvalidatordJob{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordJob {
				name := "balance-sync"
				pending := "5"
				successes := "100"
				failures := "2"
				return &openapi.TgvalidatordJob{
					Name: &name,
					Statistics: &openapi.TgvalidatordJobStatistics{
						Pending:   &pending,
						Successes: &successes,
						Failures:  &failures,
					},
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := JobFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("JobFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("JobFromDTO() returned nil for non-nil input")
			}
			// Verify fields if set
			if tt.dto.Name != nil && got.Name != *tt.dto.Name {
				t.Errorf("Name = %v, want %v", got.Name, *tt.dto.Name)
			}
			// Verify statistics is mapped if present
			if tt.dto.Statistics != nil {
				if got.Statistics == nil {
					t.Error("Statistics should not be nil when DTO has statistics")
				}
			}
		})
	}
}

func TestJobsFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dtos []openapi.TgvalidatordJob
		want int
	}{
		{
			name: "nil slice returns nil",
			dtos: nil,
			want: -1, // Special value to indicate nil check
		},
		{
			name: "empty slice returns empty slice",
			dtos: []openapi.TgvalidatordJob{},
			want: 0,
		},
		{
			name: "converts multiple jobs",
			dtos: func() []openapi.TgvalidatordJob {
				name1 := "job1"
				name2 := "job2"
				return []openapi.TgvalidatordJob{
					{Name: &name1},
					{Name: &name2},
				}
			}(),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := JobsFromDTO(tt.dtos)
			if tt.want == -1 {
				if got != nil {
					t.Errorf("JobsFromDTO() = %v, want nil", got)
				}
				return
			}
			if len(got) != tt.want {
				t.Errorf("JobsFromDTO() length = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func TestJobStatisticsFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordJobStatistics
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns statistics with zero values",
			dto:  &openapi.TgvalidatordJobStatistics{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordJobStatistics {
				pending := "10"
				successes := "500"
				failures := "5"
				avgDuration := "1.5s"
				maxDuration := "5s"
				minDuration := "0.5s"
				statusID := "status-123"
				return &openapi.TgvalidatordJobStatistics{
					Pending:     &pending,
					Successes:   &successes,
					Failures:    &failures,
					AvgDuration: &avgDuration,
					MaxDuration: &maxDuration,
					MinDuration: &minDuration,
					LastSuccess: &openapi.TgvalidatordJobStatus{
						Id: &statusID,
					},
					LastFailure: &openapi.TgvalidatordJobStatus{
						Id: &statusID,
					},
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := JobStatisticsFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("JobStatisticsFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("JobStatisticsFromDTO() returned nil for non-nil input")
			}
			// Verify duration fields if set
			if tt.dto.AvgDuration != nil && got.AvgDuration != *tt.dto.AvgDuration {
				t.Errorf("AvgDuration = %v, want %v", got.AvgDuration, *tt.dto.AvgDuration)
			}
			if tt.dto.MaxDuration != nil && got.MaxDuration != *tt.dto.MaxDuration {
				t.Errorf("MaxDuration = %v, want %v", got.MaxDuration, *tt.dto.MaxDuration)
			}
			if tt.dto.MinDuration != nil && got.MinDuration != *tt.dto.MinDuration {
				t.Errorf("MinDuration = %v, want %v", got.MinDuration, *tt.dto.MinDuration)
			}
			// Verify last success/failure are mapped if present
			if tt.dto.LastSuccess != nil && got.LastSuccess == nil {
				t.Error("LastSuccess should not be nil when DTO has LastSuccess")
			}
			if tt.dto.LastFailure != nil && got.LastFailure == nil {
				t.Error("LastFailure should not be nil when DTO has LastFailure")
			}
		})
	}
}

func TestJobStatisticsFromDTO_NumericParsing(t *testing.T) {
	tests := []struct {
		name          string
		pending       *string
		successes     *string
		failures      *string
		wantPending   int64
		wantSuccesses int64
		wantFailures  int64
	}{
		{
			name:          "nil values default to zero",
			pending:       nil,
			successes:     nil,
			failures:      nil,
			wantPending:   0,
			wantSuccesses: 0,
			wantFailures:  0,
		},
		{
			name:          "valid numeric strings are parsed",
			pending:       stringPtr("10"),
			successes:     stringPtr("100"),
			failures:      stringPtr("5"),
			wantPending:   10,
			wantSuccesses: 100,
			wantFailures:  5,
		},
		{
			name:          "invalid numeric strings default to zero",
			pending:       stringPtr("invalid"),
			successes:     stringPtr(""),
			failures:      stringPtr("not-a-number"),
			wantPending:   0,
			wantSuccesses: 0,
			wantFailures:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dto := &openapi.TgvalidatordJobStatistics{
				Pending:   tt.pending,
				Successes: tt.successes,
				Failures:  tt.failures,
			}
			got := JobStatisticsFromDTO(dto)
			if got.Pending != tt.wantPending {
				t.Errorf("Pending = %v, want %v", got.Pending, tt.wantPending)
			}
			if got.Successes != tt.wantSuccesses {
				t.Errorf("Successes = %v, want %v", got.Successes, tt.wantSuccesses)
			}
			if got.Failures != tt.wantFailures {
				t.Errorf("Failures = %v, want %v", got.Failures, tt.wantFailures)
			}
		})
	}
}

func TestJobStatusFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordJobStatus
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns status with zero values",
			dto:  &openapi.TgvalidatordJobStatus{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordJobStatus {
				id := "status-123"
				message := "Job completed successfully"
				status := "completed"
				now := time.Now()
				return &openapi.TgvalidatordJobStatus{
					Id:        &id,
					Message:   &message,
					Status:    &status,
					StartedAt: &now,
					UpdatedAt: &now,
					TimeoutAt: &now,
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := JobStatusFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("JobStatusFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("JobStatusFromDTO() returned nil for non-nil input")
			}
			// Verify fields if set
			if tt.dto.Id != nil && got.ID != *tt.dto.Id {
				t.Errorf("ID = %v, want %v", got.ID, *tt.dto.Id)
			}
			if tt.dto.Message != nil && got.Message != *tt.dto.Message {
				t.Errorf("Message = %v, want %v", got.Message, *tt.dto.Message)
			}
			if tt.dto.Status != nil && got.Status != *tt.dto.Status {
				t.Errorf("Status = %v, want %v", got.Status, *tt.dto.Status)
			}
			if tt.dto.StartedAt != nil && !got.StartedAt.Equal(*tt.dto.StartedAt) {
				t.Errorf("StartedAt = %v, want %v", got.StartedAt, *tt.dto.StartedAt)
			}
			if tt.dto.UpdatedAt != nil && !got.UpdatedAt.Equal(*tt.dto.UpdatedAt) {
				t.Errorf("UpdatedAt = %v, want %v", got.UpdatedAt, *tt.dto.UpdatedAt)
			}
			if tt.dto.TimeoutAt != nil && !got.TimeoutAt.Equal(*tt.dto.TimeoutAt) {
				t.Errorf("TimeoutAt = %v, want %v", got.TimeoutAt, *tt.dto.TimeoutAt)
			}
		})
	}
}

func TestJobStatusFromDTO_NilTimeFields(t *testing.T) {
	id := "status-123"
	dto := &openapi.TgvalidatordJobStatus{
		Id:        &id,
		StartedAt: nil,
		UpdatedAt: nil,
		TimeoutAt: nil,
	}

	got := JobStatusFromDTO(dto)
	if got == nil {
		t.Fatal("JobStatusFromDTO() returned nil for non-nil input")
	}
	// When time fields are nil, they should be the zero time value
	if !got.StartedAt.IsZero() {
		t.Errorf("StartedAt should be zero time when nil, got %v", got.StartedAt)
	}
	if !got.UpdatedAt.IsZero() {
		t.Errorf("UpdatedAt should be zero time when nil, got %v", got.UpdatedAt)
	}
	if !got.TimeoutAt.IsZero() {
		t.Errorf("TimeoutAt should be zero time when nil, got %v", got.TimeoutAt)
	}
}

func TestJobFromDTO_NilStatistics(t *testing.T) {
	name := "test-job"
	dto := &openapi.TgvalidatordJob{
		Name:       &name,
		Statistics: nil,
	}

	got := JobFromDTO(dto)
	if got == nil {
		t.Fatal("JobFromDTO() returned nil for non-nil input")
	}
	if got.Statistics != nil {
		t.Errorf("Statistics should be nil when DTO statistics is nil, got %v", got.Statistics)
	}
	if got.Name != "test-job" {
		t.Errorf("Name = %v, want test-job", got.Name)
	}
}
