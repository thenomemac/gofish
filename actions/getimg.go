package actions

import (
	"fmt"
	"io/ioutil"

	"github.com/gobuffalo/buffalo"
	"github.com/pkg/errors"
	"github.com/thenomemac/gofish/models"
)

func GetImgHandler(c buffalo.Context) error {
	var data []byte
	var err error

	id := c.Param("fishpic_id")

	if models.LOCAL_CACHE == "true" {
		data, err = ioutil.ReadFile(fmt.Sprintf("imgs/%s", id))
	} else {
		data, err = models.Downloader(fmt.Sprintf("fishpics/%s", id))
	}

	if err != nil {
		errors.WithStack(err)
	}

	w := c.Response()
	_, err = w.Write(data)
	if err != nil {
		errors.WithStack(err)
	}

	return nil // c.Render(200, r.String("done"))
}
