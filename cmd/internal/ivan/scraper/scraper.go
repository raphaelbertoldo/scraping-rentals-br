package scraper

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
	"github.com/raphaelbertoldo/scraping-rentals-br/cmd/internal/models"
)

func NewService() *Service {
	return &Service{}
}

type Service struct{}

// type Imovel struct {
// 	Url      string   `json:"url"`
// 	Title    string   `json:"title"`
// 	Type     string   `json:"type"`
// 	Subtitle string   `json:"subtitle"`
// 	Info     string   `json:"info"`
// 	Price    string   `json:"price"`
// 	Imgs     []string `json:"imgs"`
// }

func (s *Service) Scraper(urls []string) []models.Imovel {
	var imoveis []models.Imovel

	for _, url := range urls {
		fmt.Println("Processando URL:", url)

		collector := colly.NewCollector()

		collector.OnError(func(r *colly.Response, e error) {
			fmt.Println("Ocorreu um erro:", e)
		})

		imovel := models.Imovel{
			Url: url,
		}

		collector.OnHTML("section", func(section *colly.HTMLElement) {
			if section.ChildText(".row") != "" && section.ChildText(".titulo-imovel") != "" {
				imovel.Title = section.ChildText(".titulo-imovel")
				imovel.Subtitle = section.ChildText(".subtitulo-imovel")
				section.ForEach(".valores_imovel", func(_ int, elements *colly.HTMLElement) {
					elements.ForEach(".row", func(_ int, elements_ *colly.HTMLElement) {
						if strings.Contains(elements_.Text, "Total / Mês") {
							imovel.Price = elements_.ChildText("strong")
						}
					})
				})
			}

			section.ForEach(".tipo-prop", func(_ int, section_ *colly.HTMLElement) {
				imovel.Type = section_.ChildText("strong")
			})

			section.ForEach(".card-imo-radius", func(_ int, section_ *colly.HTMLElement) {
				if section_.ChildText("p") != "" && section_.ChildText(".descricao-imovel") != "" {
					imovel.Info = section_.ChildText("p")
				}
			})

			section.ForEach("#slide_fotos", func(_ int, section *colly.HTMLElement) {
				section.ForEach(".img-slider", func(i int, elements *colly.HTMLElement) {
					imovel.Imgs = append(imovel.Imgs, elements.Attr("src"))
				})
			})
		})

		collector.Visit(imovel.Url)

		imoveis = append(imoveis, imovel)
	}

	return imoveis
}
