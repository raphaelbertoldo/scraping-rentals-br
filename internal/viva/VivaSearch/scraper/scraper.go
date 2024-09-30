package scraper

type Imovel struct {
	Url      string   `json:"url"`
	Title    string   `json:"title"`
	Type     string   `json:"type"`
	Subtitle string   `json:"subtitle"`
	Info     string   `json:"info"`
	Price    string   `json:"price"`
	Imgs     []string `json:"imgs"`
}

func NewService() *Service {
	return &Service{}
}

type Service struct{}

func (s *Service) ScraperMocked() []Imovel {
	imoveis := []Imovel{
		{
			Url:      "https://exemplo.com/imovel1",
			Title:    "Apartamento Moderno",
			Type:     "Apartamento",
			Subtitle: "3 quartos, 2 banheiros, 1 vaga de garagem",
			Info:     "Apartamento espaçoso com excelente localização.",
			Price:    "R$ 2.500,00 / Mês",
			Imgs:     []string{"https://exemplo.com/img1.jpg", "https://exemplo.com/img2.jpg"},
		},
		{
			Url:      "https://exemplo.com/imovel2",
			Title:    "Casa de Campo",
			Type:     "Casa",
			Subtitle: "4 quartos, 3 banheiros, piscina",
			Info:     "Linda casa de campo com vista para as montanhas.",
			Price:    "R$ 5.000,00 / Mês",
			Imgs:     []string{"https://exemplo.com/img3.jpg", "https://exemplo.com/img4.jpg"},
		},
		{
			Url:      "https://exemplo.com/imovel3",
			Title:    "Cobertura de Luxo",
			Type:     "Cobertura",
			Subtitle: "2 quartos, 1 suíte, jacuzzi",
			Info:     "Cobertura com design moderno e acabamento de alto padrão.",
			Price:    "R$ 8.500,00 / Mês",
			Imgs:     []string{"https://exemplo.com/img5.jpg", "https://exemplo.com/img6.jpg"},
		},
	}
	return imoveis
}
