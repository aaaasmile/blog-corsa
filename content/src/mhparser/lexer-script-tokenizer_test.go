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

func TestParseHtmlFigStack(t *testing.T) {
	str := `title: Un altro post entusiasmante
datetime: 2024-12-23
id: 20241108-00
---
<p>Ciao</p>
[figstack
  'AustriaBackyardUltra2024011.jpg', 'Partenza mondiale Backyard',
  'backyard_award.png', 'Certificato finale'
]
<p>hello</p>`

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
	secline := ll.ArrayValue[0]
	if !strings.Contains(secline, "<p>Ciao</p>") {
		t.Errorf("expected <p>Ciao</p> in generated  html, but %s ", secline)
		return
	}
	secline = ll.ArrayValue[1]
	if !strings.Contains(secline, `<a id="00" onclick="appGallery.displayImage`) {
		t.Errorf("expected AustriaBackyardUltra2024011 in generated  html, but %s ", secline)
		return
	}
	if !strings.Contains(secline, `<a id="02" onclick="appGallery.displayImage`) {
		t.Errorf("expected backyard_award.png in generated  html, but %s ", secline)
		return
	}

	secline = ll.ArrayValue[2]
	if !strings.Contains(secline, "<p>hello</p>") {
		t.Errorf("expected  <p>hello</p> in generated  html, but %s ", secline)
	}
}

func TestParseZeroTags(t *testing.T) {
	str := `title: Un altro post entusiasmante
datetime: 2024-12-23
id: 20241108-00
tags:
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

	if len(lex.Tags) != 0 {
		t.Error("expected zero tags")
	}
}

func TestParseTwoTags(t *testing.T) {
	str := `title: Un altro post entusiasmante
datetime: 2024-12-23
id: 20241108-00
tags: ultra,adamello
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

	if len(lex.Tags) != 2 {
		t.Error("expected two tags")
		return
	}
	if strings.Compare(lex.Tags[0], "ultra") != 0 {
		t.Error("expected first tag ultra, but ", lex.Tags[0])
		return
	}
	if strings.Compare(lex.Tags[1], "adamello") != 0 {
		t.Error("expected second tag adamello, but ", lex.Tags[0])
		return
	}
}

func TestSimpleHtmlZeroTags(t *testing.T) {
	str := `title: Un altro post entusiasmante
datetime: 2024-12-23
id: 20241108-00
tags:
---
<header class="withimg">
  <div>
    <h1>Quo Vadis</h1>
    <time>4 Gennaio 2023</time>
  </div>
  <img src="bestage.jpg" /> 
</header>
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

	if len(lex.Tags) != 0 {
		t.Error("expected zero tags")
	}
}

func TestParseHtmlLinkWithCaptionBlock(t *testing.T) {
	str := `title: Un altro post entusiasmante
datetime: 2024-12-23
id: 20241108-00
---
<p>Pa</p>
<p>Tracker: [linkcap 'Tracker', 'https://wien-rundumadum-2024-130k.legendstracking.com/']</p>`

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
	if !strings.Contains(secline, "<p>Tracker: <a href=\"https://wien-rundumadum-2024-130k.legendstracking.com/\" target=\"_blank\">Tracker</a></p>") {
		t.Errorf("expected  <a> in generated  html, but %s ", secline)
	}
}
