package main

import (
	"log"
)

func main() {
	log.Println("Starting cron");
	/*
	- one by one filter the vendor by its id in supply-vendor-table
	- avg the price
	- send the the one tuple to ML model to give score for the given vendor-id.
	- this same process, goes to other vendors as well.
	- after every update have a delay... 
	*/
}