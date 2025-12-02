package items

import "UltimateDesktopPet/internal/attributes"

type Item struct {
	ID                    uint   `json:"id"`
	Path                  string `json:"path"`
	Name                  string `json:"name"`
	Type                  string `json:"type"`
	attributes.Attributes `json:"attributes" gorm:"embedded"`
	Description           string `json:"description"`
}
