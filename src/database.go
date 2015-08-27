package main

import (
	"encoding/base64"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Database struct {
	session  *mgo.Session
	database *mgo.Database
}

func InitDatabase(servers, dbname string) Database {
	var db Database
	var err error
	db.session, err = mgo.Dial(servers)
	if err != nil {
		panic(err)
	}
	db.database = db.session.DB(dbname)
	return db
}

func (db Database) Close() {
	db.session.Close()
}

func (db Database) GetMeetups() (meetups []Meetup, err error) {
	err = db.database.C("meetups").Find(nil).All(&meetups)
	if err != nil {
		return
	}
	return
}

func (db Database) NewMeetup(meetup Meetup) (err error, id bson.ObjectId) {
	meetup.Id = bson.NewObjectId()
	meetup.Url = base64.URLEncoding.EncodeToString([]byte(meetup.Id[:]))[:10]
	err = db.database.C("meetups").Insert(meetup)
	return
}

func (db Database) GetMeetup(url string) (meetup Meetup, err error) {
	err = db.database.C("meetups").Find(bson.M{"url": url}).One(&meetup)
	if err != nil {
		return
	}
	return
}
