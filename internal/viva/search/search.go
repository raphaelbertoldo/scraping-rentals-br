package search

func NewService() *Service {
	return &Service{}
}

type Service struct{}

type Imovel struct {
	Url      string   `json:"url"`
	Title    string   `json:"title"`
	Type     string   `json:"type"`
	Subtitle string   `json:"subtitle"`
	Info     string   `json:"info"`
	Price    string   `json:"price"`
	Imgs     []string `json:"imgs"`
}
