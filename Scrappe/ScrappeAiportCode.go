package Scrappe

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Aiport struct {
    Name, AiportCode string
}

func ScrappeAiportCode(value string) []Aiport {
	return GetAllAiport(value)
}

func GetAllAiport(value string) (Aiports []Aiport) {

	baseUrl := "https://www.air-port-codes.com/search/results/" + value 

	doc := GetDocument(baseUrl)
	doc.Find(".small>tbody>tr").Each(func(i int, oneFight *goquery.Selection) {
		if i != 0 && i < 9 {
			oneAiport := make([]string, 2)
			skip:= false
			oneFight.Find("td").Each(func(i int, element *goquery.Selection) {

				if i == 0 {
					if contains(element.Text(), "All Airports"){
						skip = true
					} else {
						oneAiport[0] = CleanString(element.Text())
					}
				}
				if i == 1 {
					if skip != true {
						oneAiport[1] = CleanString(element.Text())
						Aiports = append(Aiports, Aiport{oneAiport[0],oneAiport[1]})
					} else {
						skip = false
					}
				}

			})
		}	
	})
	return Aiports
}


func contains(str, substr string) bool{
	return strings.Contains(str, substr)
}



