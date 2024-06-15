package main

import (
    "encoding/csv"
    "fmt"
    "os"
)

type Cep struct {
	CEP string
}

func main() {
	file, err := os.Open("ceps.csv")
    if err != nil {
        panic(err)
    }
	defer file.Close()

	records, err := csv.NewReader(file).ReadAll()
    if err != nil {
        panic(err)
    }

	ceps := make([]Cep, len(records))
	for k, record := range records {
		cep := Cep{
			CEP: record[0],
		}
		ceps[k] = cep
	}
	fmt.Println(ceps)
}