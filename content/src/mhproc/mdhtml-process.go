package mhproc

import (
	"corsa-blog/content/src/mhparser"
	"fmt"
	"log"
)

type MdHtmlProcess struct {
	debug    bool
	scrGramm mhparser.ScriptGrammar
	HtmlGen  string
	ident    int
}

func (mp *MdHtmlProcess) ProcessToHtml(script string) error {
	log.Println("[ProcessToHtml] is called with a script ", len(script))
	if script == "" {
		return fmt.Errorf("[ProcessToHtml] script is empty")
	}
	mp.scrGramm = mhparser.ScriptGrammar{
		Debug: mp.debug,
	}
	if err := mp.scrGramm.ParseScript(script); err != nil {
		log.Println("[ProcessToHtml] Parser error")
		return err
	}
	if err := mp.scrGramm.CheckNorm(); err != nil {
		log.Println("[ProcessToHtml] Script structure error")
		return err
	}
	if err := mp.scrGramm.EvaluateParams(); err != nil {
		log.Println("[ProcessToHtml] EvaluateParams error")
		return err
	}
	if mp.scrGramm.Title == "" {
		return fmt.Errorf("[ProcessToHtml] field 'title' in mdhtml is empty")
	}
	if mp.scrGramm.PostId == "" {
		return fmt.Errorf("[ProcessToHtml] field 'id' in mdhtml is empty")
	}
	if mp.scrGramm.Datetime.Year() < 2010 {
		return fmt.Errorf("[ProcessToHtml] field 'datetime' is empty or invalid")
	}
	if mp.debug {
		log.Println("[ProcessToHtml]  Parser found Nodes: ", mp.scrGramm.Norm)
	}
	return mp.parsedToHtml()
}

func (mp *MdHtmlProcess) parsedToHtml() error {
	normPrg := mp.scrGramm.Norm["main"]
	firstStep := []string{}
	for _, stItem := range normPrg.FnsList {
		if stItem.Type == mhparser.TtHtmlVerbatim {
			firstStep = append(firstStep, stItem.Params[0].ArrayValue...)
		}
	}
	return nil
}
