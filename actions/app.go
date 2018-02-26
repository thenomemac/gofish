package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/middleware"
	"github.com/gobuffalo/buffalo/middleware/csrf"
	"github.com/gobuffalo/buffalo/middleware/ssl"
	"github.com/gobuffalo/envy"
	"github.com/unrolled/secure"

	"github.com/gobuffalo/buffalo/middleware/i18n"
	"github.com/gobuffalo/packr"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App
var T *i18n.Translator

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:         ENV,
			SessionName: "_gofish_session",
			LooseSlash:  true,
		})
		// Automatically redirect to SSL
		app.Use(ssl.ForceSSL(secure.Options{
			SSLRedirect:     ENV == "production",
			SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
		}))

		if ENV == "development" {
			app.Use(middleware.ParameterLogger)
		}

		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		app.Use(csrf.New)

		// Wraps each request in a transaction.
		//  c.Value("tx").(*pop.PopTransaction)
		// Remove to disable this.
		// app.Use(middleware.PopTransaction(models.DB))

		// Setup and use translations:
		var err error
		if T, err = i18n.New(packr.NewBox("../locales"), "en-US"); err != nil {
			app.Stop(err)
		}
		app.Use(T.Middleware())

		app.GET("/", HomeHandler)
		app.GET("/routes", RoutesHandler)
		app.GET("/upload", UploadHandler)
		app.GET("/regulations", RegulationsHandler)
		app.GET("/label", LabelHandler)
		app.GET("/dataexport", DataexportHandler)
		app.GET("/about", AboutHandler)
		app.GET("/legal", LegalHandler)

		app.ServeFiles("/assets", assetsBox)

		// app.ServeFiles("/imgs", packr.NewBox("../imgs"))
		// app.ServeFiles("/imgs", http.Dir("./imgs"))
		app.GET("/imgs/{fishpic_id}", GetImgHandler)

		app.Resource("/fishpics", FishpicsResource{})

		app.POST("/postimg", PostImgHandler)

		app.GET("/fishpic-results/{fishpic_id}", FishpicResultsHandler)
	}

	return app
}
