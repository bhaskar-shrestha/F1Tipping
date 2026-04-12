package models

// Team represents a constructor/team with their race positions
// ID field is the primary identifier (matches constructor_id)
type Team struct {
	ID                 string  `json:"id"`                    // Primary identifier, references constructors.constructor_id
	ConstructorName    string  `json:"constructor_name"`
	RaceCar1Position   *int    `json:"race_car1_position"`    // NULL if not finished
	RaceCar2Position   *int    `json:"race_car2_position"`    // NULL if not finished
	SprintCar1Position *int    `json:"sprint_car1_position"`  // NULL if not finished
	SprintCar2Position *int    `json:"sprint_car2_position"`  // NULL if not finished
}
