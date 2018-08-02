package starter

type Properties struct {
	Enabled   bool   `json:"enabled"`
	DependsOn string `json:"depends_on"`
	Before    string `json:"before"`
	After     string `json:"after"`
}
