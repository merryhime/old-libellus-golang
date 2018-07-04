package wikidata

import (
	"encoding/json"
)

const (
	ErrorKnowledgeType    string = "error"
	MarkdownKnowledgeType string = "markdown"
)

type Knowledge interface {
	knowledgeTag()
	GetCards() []string
	GetInfo() KnowledgeInfo
}

type KnowledgeInfo struct {
	Identifier KnowledgeId
	Type       string
	Cards      []string
}

func (ki KnowledgeInfo) GetCards() []string {
	return ki.Cards
}

type ErrorKnowledge struct {
	KnowledgeInfo
	Message string
}

func (ErrorKnowledge) knowledgeTag() {}
func (k ErrorKnowledge) GetInfo() KnowledgeInfo {
	return k.KnowledgeInfo
}

func NewErrorKnowledge(msg string) ErrorKnowledge {
	return ErrorKnowledge{
		KnowledgeInfo: KnowledgeInfo{Type: ErrorKnowledgeType},
		Message:       msg,
	}
}

type MarkdownKnowledge struct {
	KnowledgeInfo
	Markdown string
}

func (MarkdownKnowledge) knowledgeTag() {}
func (k MarkdownKnowledge) GetInfo() KnowledgeInfo {
	return k.KnowledgeInfo
}

func (wd *WikiData) parseKnowledge(meta KnowledgeMeta) Knowledge {
	path := meta.ParentPath + "/_page/" + string(meta.Identifier)

	ki, err := wd.parseKnowledgeInfo(meta)

	if err != nil {
		return ErrorKnowledge{KnowledgeInfo: ki, Message: path + ": while parsing info: " + err.Error()}
	}

	switch ki.Type {
	case MarkdownKnowledgeType:
		return wd.parseMarkdownKnowledge(path, ki, meta)

	case ErrorKnowledgeType:
		return ErrorKnowledge{KnowledgeInfo: ki, Message: "wild ErrorKnowledgeType found at " + path}
	}

	return ErrorKnowledge{KnowledgeInfo: ki, Message: "unknown knowledge type found at " + path}
}

func (wd *WikiData) parseKnowledgeInfo(meta KnowledgeMeta) (KnowledgeInfo, error) {
	infoRaw, err := wd.repo.ReadBlobFromTreeOid(meta.TreeOid, "_info")
	if err != nil {
		return KnowledgeInfo{}, err
	}

	var kid KnowledgeInfo
	err = json.Unmarshal(infoRaw, &kid)
	if err != nil {
		return KnowledgeInfo{}, err
	}
	kid.Identifier = meta.Identifier

	return kid, nil
}

func (wd *WikiData) parseMarkdownKnowledge(path string, ki KnowledgeInfo, meta KnowledgeMeta) Knowledge {
	raw, err := wd.repo.ReadBlobFromTreeOid(meta.TreeOid, "_data.md")
	if err != nil {
		return ErrorKnowledge{KnowledgeInfo: ki, Message: "could not read _data.md at " + path}
	}
	return MarkdownKnowledge{
		KnowledgeInfo: ki,
		Markdown:      string(raw),
	}
}
