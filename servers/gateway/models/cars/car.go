package cars

// Car represents a single Car that has been registered by a user in the database
type Car struct {
	ID int64 `json:"id"`
	LicensePlate string `json:"licensePlate"`
	UserID int64 `json:"userId"`
	Make string `json:"make"`
	Model string `json:"model"`
	ModelYear string `json:"year"`
	Color string `json:"color"`

}

// UpdateCar represents possible updates to a specific vehicle's information
type UpdateCar struct {
	LicensePlate string `json:"licensePlate"`
	UserID int64 `json:"userId"`
	Make string `json:"make"`
	Model string `json:"model"`
	ModelYear string `json:"year"`
	Color string `json:"color"`

}