package params

type Login struct {
	Email string `json:"email"`
	Pwd   string `json:"pwd"`
}

type User struct {
	Token  string `json:"token"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Pwd    string `json:"pwd"`
	Scope  string `json:"scope"`
	Access string `json:"access"`
}
