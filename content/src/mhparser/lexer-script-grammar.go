package mhparser

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type ParamItem struct {
	Label         string
	Value         string
	VariableName  string
	ArrayValue    []string
	IsVariable    bool
	IsUnset       bool
	IsArgument    bool
	IsModificable bool
	IsOverwrite   bool
	IsArray       bool
	_touched      bool
}

func (pi *ParamItem) setVariableName(name string) {
	pi.VariableName = name
	pi.IsVariable = true
	pi._touched = true
}

func (pi *ParamItem) setStringValue(value string) {
	pi.Value = value
	pi._touched = true
}

func (pi *ParamItem) touched() bool {
	return pi._touched
}

func (pi *ParamItem) clear() {
	pi._touched = false
	pi.IsVariable = false
	pi.IsArgument = false
	pi.Value = ""
	pi.Label = ""
	pi.VariableName = ""
}

type FnStatement struct {
	FnName          string
	Params          []ParamItem
	ResHolder       string
	IsAssign        bool // used for aa : va1 in this case FnName is :
	IsInternal      bool // function is an internal builtin function e.g. Printf,..
	HasVariableArgs bool // funtion with variable argument  e.g. Printf
	IsArray         bool // array of strings
}

type ScriptGrammar struct {
	Html     string
	Title    string
	Datetime time.Time
	Norm     map[string]*NormPrg
	st_id    int
	Debug    bool
}

func (sn *ScriptGrammar) ParseScript(source string) error {
	ll := NewL(source, lexStateInit)

	buildDescrInLex(ll)

	defer close(ll.tokens)
	fnstlex := NewFnStatLex()
	sn.Norm = make(map[string]*NormPrg)
	nrmPrg := NewProgNorm("main", false, false)
	sn.Norm[nrmPrg.Name] = nrmPrg
	var err error
	for {
		item := ll.nextItem()
		if sn.Debug {
			fmt.Println("*** type: ", item.Type.String(), item.String())
		}
		switch {
		case item.Type == itemError:
			return errors.New(item.Value)
		case item.Type == itemEndOfStatement:
			if err = sn.storeStatement(ll, fnstlex, nrmPrg); err != nil {
				return err
			}
			fnstlex = NewFnStatLex()
		case item.Type == itemVariable:
			fnstlex.varName = item.Value
		case item.Type == itemAssign:
			fnstlex.isAssign = true
		case item.Type == itemArrayBegin:
			if err := fnstlex.ItemArrayBeginStatement(item); err != nil {
				return err
			}
		case item.Type == itemArrayEnd:
			if err := fnstlex.ItemArrayEndStatement(item); err != nil {
				return err
			}
		case item.Type == itemStringValue:
			if err := fnstlex.ItemStringValueAssignStatement(item); err != nil {
				return err
			}
		case item.Type == itemParamString:
			if err := fnstlex.ItemStringParamAsStatement(ll, item.Value); err != nil {
				return err
			}
		case item.Type == itemEmptyString:
			if err := fnstlex.ItemStringParamAsStatement(ll, ""); err != nil {
				return err
			}
		case item.Type == itemBuiltinFunction:
			if !isLexfnKey(ll, item.ID) {
				return fmt.Errorf("[ParseScript] function %s is not defined", item.Value)
			}
			fnstlex.fnName = item.Value

		case item.Type == itemEOF:
			if err = sn.storeStatement(ll, fnstlex, nrmPrg); err != nil {
				return err
			}
			if nrmPrg.Name != "main" {
				return fmt.Errorf("[ParseScript] missed end of function. Do you forget } at the end of a custom function?")
			}
			return nil
		}
	}
}

func (sn *ScriptGrammar) storeStatement(l *L, fnstlex *FnStatLex, nrmPrg *NormPrg) error {
	if fnstlex.fnName == "" {
		return storeWithEmptyFunction(fnstlex, nrmPrg, sn)
	}
	for _, v := range l.descrFns {
		if strings.Compare(fnstlex.fnName, v.KeyName) == 0 {
			if len(fnstlex.params) != v.NumParam {
				if len(fnstlex.params) < v.NumParam || !v.VariableArgs {
					return fmt.Errorf("[storeStatement]  paramter in %s are %d instead of %d", fnstlex.fnName, len(fnstlex.params), v.NumParam)
				}
			}
			return storeWithFunction(fnstlex, v, sn, nrmPrg)
		}
	}

	return fmt.Errorf("[storeStatement]  function not supported %s", fnstlex.fnName)
}

func (sn *ScriptGrammar) GetNextId() int {
	sn.st_id += 1
	return sn.st_id
}

func (sn *ScriptGrammar) CheckNorm() error {
	varParArr := make([]ParamItem, 0)
	varAssignMain := make(map[string]bool)
	var err error
	if normMain, ok := sn.Norm["main"]; ok {
		varParArr, err = normMain.checkNormItem("main", varAssignMain, varParArr)
		if err != nil {
			return err
		}
	}
	for kname, normItem := range sn.Norm {
		if normItem.Name == "main" {
			continue
		}
		varParArr, err = normItem.checkNormItem(kname, varAssignMain, varParArr)
		if err != nil {
			return err
		}
	}
	if sn.Debug {
		fmt.Println("*** Pararr", varParArr)
		fmt.Println("*** varAssignArr", varAssignMain)
	}

	return nil
}

func storeWithFunction(fnstlex *FnStatLex, v DescrFnItem, sn *ScriptGrammar, nrmPrg *NormPrg) error {
	if sn.Debug {
		fmt.Printf("*** storeWithFunction [norm %s], %v\n", nrmPrg.Name, fnstlex)
	}
	if fnstlex.isAssign && fnstlex.varName != "" {
		// something like res = Sprintf('%s', 'ciao')
		// we need an extra assignement statement
		varfnlex := NewFnStatLex()
		varfnlex.isAssign = true
		varfnlex.varName = fnstlex.varName
		varfnlex.isArray = fnstlex.isArray
		varfnlex.AddParamForVariableAssign()
		if err := storeWithEmptyFunction(varfnlex, nrmPrg, sn); err != nil {
			return err
		}
	}

	fncopy := FnStatement{
		FnName:          fnstlex.fnName,
		ResHolder:       fnstlex.varName,
		IsAssign:        false, //fnstlex.isAssign,
		IsInternal:      v.Internal,
		HasVariableArgs: v.VariableArgs,
		Params:          make([]ParamItem, len(fnstlex.params)),
	}
	copy(fncopy.Params, fnstlex.params)

	_, err := nrmPrg.statementInNormMap(fnstlex.fnName, sn, len(nrmPrg.FnsList)-1)
	if err != nil {
		return err
	}
	nrmPrg.FnsList = append(nrmPrg.FnsList, fncopy)
	if sn.Debug {
		fmt.Printf("*** storeNorm %s append. Count %d \n", nrmPrg.Name, len(nrmPrg.FnsList))
	}
	return nil
}

func storeWithEmptyFunction(fnstlex *FnStatLex, nrmPrg *NormPrg, sn *ScriptGrammar) error {
	if fnstlex.isAssign && fnstlex.varName != "" {
		fncopy := FnStatement{
			IsAssign: true,
			IsArray:  fnstlex.isArray,
			Params:   make([]ParamItem, 1),
		}
		copy(fncopy.Params, fnstlex.params)
		if fncopy.IsArray {
			if len(fncopy.Params) == 1 {
				copy(fncopy.Params[0].ArrayValue, fnstlex.params[0].ArrayValue)
			} else {
				return fmt.Errorf("[storeStatement] array param len is wrong %d", len(fncopy.Params))
			}
		}
		nrmPrg.FnsList = append(nrmPrg.FnsList, fncopy)
		nrmPrg.statementInNormMap(fnstlex.varName, sn, len(nrmPrg.FnsList)-1)
		if sn.Debug {
			fmt.Printf("*** storeWithEmptyFunction [norm %s], %v\n", nrmPrg.Name, fnstlex)
		}
		return nil
	} else if fnstlex.varName != "" {
		return fmt.Errorf("[storeStatement] variable %s without statement", fnstlex.varName)
	}

	return nil
}

func isLexfnKey(l *L, customID int) bool {
	for _, v := range l.descrFns {
		if v.CustomID == customID {
			return true
		}
	}
	return false
}
