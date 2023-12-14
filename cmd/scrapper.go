package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"jonathantyar/tokopedia-crawler/src"
	"jonathantyar/tokopedia-crawler/src/service"

	"github.com/spf13/cobra"
)

var ScrapperCmd = &cobra.Command{
	Use:   "scrapper",
	Short: "Use to scrapping data from website",
	Run:   scrapperCommand,
}

var category string
var total, thread uint64
var db, file bool

func init() {
	ScrapperCmd.Flags().StringVarP(&category, "category", "c", "", "Category name for scraping")
	ScrapperCmd.Flags().Uint64VarP(&total, "total", "t", 100, "Total of data")
	ScrapperCmd.Flags().Uint64VarP(&thread, "thread", "", 5, "MultiThread scraping")
	ScrapperCmd.Flags().BoolVarP(&db, "db", "", false, "Save data to DB")
	ScrapperCmd.Flags().BoolVarP(&file, "file", "", false, "Save data to File")
}

func scrapperCommand(cmd *cobra.Command, args []string) {
	tokpedService, err := src.InitializeTokopediaService()
	if err != nil {
		panic(fmt.Errorf("error initialize tokpedService service: %w", err))
	}

	cfg := service.Config{
		Url:      category,
		Total:    total,
		SizePool: thread,
		SaveDB:   db,
		SaveFile: file,
	}

	b, _ := json.Marshal(cfg)
	fmt.Println("Running tokopedia service with config : ", string(b))
	tokpedService.HandleScraping(
		context.Background(),
		cfg,
	)
}
