package mhproc

import (
	"bytes"
	"corsa-blog/content/src/mhparser"
	"fmt"
	"log"
	"text/template"
	"time"
)

type MdHtmlProcess struct {
	debug         bool
	scrGramm      mhparser.ScriptGrammar
	HtmlGen       string
	pageTemplName string
}

func NewMdHtmlProcess(debug bool) *MdHtmlProcess {
	res := MdHtmlProcess{
		debug:         debug,
		pageTemplName: "templates/htmlgen/post.html",
	}
	return &res
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
		main_norm := mp.scrGramm.Norm["main"]
		log.Println("[ProcessToHtml] Parser nodes found: ", len(main_norm.FnsList))
	}
	return mp.parsedToHtml()
}

func (mp *MdHtmlProcess) parsedToHtml() error {
	if mp.debug {
		log.Println("create the HTML using parsed info")
	}
	normPrg := mp.scrGramm.Norm["main"]
	lines := []string{}
	for _, stItem := range normPrg.FnsList {
		if stItem.Type == mhparser.TtHtmlVerbatim {
			lines = append(lines, stItem.Params[0].ArrayValue...)
		}
	}

	templName := mp.pageTemplName
	var partFirst, partSecond, partMerged bytes.Buffer
	tmplPage := template.Must(template.New("Page").ParseFiles(templName))
	CtxFirst := struct {
		Title string
		Lines []string
	}{
		Title: mp.scrGramm.Title,
		Lines: lines,
	}

	if err := tmplPage.ExecuteTemplate(&partFirst, "postbeg", CtxFirst); err != nil {
		return err
	}

	CtxSecond := struct {
		DateFormatted string
		DateTime      string
		PostId        string
	}{
		DateTime:      mp.scrGramm.Datetime.Format("2006-01-02 15:00"),
		DateFormatted: formatPostDate(mp.scrGramm.Datetime),
		PostId:        mp.scrGramm.PostId,
	}
	if err := tmplPage.ExecuteTemplate(&partSecond, "postfinal", CtxSecond); err != nil {
		return err
	}
	partFirst.WriteTo(&partMerged)
	partSecond.WriteTo(&partMerged)
	mp.HtmlGen = partMerged.String()
	if mp.debug {
		fmt.Println("***HTML***\n", mp.HtmlGen)
	}
	return nil
}

func formatPostDate(tt time.Time) string {
	res := fmt.Sprintf("%d %s %d", tt.Day(), monthToStringIT(tt.Month()), tt.Year())
	return res
}

func monthToStringIT(month time.Month) string {
	switch month {
	case time.January:
		return "Gennaio"
	case time.February:
		return "Febbraio"
	case time.March:
		return "Marzo"
	case time.April:
		return "Aprile"
	case time.May:
		return "Maggio"
	case time.June:
		return "Giugno"
	case time.July:
		return "Luglio"
	case time.August:
		return "Agosto"
	case time.September:
		return "Settembre"
	case time.October:
		return "Ottobre"
	case time.November:
		return "Novembre"
	case time.December:
		return "Dicembre"
	default:
		return ""
	}
}
