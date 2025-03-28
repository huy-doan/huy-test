package ExampleShell

import (
	"log"

	"github.com/huydq/demo/cmd/service"
)

func Execute() {
	log.Println("======= Start Example Shell ======= ")
	defer log.Println("======= Stop Example Shell ======= ")

	batchService, err := service.NewBatchService()
	if err != nil {
		log.Fatalf("Failed to initialize batch service: %v", err)
	}
	defer batchService.Close()
}
