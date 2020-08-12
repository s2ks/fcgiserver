package main

import (
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/s2ks/fcgiserver"
	"github.com/s2ks/fcgiserver/config"
	"github.com/s2ks/fcgiserver/logger"
	"github.com/s2ks/fcgiserver/util"
)

type MyConfigPage struct {
	Path    string `xml:"path,attr"`
	Title   string `xml:"title"`
	SrcFile string `xml:"serve"`
}

type MyConfig struct {
	Name xml.Name       `xml:"user"`
	Page []MyConfigPage `xml:"page"`
}

type MyPage struct {
	Path string

	config *MyConfig
	page   *MyConfigPage
}

func (conf *MyConfig) GetPageFor(path string) *MyConfigPage {
	for _, page := range conf.Page {
		if page.Path == path {
			return &page
		}
	}

	return nil
}

func (p *MyPage) Setup(path string) error {
	p.Path = path
	p.page = p.config.GetPageFor(path)

	if p.page == nil {
		return fmt.Errorf("No configuration entry found for page \"%s\"", path)
	}

	return nil
}

func (p *MyPage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	conf, err := config.GetServerConfFromXmlFile("example.xml")

	if err != nil {
		logger.Fatal(err)
	}

	server, err := fcgiserver.New(conf.Net.Address, conf.Net.Port, conf.Net.Protocol)

	if err != nil {
		logger.Fatal(err)
	}

	buf, err := config.GetUserConfFromXmlFile("example.xml")

	if err != nil {
		logger.Fatal(err)
	}

	myconfig := new(MyConfig)

	err = xml.Unmarshal(buf, myconfig)

	if err != nil {
		logger.Fatal(err)
	}

	server.Register("/", &MyPage{config: myconfig})

	logger.Fatal(server.Serve())
}
