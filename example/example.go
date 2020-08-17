package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"flag"

	"github.com/s2ks/fcgiserver"
	"github.com/s2ks/fcgiserver/config"
	"github.com/s2ks/fcgiserver/logger"
	"github.com/s2ks/fcgiserver/util"
)

type MyPageConfig struct {
	Path    string `xml:"path,attr"`
	Title   string `xml:"title"`
	SrcFile string `xml:"serve"`
}

type MyConfig struct {
	Name xml.Name       `xml:"user"`
	Page []MyPageConfig `xml:"page"`
}

type MyPageHandler struct {
	Path string

	config *MyConfig
	page   *MyPageConfig
}

func (conf *MyConfig) GetPageFor(path string) *MyPageConfig {
	for _, page := range conf.Page {
		if page.Path == path {
			return &page
		}
	}

	return nil
}

func (p *MyPageHandler) Setup(path string) error {
	p.Path = path
	p.page = p.config.GetPageFor(path)

	if p.page == nil {
		return fmt.Errorf("No configuration entry found for page \"%s\"", path)
	}

	return nil
}

func (p *MyPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if p.Path != r.URL.Path {
		http.NotFound(w, r)
		return
	}

	buf, err := util.ReadFromFile(p.page.SrcFile)

	if err != nil {
		fcgiserver.InternalServerError(w, r, err)
		return
	}

	w.Write(buf)
}

func main() {
	var debug = flag.Bool("debug", false, "Enable debug logging")
	var configpath = flag.String("config", "example.xml", "Path to the configuration file")

	flag.Parse()

	if *debug == true {
		logger.LogLevel(logger.LogLevelDebug)
	}

	conf, err := config.GetServerConfFromXmlFile(*configpath)

	if err != nil {
		logger.Fatal(err)
	}

	server, err := fcgiserver.New(conf.Net.Address, conf.Net.Port, conf.Net.Protocol)

	if err != nil {
		logger.Fatal(err)
	}

	raw, err := config.GetUserXmlFromFile(*configpath)

	if err != nil {
		logger.Fatal(err)
	}

	myconfig := new(MyConfig)

	err = xml.Unmarshal(raw, myconfig)

	if err != nil {
		logger.Fatal(err)
	}

	server.Register("/", &MyPageHandler{config: myconfig})

	logger.Fatal(server.Serve())
}
