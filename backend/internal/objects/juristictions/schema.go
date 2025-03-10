package juristictions

type JuristictionInformation struct {
	Country        string                 `json:"country"`
	State          string                 `json:"state"`
	Municipality   string                 `json:"municipality"`
	Agency         string                 `json:"agency"`
	ProceedingName string                 `json:"proceeding_name"`
	ExtraObject    map[string]interface{} `json:"extra_object"`
}
