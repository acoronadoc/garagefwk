package garagefwk

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type GarageRequest struct {
	Params      *map[string]string
	Body        *map[string]interface{}
	VarsUrl     *map[string]interface{}
	Db          *sql.DB
	Screen      *Screen
	ScreenTypes *map[string]interface{}
	User        string
}

func InitGarageFWK(app *map[string]interface{}) *sql.DB {
	screens := (*app)["screens"].(map[string]interface{})
	fwkFolder := modulePath()
	config := readConfig()
	scrTypes := initScreensTypes(&screens)
	db := connectDB(config)
	defer db.Close()

	(*app)["db"] = db

	http.HandleFunc("/api/admin", func(w http.ResponseWriter, req *http.Request) {
		parts := make([]map[string]interface{}, 0)
		request := GarageRequest{
			Body:        ReadForm(req),
			Params:      ReadParams(req),
			Db:          db,
			Screen:      getScreenByUrl(config, req.URL.Query().Get("url")),
			ScreenTypes: scrTypes,
			User:        "anonymous",
		}

		if request.Screen != nil {
			request.VarsUrl = getVarsURL(request.Screen.Url, req.URL.Query().Get("url"))

			execScreenType(&parts, &request)
		}

		w.Header().Set("Content-Type", "application/json")

		ret := make(map[string]interface{})
		ret["parts"] = parts
		json.NewEncoder(w).Encode(ret)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		_, err := os.Stat(fwkFolder + "/static" + req.URL.Path)

		if err == nil && req.URL.Path != "/" {
			http.ServeFile(w, req, fwkFolder+"/static"+req.URL.Path)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		jsonMenu, errJson := json.Marshal(config.Menus["sidebarmenu"])

		if errJson != nil {
			panic(errJson)
		}

		serveTemplate(fwkFolder+"/templates/index.html.gohtml", map[string]interface{}{
			"css":         template.CSS(mergeFiles(fwkFolder+"/webcomponents", ".css")),
			"js":          template.JS(mergeFiles(fwkFolder+"/webcomponents", ".js")),
			"sidebarmenu": string(jsonMenu),
		}, w)
	})

	http.ListenAndServe("0.0.0.0:8080", nil)

	return db
}

func serveTemplate(templateFile string, params map[string]interface{}, w http.ResponseWriter) {
	var content bytes.Buffer
	tmpl := template.Must(template.ParseFiles(templateFile))

	err := tmpl.Execute(&content, params)

	if err != nil {
		panic(err)
	}

	w.Write(content.Bytes())
}

func modulePath() string {
	pc, _, _, _ := runtime.Caller(1)
	f, _ := runtime.FuncForPC(pc).FileLine(pc)

	return filepath.Dir(f)
}

func mergeFiles(folder string, suffix string) []byte {
	var buffer bytes.Buffer

	filesDir, err := os.ReadDir(folder)

	if err != nil {
		panic(err)
	}

	for _, info := range filesDir {
		if strings.HasSuffix(info.Name(), suffix) {
			data, err := os.ReadFile(folder + "/" + info.Name())

			if err != nil {
				panic(err)
			}

			buffer.Write(data)
			buffer.WriteByte('\n')
		}
	}

	return buffer.Bytes()
}
