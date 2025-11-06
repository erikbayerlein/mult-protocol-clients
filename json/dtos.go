package json

type Auth struct {
	Type      string `json:"tipo"`
	StudentId string `json:"aluno_id"`
}

type Logout struct {
	Type  string `json:"tipo"`
	Token string `json:"token"`
}

type Operation struct {
	Type      string `json:"tipo"`
	Operation string `json:"operacao"`
	Token     string `json:"token"`
	Params    any    `json:"parametros"`
}

type EchoParams struct {
	Message string `json:"mensagem"`
}

type SumParams struct {
	Numeros []int `json:"numeros"`
}

type TimestampParams struct{}

type StatusParams struct {
	Detalhado bool `json:"detalhado"`
}

type HistoryParams struct {
	Limite int `json:"limite"`
}
