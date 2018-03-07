package models

import (
	"io/ioutil"
	"path"
	"strings"

	"github.com/salsita/go-pivotaltracker/v5/pivotal"
)

// OutRequest -
type OutRequest struct {
	Source Source    `json:"source"`
	Params OutParams `json:"params"`
}

// Source -
type Source struct {
	BaseURL   string `json:"url"`
	Token     string `json:"token"`
	ProjectID int    `json:"project_id"`
}

// OutParams -
type OutParams struct {
	NameFile string `json:"name_file"`
	pivotal.StoryRequest
}

// ReplaceName -
func (o *OutParams) ReplaceName(filepath string) error {
	if o.NameFile != "" {

		dir := path.Join(filepath, o.NameFile)
		dat, err := ioutil.ReadFile(dir)

		if err != nil {
			return err
		}

		o.Name = strings.Replace(o.Name, "$NAME_FILE", string(dat), -1)
	}

	return nil
}
