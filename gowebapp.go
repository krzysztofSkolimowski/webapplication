package main

import (
	"encoding/json"
	"log"
	"os"
	"runtime"

	"github.com/krzysztofSkolimowski/webapplication/service/shared/jsonconfig"
	"github.com/krzysztofSkolimowski/webapplication/service/shared/session"
	"github.com/krzysztofSkolimowski/webapplication/service/shared/database"
	"github.com/krzysztofSkolimowski/webapplication/service/shared/recaptcha"
	"github.com/krzysztofSkolimowski/webapplication/service/shared/view"
	"github.com/krzysztofSkolimowski/webapplication/service/shared/server"
	"github.com/krzysztofSkolimowski/webapplication/service/route"
	"github.com/krzysztofSkolimowski/webapplication/service/shared/email"
	"github.com/krzysztofSkolimowski/webapplication/service/shared/view/plugin"
)

func init() {
	log.SetFlags(log.Lshortfile)
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	jsonconfig.Load("config"+string(os.PathSeparator)+"config.json", config)
	session.Configure(config.Session)
	database.Connect(config.Database)
	recaptcha.Configure(config.Recaptcha)
	view.Configure(config.View)
	view.LoadTemplates(config.Template.Root, config.Template.Children)
	view.LoadPlugins(
		plugin.TagHelper(config.View),
		plugin.NoEscape(),
		plugin.PrettyTime(),
		recaptcha.Plugin())
	server.Run(route.LoadHTTP(), route.LoadHTTPS(), config.Server)
}

var config = &configuration{}

type configuration struct {
	Database  database.Info   `json:"Database"`
	Email     email.SMTPInfo  `json:"Email"`
	Recaptcha recaptcha.Info  `json:"Recaptcha"`
	Server    server.Server   `json:"Server"`
	Session   session.Session `json:"Session"`
	Template  view.Template   `json:"Template"`
	View      view.View       `json:"View"`
}

func (c *configuration) ParseJSON(b []byte) error {
	return json.Unmarshal(b, &c)
}
