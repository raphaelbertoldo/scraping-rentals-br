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

type ngss struct {
}

var neighborhoodList = map[string]string{
	"Aclimação":                   "146",
	"Altamira":                    "154",
	"Alto Umuarama":               "1456",
	"Aparecida":                   "1308",
	"Área Rural":                  "1591",
	"Bom Jesus":                   "174",
	"Bosque dos Buritis":          "310",
	"Brasil":                      "182",
	"Carajás":                     "186",
	"Cazeca":                      "190",
	"Centro":                      "6",
	"Centro Empresarial Leste":    "194",
	"Chácaras Douradinho":         "198",
	"Chácaras Rancho Alegre":      "1798",
	"Chácaras Tubalina E Quartel": "1307",
	"Cidade Jardim":               "642",
	"City Uberlândia":             "238",
	"Cond. Cyrela Buritis":        "868",
	"Cond. Gávea Paradiso":        "1482",
	"Cond. GSP Arts":              "1452",
	"Cond. Jardins Barcelona":     "1510",
	"Cond. Jardins Roma":          "1518",
	"Cond. Royal Park":            "1530",
	"Cond. Varanda Sul":           "1538",
	"Condomínio Alphaville I":     "862",
	"Conjunto Alvorada":           "1345",
	"Copacabana":                  "398",
	"Custodio Pereira":            "406",
	"Daniel Fonseca":              "1636",
	"Distrito Industrial":         "1339",
	"Dona Zulmira":                "422",
	"Erlan":                       "1390",
	"Fundinho":                    "1580",
	"Gávea":                       "274",
	"Gávea Hill I":                "992",
	"Gávea Sul":                   "314",
	"General Osorio":              "1001",
	"Granada":                     "1402",
	"Grand Ville":                 "1672",
	"Granja Marileusa":            "849",
	"Guarani":                     "450",
	"Itapema Sul":                 "454",
	"Jaraguá":                     "458",
	"Jardim América I":            "470",
	"Jardim Botânico":             "482",
	"Jardim Brasília":             "10",
	"Jardim Califórnia":           "334",
	"Jardim Canaã":                "1351",
	"Jardim Celia":                "494",
	"Jardim das Acácias":          "506",
	"Jardim das Palmeiras":        "510",
	"Jardim Europa":               "1614",
	"Jardim Finotti":              "18",
	"Jardim Holanda":              "1334",
	"Jardim Inconfidência":        "326",
	"Jardim Ipanema":              "606",
	"Jardim Karaiba":              "962",
	"Jardim Ozanan":               "1019",
	"Jardim Patrícia":             "710",
	"Jardim Sucupira":             "554",
	"Jardim Sul":                  "846",
	"Lagoinha":                    "570",
	"Laranjeiras":                 "1010",
	"Lidice":                      "582",
	"Loteamento Centro Empresarial Leste III": "1737",
	"Loteamento empresarial Taiaman":          "1779",
	"Loteamento Luizote de Freitas IV":        "1803",
	"Loteamento Monte Hebron":                 "1522",
	"Loteamento Portal do Vale II":            "1462",
	"Loteamento Residencial Pequis":           "1040",
	"Lourdes":                                 "1046",
	"Luizote de Freitas":                      "1731",
	"Luizote de Freitas IV":                   "1539",
	"Mansões Aeroporto":                       "622",
	"Mansour":                                 "626",
	"Mansour III":                             "1598",
	"Maracanã":                                "630",
	"Maravilha":                               "634",
	"Marta Helena":                            "26",
	"Martins":                                 "30",
	"Minas Gerais":                            "638",
	"Miranda":                                 "1052",
	"Monte Hebron":                            "1058",
	"Morada da Colina":                        "386",
	"Morada do Sol":                           "650",
	"Morada dos Passaros":                     "654",
	"Morada Nova":                             "662",
	"Morumbi":                                 "1637",
	"Nossa Senhora Aparecida":                 "953",
	"Nossa Senhora das Graças":                "1064",
	"Nova Uberlândia":                         "302",
	"Novo Mundo":                              "322",
	"Osvaldo Rezende":                         "1109",
	"Pacaembu":                                "698",
	"Pampulha":                                "702",
	"Panorama":                                "968",
	"Patrimonio":                              "722",
	"Planalto":                                "34",
	"Presidente Roosevelt":                    "1091",
	"Residencial Fruta do Conde":              "1442",
	"Residencial Gramado":                     "742",
	"Residencial Lago Azul":                   "746",
	"Residencial Liberdade":                   "1037",
	"Residencial Pequis":                      "726",
	"Rezende":                                 "830",
	"Santa Luzia":                             "750",
	"Santa Maria":                             "42",
	"Santa Mônica":                            "1553",
	"Santa Rosa":                              "578",
	"São Jorge":                               "1082",
	"São José":                                "770",
	"Saraiva":                                 "778",
	"Segismundo Pereira":                      "1631",
	"Shopping Park":                           "1616",
	"Tabajaras":                               "798",
	"Taiaman":                                 "1640",
	"Tibery":                                  "802",
	"Tocantins":                               "806",
	"Tubalina":                                "1669",
	"Umuarama":                                "1621",
	"Varanda Sul":                             "1628",
	"Vigilato Pereira":                        "1658",
	"Vila Oswaldo":                            "1314",
	"Vila Póvoa":                              "1333",
	"Zona Rural":                              "1350",
}

// todo pass query
func (s *Service) Search(neighborhood string, min string, max string) ([]Imovel, error) {
	scraperService := scraper.NewService()
	neighborhoodId := neighborhoodList[neighborhood]
	fmt.Println("🚀 ~ file: search.go ~ line 167 ~ func ~ nh : ", neighborhoodId)
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
		chromedp.WaitVisible("#busca_detalhada", chromedp.ByID),
		chromedp.SetValue("#locacao_venda", "L", chromedp.ByID),
		chromedp.Evaluate(`document.querySelector("#locacao_venda").dispatchEvent(new Event("change"))`, nil),
		chromedp.Sleep(2*time.Second),
		chromedp.SetValue("#id_cidade", "2", chromedp.ByID),
		chromedp.Sleep(2*time.Second),
		chromedp.SetValue("#id_bairro", neighborhoodId, chromedp.ByID),
		chromedp.Sleep(2*time.Second),
		chromedp.SetValue("#vmi", min, chromedp.ByID),
		chromedp.Sleep(2*time.Second),
		chromedp.SetValue("#vma", max, chromedp.ByID),
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
