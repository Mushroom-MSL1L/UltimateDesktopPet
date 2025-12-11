package activities

import "github.com/Mushroom-MSL1L/UltimateDesktopPet/desktop_pet/internal/attributes"

type Activity struct {
	ID                    uint   `json:"id"`
	Path                  string `json:"path"`
	Name                  string `json:"name"`
	Type                  string `json:"type"`
	attributes.Attributes `json:"attributes" gorm:"embedded"`
	DurationMinute        int16  `json:"duration_minute"`
	Description           string `json:"description"`
}
