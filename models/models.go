package models

import (
	"io/ioutil"
	"path"
	"regexp"
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

// MergeName -
func (o *OutParams) MergeName(filepath string) error {
	if o.NameFile == "" {
		return nil
	}

	dir := path.Join(filepath, o.NameFile)
	data, err := ioutil.ReadFile(dir)
	if err != nil {
		return err
	}

	replaced := strings.Replace(o.Name, "$NAME_FILE", string(data), -1)
	trimmed := strings.TrimSpace(replaced)
	re := regexp.MustCompile(`\r?\n`)
	formatted := re.ReplaceAllString(trimmed, "")

	o.Name = formatted

	return nil
}
