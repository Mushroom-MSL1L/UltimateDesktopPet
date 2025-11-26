package items

import "UltimateDesktopPet/pkg/attributes.go"

type Item struct {
	ID                    uint   `json:"id"`
	Path                  string `json:"path"`
	Name                  string `json:"name"`
	Type                  string `json:"type"`
	attributes.Attributes `json:"attributes" gorm:"embedded"`
	Description           string `json:"description"`
}
