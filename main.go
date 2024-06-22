package main

import (
	"encoding/csv"
	"fmt"
	"github.com/gofor-little/env"
	"log"
	"os"
)

func main() {
	loadEnv()
	fmt.Println("Iniciando script de busca de endereços")
	records := getCEPs()

	var ceps []string

	for _, record := range records {
		ceps = append(ceps, record[0])
	}

	addressesFile, err := os.Create("./files/addresses.csv")
	if err != nil {
		panic(err)
	}
	defer func(addressesFile *os.File) {
		err := addressesFile.Close()
		if err != nil {
			log.Fatalf("Error closing file: %v", err)
		}
	}(addressesFile)

	writer := csv.NewWriter(addressesFile)
	defer writer.Flush()

	ceptowrite := make([][]string, len(ceps))
	for i, c := range ceps {
		// TODO: Implementar busca no google maps
		ceptowrite[i] = []string{c}
		if err := writer.Write(ceptowrite[i]); err != nil {
			panic(err)
		}
	}
	fmt.Println("Script finalizado")
}

func getCEPs() [][]string {
	file, err := os.Open("./files/ceps.csv")
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("failed to close file: %s", err)
		}
	}(file)

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		panic(err)
	}
	return records
}

func loadEnv() {
	err := env.Load(".env")
	if err != nil {
		log.Fatalf("Falha ao localizar arquivo .env encerrando script")
	}
}

//func readJson() {
//	// Your JSON string
//	var jsonStr = []byte(`{
//		"name": "John",
//		"age": 30,
//		"city": "New York",
//		"social": {
//			"twitter": "@john",
//			"email": "john@example.com"
//		}
//	}`)
//
//	// Use a map of string keys to values of any type
//	var result map[string]interface{}
//
//	// Unmarshal the JSON string into the map
//	if err := json.Unmarshal(jsonStr, &result); err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(result)
//
//	// Print the map
//	for key, value := range result {
//		fmt.Printf("%s: %v\n", key, value)
//	}
//}
