package testing

import (
	"testing"
)

// TestSprintPoints validates sprint race points for top 8 finishers
func TestSprintPoints(t *testing.T) {
	tests := []struct {
		position int
		expected int
	}{
		{1, 8},
		{2, 7},
		{3, 6},
		{4, 5},
		{5, 4},
		{6, 3},
		{7, 2},
		{8, 1},
		{9, 0},
		{10, 0},
		{22, 0},
	}

	for _, test := range tests {
		result := SprintPoints[test.position]
		if result != test.expected {
			t.Errorf("Expected sprint points %d for position %d, got %d", test.expected, test.position, result)
		}
	}
}

// TestRacePoints validates main race points for top 10 finishers
func TestRacePoints(t *testing.T) {
	tests := []struct {
		position int
		expected int
	}{
		{1, 25},
		{2, 18},
		{3, 15},
		{4, 12},
		{5, 10},
		{6, 8},
		{7, 6},
		{8, 4},
		{9, 2},
		{10, 1},
		{11, 0},
		{22, 0},
	}

	for _, test := range tests {
		result := RacePoints[test.position]
		if result != test.expected {
			t.Errorf("Expected race points %d for position %d, got %d", test.expected, test.position, result)
		}
	}
}

// TestPointsValidation validates the 2026 F1 points rules
func TestPointsValidation(t *testing.T) {
	// Sprint: Top 8 finishers
	// Race: Top 10 finishers
	// Team: Both cars accumulate points

	expectedMaxSprint := 24 // 8+7+6+5+4+3+2+1 for all 8 positions
	expectedMaxRace := 145 // 25+18+15+12+10+8+6+4+2+1 for top 10
	expectedTeamMax := 260 // Both cars scoring max race points

	if len(SprintPoints) < 22 {
		t.Error("SprintPoints should have 22 entries")
	}

	if len(RacePoints) < 22 {
		t.Error("RacePoints should have 22 entries")
	}

	// Verify max race points
	sum := 0
	for i, p := range RacePoints {
		sum += p
		if i >= 10 {
			break
		}
	}
	if sum != expectedMaxRace {
		t.Errorf("Expected max race points %d, got %d", expectedMaxRace, sum)
	}
}
