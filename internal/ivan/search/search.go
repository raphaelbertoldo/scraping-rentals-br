package search

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/raphaelbertoldo/scraping-rentals-br/internal/ivan/scraper"
)

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

func (s *Service) Search(query string) ([]Imovel, error) {
	scraperService := scraper.NewService()

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
		chromedp.Navigate("https://www.ivannegocios.com.br/"),
		chromedp.WaitVisible("#locacao_venda", chromedp.ByID),
		chromedp.SetValue("#locacao_venda", "L", chromedp.ByID),
		chromedp.Evaluate(`document.querySelector("#locacao_venda").dispatchEvent(new Event("change"))`, nil),
		chromedp.Sleep(2*time.Second),

		chromedp.Click("button.btn.btn-primary.cantos-arredondados.loading[type='submit']", chromedp.ByQuery),

		chromedp.Sleep(5*time.Second),
		chromedp.WaitVisible(".muda_card1", chromedp.ByQuery),

		chromedp.Evaluate(`
			(() => {
				const uniqueHrefs = new Set();
				document.querySelectorAll('.muda_card1 .carousel-cell').forEach(el => {
					if (el.href) {
						uniqueHrefs.add(el.href);
					}
				});
				return Array.from(uniqueHrefs);
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
		fmt.Println("Nenhum href único encontrado. Verifique se os seletores estão corretos.")
	} else {
		fmt.Printf("Total de hrefs únicos encontrados: %d\n", len(hrefs))
	}

	scrapedImoveis := scraperService.Scraper(hrefs)

	var imoveis []Imovel
	for _, scrapedImovel := range scrapedImoveis {
		imoveis = append(imoveis, Imovel{
			Url:      scrapedImovel.Url,
			Title:    scrapedImovel.Title,
			Type:     scrapedImovel.Type,
			Subtitle: scrapedImovel.Subtitle,
			Info:     scrapedImovel.Info,
			Price:    scrapedImovel.Price,
			Imgs:     scrapedImovel.Imgs,
		})
	}

	return imoveis, nil
}
