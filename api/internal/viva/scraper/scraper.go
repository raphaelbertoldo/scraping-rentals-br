package scraper

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/raphaelbertoldo/scraping-rentals-br/api/internal/models"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Scraper(urls []string) []models.Imovel {
	var imoveis []models.Imovel

	for _, url := range urls {
		opts := append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("no-sandbox", true),
			chromedp.Flag("disable-setuid-sandbox", true),
			chromedp.Flag("disable-blink-features", "AutomationControlled"),
			chromedp.Flag("window-size", "1920,1080"),
			chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.0.0 Safari/537.36"),
		)

		allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
		defer cancel()

		ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
		defer cancel()

		timeout := 60 * time.Second
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()

		fmt.Println("[VIVA - SCRAPER] Processando URL:", url)

		var title, propertyType, subtitle, info, address, price string
		var imgUrls []string
		var body string
		err := chromedp.Run(ctx,
			chromedp.Navigate(url),
			chromedp.WaitReady("body", chromedp.ByQuery),
			chromedp.OuterHTML("html", &body, chromedp.ByQuery),
			chromedp.Text(".price-info-value", &price, chromedp.ByQuery),
			chromedp.Text(".description__title", &title, chromedp.ByQuery),
			chromedp.Text(".description__title", &subtitle, chromedp.ByQuery),
			chromedp.Text(".address-info-value", &address, chromedp.ByQuery),
			chromedp.Text(".description__content--text", &info, chromedp.ByQuery),

			chromedp.Evaluate(`
				Array.from(document.querySelectorAll('.carousel-photos--item > picture > source')).map(src => {
				console.log(">>>>")
				return src.srcset
				})
			`, &imgUrls),
		)

		// fmt.Println("body ::::", body)

		if err != nil {
			log.Printf("Erro ao processar %s: %v", url, err)
			continue
		}

		title = strings.TrimSpace(title)
		propertyType = strings.TrimSpace(propertyType)
		subtitle = strings.TrimSpace(subtitle)
		info = strings.TrimSpace(info)
		address = strings.TrimSpace(address)
		price = strings.TrimSpace(price)

		imovel := models.Imovel{
			Url:      url,
			Title:    title,
			Type:     propertyType,
			Subtitle: subtitle,
			Info:     info,
			Address:  address,
			Price:    price,
			Imgs:     imgUrls,
		}

		imoveis = append(imoveis, imovel)

		fmt.Printf("[VIVA - SCRAPER] Imóvel processado: %+v\n", imovel)

		time.Sleep(5 * time.Second)
	}

	return imoveis
}
