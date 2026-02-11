package mapper

import (
	"testing"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
)

func TestRefreshScoreResultFromDTO(t *testing.T) {
	tests := []struct {
		name string
		dto  *openapi.TgvalidatordRefreshScoreReply
	}{
		{
			name: "nil input returns nil",
			dto:  nil,
		},
		{
			name: "empty DTO returns result with nil scores",
			dto:  &openapi.TgvalidatordRefreshScoreReply{},
		},
		{
			name: "complete DTO maps all fields",
			dto: func() *openapi.TgvalidatordRefreshScoreReply {
				id := "score-123"
				provider := "chainalysis"
				scoreType := "risk"
				score := "75"
				updateDate := time.Now()
				return &openapi.TgvalidatordRefreshScoreReply{
					Scores: []openapi.TgvalidatordScore{
						{
							Id:         &id,
							Provider:   &provider,
							Type:       &scoreType,
							Score:      &score,
							UpdateDate: &updateDate,
						},
					},
				}
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RefreshScoreResultFromDTO(tt.dto)
			if tt.dto == nil {
				if got != nil {
					t.Errorf("RefreshScoreResultFromDTO() = %v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("RefreshScoreResultFromDTO() returned nil for non-nil input")
			}
			// Verify scores are mapped
			if tt.dto.Scores != nil && len(got.Scores) != len(tt.dto.Scores) {
				t.Errorf("Scores length = %v, want %v", len(got.Scores), len(tt.dto.Scores))
			}
		})
	}
}

func TestRefreshScoreResultFromDTO_MultipleScores(t *testing.T) {
	provider1 := "chainalysis"
	provider2 := "elliptic"
	score1 := "75"
	score2 := "80"

	dto := &openapi.TgvalidatordRefreshScoreReply{
		Scores: []openapi.TgvalidatordScore{
			{Provider: &provider1, Score: &score1},
			{Provider: &provider2, Score: &score2},
		},
	}

	got := RefreshScoreResultFromDTO(dto)
	if got == nil {
		t.Fatal("RefreshScoreResultFromDTO() returned nil for non-nil input")
	}
	if len(got.Scores) != 2 {
		t.Errorf("Scores length = %v, want 2", len(got.Scores))
	}
	if got.Scores[0].Provider != "chainalysis" {
		t.Errorf("First score provider = %v, want chainalysis", got.Scores[0].Provider)
	}
	if got.Scores[1].Provider != "elliptic" {
		t.Errorf("Second score provider = %v, want elliptic", got.Scores[1].Provider)
	}
}

func TestRefreshScoreResultFromDTO_EmptyScores(t *testing.T) {
	dto := &openapi.TgvalidatordRefreshScoreReply{
		Scores: []openapi.TgvalidatordScore{},
	}

	got := RefreshScoreResultFromDTO(dto)
	if got == nil {
		t.Fatal("RefreshScoreResultFromDTO() returned nil for non-nil input")
	}
	if len(got.Scores) != 0 {
		t.Errorf("Scores should be empty, got %v", got.Scores)
	}
}
