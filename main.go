package main

import (
    "encoding/csv"
    "os"
)

func main() {
	file, err := os.Open("./files/ceps.csv")
    if err != nil {
        panic(err)
    }
	defer file.Close()

	records, err := csv.NewReader(file).ReadAll()
    if err != nil {
        panic(err)
    }

	var cep []string

	for _, record := range records {
		cep = append(cep, record[0])
	}

	addressesFile, err := os.Create("./files/addresses.csv")
    if err != nil {
        panic(err)
    }
	defer addressesFile.Close()

	writer := csv.NewWriter(addressesFile)
	defer writer.Flush()

	ceptowrite := make([][]string, len(cep))
	for i, c := range cep {
		ceptowrite[i] = []string{c}
		if err := writer.Write(ceptowrite[i]); err != nil {
			panic(err)
		}
	}
}