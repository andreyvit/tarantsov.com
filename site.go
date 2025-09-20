package main

import (
	"bytes"
	"cmp"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"maps"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/andreyvit/jsonfix"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

type Roots struct {
	ContentDir string
	ThemeDir   string
	AssetsDir  string
	DataDir    string
}

type Library struct {
	Items            []*Item
	ItemsByServePath map[string]*Item
	ItemsByName      map[string]*Item

	Templates map[string]*Template
	Layouts   map[string]*Template
	Partials  map[string]*Template

	MainNavItems []*MainNavItem

	Errors []error
}

func (lib *Library) AddError(err error) {
	log.Printf("** %v", err)
	lib.Errors = append(lib.Errors, err)
}

type ErrSink interface {
	AddError(err error)
}

type Item struct {
	Name        string
	Ext         string
	Error       error
	ServePath   string
	SourcePath  string
	Frontmatter *Frontmatter
	Raw         []byte
	Rendered    []byte
	LinkURL     string
}

type MainNavItem struct {
	RawText  string `json:"text"`
	ItemName string `json:"item"`
	Item     *Item  `json:"-"`

	TextPrefix string `json:"-"`
	TextCore   string `json:"-"`
	TextSuffix string `json:"-"`
	FullText   string `json:"-"`
}

type Template struct {
	Name  string
	Templ *template.Template
}

type Frontmatter struct {
	Title       string          `json:"title"`
	Template    string          `json:"template"`
	Layout      string          `json:"layout"`
	PageClasses []string        `json:"page_classes"`
	CTAs        map[string]*CTA `json:"cta"`
}

type CTA struct {
	Title   string `json:"title"`
	LinkURL string `json:"href"`
}

type RenderContext struct {
	Item *Item
}

type LayoutInput struct {
	Title       string
	PageClasses []string
	Content     template.HTML
}

const (
	mdExt   = ".md"
	htmlExt = ".html"
	none    = "none"
)

var contentExts = []string{mdExt, htmlExt}

func main() {
	log.SetOutput(os.Stderr)

	var rootDir string
	var listenAddr string
	var isDevMode bool
	flag.StringVar(&rootDir, "root", ".", "root directory")
	flag.StringVar(&listenAddr, "listen", ":8080", "listen address")
	flag.BoolVar(&isDevMode, "dev", false, "development mode (reload content from disk)")
	flag.Parse()

	roots := &Roots{
		ContentDir: filepath.Join(rootDir, "content"),
		ThemeDir:   filepath.Join(rootDir, "theme"),
		AssetsDir:  filepath.Join(rootDir, "assets"),
		DataDir:    filepath.Join(rootDir, "data"),
	}
	sharedLib := loadLibrary(roots)

	http.HandleFunc("GET /assets/{path...}", func(w http.ResponseWriter, r *http.Request) {
		path := r.PathValue("path")
		log.Printf("serving asset: %q", path)
		http.ServeFileFS(w, r, os.DirFS(roots.AssetsDir), path)
	})
	http.HandleFunc("GET /{path...}", func(w http.ResponseWriter, r *http.Request) {
		lib := sharedLib
		if isDevMode {
			lib = loadLibrary(roots)
		}

		path := r.PathValue("path")
		if path == "" {
			path = "/"
		}
		log.Printf("serving: %q", path)
		item := lib.ItemsByServePath[path]
		if item == nil {
			http.NotFound(w, r)
			return
		}
		if item.Error != nil {
			fmt.Fprintf(w, "Error: %v", item.Error)
		} else {
			w.Write(item.Rendered)
		}
	})
	log.Printf("Listening on %s", listenAddr)
	http.ListenAndServe(listenAddr, nil)
}

var tags = map[string]func(*Tag, *Library, *RenderContext) (string, error){
	"x-textnav": renderTextNav,
	"x-cta":     renderCTA,
}

func renderTextNav(tag *Tag, lib *Library, renderCtx *RenderContext) (string, error) {
	type MainNavItemVM struct {
		*MainNavItem
		IsLink    bool
		IsCurrent bool
		LinkURL   string
	}
	type TextNavInput struct {
		Variant      string
		MainNavItems []*MainNavItemVM
	}

	items := make([]*MainNavItemVM, 0, len(lib.MainNavItems))
	for _, ni := range lib.MainNavItems {
		vm := &MainNavItemVM{
			MainNavItem: ni,
		}
		if ni.Item != nil {
			vm.IsLink = true
			vm.IsCurrent = (ni.Item == renderCtx.Item)
			vm.LinkURL = ni.Item.LinkURL
		}
		items = append(items, vm)
	}

	in := &TextNavInput{
		Variant:      tag.Attrs["variant"],
		MainNavItems: items,
	}
	return renderPartial(lib, "textnav", in)
}

func renderCTA(tag *Tag, lib *Library, renderCtx *RenderContext) (string, error) {
	// {{ $cta := index $.Page.Params.cta (default "primary" (.Get 0)) -}}
	// <div class="cta">
	//     <header><a href="{{ $cta.href }}">{{ $cta.title }}</a></header>
	// </div>

	type CTAInput struct {
		*CTA
	}

	name := cmp.Or(tag.Attrs["name"], "primary")
	cta := renderCtx.Item.Frontmatter.CTAs[name]
	if cta == nil {
		return "", fmt.Errorf("cta not found: %q", name)
	}

	in := &CTAInput{
		CTA: cta,
	}
	return renderPartial(lib, "cta", in)
}

var navItemTextRe = regexp.MustCompile(`\[(.*)\]`)

func loadLibrary(roots *Roots) *Library {
	lib := &Library{}
	loadContent(roots.ContentDir, lib)
	lib.Templates = loadTemplates(filepath.Join(roots.ThemeDir, "templates"))
	lib.Layouts = loadTemplates(filepath.Join(roots.ThemeDir, "layouts"))
	lib.Partials = loadTemplates(filepath.Join(roots.ThemeDir, "partials"))

	loadDataFile(filepath.Join(roots.DataDir, "nav/main.json"), &lib.MainNavItems, lib)
	log.Printf("Loaded %d main nav items", len(lib.MainNavItems))

	for _, ni := range lib.MainNavItems {
		if m := navItemTextRe.FindStringSubmatchIndex(ni.RawText); m != nil {
			ni.TextPrefix = ni.RawText[:m[0]]
			ni.TextCore = ni.RawText[m[2]:m[3]]
			ni.TextSuffix = ni.RawText[m[1]:]
		} else {
			ni.TextCore = ni.RawText
		}
		ni.FullText = ni.TextPrefix + ni.TextCore + ni.TextSuffix

		if ni.ItemName != "" {
			item := lib.ItemsByName[ni.ItemName]
			if item == nil {
				lib.AddError(fmt.Errorf("nav item %q references unknown content item %q", ni.RawText, ni.ItemName))
				continue
			}
			ni.Item = item
			// log.Printf("nav item %q references content item %q", ni.FullText, item.Name)
		}
	}

	for _, item := range lib.Items {
		if item.Error != nil {
			continue
		}
		var err error
		item.Rendered, err = renderItem(item, lib)
		if err != nil {
			failed(item, err)
		}
	}

	return lib
}

func loadContent(contentDir string, lib *Library) {
	var items []*Item
	walkDir(contentDir, func(fullPath, relPath string, d fs.DirEntry) {
		item, err := loadContentFile(fullPath, relPath)
		if err != nil {
			failed(item, err)
			return
		}
		items = append(items, item)
	})

	lib.Items = items
	lib.ItemsByServePath = make(map[string]*Item)
	lib.ItemsByName = make(map[string]*Item)
	for _, item := range items {
		if item.ServePath != "" {
			lib.ItemsByServePath[item.ServePath] = item
		}
		lib.ItemsByName[item.Name] = item
	}
}

func loadContentFile(fullPath, relPath string) (*Item, error) {
	item := &Item{
		SourcePath: fullPath,
	}

	raw := must(os.ReadFile(fullPath))
	base, ext := parseExt(relPath, contentExts)
	if ext == "" {
		return item, fmt.Errorf("unknown file extension")
	}
	item.Name = base
	item.Ext = ext

	raw, fm, err := extractFrontmatter(raw)
	if err != nil {
		return item, err
	}
	item.Raw = raw
	item.Frontmatter = fm

	var servePath string
	if s, ok := strings.CutSuffix(base, "index"); ok {
		if s == "" {
			servePath = "/"
		} else {
			servePath = s + "/"
		}
	} else {
		servePath = base + "/"
	}
	item.ServePath = servePath

	item.LinkURL = "/" + strings.TrimPrefix(item.ServePath, "/")

	log.Printf("path = %q, name = %q, frontmatter = %s", item.ServePath, item.Name, jsonstr(item.Frontmatter))
	return item, nil
}

func loadTemplates(dir string) map[string]*Template {
	result := make(map[string]*Template)
	walkDir(dir, func(fullPath, relPath string, d fs.DirEntry) {
		t, err := loadTemplateFile(fullPath, relPath)
		if err != nil {
			log.Printf("** %s: %v", fullPath, err)
			return
		}
		result[t.Name] = t
	})
	log.Printf("loaded templates: %s", strings.Join(slices.Sorted(maps.Keys(result)), ", "))
	return result
}

func loadTemplateFile(fullPath, relPath string) (*Template, error) {
	raw := must(os.ReadFile(fullPath))

	base, ext := parseExt(relPath, []string{".html"})
	if ext == "" {
		return nil, fmt.Errorf("unknown file extension")
	}

	t := template.New(base)
	_, err := t.Parse(string(raw))
	if err != nil {
		return nil, err
	}

	return &Template{
		Name:  base,
		Templ: t,
	}, nil
}

func loadDataFile(fullPath string, v any, sink ErrSink) {
	raw := must(os.ReadFile(fullPath))
	err := json.Unmarshal(jsonfix.Bytes(raw), v)
	if err != nil {
		sink.AddError(fmt.Errorf("%s: %w", fullPath, err))
	}
}

func renderItem(item *Item, lib *Library) ([]byte, error) {
	var content template.HTML
	switch item.Ext {
	case mdExt:
		p := parser.NewWithExtensions(parser.CommonExtensions | parser.AutoHeadingIDs | parser.Mmark)
		content = template.HTML(markdown.ToHTML(item.Raw, p, nil))
	case htmlExt:
		content = template.HTML(item.Raw)
	default:
		return nil, fmt.Errorf("unknown file extension")
	}

	renderCtx := &RenderContext{
		Item: item,
	}

	in := &LayoutInput{
		Title:       item.Frontmatter.Title,
		PageClasses: item.Frontmatter.PageClasses,
	}

	in.Content = content
	content, err := renderTemplate(cmp.Or(item.Frontmatter.Template, none), lib.Templates, in, content)
	if err != nil {
		return nil, err
	}

	in.Content = content
	content, err = renderTemplate(cmp.Or(item.Frontmatter.Layout, "default"), lib.Layouts, in, content)
	if err != nil {
		return nil, err
	}

	content, err = renderPseudoTags(content, lib, renderCtx)
	if err != nil {
		return nil, err
	}

	return []byte(content), nil
}

func renderTemplate(templName string, avail map[string]*Template, in any, content template.HTML) (template.HTML, error) {
	if templName == none {
		return content, nil
	}
	t := avail[templName]
	if t == nil {
		return "", fmt.Errorf("unknown template %q", templName)
	}

	var buf bytes.Buffer
	err := t.Templ.Execute(&buf, in)
	if err != nil {
		return "", fmt.Errorf("%s: %w", templName, err)
	}
	s := template.HTML(buf.Bytes())
	return s, nil
}

func renderPartial(lib *Library, templName string, in any) (string, error) {
	s, err := renderTemplate(templName, lib.Partials, in, "")
	return string(s), err
}

func renderPseudoTags(content template.HTML, lib *Library, renderCtx *RenderContext) (template.HTML, error) {
	result, err := findAndRenderTags(string(content), func(tag *Tag) (string, error) {
		f := tags[tag.Name]
		if f != nil {
			return f(tag, lib, renderCtx)
		}
		return "", fmt.Errorf("unknown tag <%s>", tag.Name)
	})
	return template.HTML(result), err
}

type Tag struct {
	Name  string
	Attrs map[string]string
}

var startRe = regexp.MustCompile(`<(x-\w+)`)

func findAndRenderTags(content string, f func(tag *Tag) (string, error)) (string, error) {
	var result strings.Builder
	for {
		m := startRe.FindStringSubmatchIndex(content)
		if m == nil {
			result.WriteString(content)
			break
		}
		result.WriteString(content[:m[0]])

		tag := &Tag{Name: content[m[2]:m[3]]}

		endStr := fmt.Sprintf("</%s>", tag.Name)
		i := strings.Index(content[m[1]:], endStr)
		if i < 0 {
			return "", fmt.Errorf("unclosed <%s>", tag.Name)
		}
		i += m[1] + len(endStr)
		tagXML := content[m[0]:i]
		content = content[i:]

		decoder := xml.NewDecoder(strings.NewReader(tagXML))
		var attrs map[string]string
		for {
			tok, err := decoder.Token()
			if err != nil {
				return "", fmt.Errorf("failed to parse tag %s: %w", tag.Name, err)
			}
			if startElem, ok := tok.(xml.StartElement); ok {
				attrs = make(map[string]string)
				for _, attr := range startElem.Attr {
					attrs[attr.Name.Local] = attr.Value
				}
				break
			}
		}
		tag.Attrs = attrs

		rendered, err := f(tag)
		if err != nil {
			return "", err
		}
		result.WriteString(rendered)
	}
	return result.String(), nil
}

func failed(item *Item, err error) {
	if item.Error == nil {
		item.Error = err
	}
	log.Printf("** %s: %v", item.SourcePath, err)
}

func walkDir(dir string, fn func(fullPath, relPath string, d fs.DirEntry)) {
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if strings.HasPrefix(d.Name(), ".") {
			return nil
		}
		if path == dir {
			return nil
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		fn(path, rel, d)
		return nil
	})
	if err != nil {
		log.Fatalf("** %v", fmt.Errorf("failed to walk %s: %w", dir, err))
	}
}

func extractFrontmatter(raw []byte) ([]byte, *Frontmatter, error) {
	const endSep = "\n}\n"
	raw = bytes.TrimSpace(raw)
	fm := new(Frontmatter)
	if len(raw) > 0 && raw[0] == '{' {
		end := bytes.Index(raw, []byte(endSep))
		if end < 0 {
			return nil, nil, fmt.Errorf("frontmatter: missing end")
		}
		end += len(endSep)

		err := json.Unmarshal(jsonfix.Bytes(raw[:end]), fm)
		if err != nil {
			return nil, nil, fmt.Errorf("frontmatter: %w", err)
		}

		raw = bytes.TrimSpace(raw[end:])
	}
	return raw, fm, nil
}

func parseExt(fn string, validExts []string) (string, string) {
	for _, ext := range validExts {
		if base, ok := strings.CutSuffix(fn, ext); ok {
			return base, ext
		}
	}
	return "", ""
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func ensure(err error) {
	if err != nil {
		panic(err)
	}
}

func jsonstr(v any) string {
	return string(must(json.Marshal(v)))
}
