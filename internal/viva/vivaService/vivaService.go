package vivaService

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/raphaelbertoldo/scraping-rentals-br/internal/models"
	"github.com/raphaelbertoldo/scraping-rentals-br/internal/viva/scraper"
)

func NewService() *Service {
	return &Service{}
}

type Service struct{}

func (s *Service) Search(neighborhood string, min string, max string) ([]models.Imovel, error) {

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
		chromedp.Sleep(2*time.Second),

		chromedp.SendKeys("#filter-location-search-input", neighborhood, chromedp.ByID),
		chromedp.Sleep(3*time.Second),

		chromedp.Click(`.autocomplete__list > li[data-type="neighborhood"]`, chromedp.ByQuery),
		chromedp.Sleep(2*time.Second),

		chromedp.SendKeys("#filter-range-from-price", min, chromedp.ByID),
		chromedp.Sleep(2*time.Second),

		chromedp.SendKeys("#filter-range-to-price", max+"\n", chromedp.ByID),

		chromedp.Sleep(6*time.Second),

		chromedp.Evaluate(`
			(() => {
				const uniqueHrefs = new Set();
				document.querySelectorAll('.property-card__content-link').forEach(el => {
					if (el.href) {
						uniqueHrefs.add(el.href);
					}
				});
				return Array.from(uniqueHrefs)
			})()
		`, &hrefs),
	)

	if err != nil {
		log.Fatal(err)
	}

	for i, href := range hrefs {
		fmt.Printf("Card %d: %s\n", i+1, href)
	}

	if len(hrefs) == 0 {
		fmt.Println("[VIVA] - Nenhum href único encontrado. Verifique se os seletores estão corretos.")
	} else {
		fmt.Printf("Total de hrefs únicos encontrados: %d\n", len(hrefs))
	}
	scraperService := scraper.NewService()
	scrapedImoveis := scraperService.Scraper(hrefs)

	var imoveis []models.Imovel
	for _, scrapedImovel := range scrapedImoveis {
		imoveis = append(imoveis, models.Imovel{
			Url:      scrapedImovel.Url,
			Title:    scrapedImovel.Title,
			Subtitle: scrapedImovel.Subtitle,
			Address:  scrapedImovel.Address,
			Type:     scrapedImovel.Type,
			Info:     scrapedImovel.Info,
			Price:    scrapedImovel.Price,
			Imgs:     scrapedImovel.Imgs,
		})
	}

	return imoveis, nil
}
