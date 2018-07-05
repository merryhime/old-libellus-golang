package wikidata

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/MerryMage/libellus/objstore"
	"github.com/MerryMage/libellus/objstore/filemode"
	"github.com/MerryMage/libellus/objstore/objid"
	"github.com/MerryMage/libellus/objstore/tree"
)

type KnowledgeId string
type CardId string

type KnowledgeMeta struct {
	ParentPath string
	Identifier KnowledgeId
	TreeOid    objid.Oid
	Cards      []CardId
}

type CardMeta struct {
	ParentParentPath string
	ParentIdentifier KnowledgeId
	BlobOid          objid.Oid
}

type PageInfo struct {
	Title      string
	Knowledges []KnowledgeId
}

type Page struct {
	PageInfo
	ActualKnowledges []KnowledgeId
	Path             string
	Children         []string
	NoInfo           bool
}

type WikiData struct {
	repo *objstore.Repository
	ref  string

	knowledges map[KnowledgeId]KnowledgeMeta
	cards      map[CardId]CardMeta
	pages      map[string]Page
}

type RefreshStateErrorInfo struct {
	Path string
	Err  error
}

func New(repo *objstore.Repository, ref string) *WikiData {
	wd := &WikiData{
		repo:       repo,
		ref:        ref,
		knowledges: make(map[KnowledgeId]KnowledgeMeta),
		pages:      make(map[string]Page),
	}

	wd.RefreshState()

	return wd
}

func (wd *WikiData) addError(path string, err error) {
	log.Println(path, "-", err)
}

func (wd *WikiData) parsePageInfo(currentPage *Page, pageTreeEntry *tree.Entry) {
	pageTree, err := wd.repo.Tree(pageTreeEntry.Oid)
	if err != nil {
		wd.addError(currentPage.Path+"/_page", err)
		return
	}

	for _, e := range pageTree.Entries {
		if e.Name[0] == '_' {
			continue
		}

		if e.Mode == filemode.Dir {
			kid := KnowledgeId(e.Name)
			currentPage.ActualKnowledges = append(currentPage.ActualKnowledges, kid)

			km := KnowledgeMeta{
				ParentPath: currentPage.Path,
				Identifier: kid,
				TreeOid:    e.Oid,
			}

			if cardsTreeEntry, err := tree.Lookup(wd.repo, e.Oid, "_cards"); err == nil {
				wd.parseCardInfo(currentPage, &km, cardsTreeEntry)
			}

			wd.knowledges[kid] = km

			continue
		}

		wd.addError(currentPage.Path+"/_page/"+e.Name, errors.New("wikidata/parsePageInfo: unexpected loose file"))
	}

	infoRaw, err := wd.repo.ReadBlobFromTree(pageTree, "_info")
	if err != nil {
		wd.addError(currentPage.Path+"/_page/_info", err)
		return
	}

	err = json.Unmarshal(infoRaw, &currentPage.PageInfo)
	if err != nil {
		wd.addError(currentPage.Path+"/_page/_info", err)
		return
	}

	currentPage.NoInfo = false
}

func (wd *WikiData) parseCardInfo(currentPage *Page, km *KnowledgeMeta, cardsTreeEntry *tree.Entry) {
	cardsTree, err := wd.repo.Tree(cardsTreeEntry.Oid)
	if err != nil {
		wd.addError(currentPage.Path+"/_page/"+string(km.Identifier)+"/_cards", err)
		return
	}

	for _, e := range cardsTree.Entries {
		if e.Name[0] == '_' {
			continue
		}

		if e.Mode == filemode.Regular {
			cid := CardId(e.Name)
			km.Cards = append(km.Cards, cid)

			wd.cards[cid] = CardMeta{
				ParentParentPath: currentPage.Path,
				ParentIdentifier: km.Identifier,
				BlobOid:          e.Oid,
			}

			continue
		}

		wd.addError(currentPage.Path+"/_page/"+string(km.Identifier)+"/_cards/"+e.Name, errors.New("wikidata/parseCardInfo: unexpected non-regular entry"))
	}
}

func (wd *WikiData) refreshStateHelper(currentPath string, currentTree objid.Oid) {
	tree, err := wd.repo.Tree(currentTree)
	if err != nil {
		wd.addError(currentPath, err)
		return
	}

	currentPage := &Page{
		Path:   currentPath,
		NoInfo: true,
	}

	if currentPath == "" {
		currentPage.Path = "/"
	}

	if pageTreeEntry := tree.Find("_page"); pageTreeEntry != nil {
		wd.parsePageInfo(currentPage, pageTreeEntry)
	}

	for _, e := range tree.Entries {
		path := currentPath + "/" + e.Name

		if e.Name[0] == '_' {
			continue
		}

		currentPage.Children = append(currentPage.Children, path)

		if e.Mode == filemode.Dir {
			wd.refreshStateHelper(path, e.Oid)
			continue
		}

		wd.addError(path, errors.New("wikidata: unexpected loose file"))
	}

	wd.pages[currentPage.Path] = *currentPage
}

func (wd *WikiData) RefreshState() {
	rootTreeEntry, err := wd.repo.LookupEntryByPath(wd.ref, "_wiki")
	if err != nil {
		wd.addError("/", err)
		return
	}
	wd.refreshStateHelper("", rootTreeEntry.Oid)
}

func (wd *WikiData) LookupPage(path string) (Page, bool) {
	page, ok := wd.pages[path]
	return page, ok
}

func (wd *WikiData) LookupKnowledgeMeta(kid KnowledgeId) (KnowledgeMeta, bool) {
	k, ok := wd.knowledges[kid]
	return k, ok
}

func (wd *WikiData) LookupKnowledge(kid KnowledgeId) (KnowledgeMeta, Knowledge) {
	k, ok := wd.knowledges[kid]
	if !ok {
		return KnowledgeMeta{
			ParentPath: "/_error",
			Identifier: kid,
		}, NewErrorKnowledge("kid \"" + string(kid) + "\" not found")
	}

	return k, wd.parseKnowledge(k)
}

func (wd *WikiData) LookupCardMeta(cid CardId) (CardMeta, bool) {
	c, ok := wd.cards[cid]
	return c, ok
}
