package mhparser

import "testing"

func TestParseScript(t *testing.T) {
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
		t.Error("Error in parser ", err)
	}
	err = lex.EvaluateParams()
	if err != nil {
		t.Error("Error in evaluate ", err)
	}
	if lex.PostId != "20241108-00" {
		t.Error("unexpected id", lex.PostId)
	}
	if lex.Title != "Prossima gara Wien Rundumadum" {
		t.Error("unexpected Title", lex.Title)
	}
}
