package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gocolly/colly"
)

// type Images map[string]string
type Imovel struct {
	Url      string   `json:"url"`
	Title    string   `json:"title"`
	Subtitle string   `json:"subtitle"`
	Price    string   `json:"price"`
	Imgs     []string `json:"Imgs"`
}

func main() {
	collector := colly.NewCollector()
	var imoveis []Imovel

	collector.OnError(func(r *colly.Response, e error) {
		fmt.Println("Blimey, an error occurred!:", e)
	})

	collector.OnHTML(".container", func(section *colly.HTMLElement) {
		imovel := Imovel{}
		if section.ChildText(".row") != "" && section.ChildText(".titulo-imovel") != "" {
			imovel.Title = section.ChildText(".titulo-imovel")
			imovel.Subtitle = section.ChildText(".subtitulo-imovel")
			section.ForEach(".text-end", func(i int, elements *colly.HTMLElement) {
				if elements.ChildText("strong") != "" {
					imovel.Price = elements.ChildText("strong")
				}
			})
			jsonData, err := json.MarshalIndent(imovel, "", "  ")
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(jsonData))
			imoveis = append(imoveis, imovel)
		}
	})
	//TODO - GET IMAGES IMVOEL

	collector.Visit("https://www.ivannegocios.com.br/alugar/Uberlandia/Comercial/Loja/Centro/66345")

	// Convertendo os dados para JSON formatado
	// jsonData, err := json.MarshalIndent(imoveis, "", "  ")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Imprimindo os dados formatados
	// fmt.Println(string(jsonData))
}
