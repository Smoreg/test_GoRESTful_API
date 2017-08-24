package test_GoRESTful_API

import (
	"gopkg.in/mgo.v2"
	"log"
	"gopkg.in/mgo.v2/bson"
	. "github.com/Smoreg/test_GoRESTful_API/model"
)

type memesDAO struct {
	Server   string `json:"server"`
	Database string `json:"database"`
}

var db *mgo.Database

func (m memesDAO) connect()  {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Panic(err)
	}
	db = session.DB(m.Database)
	return
}

func getClt() *mgo.Collection {
	return db.C("memes")
}

func (m *memesDAO) FindAll() ([]Meme, error) {
	var memes []Meme
	err := getClt().Find(bson.M{}).All(&memes)
	return memes, err
}

func (m *memesDAO) FindById(id string) (Meme, error) {
	var meme Meme
	err := getClt().FindId(bson.ObjectIdHex(id)).One(&meme)
	return meme, err
}

func (m *memesDAO) Insert(meme Meme) error {
	err := getClt().Insert(&meme)
	return err
}

func (m *memesDAO) Delete(meme Meme) error {
	err := getClt().Remove(&meme)
	return err
}

func (m *memesDAO) Update(meme Meme) error {
	err := getClt().UpdateId(meme.ID, &meme)
	return err
}

