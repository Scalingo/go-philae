package mongoprobe

import (
	errgo "gopkg.in/errgo.v1"
	mgo "gopkg.in/mgo.v2"
)

type MongoProbe struct {
	name string
	url  string
}

func NewMongoProbe(name, url string) MongoProbe {
	return MongoProbe{
		name: name,
		url:  url,
	}
}

func (p MongoProbe) Name() string {
	return p.name
}

func (p MongoProbe) Check() error {
	session, err := mgo.Dial(p.url)
	if err != nil {
		return errgo.Notef(err, "Unable to contact server")
	}
	defer session.Close()

	_, err = session.DatabaseNames()
	if err != nil {
		return errgo.Notef(err, "Unable to send query")
	}
	return nil
}
