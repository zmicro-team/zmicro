package docs

import (
	_ "embed"
)

//go:embed swagger.swagger.json
var docs string

type Docs struct{}

func (*Docs) ReadDoc() string {
	return docs
}
