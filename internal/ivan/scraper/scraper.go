package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

type Imovel struct {
	Url      string   `json:"url"`
	Title    string   `json:"title"`
	Type     string   `json:"type"`
	Subtitle string   `json:"subtitle"`
	Info     string   `json:"Info"`
	Price    string   `json:"price"`
	Imgs     []string `json:"imgs"`
}

func main() {

	collector := colly.NewCollector()
	collector.OnError(func(r *colly.Response, e error) {
		fmt.Println("Blimey, an error occurred!:", e)
	})

	imovel := Imovel{}
	url := os.Args[1]
	imovel.Url = url

	collector.OnHTML("section", func(section *colly.HTMLElement) {
		if section.ChildText(".row") != "" && section.ChildText(".titulo-imovel") != "" {
			imovel.Title = section.ChildText(".titulo-imovel")
			imovel.Subtitle = section.ChildText(".subtitulo-imovel")
			section.ForEach(".text-end", func(i int, elements *colly.HTMLElement) {
				if elements.ChildText("strong") != "" {
					imovel.Price = elements.ChildText("strong")
				}
			})

		}
		fmt.Println(section.Attr("title"))
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
		jsonData, err := json.MarshalIndent(imovel, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(jsonData))
	})

	collector.Visit(imovel.Url)
}
