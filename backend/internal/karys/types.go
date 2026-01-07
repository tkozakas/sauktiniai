package karys

type Person struct {
	Pos        string `json:"pos"`
	Number     string `json:"-"`
	Name       string `json:"name"`
	Lastname   string `json:"lastname"`
	Bdate      string `json:"bdate"`
	Department string `json:"-"`
	Info       string `json:"-"`
}

type Region struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
