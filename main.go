package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type SymbolDetail struct {
	ID            string `json:"id"`
	BaseCurrency  string `json:"baseCurrency"`
	QuoteCurrency string `json:"quoteCurrency"`
}
type ArrSymbolDetail []SymbolDetail

type CurrDetail struct {
	ID       string `json:"id"`
	FullName string `json:"fullName"`
}
type ArrCurrDetail []CurrDetail
type CurrenciesData struct {
	Symbol      string `json:"symbol"`
	Ask         string `json:"ask"`
	Bid         string `json:"bid"`
	Last        string `json:"last"`
	Open        string `json:"open"`
	Low         string `json:"low"`
	High        string `json:"high"`
	FeeCurrency string `json:"feeCurrency"`
}
type ArrCurrenciesData []CurrenciesData

type responseData struct {
	CryptoId    string
	ID          string `json:"id"`
	FullName    string `json:"fullName"`
	Ask         string `json:"ask"`
	Bid         string `json:"bid"`
	Last        string `json:"last"`
	Open        string `json:"open"`
	Low         string `json:"low"`
	High        string `json:"high"`
	FeeCurrency string `json:"feeCurrency"`
}
type AllData []responseData

func getCurrencyDetails(w http.ResponseWriter, r *http.Request) {
	id := strings.Split(r.URL.Path, "/")
	switch r.Method {
	case "GET":
		w.WriteHeader(http.StatusOK)
		if id[len(id)-1] != "all" {
			if id[len(id)-1] != "" {
				s, err := checkSymbol(w, id[len(id)-1])
				if err != nil {
					w.Write([]byte(`{"message": "` + string(err.Error()) + `"}`))
				}
				a, _ := getSymboldetail(w, s.BaseCurrency)
				c, _ := getImpotantData(w, s.ID)
				r := responseData{a.ID, s.BaseCurrency, a.FullName, c.Ask, c.Bid, c.Last, c.Open, c.Low, c.High, c.FeeCurrency}
				data, err := json.Marshal(r)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Fprint(w, string(data))
			}
		} else {
			//getSymboldetail(w, "")
			ASD, err := AllSymbol(w)
			//fmt.Println("Curr Detail ", ASD)
			if err != nil {
				w.Write([]byte(`{"message": "` + string(err.Error()) + `"}`))
			}
			allData := AllData{}
			CD, err := AllSymboldetail()
			//fmt.Println("Curr Detail ", CD)
			if err != nil {
				w.Write([]byte(`{"message": "` + string(err.Error()) + `"}`))
			}
			AID, err := AllImpotantData()
			fmt.Println("Curr Detail ", AID)
			if err != nil {
				w.Write([]byte(`{"message": "` + string(err.Error()) + `"}`))
			}
			for _, val := range ASD {
				//r := responseData{val.BaseCurrency, CD[0].FullName, AID[0].Ask, AID[0].Bid, AID[0].Last, AID[0].Open, AID[0].Low, AID[0].High, AID[0].FeeCurrency}
				r := responseData{}
				r.ID = val.BaseCurrency
				r.CryptoId = val.ID
				for _, val2 := range CD {
					if val2.ID == val.BaseCurrency {
						r.FullName = val2.FullName
						break
					}
				}
				for _, val2 := range AID {
					if val2.Symbol == val.ID {
						r.Ask = val2.Ask
						r.FeeCurrency = val2.FeeCurrency
						r.Low = val2.Low
						r.High = val2.High
						r.Open = val2.Open
						r.Bid = val2.Bid
						r.Last = val2.Last
						break
					}
				}
				allData = append(allData, r)
				//break
			}
			data, err := json.Marshal(allData)
			//fmt.Println(allData)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Fprint(w, string(data))
		}
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "not found"}`))
	}

}

func syncRealTimeData() {

}

func handleRequests() {
	http.HandleFunc("/currency/", getCurrencyDetails)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	handleRequests()
}

func AllSymboldetail() (ArrCurrDetail, error) {
	url := "https://api.hitbtc.com/api/2/public/currency"
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return ArrCurrDetail{}, err
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	//fmt.Println(string(responseData))
	if err != nil {
		log.Fatal(err)
		return ArrCurrDetail{}, err
	}
	r := []CurrDetail{}
	err = json.Unmarshal(responseData, &r)
	if err != nil {
		log.Fatal(err)
	}
	if len(r) == 0 {
		fmt.Println("No Data Found")
		return ArrCurrDetail{}, errors.New("No Data Found")
	}
	fmt.Println(r)
	return r, nil
}

func getSymboldetail(w http.ResponseWriter, curr string) (CurrDetail, error) {
	fmt.Println("I am here ", curr)
	url := "https://api.hitbtc.com/api/2/public/currency/" + curr
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return CurrDetail{}, err
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	fmt.Println(string(responseData))
	if err != nil {
		log.Fatal(err)
		return CurrDetail{}, err
	}
	r := CurrDetail{}
	err = json.Unmarshal(responseData, &r)
	if err != nil {
		log.Fatal(err)
	}
	if (CurrDetail{}) == r {
		fmt.Println("Invalid Symbol")
		return CurrDetail{}, errors.New("Invalid Symbol")
	}
	return r, nil
}

func getImpotantData(w http.ResponseWriter, ID string) (CurrenciesData, error) {
	url := "https://api.hitbtc.com/api/2/public/ticker/" + ID
	if ID == "" {
		url = "https://api.hitbtc.com/api/2/public/ticker"
	}
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return CurrenciesData{}, err
	}
	defer response.Body.Close()

	responseFromURL, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		return CurrenciesData{}, err
	}
	r := CurrenciesData{}

	err = json.Unmarshal(responseFromURL, &r)
	if err != nil {
		log.Fatal(err)
	}
	if (CurrenciesData{}) == r {
		fmt.Println("Unable to unmarshal")
		return CurrenciesData{}, errors.New("Invalid Symbol")
	}
	return r, nil
}
func AllImpotantData() (ArrCurrenciesData, error) {
	url := "https://api.hitbtc.com/api/2/public/ticker"
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return ArrCurrenciesData{}, err
	}
	defer response.Body.Close()

	responseFromURL, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		return ArrCurrenciesData{}, err
	}
	r := ArrCurrenciesData{}

	err = json.Unmarshal(responseFromURL, &r)
	if err != nil {
		log.Fatal(err)
	}
	if len(r) == 0 {
		fmt.Println("Data Not found")
		return ArrCurrenciesData{}, errors.New("Data Not found")
	}
	return r, nil
}

func AllSymbol(w http.ResponseWriter, ) (ArrSymbolDetail, error) {
	url := "https://api.hitbtc.com/api/2/public/symbol"
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return ArrSymbolDetail{}, err
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		return ArrSymbolDetail{}, err
	}
	//fmt.Fprint(w, string(responseData))
	s := ArrSymbolDetail{}
	err = json.Unmarshal(responseData, &s)
	if err != nil {
		log.Fatal(err)
		return ArrSymbolDetail{}, err
	}
	if len(s) == 0 {
		return ArrSymbolDetail{}, errors.New("No sysmbol found")
	}
	return s, nil
}

func checkSymbol(w http.ResponseWriter, symbol string) (SymbolDetail, error) {
	url := "https://api.hitbtc.com/api/2/public/symbol/" + symbol
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return SymbolDetail{}, err
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		return SymbolDetail{}, err
	}
	s := SymbolDetail{}
	err = json.Unmarshal(responseData, &s)
	if err != nil {
		log.Fatal(err)
		return SymbolDetail{}, err
	}
	if (SymbolDetail{}) == s {
		return SymbolDetail{}, errors.New("Invalid Symbol")
	}
	return s, nil
}
