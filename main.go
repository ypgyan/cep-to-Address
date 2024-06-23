package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gofor-little/env"
)

const (
	addressFilePath = "./files/addresses.csv"
	cepFilePath     = "./files/ceps.csv"
)

func main() {
	loadEnv()

	log.Println("Starting address search script")

	records := getCEPs()
	var ceps []string
	for _, record := range records {
		ceps = append(ceps, record[0])
	}

	writer, addressesFile := initAddressFile()
	defer func(addressesFile *os.File) {
		err := addressesFile.Close()
		if err != nil {
			log.Fatalf("Error closing file: %v", err)
		}
	}(addressesFile)
	defer writer.Flush()

	log.Println("Starting search on google")

	for _, cep := range ceps {
		addressData := searchOnGoogle(cep)
		if len(addressData.Results) == 0 {
			log.Println("No results found for CEP", cep)
			continue
		}
		address := extractAddress(cep, addressData)
		writeToFile(writer, address)
	}

	log.Println("Finished address search script")
}

func loadEnv() {
	err := env.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func getCEPs() [][]string {
	log.Println("Loading CEPs from file")
	file, err := os.Open(cepFilePath)
	if err != nil {
		log.Fatalf("Error occurred while opening file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("failed to close file: %s", err)
		}
	}(file)
	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		log.Fatalf("Error occurred while reading file: %v", err)
	}
	return records
}

func initAddressFile() (writer *csv.Writer, addressesFile *os.File) {
	addressesFile, err := os.Create(addressFilePath)
	if err != nil {
		log.Fatalf("Error occurred while creating file: %v", err)
	}
	writer = csv.NewWriter(addressesFile)
	writer.Comma = ';'
	return
}

func searchOnGoogle(c string) GoogleResponse {
	url := env.Get("GOOGLE_MAPS_URL", "")
	resp, err := http.Get(url + "?address=" + c + "&key=" + env.Get("GOOGLE_MAPS_API_KEY", ""))
	if err != nil {
		log.Fatalf("Error occurred while fetching data from Google: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalf("Error closing body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Request failed with status code: %d", resp.StatusCode)
	}

	var response GoogleResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Fatalf("Error occurred while decoding response: %v", err)
	}

	return response
}

func extractAddress(cep string, response GoogleResponse) *Address {
	result := response.Results[0]
	lat, lng := result.Geometry.Location.Lat, result.Geometry.Location.Lng

	return &Address{
		Cep:              cep,
		Latitude:         lat,
		Longitude:        lng,
		FormattedAddress: result.FormattedAddress,
	}
}

func writeToFile(writer *csv.Writer, address *Address) {
	lineToWrite := addressToStrings(address)

	err := writer.Write(lineToWrite)
	if err != nil {
		log.Println("Error occurred while writing to file:", err)
	}
}

func addressToStrings(address *Address) []string {
	return []string{address.Cep, fmt.Sprintf("%f", address.Latitude), fmt.Sprintf("%f", address.Longitude), address.FormattedAddress}
}
