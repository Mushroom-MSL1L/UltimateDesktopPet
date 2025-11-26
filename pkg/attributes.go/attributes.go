package attributes

type Attributes struct {
	Experience int   `json:"experience"`
	Water      int16 `json:"water"`
	Hunger     int16 `json:"hunger"`
	Health     int16 `json:"health"`
	Mood       int16 `json:"mood"`
	Energy     int16 `json:"energy"`
	Money      int   `json:"money"`
}
