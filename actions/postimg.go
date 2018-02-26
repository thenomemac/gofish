package actions

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/gobuffalo/buffalo"
	"github.com/pkg/errors"
	"github.com/thenomemac/gofish/models"
)

func PostImgHandler(c buffalo.Context) error {
	var body []byte

	body, errBody := ioutil.ReadAll(c.Request().Body)

	// if there is a form file the body file will be overwritten
	file, errForm := c.File("fish-pic-input")
	if errForm == nil {
		body, errForm = ioutil.ReadAll(file)
	}

	if errBody != nil && errForm != nil {
		log.Println("upload errors::", errForm, errBody)
		return c.Render(500, r.String("Upload an image fool!"))
	}

	fishpic := models.Fishpic{}
	err := fishpic.CreateAndSave()
	if err != nil {
		errors.WithStack(err)
	}

	if models.LOCAL_CACHE == "true" {
		err = ioutil.WriteFile(fmt.Sprintf("fishpics/%s.jpg", fishpic.ID), body, 0600)
	} else {
		err = models.Uploader(fmt.Sprintf("fishpics/%s.jpg", fishpic.ID), body)
	}

	if err != nil {
		errors.WithStack(err)
	}

	return c.Redirect(301, "/fishpic-results/%s", fishpic.ID)
}
