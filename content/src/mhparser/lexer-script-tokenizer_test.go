package mhparser

import (
	"strings"
	"testing"
)

func TestParseData(t *testing.T) {
	str := `title: Prossima gara Wien Rundumadum
datetime: 2024-11-08 19:00
id: 20241108-00`

	lex := ScriptGrammar{
		Debug: true,
	}
	err := lex.ParseScript(str)
	if err != nil {
		t.Error("Error is: ", err)
		return
	}

	err = lex.CheckNorm()
	if err != nil {
		t.Error("Error in parser norm ", err)
		return
	}
	err = lex.EvaluateParams()
	if err != nil {
		t.Error("Error in evaluate ", err)
		return
	}
	if lex.PostId != "20241108-00" {
		t.Error("unexpected id", lex.PostId)
	}
	if lex.Title != "Prossima gara Wien Rundumadum" {
		t.Error("unexpected Title", lex.Title)
	}
	if lex.Datetime.Year() != 2024 {
		t.Error("unexpected Year", lex.Datetime)
	}
	if lex.Datetime.Hour() != 19 {
		t.Error("unexpected Hour", lex.Datetime)
	}
}

func TestParseCustomData(t *testing.T) {
	str := `title: Un altro post entusiasmante
datetime: 2024-12-23
id: 20241108-00
frasefamosa : non dire gatto
`

	lex := ScriptGrammar{
		Debug: true,
	}
	err := lex.ParseScript(str)
	if err != nil {
		t.Error("Error is: ", err)
		return
	}

	err = lex.CheckNorm()
	if err != nil {
		t.Error("Error in parser norm ", err)
		return
	}
	err = lex.EvaluateParams()
	if err != nil {
		t.Error("Error in evaluate ", err)
		return
	}

	if frfam, ok := lex.CustomData["frasefamosa"]; ok {
		if frfam != "non dire gatto" {
			t.Error("unexpected custom data", lex.CustomData)
		}
	} else {
		t.Error("custom data missed", lex.CustomData)
	}
}

func TestParseSimpleHtml(t *testing.T) {
	str := `title: Un altro post entusiasmante
datetime: 2024-12-23
id: 20241108-00
---
<p>Pa</p>
il nuovo`

	lex := ScriptGrammar{
		Debug: true,
	}
	err := lex.ParseScript(str)
	if err != nil {
		t.Error("Error is: ", err)
		return
	}

	err = lex.CheckNorm()
	if err != nil {
		t.Error("Error in parser norm ", err)
		return
	}
	err = lex.EvaluateParams()
	if err != nil {
		t.Error("Error in evaluate ", err)
		return
	}
	nrm := lex.Norm["main"]
	lastFns := len(nrm.FnsList) - 1
	stFns := nrm.FnsList[lastFns]
	if len(stFns.Params) != 1 && !stFns.Params[0].IsArray {
		t.Error("expected one array param with lines")
		return
	}
	ll := &stFns.Params[0]
	if len(ll.ArrayValue) != 2 {
		t.Errorf("expected two html lines, but %d", len(ll.ArrayValue))
		return
	}
}

func TestParseHtmlLinkBlock(t *testing.T) {
	str := `title: Un altro post entusiasmante
datetime: 2024-12-23
id: 20241108-00
---
<p>Pa</p>
<p>Tracker: [link 'https://wien-rundumadum-2024-130k.legendstracking.com/']</p>`

	lex := ScriptGrammar{
		Debug:    true,
		TemplDir: "../templates/htmlgen",
	}
	err := lex.ParseScript(str)
	if err != nil {
		t.Error("Error is: ", err)
		return
	}

	err = lex.CheckNorm()
	if err != nil {
		t.Error("Error in parser norm ", err)
		return
	}
	err = lex.EvaluateParams()
	if err != nil {
		t.Error("Error in evaluate ", err)
		return
	}
	nrm := lex.Norm["main"]
	lastFns := len(nrm.FnsList) - 1
	stFns := nrm.FnsList[lastFns]
	if len(stFns.Params) != 1 && !stFns.Params[0].IsArray {
		t.Error("expected one array param with lines")
		return
	}
	ll := &stFns.Params[0]
	if len(ll.ArrayValue) != 2 {
		t.Errorf("expected two html lines, but %d", len(ll.ArrayValue))
		return
	}
	secline := ll.ArrayValue[1]
	if !strings.Contains(secline, "<p>Tracker: <a href=\"https://wien-rundumadum-2024-130k.legendstracking.com/\" target=\"_blank\">https://wien-rundumadum-2024-130k.legendstracking.com/</a></p>") {
		t.Errorf("expected  <a> in generated  html, but %s ", secline)
	}
}

func TestParseHtmlLinkBlockThreeLines(t *testing.T) {
	str := `title: Un altro post entusiasmante
datetime: 2024-12-23
id: 20241108-00
---
<p>Pa</p>
<p>Tracker: [link 'https://wien-rundumadum-2024-130k.legendstracking.com/']</p>
<div>hello</div>`

	lex := ScriptGrammar{
		Debug:    true,
		TemplDir: "../templates/htmlgen",
	}
	err := lex.ParseScript(str)
	if err != nil {
		t.Error("Error is: ", err)
		return
	}

	err = lex.CheckNorm()
	if err != nil {
		t.Error("Error in parser norm ", err)
		return
	}
	err = lex.EvaluateParams()
	if err != nil {
		t.Error("Error in evaluate ", err)
		return
	}
	nrm := lex.Norm["main"]
	lastFns := len(nrm.FnsList) - 1
	stFns := nrm.FnsList[lastFns]
	if len(stFns.Params) != 1 && !stFns.Params[0].IsArray {
		t.Error("expected one array param with lines")
		return
	}
	ll := &stFns.Params[0]
	if len(ll.ArrayValue) != 3 {
		t.Errorf("expected 3 html lines, but %d", len(ll.ArrayValue))
		return
	}
	secline0 := ll.ArrayValue[0]
	if !strings.Contains(secline0, "<p>Pa</p>") {
		t.Errorf("expected <p>Pa</p> in generated  html, but %s ", secline0)
	}
	secline := ll.ArrayValue[1]
	if !strings.Contains(secline, "<p>Tracker: <a href=\"https://wien-rundumadum-2024-130k.legendstracking.com/\" target=\"_blank\">https://wien-rundumadum-2024-130k.legendstracking.com/</a></p>") {
		t.Errorf("expected  <a href> in generated  html, but %s ", secline)
	}
	secline = ll.ArrayValue[2]
	if !strings.Contains(secline, "<div>hello</div>") {
		t.Errorf("expected  <div>hello</div> in generated  html, but %s ", secline)
	}
}

func TestParseHtmlLinkBlockOneLine(t *testing.T) {
	str := `title: Un altro post entusiasmante
datetime: 2024-12-23
id: 20241108-00
---
[link 'https://wien-rundumadum-2024-130k.legendstracking.com/']`

	lex := ScriptGrammar{
		Debug:    true,
		TemplDir: "../templates/htmlgen",
	}
	err := lex.ParseScript(str)
	if err != nil {
		t.Error("Error is: ", err)
		return
	}

	err = lex.CheckNorm()
	if err != nil {
		t.Error("Error in parser norm ", err)
		return
	}
	err = lex.EvaluateParams()
	if err != nil {
		t.Error("Error in evaluate ", err)
		return
	}
	nrm := lex.Norm["main"]
	lastFns := len(nrm.FnsList) - 1
	stFns := nrm.FnsList[lastFns]
	if len(stFns.Params) != 1 && !stFns.Params[0].IsArray {
		t.Error("expected one array param with lines")
		return
	}
	ll := &stFns.Params[0]
	if len(ll.ArrayValue) != 1 {
		t.Errorf("expected one html lines, but %d", len(ll.ArrayValue))
		return
	}
	secline := ll.ArrayValue[0]
	if !strings.Contains(secline, "<a href=\"https://wien-rundumadum-2024-130k.legendstracking.com/\" target=\"_blank\">https://wien-rundumadum-2024-130k.legendstracking.com/</a>") {
		t.Errorf("expected  <a> in generated  html, but %s ", secline)
	}
}

func TestParseHtmlLinkBlockTwoLines(t *testing.T) {
	str := `title: Un altro post entusiasmante
datetime: 2024-12-23
id: 20241108-00
---
[link 'https://wien-rundumadum-2024-130k.legendstracking.com/']<p>
hello</p>`

	lex := ScriptGrammar{
		Debug:    true,
		TemplDir: "../templates/htmlgen",
	}
	err := lex.ParseScript(str)
	if err != nil {
		t.Error("Error is: ", err)
		return
	}

	err = lex.CheckNorm()
	if err != nil {
		t.Error("Error in parser norm ", err)
		return
	}
	err = lex.EvaluateParams()
	if err != nil {
		t.Error("Error in evaluate ", err)
		return
	}
	nrm := lex.Norm["main"]
	lastFns := len(nrm.FnsList) - 1
	stFns := nrm.FnsList[lastFns]
	if len(stFns.Params) != 1 && !stFns.Params[0].IsArray {
		t.Error("expected one array param with lines")
		return
	}
	ll := &stFns.Params[0]
	if len(ll.ArrayValue) != 2 {
		t.Errorf("expected 2 html lines, but %d", len(ll.ArrayValue))
		return
	}
	secline := ll.ArrayValue[0]
	if !strings.Contains(secline, "<a href=\"https://wien-rundumadum-2024-130k.legendstracking.com/\" target=\"_blank\">https://wien-rundumadum-2024-130k.legendstracking.com/</a><p>") {
		t.Errorf("expected  <a> in generated  html, but %s ", secline)
	}
	secline = ll.ArrayValue[1]
	if !strings.Contains(secline, "hello</p>") {
		t.Errorf("expected  hello</p> in generated  html, but %s ", secline)
	}
}
