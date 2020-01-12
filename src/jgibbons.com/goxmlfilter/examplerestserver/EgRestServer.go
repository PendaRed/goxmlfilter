package main

import (
	//	"encoding/json"
	"jgibbons.com/goxmlfilter/examplerestserver/data"
	"log"
	"net/http"
)

// Cribbed from https://dev.to/moficodes/build-your-first-rest-api-with-go-2gcj
type server struct{}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	recs := new(data.EgRecords)
	recs.Add(&data.EgDataRecord{Forename: "Jonathan", Surname: "Gibbons", Gender: "Male", House: "12 A Street"})
	recs.Add(&data.EgDataRecord{Forename: "Caladan", Surname: "Brood", Gender: "Alien", House: "20 A Street"})
	recs.Add(&data.EgDataRecord{Forename: "Bob", Surname: "Builder", Gender: "Male", House: "21 A Street"})
	recs.Add(&data.EgDataRecord{Forename: "Freddy", Surname: "Mouse", Gender: "Male", House: "22 A Street"})
	recs.Add(&data.EgDataRecord{Forename: "Father", Surname: "Christmas", Gender: "Male", House: "North Pole"})

	//	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK) // do this last
	//	json.NewEncoder(w).Encode(recs)
	//	w.Write([]byte(recs.AsJson()))
	w.Write([]byte(recs.AsXml()))
}

func main() {
	s := &server{}
	http.Handle("/", s)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
