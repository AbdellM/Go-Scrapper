package Scrappe

import (
	"encoding/csv"
	"encoding/json"
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

func Scrappe(airportCode string, airportQueryType string) []flight {
	var pageNbr int = 0

	baseUrl := "https://www.flightstats.com/go/weblet?guid=4ee8679bc940b230:4189663e:145719871fa:-173&weblet=status&action=AirportFlightStatus&airportCode="+airportCode+"&airportQueryType="+airportQueryType 
	GetPageNum(&baseUrl, &pageNbr)
	baseUrl = baseUrl + "&airportQueryTimePeriod="
	// WriteFlightsCSV(GetAllFlight(&pageNbr, &baseUrl))
	// return WriteFlightsJSON(GetAllFlight(&pageNbr, &baseUrl))
	return GetAllFlight(&pageNbr, &baseUrl)
}

func WriteFlightsJSON(flights []flight) *os.File {

	json_data, err := json.Marshal(flights)
	throwError(err)

	json_file, err := os.Create("sample.json")
	throwError(err)

	defer json_file.Close()
	
	json_file.Write(json_data)

	return json_file
}

func WriteFlightsCSV(flights []flight) bool {

	file, err := os.Create("flights.csv")
	throwError(err)

	w := csv.NewWriter(file)

	header := []string{"FlightCode", "Company", "Destination", "Time", "Situation" }

	err = w.Write(header)
	throwError(err)

	for _, flight := range flights {
		err = w.Write([]string{flight.FlightCode, flight.Company, flight.Destination, flight.Time, flight.Situation })
		throwError(err)
	}
	defer w.Flush()

	return true
}

func GetAllFlight(pageNbr *int, baseUrl *string) (flights []flight) {
	if *pageNbr == 0 {
		*pageNbr = 1
	}
	for i := 1; i <= *pageNbr; i++ {
		flights = append(flights,GetFlightSinglePage(*baseUrl + strconv.Itoa(i))...)
	}
	return flights
}

func GetFlightSinglePage(baseUrl string) (flights []flight) {

	doc := GetDocument(baseUrl)
	doc.Find(".tableListingTable>tbody>tr").Each(func(i int, oneFight *goquery.Selection) {
		if i!=0 {
			flights = append(flights, ExtractAFlight(oneFight))

		}		
	})
	return flights
}

func ExtractAFlight(oneFight *goquery.Selection) flight {
	oneFlight := make([]string, 5)

	oneFight.Find("td").Each(func(i int, element *goquery.Selection) {
		oneFlight[i] = CleanString(element.Text())
	})
	
	return flight{oneFlight[0],oneFlight[1],oneFlight[2],oneFlight[3],oneFlight[4]}
}

func GetPageNum(baseUrl *string, pageNbr *int) {
	doc := GetDocument(*baseUrl)

	doc.Find("option").Each(func(i int, s *goquery.Selection) {
		*pageNbr++
	})
}

func GetDocument(baseUrl string) *goquery.Document {
	res, err := http.Get(baseUrl)
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