package mhproc

import (
	"bytes"
	"corsa-blog/content/src/mhparser"
	"corsa-blog/idl"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"text/template"
	"time"
)

type MdHtmlProcess struct {
	debug             bool
	scrGramm          mhparser.ScriptGrammar
	HtmlGen           string
	ImgJsonGen        string
	templDir          string
	validateMandatory bool
	SourceDir         string
	RootStaticDir     string
	TargetDir         string
}

func NewMdHtmlProcess(debug bool) *MdHtmlProcess {
	res := MdHtmlProcess{
		debug:             debug,
		validateMandatory: true,
		templDir:          "templates/htmlgen",
	}
	return &res
}

func (mp *MdHtmlProcess) ProcessToHtml(script string) error {
	log.Println("[ProcessToHtml] is called with a script len ", len(script))
	if script == "" {
		return fmt.Errorf("[ProcessToHtml] script is empty")
	}
	mp.scrGramm = mhparser.ScriptGrammar{
		Debug:    mp.debug,
		TemplDir: mp.templDir,
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
	if mp.validateMandatory {
		if mp.scrGramm.Title == "" {
			return fmt.Errorf("[ProcessToHtml] field 'title' in mdhtml is empty")
		}
		if mp.scrGramm.PostId == "" {
			return fmt.Errorf("[ProcessToHtml] field 'id' in mdhtml is empty")
		}
		if mp.scrGramm.Datetime.Year() < 2010 {
			return fmt.Errorf("[ProcessToHtml] field 'datetime' is empty or invalid")
		}
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
	img_data_items := []idl.ImgDataItem{}
	for _, stItem := range normPrg.FnsList {
		if stItem.Type == mhparser.TtHtmlVerbatim {
			lines = append(lines, stItem.Params[0].ArrayValue...)
		}
		if stItem.Type == mhparser.TtJsonBlock {
			labelJson := stItem.Params[0].Label
			if labelJson == "TtJsonImgs" {
				img_arr := idl.ImgDataItems{}
				bb := []byte(stItem.Params[0].Value)
				if err := json.Unmarshal(bb, &img_arr); err != nil {
					log.Println("[parsedToHtml] Unmarshal error")
					return err
				}
				img_data_items = append(img_data_items, img_arr.Images...)
			} else {
				return fmt.Errorf("[parsedToHtml] %s json block not supported", labelJson)
			}
			//fmt.Println("*** json item", stItem.Params[0].Value)
		}
	}
	if len(img_data_items) > 0 {
		imgs := idl.ImgDataItems{Images: img_data_items}
		data_img, err := json.Marshal(imgs)
		if err != nil {
			log.Println("[parsedToHtml] Marshal error")
			return err
		}
		mp.ImgJsonGen = string(data_img)
	}

	if mp.templDir != "" {
		return mp.htmlFromTemplate(lines)
	}
	mp.HtmlGen = strings.Join(lines, "\n")
	mp.printGenHTML()

	return nil
}

func (mp *MdHtmlProcess) printGenHTML() {
	if mp.debug {
		fmt.Printf("***HTML***\n%s\n", mp.HtmlGen)
	}
}

func (mp *MdHtmlProcess) htmlFromTemplate(lines []string) error {
	templName := path.Join(mp.templDir, "post.html")
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
	mp.printGenHTML()
	return nil
}

func (mp *MdHtmlProcess) CreateOrUpdateStaticHtml(sourceName string) error {
	arr := strings.Split(sourceName, "\\")
	if len(arr) < 4 {
		return fmt.Errorf("soure filename is not conform to expected path: <optional/>yyyy/mm/dd/fname.mdhtml")
	}
	log.Println("Processing stack from source ", arr)
	last_ix := len(arr) - 1
	ext := path.Ext(arr[last_ix])
	last_dir := strings.Replace(arr[last_ix], ext, "", -1)
	arr[last_ix] = last_dir
	dir_stack := []string{arr[last_ix-3], arr[last_ix-2], arr[last_ix-1], arr[last_ix]}
	if mp.debug {
		log.Println("dir structure for output ", dir_stack)
	}
	if err := mp.checkOrCreateOutDir(dir_stack); err != nil {
		return err
	}
	log.Println("target dir", mp.TargetDir)
	if err := mp.createIndexHtml(); err != nil {
		return err
	}
	if err := mp.createImageGalleryJson(); err != nil {
		return err
	}

	src_arr := make([]string, 0)
	src_arr = append(src_arr, arr[0:last_ix]...)
	mp.SourceDir = strings.Join(src_arr, "\\")
	log.Println("source dir", mp.SourceDir)

	return nil
}

func (mp *MdHtmlProcess) createIndexHtml() error {
	fname := path.Join(mp.TargetDir, "index.html")
	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(mp.HtmlGen); err != nil {
		return err
	}
	log.Println("file created ", fname)
	return nil
}

func (mp *MdHtmlProcess) createImageGalleryJson() error {
	if mp.ImgJsonGen == "" {
		return nil
	}
	fname := path.Join(mp.TargetDir, "photos.json")
	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(mp.ImgJsonGen); err != nil {
		return err
	}
	log.Println("[createImageGalleryJson] file created ", fname)
	return nil
}

func (mp *MdHtmlProcess) checkOrCreateOutDir(dir_stack []string) error {
	dir_path := mp.RootStaticDir
	for _, item := range dir_stack {
		dir_path = path.Join(dir_path, item)
		//log.Println("check if out dir is here ", dir_path)
		if info, err := os.Stat(dir_path); err == nil && info.IsDir() {
			if mp.debug {
				log.Println("dir exist", dir_path)
			}
		} else {
			if mp.debug {
				log.Println("create dir", dir_path)
			}
			os.MkdirAll(dir_path, 0700)
		}
	}
	mp.TargetDir = dir_path
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
