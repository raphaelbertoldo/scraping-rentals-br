package search

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/raphaelbertoldo/scraping-rentals-br/internal/models"
)

func NewVivaService() *VivaService {
	return &VivaService{}
}

type VivaService struct{}

func (s *VivaService) Search(neighborhood string, min string, max string) (string, error) {
	fmt.Println("🚀 ~ file: VivaSearch.go ~ line 18 ~ func ~ max : ", max)
	fmt.Println("🚀 ~ file: VivaSearch.go ~ line 18 ~ func ~ min : ", min)
	fmt.Println("🚀 ~ file: VivaSearch.go ~ line 18 ~ func ~ neighborhood : ", neighborhood)
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("headless", false),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var hrefs []string

	err := chromedp.Run(ctx,
		chromedp.Navigate("https://www.vivareal.com.br/aluguel/minas-gerais/uberlandia/casa_residencial/"),
		chromedp.WaitVisible(".filters-panel__main-filters", chromedp.ByID),
		chromedp.SendKeys(".filters-panel__main-filters", "cu de egua.. ?", chromedp.ByID),
	)

	if err != nil {
		log.Fatal(err)
	}

	for i, href := range hrefs {
		fmt.Printf("Card %d: %s\n", i+1, href)
	}

	if len(hrefs) == 0 {
		fmt.Println("Nenhum href único encontrado. Verifique se os seletores estão corretos.")
	} else {
		fmt.Printf("Total de hrefs únicos encontrados: %d\n", len(hrefs))
	}

	scrapedImoveis := []models.Imovel{
		{
			Url:      "https://exemplo.com/imovel1",
			Title:    "Apartamento Moderno",
			Type:     "Apartamento",
			Subtitle: "3 quartos, 2 banheiros, 1 vaga de garagem",
			Info:     "Apartamento espaçoso com excelente localização.",
			Price:    "R$ 2.500,00 / Mês",
			Imgs:     []string{"https://exemplo.com/img1.jpg", "https://exemplo.com/img2.jpg"},
		}}

	var imoveis []models.Imovel
	for _, scrapedImovel := range scrapedImoveis {
		imoveis = append(imoveis, models.Imovel{
			Url:      scrapedImovel.Url,
			Title:    scrapedImovel.Title,
			Type:     scrapedImovel.Type,
			Subtitle: scrapedImovel.Subtitle,
			Info:     scrapedImovel.Info,
			Price:    scrapedImovel.Price,
			Imgs:     scrapedImovel.Imgs,
		})
	}
	res := "ok"
	return res, nil
}
