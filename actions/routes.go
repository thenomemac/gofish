package actions

import (
	"fmt"
	"log"

	"github.com/gobuffalo/buffalo"
	"github.com/pkg/errors"
	"github.com/thenomemac/gofish/models"
)

// RoutesHandler
func RoutesHandler(c buffalo.Context) error {
	return c.Render(200, r.HTML("routes.html"))
}

// UploadHandler
func UploadHandler(c buffalo.Context) error {
	return c.Render(200, r.HTML("upload.html"))
}

// RegulationsHandler
func RegulationsHandler(c buffalo.Context) error {
	return c.Render(200, r.HTML("regulations.html"))
}

// LabelHandler
func LabelHandler(c buffalo.Context) error {
	return c.Render(200, r.HTML("label.html"))
}

// DataexportHandler
func DataexportHandler(c buffalo.Context) error {
	return c.Render(200, r.HTML("dataexport.html"))
}

// AboutHandler
func AboutHandler(c buffalo.Context) error {
	return c.Render(200, r.HTML("about.html"))
}

// LegalHandler
func LegalHandler(c buffalo.Context) error {
	return c.Render(200, r.HTML("legal.html"))
}

// FishpicResultsHandler
func FishpicResultsHandler(c buffalo.Context) error {
	id := c.Param("fishpic_id")

	log.Println("fishpic_id:", id)

	fishpic := &models.Fishpic{}
	err := fishpic.Find(id)
	if err != nil {
		errors.WithStack(err)
	}

	fishpicURL := fmt.Sprintf("/imgs/%s.jpg", fishpic.ID)
	c.Set("fishpic_url", fishpicURL)
	return c.Render(200, r.HTML("fishpic-results.html"))
}
