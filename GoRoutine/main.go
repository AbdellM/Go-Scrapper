package main

import (
	"encoding/csv"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type flight struct {
    FlightCode, Company, Destination, Time, Situation string
}

// func main(airportCode string, airportQueryType string) {
func main() {
	airportCode :="CMN"
	airportQueryType :="1"

	baseUrl := "https://www.flightstats.com/go/weblet?guid=4ee8679bc940b230:4189663e:145719871fa:-173&weblet=status&action=AirportFlightStatus&airportCode="+airportCode+"&airportQueryType="+airportQueryType 
	pageNbr := GetPageNum(&baseUrl)

	flights := GetAllFlight(&pageNbr, baseUrl + "&airportQueryTimePeriod=")
	WriteFlights(&flights)
}

func WriteFlights(flights *[]flight) bool {

	file, err := os.Create("flights.csv")
	throwError(err)

	w := csv.NewWriter(file)

	header := []string{"FlightCode", "Company", "Destination", "Time", "Situation" }

	err = w.Write(header)
	throwError(err)

	for _, flight := range *flights {
		err = w.Write([]string{flight.FlightCode, flight.Company, flight.Destination, flight.Time, flight.Situation })
		throwError(err)
	}

	defer w.Flush()

	return true
}

func GetAllFlight(pageNbr *int, baseUrl string) (flights []flight) {
	singlePageChan := make(chan []flight)

	for i := 1; i <= *pageNbr; i++ {
		go GetFlightSinglePage(baseUrl + strconv.Itoa(i), singlePageChan)
	}
	for i := 1; i <= *pageNbr; i++ {
		flights = append(flights, <- singlePageChan...)
	}
	return flights
}

func GetFlightSinglePage(baseUrl string, singlePageChan chan []flight) {
	var flights []flight
	doc := GetDocument(&baseUrl)

	doc.Find(".tableListingTable>tbody>tr").Each(func(i int, oneFight *goquery.Selection) {
		if i!=0 {
			flights = append(flights, ExtractAFlight(oneFight))
		}		
	})
	singlePageChan <- flights
}

func ExtractAFlight(oneFight *goquery.Selection) flight {
	oneFlight := make([]string, 5)

	oneFight.Find("td").Each(func(i int, element *goquery.Selection) {
		oneFlight[i] = CleanString(element.Text())
	})
	
	return flight{oneFlight[0],oneFlight[1],oneFlight[2],oneFlight[3],oneFlight[4]}
}

func GetPageNum(baseUrl *string) (pageNbr int)  {

	doc := GetDocument(baseUrl)

	doc.Find("option").Each(func(i int, s *goquery.Selection) {
		pageNbr++
	})
	return pageNbr
}

func GetDocument(baseUrl *string) *goquery.Document {
	res, err := http.Get(*baseUrl)
	throwError(err)

	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	throwError(err)

	return  doc
}

func CleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

func throwError(err error)  {
	if err != nil {
		log.Fatal(err)
	}
}