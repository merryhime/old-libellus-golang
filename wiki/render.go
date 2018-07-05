package wiki

import (
	"html/template"
	"strings"

	"github.com/MerryMage/libellus/wikidata"
)

type RenderedPath []string

func (rp RenderedPath) IsLast(i int) bool {
	return i == len(rp)-1
}

func (rp RenderedPath) Partial(i int) string {
	return "/" + strings.Join(rp[:i+1], "/")
}

func (rp RenderedPath) NotRoot() bool {
	return !(len(rp) == 1 && rp[0] == "")
}

type RenderedSubpage struct {
	Path  string
	Title string
}

type RenderedKnowledge struct {
	Identifier   string
	RenderedHTML template.HTML
	CardCount    int
}

type RenderedPage struct {
	Authorized bool

	Title      string
	Path       RenderedPath
	Subpages   []RenderedSubpage
	Knowledges []RenderedKnowledge
}

func (wiki *Wiki) RenderKnowledge(kid wikidata.KnowledgeId) RenderedKnowledge {
	km, k := wiki.config.WikiData.LookupKnowledge(kid)

	var rendered RenderedKnowledge
	rendered.CardCount = len(km.Cards)

	switch k := k.(type) {
	case wikidata.ErrorKnowledge:
		rendered.Identifier = string(k.Identifier)
		rendered.RenderedHTML = template.HTML(k.Message)
	case wikidata.MarkdownKnowledge:
		rendered.Identifier = string(k.Identifier)
		rendered.RenderedHTML = template.HTML(k.Markdown)
	default:
		rendered.Identifier = "invalid"
		rendered.RenderedHTML = template.HTML("no idea what's going on")
	}

	return rendered
}
