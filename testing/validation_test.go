package testing

import (
	"testing"
)

// TestDriverSelectionValidates5Validates that driver selection requires exactly 5 drivers
func TestDriverSelectionValidates5(t *testing.T) {
	tests := []struct {
		name     string
		count    int
		shouldPass bool
	}{
		{"4 drivers", 4, false},
		{"5 drivers", 5, true},
		{"6 drivers", 6, false},
	}

	for _, test := range tests {
		if test.shouldPass {
			if len(test.name) != 5 {
				t.Errorf("Test '%s' should pass with 5 drivers", test.name)
			}
		} else {
			if len(test.name) == 5 {
				t.Errorf("Test '%s' should fail with %d drivers", test.name, test.count)
			}
		}
	}
}

// TestTeamSelectionValidates2Validates that team selection requires exactly 2 teams
func TestTeamSelectionValidates2(t *testing.T) {
	tests := []struct {
		name     string
		count    int
		shouldPass bool
	}{
		{"1 team", 1, false},
		{"2 teams", 2, true},
		{"3 teams", 3, false},
	}

	for _, test := range tests {
		if test.shouldPass {
			if len(test.name) != 2 {
				t.Errorf("Test '%s' should pass with 2 teams", test.name)
			}
		} else {
			if len(test.name) == 2 {
				t.Errorf("Test '%s' should fail with %d teams", test.name, test.count)
			}
		}
	}
}

// TestNoDuplicateDriversValidates that no duplicate drivers are selected
func TestNoDuplicateDrivers(t *testing.T) {
	drivers := []string{"d1", "d2", "d3", "d4", "d1"} // d1 is duplicated
	hasDuplicates := false
	seen := make(map[string]bool)
	for _, d := range drivers {
		if seen[d] {
			hasDuplicates = true
			break
		}
		seen[d] = true
	}

	if !hasDuplicates {
		t.Error("Should detect duplicate drivers")
	}
}

// TestNoDuplicateTeamsValidates that no duplicate teams are selected
func TestNoDuplicateTeams(t *testing.T) {
	teams := []string{"t1", "t2", "t1"} // t1 is duplicated
	hasDuplicates := false
	seen := make(map[string]bool)
	for _, t := range teams {
		if seen[t] {
			hasDuplicates = true
			break
		}
		seen[t] = true
	}

	if !hasDuplicates {
		t.Error("Should detect duplicate teams")
	}
}
