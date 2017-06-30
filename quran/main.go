//go:generate goagen bootstrap -d github.com/ccdatatraits/quran/design

package main

import (
	"github.com/ccdatatraits/quran/app"
	"github.com/goadesign/goa"
	"github.com/goadesign/goa/middleware"
	"gopkg.in/mgo.v2"
)

func ensureIndex(s *mgo.Session) {
	session := s.Copy()
	defer session.Close()

	c := session.DB("quran").C("quran")

	index := mgo.Index{
		Key:        []string{"index"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

func main() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	ensureIndex(session)

	// Create service
	service := goa.New("quran")

	// Mount middleware
	service.Use(middleware.RequestID())
	service.Use(middleware.LogRequest(true))
	service.Use(middleware.ErrorHandler(service, true))
	service.Use(middleware.Recover())

	// Mount "quran" controller
	c := NewQuranController(session, service)
	app.MountQuranController(service, c)

	// Start service
	if err := service.ListenAndServe(":8888"); err != nil {
		service.LogError("startup", "err", err)
	}

}
