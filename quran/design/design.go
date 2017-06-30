package design // The convention consists of naming the design
// package "design"
import (
	. "github.com/goadesign/goa/design" // Use . imports to enable the DSL
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = API("quran", func() { // API defines the microservice endpoint and
	Title("Quranic API")                                          // other global properties. There should be one
	Description("A simple microservice for 'sura's, 'aya's etc.") // and exactly one API definition appearing in
	Scheme("http")                                                // the design.
	Host("localhost:8888")
})

var _ = Resource("quran", func() { // Resources group related API endpoints
	// All resource actions require "api:read" scope
	Security(OAuth2, func() { Scope("api:read") })
	BasePath("/quran")       // together. They map to REST resources for REST
	DefaultMedia(QuranMedia) // services.

	/*Action("show", func() { // Actions define a single API endpoint together
		Description("Get aya by index") // with its path, parameters (both path
		Routing(GET("/:indexID"))       // parameters and querystring values) and payload
		Params(func() {                 // (shape of the request body).
			Param("indexID", Integer, "index ID")
		})
		Response(OK)       // Responses define the shape and status code
		Response(NotFound) // of HTTP responses.
	})*/

	Action("show", func() { // Actions define a single API endpoint together
		Description("Get aya by sura & aya") // with its path, parameters (both path
		Routing(GET("/:suraID/:ayaID"))      // parameters and querystring values) and payload
		Params(func() {                      // (shape of the request body).
			Param("suraID", Integer, "sura ID")
			Param("ayaID", Integer, "aya ID")
		})
		Response(OK)       // Responses define the shape and status code
		Response(NotFound) // of HTTP responses.
	})
})

// OAuth2 for security
var OAuth2 = OAuth2Security("OAuth2", func() {
	Description("Use OAuth2 to authenticate")
	AccessCodeFlow("/authorization", "/token")
	Scope("api:read", "Provides read access")
	//Scope("api:write", "Provides write access")
})

// QuranMedia defines the media type used to fetch Quran details.
var QuranMedia = MediaType("application/vnd.quran+json", func() {
	Description("Information from Quran")
	Attributes(func() { // Attributes define the media type shape.
		Attribute("index", Integer, "index to each aya")
		Attribute("sura", Integer, "sura #")
		Attribute("aya", Integer, "aya # within sura")
		Attribute("href", String, "API href for making requests to the aya")
		Attribute("text", String, "Arabic text")
		Required("index", "href", "text")
	})
	View("default", func() { // View defines a rendering of the media type.
		Attribute("index") // Media types may have multiple views and must
		Attribute("sura")  // have a "default" view.
		Attribute("aya")
		Attribute("href")
		Attribute("text")
	})
})
