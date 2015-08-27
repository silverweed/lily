package main

import (
	"../mustache"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

func httpTemplate(template string) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		send(rw, req, template, template, nil)
	}
}

func httpIndex(rw http.ResponseWriter, req *http.Request) {

	meetups, err := db.GetMeetups()
	if err != nil {
		http.Error(rw, "Server error", 500)
		return
	}

	for i := 0; i < len(meetups); i++ {
		fillFriendlyTime(&meetups[i])
	}

	send(rw, req, "index", "Index", struct {
		Meetups []Meetup
	}{
		meetups,
	})
}

func httpMeetup(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := vars["id"]

	meetup, err := db.GetMeetup(id)
	if err != nil {
		http.Error(rw, "Meetup not found", 404)
		return
	}

	fillFriendlyTime(&meetup)

	var days []Day

	if !meetup.DateConfirmed {
		duration := meetup.To.Sub(meetup.From)
		diff := int(duration / time.Hour / 24)
		days = getDays(&meetup, meetup.From, diff)
	} else {
		days = getDays(&meetup, meetup.When, 1)
	}

	send(rw, req, "meetup", "Meetup", struct {
		Meetup Meetup
		Days   []Day
	}{
		meetup,
		days,
	})
}

func send(rw http.ResponseWriter, req *http.Request,
	name string, title string, context interface{}) {
	if len(title) > 0 {
		title = " ~ " + title
	}
	fmt.Fprintln(rw,
		mustache.RenderFileInLayout(
			rootDir+"/template/"+name+".html",
			rootDir+"/template/layout.html",
			struct {
				Title string
				Data  interface{}
			}{
				title,
				context,
			}))
}

func fillFriendlyTime(m *Meetup) {
	if m.DateConfirmed {
		m.FriendlyTime.When = m.When.Format(dateFormat)
	} else {
		m.FriendlyTime.From = m.From.Format(dateFormat)
		m.FriendlyTime.To = m.To.Format(dateFormat)
	}
}

func getDays(meetup *Meetup, start time.Time, n int) []Day {
	days := make([]Day, n)
	for i := 0; i < n; i++ {
		d, err := time.ParseDuration(strconv.Itoa(i*24) + "h")
		if err != nil {
			panic("parseDuration error: " + err.Error())
		}
		day := start.Add(d)
		days[i].Day = day.Day()
		days[i].Month = int(day.Month())
		days[i].Year, _ = strconv.Atoi(day.Format("06"))
		for _, p := range meetup.People {
			for _, d := range p.Days {
				if d.Equal(day) {
					days[i].People = append(days[i].People, p)
				}
			}
		}
	}
	return days
}
