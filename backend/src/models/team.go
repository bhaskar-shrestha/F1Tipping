package models

// Team represents a constructor/team with their race positions
type Team struct {
	ID                 string  `json:"id"`
	ConstructorID      string  `json:"constructor_id"`
	ConstructorName    string  `json:"constructor_name"`
	RaceCar1Position   *int    `json:"race_car1_position"`
	RaceCar2Position   *int    `json:"race_car2_position"`
	SprintCar1Position *int    `json:"sprint_car1_position"`
	SprintCar2Position *int    `json:"sprint_car2_position"`
}
