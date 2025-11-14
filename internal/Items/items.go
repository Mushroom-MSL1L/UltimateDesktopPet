package items

type Item struct {
	Path        string `json:"path"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Water       int16  `json:"water"`
	Hunger      int16  `json:"hunger"`
	Health      int16  `json:"health"`
	Mood        int16  `json:"mood"`
	Energy      int16  `json:"energy"`
	MoneyCost   int  `json:"moneyCost"`
	Description string `json:"description"`
}
