package models

type Imovel struct {
	Url      string   `json:"url"`
	Title    string   `json:"title"`
	Type     string   `json:"type"`
	Subtitle string   `json:"subtitle"` // describ
	Info     string   `json:"info"`
	Address  string   `json:"address"`
	Price    string   `json:"price"`
	Imgs     []string `json:"imgs"`
}

type Service struct{}
