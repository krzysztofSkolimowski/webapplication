package view

import (
	"encoding/gob"
	"sync"
	"net/http"
	"github.com/krzysztofSkolimowski/webapplication/service/shared/session"
	"strings"
	"path/filepath"
	"os"
	"encoding/json"
	"net/url"
	"fmt"
	"html/template"
)

func init() {
	gob.Register(Flash{})
}

var (
	FlashError = "alert-danger"
	FlashSuccess = "alert-success"
	FlashNotice = "alert-info"
	FlashWarning = "alert-warning"

	childTemplates     []string
	rootTemplate       string
	templateCollection = make(map[string]*template.Template)
	pluginCollection   = make(template.FuncMap)
	mutex              sync.RWMutex
	mutexPlugins       sync.RWMutex
	sessionName        string
	viewInfo           View
)

type Template struct {
	Root     string   `json:"Root"`
	Children []string `json:"Children"`
}

type View struct {
	BaseURI   string
	Extension string
	Folder    string
	Name      string
	Caching   bool
	Vars      map[string]interface{}
	request   *http.Request
}

type Flash struct {
	Message string
	Class   string
}

func Configure(vi View) {
	viewInfo = vi
}

func ReadConfig() View {
	return viewInfo
}

func LoadTemplates(rootTemp string, childTemps []string) {
	rootTemplate = rootTemp
	childTemplates = childTemps
}

func LoadPlugins(fms ...template.FuncMap) {
	fm := make(template.FuncMap)
	for _, m := range fms {
		for k, v := range m {
			fm[k] = v
		}
	}
	mutexPlugins.Lock()
	pluginCollection = fm
	mutexPlugins.Unlock()
}

func (v *View) PrependBaseURI(s string) string {
	return v.BaseURI + s
}

func New(r *http.Request) *View {
	v := &View{}
	v.Vars = make(map[string]interface{})
	v.Vars["AuthLevel"] = "anon"
	v.BaseURI = viewInfo.BaseURI
	v.Extension = viewInfo.Extension
	v.Folder = viewInfo.Folder
	v.Name = viewInfo.Name
	v.Vars["BaseURI"] = v.BaseURI
	v.request = r
	s := session.Instance(v.request)
	if s.Values["id"] != nil {
		v.Vars["AuthLevel"] = "auth"
	}
	return v
}

func (v *View) AssetTimePath(s string) (string, error) {
	if strings.HasPrefix(s, "//") {
		return s, nil
	}

	s = strings.TrimLeft(s, "/")
	abs, err := filepath.Abs(s)

	if err != nil {
		return "", err
	}

	time, err2 := FileTime(abs)
	if err2 != nil {
		return "", err2
	}

	return v.PrependBaseURI(s + "?" + time), nil
}

func (v *View) RenderSingle(w http.ResponseWriter) {
	mutexPlugins.RLock()
	pc := pluginCollection
	mutexPlugins.RUnlock()

	templateList := []string{v.Name}
	for i, name := range templateList {
		path, err := filepath.Abs(v.Folder + string(os.PathSeparator) + name + "." + v.Extension)
		if err != nil {
			http.Error(w, "Template Path Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		templateList[i] = path
	}

	templates, err := template.New(v.Name).Funcs(pc).ParseFiles(templateList...)
	if err != nil {
		http.Error(w, "Template Parse Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tc := templates
	s := session.Instance(v.request)

	if flashes := s.Flashes(); len(flashes) > 0 {
		v.Vars["flashes"] = make([]Flash, len(flashes))
		for i, f := range flashes {
			switch f.(type) {
			case Flash:
				v.Vars["flashes"].([]Flash)[i] = f.(Flash)
			default:
				v.Vars["flashes"].([]Flash)[i] = Flash{f.(string), "alert-box"}
			}

		}
		s.Save(v.request, w)
	}

	err = tc.Funcs(pc).ExecuteTemplate(w, v.Name+"."+v.Extension, v.Vars)
	if err != nil {
		http.Error(w, "Template File Error: "+err.Error(), http.StatusInternalServerError)
	}
}

func (v *View) Render(w http.ResponseWriter) {

	mutex.RLock()
	tc, ok := templateCollection[v.Name]
	mutex.RUnlock()

	mutexPlugins.RLock()
	pc := pluginCollection
	mutexPlugins.RUnlock()

	if !ok || !viewInfo.Caching {

		var templateList []string
		templateList = append(templateList, rootTemplate)
		templateList = append(templateList, v.Name)
		templateList = append(templateList, childTemplates...)

		for i, name := range templateList {
			path, err := filepath.Abs(v.Folder + string(os.PathSeparator) + name + "." + v.Extension)
			if err != nil {
				http.Error(w, "Template Path Error: "+err.Error(), http.StatusInternalServerError)
				return
			}
			templateList[i] = path
		}

		templates, err := template.New(v.Name).Funcs(pc).ParseFiles(templateList...)

		if err != nil {
			http.Error(w, "Template Parse Error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		mutex.Lock()
		templateCollection[v.Name] = templates
		mutex.Unlock()

		tc = templates
	}

	s := session.Instance(v.request)

	if flashes := s.Flashes(); len(flashes) > 0 {
		v.Vars["flashes"] = make([]Flash, len(flashes))
		for i, f := range flashes {
			switch f.(type) {
			case Flash:
				v.Vars["flashes"].([]Flash)[i] = f.(Flash)
			default:
				v.Vars["flashes"].([]Flash)[i] = Flash{f.(string), "alert-box"}
			}

		}
		s.Save(v.request, w)
	}

	err := tc.Funcs(pc).ExecuteTemplate(w, rootTemplate+"."+v.Extension, v.Vars)

	if err != nil {
		http.Error(w, "Template File Error: "+err.Error(), http.StatusInternalServerError)
	}
}

func Validate(r *http.Request, required []string) (bool, string) {
	for _, v := range required {
		if r.FormValue(v) == "" {
			return false, v
		}
	}

	return true, ""
}

func (v *View) SendFlashes(w http.ResponseWriter) {
	s := session.Instance(v.request)

	flashes := peekFlashes(w, v.request)
	s.Save(v.request, w)

	js, err := json.Marshal(flashes)
	if err != nil {
		http.Error(w, "JSON Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func peekFlashes(w http.ResponseWriter, r *http.Request) []Flash {
	s := session.Instance(r)

	var v []Flash

	if flashes := s.Flashes(); len(flashes) > 0 {
		v = make([]Flash, len(flashes))
		for i, f := range flashes {
			switch f.(type) {
			case Flash:
				v[i] = f.(Flash)
			default:
				v[i] = Flash{f.(string), "alert-box"}
			}

		}
	}

	return v
}

func Repopulate(l []string, src url.Values, dst map[string]interface{}) {
	for _, v := range l {
		dst[v] = src.Get(v)
	}
}

func FileTime(path string) (string, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	mtime := fi.ModTime().Unix()
	return fmt.Sprintf("%v", mtime), nil
}
