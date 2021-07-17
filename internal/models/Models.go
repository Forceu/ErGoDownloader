package models

// Configuration contains the user credentials
type Configuration struct {
	Username string `json:"Username"`
	Token    string `json:"Token"`
	MaxPost  string `json:"MaxPost"`
}

// Application contains the application (reddit server-side) credentials
type Application struct {
	Id          string
	Useragent   string
	RedirectUrl string
}
