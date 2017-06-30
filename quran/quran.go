package main

import (
	"github.com/ccdatatraits/quran/app"
	"github.com/goadesign/goa"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// QuranController implements the quran resource.
type QuranController struct {
	*goa.Controller
	*mgo.Session
}

// NewQuranController creates a quran controller.
func NewQuranController(session *mgo.Session, service *goa.Service) *QuranController {
	return &QuranController{Controller: service.NewController("QuranController"), Session: session}
}

// Show implements the "show" action of the "quran" controller.
func (c *QuranController) Show(ctx *app.ShowQuranContext) error {
	if ctx.SuraID == 0 || ctx.AyaID == 0 {
		// Emulate a missing record with ID 0
		return ctx.NotFound()
	}

	mgoC := c.Session.DB("quran").C("quran")

	var aya app.Quran
	err := mgoC.Find(bson.M{"sura": ctx.SuraID, "aya": ctx.AyaID}).One(&aya)
	if err != nil {
		return ctx.NotFound()
	}
	// Build the resource using the generated data structure
	/*aya := app.Quran{
		ID:   ctx.IndexID,
		Text: fmt.Sprintf("Aya #%d", ctx.IndexID),
		Href: app.QuranHref(ctx.IndexID),
	}*/
	aya.Href = app.QuranHref(ctx.SuraID, ctx.AyaID)

	// Let the generated code produce the HTTP response using the
	// media type described in the design (QuranMedia).
	return ctx.OK(&aya)
}
