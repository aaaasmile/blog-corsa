package mhparser

import "sort"

func buildDescrInLex(l *L) {
	arr := make([]DescrFnItem, 0)

	fn := DescrFnItem{
		KeyName:       "link",
		Labels:        []string{"Url"},
		ItemTokenType: itemLinkSimple,
		IsMultiline:   false,
	}
	fn.NumParam = len(fn.Labels)
	arr = append(arr, fn)

	fn = DescrFnItem{
		KeyName:       "figstack",
		ItemTokenType: itemFigStack,
		Labels:        []string{},
		IsMultiline:   true,
	}
	fn.NumParam = len(fn.Labels)
	arr = append(arr, fn)

	arr2 := make([]DescrFnItem, 0, len(arr))
	for ix, v := range arr {
		v.CustomID = ix + 1
		arr2 = append(arr2, v)
	}
	sort.Slice(arr2, func(i, j int) bool {
		return len(arr2[i].KeyName) > len(arr2[j].KeyName)
	})
	l.descrFns = arr2
}
