package main

import (
	"fmt"
	"github.com/chamodshehanka/ecr-variant-hunter/internal/config"
	"github.com/chamodshehanka/ecr-variant-hunter/services"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config: ", err.Error())
		return
	}

	ecrClient := services.GetECRConfig()
	services.DeleteOutdatedImages(ecrClient)
}
