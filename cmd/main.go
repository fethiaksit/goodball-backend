package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gocolly/colly"
)

type Item struct {
	Team1  string `json:"team1"`
	Team2  string `json:"team2"`
	Score1 string `json:"score1"`
	Score2 string `json:"score2"`
}
type Response struct {
	Success bool   `json:"success"`
	Result  []Item `json:"result"`
}

func getMatches() []Item {
	var matches []Item

	c := colly.NewCollector(
		colly.AllowedDomains("www.tff.org"),
		colly.CacheDir(""),
	)

	c.OnHTML("tr.haftaninMaclariTr", func(h *colly.HTMLElement) {
		item := Item{
			Team1:  h.ChildText("td.haftaninMaclariEv span"),
			Score1: h.ChildText("td.haftaninMaclariSkor span:nth-of-type(1)"),
			Score2: h.ChildText("td.haftaninMaclariSkor span:nth-of-type(2)"),
			Team2:  h.ChildText("td.haftaninMaclariDeplasman span"),
		}
		matches = append(matches, item)

	})
	err := c.Visit("https://www.tff.org/Default.aspx?pageID=142&hafta=23")
	if err != nil {
		fmt.Println("error visiting page", err)
	}
	return matches
}

func main() {
	http.HandleFunc("/matches", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		data := getMatches()
		resp := Response{
			Success: true,
			Result:  data,
		}
		json.NewEncoder(w).Encode(resp)

	})

	fmt.Println("Sunucu çalışıyor http://0.0.0.0:8080/matches")

	http.ListenAndServe("0.0.0.0:8080", nil)

}
