package gorm

type properties struct {
	Type      string `json:"type"` // mysql, postgres, sqlite3, mssql,
	Host      string `json:"host"`
	Port      string `json:"port"`
	Database  string `json:"database"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Charset   string `json:"charset"`
	ParseTime string `json:"parse_time"`
	Loc       string `json:"loc"`
}
