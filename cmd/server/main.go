package main

import (
	"log"

	"github.com/hphp/linkvault/internal/config"
	"github.com/hphp/linkvault/internal/db"
	"github.com/hphp/linkvault/internal/router"
)

func main() {
	cfg := config.Load()
	
	dynamoClient := db.NewDynamoClient(cfg.AWSRegion, cfg.DynamoEndpoint)
	
	r := router.New(dynamoClient, cfg.TableName)
	
	log.Printf("Starting server on %s", cfg.Port)
	if err := r.Run(cfg.Port); err != nil {
		log.Fatal(err)
	}
}