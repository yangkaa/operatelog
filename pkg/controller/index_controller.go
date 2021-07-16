package controller

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path"

	restful "github.com/emicklei/go-restful"
	"github.com/sirupsen/logrus"
)

//IndexCtrl index page controll
type IndexCtrl struct {
	uiPath    http.Dir
	indexBody []byte
}

//NewIndexCtrl new index ctrl
func NewIndexCtrl() *IndexCtrl {
	uipath := os.Getenv("UI_PATH")
	if uipath == "" {
		curDir, _ := os.Getwd()
		if curDir != "" {
			uipath = path.Join(curDir, "ui", "dist")
		} else {
			uipath = "/app/ui/dist"
		}
	}
	return &IndexCtrl{uiPath: http.Dir(uipath)}
}

func (i IndexCtrl) DisablePrefix() bool {
	return true
}

// WebService returns the restful webservice
func (i IndexCtrl) WebService(ws *restful.WebService) {
	// swagger spec define
	ws.Route(ws.GET("/").To(i.index))
	ws.Route(ws.GET("/static/css/{subpath:*}").To(i.staic))
	ws.Route(ws.GET("/static/js/{subpath:*}").To(i.staic))
	ws.Route(ws.GET("/static/assets/{subpath:*}").To(i.staic))
	ws.Route(ws.GET("/static/images/{subpath:*}").To(i.staic))
	ws.Route(ws.GET("/static/ui/{subpath:*}").To(i.index))
	ws.Route(ws.GET("/static/user/{subpath:*}").To(i.index))
	ws.Route(ws.GET("/static/app-server/{subpath:*}").To(i.appProxy))
}

func (i IndexCtrl) index(request *restful.Request, response *restful.Response) {
	if i.indexBody == nil {
		body, err := ioutil.ReadFile(path.Join(string(i.uiPath), "index.html"))
		if err != nil {
			response.WriteErrorString(404, "index file not found")
			return
		}
		i.indexBody = body
	}
	fmt.Fprint(response.ResponseWriter, string(i.indexBody))
}

func (i IndexCtrl) staic(request *restful.Request, response *restful.Response) {
	http.FileServer(i.uiPath).ServeHTTP(response.ResponseWriter, request.Request)
}

func (i IndexCtrl) appProxy(request *restful.Request, response *restful.Response) {
	storeServer := "http://127.0.0.1:4000"
	url, err := url.Parse(storeServer)
	if err != nil {
		logrus.Error(err)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ServeHTTP(response, request.Request)
}
