package service

import (
	"context"
	"encoding/json"
	"fmt"
	"jonathantyar/tokopedia-crawler/src/constant"
	"jonathantyar/tokopedia-crawler/src/export"
	"jonathantyar/tokopedia-crawler/src/model"
	"jonathantyar/tokopedia-crawler/src/repository"
	"log"
	"os"
	"os/exec"
	"time"
)

type Config struct {
	Url      string
	Total    uint64
	SizePool uint64
	SaveDB   bool
	SaveFile bool
}

type ScrapperServiceTokopedia struct {
	Category          string
	Filename          string
	Product           []model.Product
	ProductRepository repository.ProductRepository
	Limit             uint64
}

// HandleScraping implements ScrapperServiceInterface.
func (s *ScrapperServiceTokopedia) HandleScraping(ctx context.Context, cfg Config) {
	s.Limit = cfg.SizePool
	s.Category = cfg.Url
	s.Filename = fmt.Sprintf("%v", time.Now().Unix())

	numCalls := cfg.Total / 50
	remaining := cfg.Total % 50

	for tmp := 1; tmp <= int(numCalls); tmp++ {
		s.ScrapingData(tmp, uint64(50))
	}
	if remaining > 0 {
		s.ScrapingData(int(numCalls)+1, uint64(remaining))
	}

	if cfg.SaveDB {
		err := s.ProductRepository.Create(ctx, s.Product)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Data has been saved")
	}

	if cfg.SaveFile {
		filename, err := export.GenerateCsv(s.Filename, s.Product)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("File created : %s", filename)
	}
}

// ScrapingData implements ScrapperServiceInterface.
func (s *ScrapperServiceTokopedia) ScrapingData(page int, total uint64) {
	fmt.Printf("scraping page : %v with data : %v", page, total)

	cmd := exec.Command("python", "./other/scrapeList.py", s.Category, fmt.Sprint(page))

	// Capture the error output
	errorOutput, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("CombinedOutput output: %s\n", errorOutput)
	}
	if len(errorOutput) > 0 {
		log.Printf("Error output: %s\n", errorOutput)
	}

	var result Result
	fmt.Println("output", string(errorOutput))
	err = json.Unmarshal(errorOutput, &result)
	if err != nil {
		log.Fatal(err)
	}

	if !result.Status {
		log.Fatal(result.Message)
	}

	s.GetData(result.File, page, total)
}

// GetData implements ScrapperServiceInterface.
func (s *ScrapperServiceTokopedia) GetData(jsonFile string, page int, total uint64) {
	var data map[string]interface{}

	jsonData, err := os.ReadFile(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		log.Fatal(err)
	}

	// recommendationData, ok := data[constant.RECOMMENDATION_PRODUCT]
	// if !ok {
	// 	log.Fatalf("Recommendation Product with key '%s' not found", constant.RECOMMENDATION_PRODUCT)
	// }

	// productBytes, err := json.Marshal(recommendationData)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// var productRecommendation Recommendation
	// err = json.Unmarshal(productBytes, &productRecommendation)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	start := ((page - 1) * 60) + 1
	aceProductData, ok := data[fmt.Sprintf(constant.ACE_SEARCH_PRODUCT, page, start, page)]
	if !ok {
		log.Fatalf("Ace Product with key '%s' not found", fmt.Sprintf(constant.ACE_SEARCH_PRODUCT, page, start, page))
	}

	aceProductBytes, err := json.Marshal(aceProductData)
	if err != nil {
		log.Fatal(err)
	}

	var productAce AceProduct
	err = json.Unmarshal(aceProductBytes, &productAce)
	if err != nil {
		log.Fatal(err)
	}

	worker := 0
	limiter := make(chan int, s.Limit)
	for _, product := range productAce.Products {
		limiter <- 1
		go func(p Product) {
			detail, ok := data[p.ID]
			if !ok {
				log.Fatalf("Detail Product with key '%s' not found", p.ID)
			}
			detailBytes, err := json.Marshal(detail)
			if err != nil {
				log.Fatal(err)
			}

			var detailProduct DetailProduct
			err = json.Unmarshal(detailBytes, &detailProduct)
			if err != nil {
				log.Fatal(err)
			}

			if worker == int(total) {
				<-limiter
				return
			}

			s.GetProducts(detailProduct)
			worker++
			<-limiter
		}(product)
	}
}

// GetProducts implements ScrapperServiceInterface.
func (s *ScrapperServiceTokopedia) GetProducts(product DetailProduct) {
	cmd := exec.Command("python", "./other/scrapeProduct.py", fmt.Sprint(product.URL))

	// Capture the error output
	errorOutput, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("CombinedOutput output: %s\n", errorOutput)
	}

	var result model.Product
	err = json.Unmarshal(errorOutput, &result)
	if err != nil {
		log.Fatal(err)
	}

	mProduct := model.Product{
		Name:        result.Name,
		Description: result.Description,
		ImageLink:   result.ImageLink,
		Price:       float64(product.PriceInt),
		Rating:      result.Rating,
		Merchant:    result.Merchant,
	}

	s.Product = append(s.Product, mProduct)
}

func NewScrapperServiceTokopedia(pr repository.ProductRepository) ScrapperServiceInterface {
	return &ScrapperServiceTokopedia{
		ProductRepository: pr,
	}
}
