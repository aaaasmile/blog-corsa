package mhparser

import "fmt"

type NormPrg struct {
	Name             string
	FnsList          []FnStatement
	statements       map[string]int
	Arguments        []string //exp: in function(aa,bb){, Arguments will be ['aa', 'bb']
	IsLeftdBlock     bool
	IsRightdBlock    bool
	IsCustomFunction bool
}

type normNode struct {
	prg  *NormPrg
	next *normNode
}

type normStack struct {
	start *normNode
}

func (ns *normStack) Pop() *NormPrg {
	if ns.start == nil {
		return nil
	} else {
		node := ns.start
		ns.start = ns.start.next
		return node.prg
	}
}

func (ns *normStack) Push(prg *NormPrg) {
	node := normNode{prg: prg}
	if ns.start == nil {
		ns.start = &node
	} else {
		node.next = ns.start
		ns.start = &node
	}
}

func (np *NormPrg) setArguments(fnst *FnStatLex) {
	for _, v := range fnst.params {
		np.Arguments = append(np.Arguments, v.VariableName)
	}
}

func NewProgNorm(name_st string, left_block, right_block bool) *NormPrg {
	ret := &NormPrg{
		Name:          name_st,
		IsRightdBlock: right_block,
		IsLeftdBlock:  left_block,
		FnsList:       make([]FnStatement, 0),
		Arguments:     make([]string, 0),
		statements:    map[string]int{},
	}
	return ret
}

func (normItem *NormPrg) checkNormItem(kname string, varAssignMain map[string]bool, varParArr []ParamItem) ([]ParamItem, error) {
	varAssign := make(map[string]int)
	for _, v := range normItem.Arguments {
		varAssign[v] = 1
	}
	for _, itemFnSt := range normItem.FnsList {
		if itemFnSt.ResHolder != "" {
			kk := itemFnSt.ResHolder
			if varAssign[kk] > 1 {
				return nil, fmt.Errorf("(CheckNorm) assignement for res holder %s (norm %s) not found", kk, normItem.Name)
			}
			if kname == "main" {
				varAssignMain[kk] = true
			} else {
				varAssign[kk] += 1
			}
		} else if itemFnSt.IsAssign {
			if len(itemFnSt.Params) != 1 {
				return nil, fmt.Errorf("(CheckNorm)assignemt is malformed on %v (norm %s)", itemFnSt, normItem.Name)
			}
			kk := itemFnSt.Params[0].VariableName
			if varAssign[kk] > 0 {
				return nil, fmt.Errorf("(CheckNorm) multiple assignement on %s (norm %s)", kk, normItem.Name)
			}
			if kname == "main" {
				varAssignMain[kk] = true
			} else {
				varAssign[kk] += 1
			}
		}
		for _, pp := range itemFnSt.Params {
			if pp.IsVariable {
				if pp.VariableName == "" {
					return nil, fmt.Errorf("(CheckNorm): variable name is not set %v (norm %s)", pp, normItem.Name)
				}
				if !pp.IsArgument {
					if (varAssign[pp.VariableName] == 0) && !varAssignMain[pp.VariableName] {
						return nil, fmt.Errorf("(CheckNorm): undefined variable %s, norm %s", pp.VariableName, normItem.Name)
					}
				}
				varParArr = append(varParArr, pp)
			}
		}
	}
	return varParArr, nil
}

func (nrmPrg *NormPrg) statementInNormMap(name string, sn *ScriptGrammar, ixst int) (string, error) {
	st_name := fmt.Sprintf("%s-%s-st%d", nrmPrg.Name, name, sn.GetNextId())
	if _, here := nrmPrg.statements[st_name]; here {
		return "", fmt.Errorf("[statementInNormMap]  statement %s is already here", st_name)
	}
	nrmPrg.statements[st_name] = ixst
	return st_name, nil
}
