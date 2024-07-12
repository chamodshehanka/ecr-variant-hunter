package main

import (
	"fmt"
	"github.com/chamodshehanka/ecr-variant-hunter/internal/config"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config: ", err.Error())
		return
	}
	// Initialize AWS session
	//sess, _ := session.NewSession(&aws.Config{
	//	Region: aws.String("us-west-2"),
	//})
	//
	//svc := ecr.New(sess)

	// Example: Iterate over repoNames and delete outdated images
	// This is a simplified example. Implement your logic for listing and deleting images.
}
