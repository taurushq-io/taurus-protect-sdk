package model

// RefreshScoreResult contains the result of refreshing an address score.
// The Score type is defined in whitelisted_address.go and represents a risk score from a scoring provider.
type RefreshScoreResult struct {
	// Scores contains the updated scores for the address.
	Scores []Score `json:"scores"`
}
