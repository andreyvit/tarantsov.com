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
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/andreyvit/jsonfix"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
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
	ItemsBySection   map[SectionName][]*Item

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
	Name            string
	Ext             string
	Error           error
	ServePath       string
	SourcePath      string
	Frontmatter     *PageFrontmatter
	Raw             []byte
	MarkdownDoc     ast.Node
	TemplateDoc     *template.Template
	Rendered        []byte
	LinkURL         string
	Section         SectionName
	Date            time.Time
	DateStr         string
	Ordinal         int // ordering of articles posted on same day
	DefaultTemplate string
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

type PageFrontmatter struct {
	Title       string          `json:"title"`
	Template    string          `json:"template"`
	Layout      string          `json:"layout"`
	PageClasses []string        `json:"page_classes"`
	CTAs        map[string]*CTA `json:"cta"`
	Path        string          `json:"path"`
}

type CTA struct {
	Title   string `json:"title"`
	LinkURL string `json:"href"`
}

type RenderContext struct {
	Item *Item
}

type PageVM struct {
	*ItemVM
	Site    *SiteVM
	Content template.HTML
}

type SiteVM struct {
	BlogItems []*ItemVM
}

type ItemVM struct {
	item *Item
}

func (vm *ItemVM) LinkURL() string       { return vm.item.LinkURL }
func (vm *ItemVM) Title() string         { return vm.item.Frontmatter.Title }
func (vm *ItemVM) DateStr() string       { return vm.item.DateStr }
func (vm *ItemVM) PageClasses() []string { return vm.item.Frontmatter.PageClasses }

type SectionName string

const (
	Blog SectionName = "blog"
)

const (
	mdExt   = ".md"
	htmlExt = ".html"
	none    = "none"
)

var contentExts = []string{mdExt, htmlExt}

var filenameDatePrefixRe = regexp.MustCompile(`^(\d{4})-(\d{2})-(\d{2})([a-z]?)-`)

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)

	var rootDir string
	var outputDir string
	var listenAddr string
	var isDevMode, isWriteMode bool
	flag.StringVar(&rootDir, "root", ".", "root directory")
	flag.StringVar(&outputDir, "o", "public", "output directory")
	flag.StringVar(&listenAddr, "listen", ":8080", "listen address")
	flag.BoolVar(&isDevMode, "dev", false, "development mode (reload content from disk)")
	flag.BoolVar(&isWriteMode, "w", false, "write content to disk instead of serving it")
	flag.Parse()

	roots := &Roots{
		ContentDir: filepath.Join(rootDir, "content"),
		ThemeDir:   filepath.Join(rootDir, "theme"),
		AssetsDir:  filepath.Join(rootDir, "assets"),
		DataDir:    filepath.Join(rootDir, "data"),
	}

	if isWriteMode {
		build(roots, func(path string, content []byte) {
			fn := filepath.Join(outputDir, path)
			dir := filepath.Dir(fn)
			log.Printf("+ %s", fn)
			if err := os.MkdirAll(dir, 0755); err != nil {
				log.Fatal(err)
			}
			if err := os.WriteFile(fn, content, 0644); err != nil {
				log.Fatal(err)
			}
		})
		log.Printf("✓ build finished")
	} else {
		serve(roots, listenAddr, isDevMode)
	}
}

func build(roots *Roots, write func(path string, content []byte)) {
	lib := loadLibrary(roots)

	var failed bool
	for _, item := range lib.Items {
		if item.Error != nil {
			failed = true
			log.Printf("** failed to build %s: %v", item.ServePath, item.Error)
		}
	}
	if failed {
		log.Fatal("** errors found.")
	}

	servePaths := slices.Sorted(maps.Keys(lib.ItemsByServePath))
	for _, sp := range servePaths {
		item := lib.ItemsByServePath[sp]

		if strings.HasSuffix(sp, "/") {
			sp = sp + "index.html"
		}
		sp = strings.TrimPrefix(sp, "/")
		write(sp, item.Rendered)
	}

	walkDir(roots.AssetsDir, func(fullPath string, relPath string, d fs.DirEntry) {
		raw := must(os.ReadFile(fullPath))
		write(path.Join("assets", relPath), raw)
	})
}

func serve(roots *Roots, listenAddr string, isDevMode bool) {
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
	err := http.ListenAndServe(listenAddr, nil)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

var tags = map[string]func(*Tag, *Library, *RenderContext) (string, error){
	"x-textnav": renderTextNav,
	"x-cta":     renderCTA,
}

func isCurrentItem(item *Item, currentItem *Item) bool {
	return item == currentItem
}

func isParentOfActiveItem(item *Item, currentItem *Item) bool {
	return string(currentItem.Section) == item.Name
}

func renderTextNav(tag *Tag, lib *Library, renderCtx *RenderContext) (string, error) {
	type MainNavItemVM struct {
		*MainNavItem
		IsLink    bool
		IsCurrent bool
		IsActive  bool
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
			vm.IsCurrent = isCurrentItem(ni.Item, renderCtx.Item)
			vm.IsActive = vm.IsCurrent || isParentOfActiveItem(ni.Item, renderCtx.Item)
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

	siteVM := &SiteVM{
		BlogItems: wrapItems(lib.ItemsBySection[Blog]),
	}

	for _, item := range lib.Items {
		if item.Error != nil {
			continue
		}
		var err error
		item.Rendered, err = renderItem(item, lib, siteVM)
		if err != nil {
			failed(item, err)
		}
	}

	return lib
}

func loadContent(contentDir string, lib *Library) {
	var items []*Item
	walkDir(contentDir, func(fullPath, relPath string, d fs.DirEntry) {
		item, err := loadContentItem(fullPath, relPath)
		if err != nil {
			failed(item, err)
			return
		}
		items = append(items, item)
	})

	lib.Items = items
	lib.ItemsByServePath = make(map[string]*Item)
	lib.ItemsByName = make(map[string]*Item)
	lib.ItemsBySection = make(map[SectionName][]*Item)
	for _, item := range items {
		if item.ServePath != "" {
			lib.ItemsByServePath[item.ServePath] = item
		}
		lib.ItemsByName[item.Name] = item

		lib.ItemsBySection[item.Section] = append(lib.ItemsBySection[item.Section], item)
	}

	for _, items := range lib.ItemsBySection {
		slices.SortFunc(items, func(a, b *Item) int {
			return cmp.Or(
				-cmpBool(a.Date.IsZero(), b.Date.IsZero()),
				-a.Date.Compare(b.Date),
				-cmp.Compare(a.Ordinal, b.Ordinal),
			)
		})
	}
}

func loadContentItem(fullPath, relPathWithExt string) (*Item, error) {
	item := &Item{
		SourcePath: fullPath,
	}

	raw := must(os.ReadFile(fullPath))
	relPath, ext := parseExt(relPathWithExt, contentExts)
	if ext == "" {
		return item, fmt.Errorf("unknown file extension")
	}
	item.Name = relPath
	item.Ext = ext

	raw, fm, err := extractFrontmatter[PageFrontmatter](raw)
	if err != nil {
		return item, err
	}
	item.Raw = raw
	item.Frontmatter = fm

	switch item.Ext {
	case mdExt:
		p := parser.NewWithExtensions(parser.CommonExtensions | parser.AutoHeadingIDs | parser.Mmark)
		item.MarkdownDoc = markdown.Parse(item.Raw, p)

		// dump(item.MarkdownDoc, "")
		title := extractHeading(item.MarkdownDoc)
		if fm.Title == "" {
			fm.Title = title
		}
	case htmlExt:
		t := template.New(relPath)
		_, err := t.Parse(string(raw))
		if err != nil {
			return nil, err
		}
		item.TemplateDoc = t
	default:
		return nil, fmt.Errorf("unknown file extension")
	}

	sectionPath, sectionRelPath, ok := strings.Cut(relPath, "/")
	if !ok {
		sectionPath, sectionRelPath = "", sectionPath
	}
	item.Section = SectionName(sectionPath)

	switch item.Section {
	case Blog:
		item.DefaultTemplate = "blog-page"
		sectionPath = "" // serve blog without a prefix
	default:
		item.DefaultTemplate = "main"
	}

	sectionSubpath, baseName, ok := cutLast(sectionRelPath, "/")
	if !ok {
		sectionSubpath, baseName = "", sectionSubpath
	}

	// log.Printf("relPath = %q => section = %q, sectionRelPath = %q, sectionSubpath = %q, baseName = %q", relPath, section, sectionRelPath, sectionSubpath, baseName)

	if m := filenameDatePrefixRe.FindStringSubmatchIndex(baseName); m != nil {
		yr := must(strconv.Atoi(baseName[m[2]:m[3]]))
		mn := must(strconv.Atoi(baseName[m[4]:m[5]]))
		dy := must(strconv.Atoi(baseName[m[6]:m[7]]))
		baseName = baseName[m[1]:]
		item.Date = time.Date(yr, time.Month(mn), dy, 12, 0, 0, 0, time.UTC)
		item.DateStr = item.Date.Format("Jan 2, 2006")

		extra := baseName[m[8]:m[9]]
		if len(extra) > 0 {
			item.Ordinal = int(extra[0] - 'a')
		}
	}

	sectionRelPath = path.Join(sectionSubpath, baseName)
	relPath = path.Join(sectionPath, sectionRelPath)

	// log.Printf("new relPath = %q, sectionRelPath = %q", relPath, sectionRelPath)

	var servePath string
	if s, ok := strings.CutSuffix(relPath, "index"); ok {
		servePath = s
	} else {
		servePath = relPath
	}
	item.ServePath = servePath + "/"

	item.LinkURL = "/" + strings.TrimPrefix(item.ServePath, "/")

	log.Printf("Item(path = %q, name = %q, section = %q, title = %q)", item.ServePath, item.Name, item.Section, item.Frontmatter.Title)
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

func wrapItems(items []*Item) []*ItemVM {
	vms := make([]*ItemVM, len(items))
	for i, item := range items {
		vms[i] = wrapItem(item)
	}
	return vms
}

func wrapItem(item *Item) *ItemVM {
	return &ItemVM{
		item: item,
	}
}

func renderItem(item *Item, lib *Library, siteVM *SiteVM) ([]byte, error) {
	in := &PageVM{
		ItemVM: wrapItem(item),
		Site:   siteVM,
	}

	var content template.HTML
	switch item.Ext {
	case mdExt:
		opts := html.RendererOptions{
			Flags: html.CommonFlags,
		}
		renderer := html.NewRenderer(opts)
		content = template.HTML(markdown.Render(item.MarkdownDoc, renderer))
	case htmlExt:
		var buf bytes.Buffer
		err := item.TemplateDoc.Execute(&buf, in)
		if err != nil {
			return nil, err
		}
		content = template.HTML(buf.Bytes())
	default:
		return nil, fmt.Errorf("unknown file extension")
	}

	renderCtx := &RenderContext{
		Item: item,
	}

	templ := cmp.Or(item.Frontmatter.Template, item.DefaultTemplate, none)

	in.Content = content
	content, err := renderTemplate(templ, lib.Templates, in, content)
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
		if d.IsDir() {
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

func extractFrontmatter[T any](raw []byte) ([]byte, *T, error) {
	const endSep = "\n}\n"
	raw = bytes.TrimSpace(raw)
	fm := new(T)
	if len(raw) > 0 && raw[0] == '{' {
		lenSep := len(endSep)
		end := bytes.Index(raw, []byte(endSep))
		if end < 0 && bytes.HasSuffix(raw, []byte(endSep[:lenSep-1])) {
			lenSep--
			end = len(raw) - lenSep
		}
		if end < 0 {
			return nil, nil, fmt.Errorf("frontmatter: missing end")
		}
		end += lenSep

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

func extractHeading(doc ast.Node) string {
	cn := doc.AsContainer()
	if cn == nil || len(cn.Children) == 0 {
		return ""
	}
	first := cn.Children[0]
	if h, ok := first.(*ast.Heading); ok {
		dump(h, "doc")
		cn.Children = cn.Children[1:]
		return collectText(h)
	}
	return ""
}

func dump(node ast.Node, prefix string) {
	if node == nil {
		return
	}
	dumpSubnode(-1, node, prefix)
}

func dumpSubnode(index int, node ast.Node, prefix string) {
	if index < 0 {
		prefix = cmp.Or(prefix, "D")
	} else {
		prefix = fmt.Sprintf("%s.%02d", prefix, index)
	}
	log.Printf("%s) %T %v", prefix, node, node)
	if cn := node.AsContainer(); cn != nil {
		for i, c := range cn.Children {
			dumpSubnode(i, c, prefix)
		}
	}
}

func collectText(root ast.Node) string {
	var buf strings.Builder
	walkMarkdown(root, func(n ast.Node) {
		// log.Printf("collectText(%T) meets %T", root, n)
		if t, ok := n.(*ast.Text); ok {
			// log.Printf("collectText(%T) = %q", t, t.Literal)
			buf.Write(t.Literal)
		}
	})
	return buf.String()
}

func walkMarkdown(node ast.Node, f func(n ast.Node)) {
	f(node)
	if cn := node.AsContainer(); cn != nil {
		for _, c := range cn.Children {
			walkMarkdown(c, f)
		}
	}
}

// cutLast slices s around the last instance of sep, returning the text before and after sep.
// The found result reports whether sep appears in s.
// If sep does not appear in s, cut returns s, "", false.
func cutLast(s, sep string) (before, after string, found bool) {
	if i := strings.LastIndex(s, sep); i >= 0 {
		return s[:i], s[i+len(sep):], true
	}
	return s, "", false
}

func cmpBool(a, b bool) int {
	if a {
		if !b {
			return 1
		}
	} else {
		if b {
			return -1
		}
	}
	return 0
}
