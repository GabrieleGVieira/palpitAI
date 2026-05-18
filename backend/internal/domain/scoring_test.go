package domain

import "testing"

func TestScorePrediction(t *testing.T) {
	tests := []struct {
		name          string
		predictedHome int
		predictedAway int
		finalHome     int
		finalAway     int
		want          int
	}{
		{
			name:          "exact score",
			predictedHome: 2,
			predictedAway: 1,
			finalHome:     2,
			finalAway:     1,
			want:          10,
		},
		{
			name:          "same winner",
			predictedHome: 3,
			predictedAway: 0,
			finalHome:     1,
			finalAway:     0,
			want:          5,
		},
		{
			name:          "same draw",
			predictedHome: 0,
			predictedAway: 0,
			finalHome:     2,
			finalAway:     2,
			want:          5,
		},
		{
			name:          "wrong result",
			predictedHome: 1,
			predictedAway: 0,
			finalHome:     0,
			finalAway:     1,
			want:          0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ScorePrediction(tt.predictedHome, tt.predictedAway, tt.finalHome, tt.finalAway)
			if got != tt.want {
				t.Fatalf("expected %d, got %d", tt.want, got)
			}
		})
	}
}
