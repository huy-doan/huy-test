package models

// Role represents a user role in the system
type Screen struct {
	ID int `json:"id"`
	BaseColumnTimestamp

	Name       string `json:"name"`
	ScreenCode string `json:"screen_code"`
	ScreenPath string `json:"screen_path"`
}

// TableName specifies the database table name
func (Screen) TableName() string {
	return "screen"
}
