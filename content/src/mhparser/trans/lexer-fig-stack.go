package trans

import (
	"bytes"
	"fmt"
	"path"
	"strings"
	"text/template"
)

type figure struct {
	FullName        string
	ReducedFullName string
	Caption         string
	FigId           string
	Sep             string
}

func (fg *figure) calcReduced() error {
	ext := path.Ext(fg.FullName)
	if ext == "" {
		return fmt.Errorf("[calcReduced] extension on %s is empty, this is not supported", fg.FullName)
	}
	bare_name := strings.Replace(fg.FullName, ext, "", -1)
	fg.ReducedFullName = fmt.Sprintf("%s_320%s", bare_name, ext)
	return nil
}

type mdhtFigStackNode struct {
	MdhtLineNode
	figItems    []string
	jsonImgPart string
}

func NewFigStackNode(preline string) *mdhtFigStackNode {
	res := mdhtFigStackNode{figItems: make([]string, 0)}
	arr := strings.Split(preline, "[")
	if len(arr) > 0 {
		res.before_link = arr[0]
	}
	return &res
}

func (ln *mdhtFigStackNode) AddParamString(parVal string) error {
	if parVal == "" {
		return fmt.Errorf("param is empty")
	}
	ln.figItems = append(ln.figItems, parVal)
	return nil
}

func (ln *mdhtFigStackNode) AddblockHtml(val string) error {
	if ln.after_link != "" {
		return fmt.Errorf("[AddblockHtml] already set")
	}
	ln.after_link = val
	return nil
}

func (ln *mdhtFigStackNode) Transform(templDir string) error {
	if templDir == "" {
		return fmt.Errorf("[Transform] templ dir is not set")
	}
	figs := make([]figure, 0)
	is_next_caption := false
	new_fig := figure{}
	for ix, item := range ln.figItems {
		if !is_next_caption {
			new_fig = figure{FullName: item, FigId: fmt.Sprintf("%02d", ix), Sep: ","}
			if err := new_fig.calcReduced(); err != nil {
				return err
			}
			is_next_caption = true
		} else {
			new_fig.Caption = item
			is_next_caption = false
			figs = append(figs, new_fig)
		}
	}
	if len(figs) > 0 {
		figs[len(figs)-1].Sep = ""
	}
	templName := path.Join(templDir, "transform.html")
	tmplPage := template.Must(template.New("FigStack").ParseFiles(templName))
	Ctx := struct {
		Figures []figure
	}{
		Figures: figs,
	}
	var partStack bytes.Buffer
	if err := tmplPage.ExecuteTemplate(&partStack, "figstack", Ctx); err != nil {
		return err
	}

	res := fmt.Sprintf("%s%s%s", ln.before_link, partStack.String(), ln.after_link)
	ln.block = res

	var partJson bytes.Buffer
	if err := tmplPage.ExecuteTemplate(&partJson, "galleryImgItem", Ctx); err != nil {
		return err
	}
	ln.jsonImgPart = partJson.String()
	fmt.Println("*** [mdhtFigStackNode] json part", ln.jsonImgPart)
	return nil
}

func (ln *mdhtFigStackNode) JsonBlock() string {
	return ln.jsonImgPart
}

func (ln *mdhtFigStackNode) HasJsonBlock() bool {
	return len(ln.jsonImgPart) > 0
}
