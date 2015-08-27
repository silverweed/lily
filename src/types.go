package main

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Meetup struct {
	Id            bson.ObjectId `_id`
	Url           string
	City          string
	When          time.Time
	From          time.Time
	To            time.Time
	DateConfirmed bool
	FriendlyTime  StringTime
	People        []Person
}

type StringTime struct {
	When string
	From string
	To   string
}

type Person struct {
	Name string
	Days []time.Time
}

type Day struct {
	Day    int
	Month  int
	Year   int
	People []Person
}
