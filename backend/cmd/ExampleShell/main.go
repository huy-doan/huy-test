package ExampleShell

import (
	"log"

	"github.com/huydq/test/cmd/service"
)

func Execute() {
	log.Println("======= Start Example Shell ======= ")
	defer log.Println("======= Stop Example Shell ======= ")

	batchService, err := service.NewBatchService()
	if err != nil {
		log.Fatalf("Failed to initialize batch service: %v", err)
	}

	// Check for errors when closing batch service
	defer func() {
		if err := batchService.Close(); err != nil {
			log.Printf("Error closing batch service: %v", err)
		}
	}()
}
