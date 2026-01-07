package karys

type Person struct {
	Pos        string `json:"pos"`
	Number     string `json:"number"`
	Name       string `json:"name"`
	Lastname   string `json:"lastname"`
	Bdate      string `json:"bdate"`
	Department string `json:"department"`
	Info       string `json:"info"`
}

type Region struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
