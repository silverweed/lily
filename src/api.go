package main

import (
	"log"
	"net/http"
	"time"
)

const dateFormat = "2/1/2006"

// params: city=string, from=dd/mm/yy, to=dd/mm/yy
func apiNew(rw http.ResponseWriter, req *http.Request) {
	city := req.PostFormValue("city")
	from := req.PostFormValue("from")
	to := req.PostFormValue("to")

	if len(city) < 1 || len(from) < 8 || len(to) < 8 {
		http.Error(rw, "Invalid parameters", http.StatusBadRequest)
		return
	}

	datestart, err := time.Parse(dateFormat, from)
	if err != nil || datestart.Before(time.Now()) {
		http.Error(rw, "Invalid starting date", http.StatusBadRequest)
		return
	}
	dateend, err := time.Parse(dateFormat, to)
	if err != nil || dateend.Before(datestart) {
		http.Error(rw, "Invalid ending date", http.StatusBadRequest)
		return
	}

	meetup := Meetup{
		City:          city,
		DateConfirmed: false,
		From:          datestart,
		To:            dateend,
	}

	if err, _ := db.NewMeetup(meetup); err != nil {
		http.Error(rw, "Server error", 500)
		log.Printf("Error inserting meetup: %s\n", err.Error())
		return
	}

	http.Redirect(rw, req, "/", http.StatusMovedPermanently)
}
