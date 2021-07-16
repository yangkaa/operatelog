package restfulutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/go-openapi/spec"
	"github.com/sirupsen/logrus"

	restful "github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
)

//Controller controller
type Controller interface {
	WebService(ws *restful.WebService)
}

//ServiceType service type
type ServiceType string

//APIForUI service for api
var APIForUI ServiceType = "APIForUI"

//OpenAPI open api
var OpenAPI ServiceType = "OpenAPI"

//Service Microservice main module
type Service struct {
	name            string
	host            string
	port            int
	serviceType     ServiceType
	ws              *restful.WebService
	noPrifixWs      *restful.WebService
	container       *restful.Container
	enableAuth      bool
	version         string
	enableAccessLog bool
}

// NewService new service
func NewService(name, version string, listenHost string, port int, serviceType ServiceType, enableAccessLog bool) *Service {
	container := restful.NewContainer()
	s := &Service{
		name:            name,
		port:            port,
		host:            listenHost,
		serviceType:     serviceType,
		container:       container,
		enableAuth:      true,
		version:         version,
		enableAccessLog: enableAccessLog,
	}
	return s
}

// DisableAuth for test, disable auth
func (s *Service) DisableAuth() {
	s.enableAuth = false
}

func (s *Service) enrichSwaggerObject(swo *spec.Swagger) {
	swo.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       s.name,
			Description: "Resource for managing " + s.name,
			Contact: &spec.ContactInfo{
				ContactInfoProps: spec.ContactInfoProps{},
			},
			License: &spec.License{
				LicenseProps: spec.LicenseProps{
					Name: "GNU",
					URL:  "http://fsf.org",
				},
			},
			Version: "1.0.0",
		},
	}
	swo.Tags = []spec.Tag{spec.Tag{TagProps: spec.TagProps{
		Name:        s.name,
		Description: "kayaking server api"}}}
}

//Add add web service
func (s *Service) Add(services ...Controller) {
	if s.ws == nil {
		s.ws = new(restful.WebService)
		s.ws.Path(fmt.Sprintf("/%s/%s/", s.name, s.version)).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
	}
	for i := range services {
		services[i].WebService(s.ws)
	}
}

//Run service run
func (s *Service) Run(errChan chan error) {
	if s.ws != nil {
		s.container.Add(s.ws)
	}
	if s.noPrifixWs != nil {
		s.container.Add(s.noPrifixWs)
	}
	s.container.Add(restfulspec.NewOpenAPIService(restfulspec.Config{
		Host:                          fmt.Sprintf("127.0.0.1:%d", s.port),
		APIVersion:                    "V1",
		WebServices:                   s.container.RegisteredWebServices(),
		APIPath:                       path.Join(s.name, "/apidocs.json"),
		PostBuildSwaggerObjectHandler: s.enrichSwaggerObject}))
	s.container.DoNotRecover(false)
	s.container.RecoverHandler(s.RecoverHandler)

	if s.serviceType == OpenAPI {
		s.container.Filter(s.permissionOpenAPIFilter)
	}
	if s.serviceType == APIForUI {
		s.container.Filter(s.permissionFilter)
	}
	if s.enableAccessLog {
		s.EnableAccessLog()
	}
	logrus.Infof("Get the %s API using http://localhost:%d/%s", s.serviceType, s.port, path.Join(s.name, "/apidocs.json"))
	go func() {
		listenAddr := fmt.Sprintf("%s:%d", s.host, s.port)
		logrus.Infof("%s restful %s api listen %s", s.name, s.serviceType, listenAddr)
		if err := http.ListenAndServe(listenAddr, s.container); err != nil {
			errChan <- err
		}
	}()
}

//EnableAccessLog enable access log
func (s *Service) EnableAccessLog() {
	// Enable NCSA common log
	s.container.Filter(ncsaCommonLogFormatLogger())
}


func (s *Service) permissionFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	//TODO:
	chain.ProcessFilter(req, resp)
}

func (s *Service) permissionOpenAPIFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	chain.ProcessFilter(req, resp)
}

//RecoverHandler recover handler
func (s *Service) RecoverHandler(panicReason interface{}, res http.ResponseWriter) {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("recover from panic situation: - %v\r\n", panicReason))
	for i := 2; ; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		buffer.WriteString(fmt.Sprintf("    %s:%d\r\n", file, line))
	}
	fmt.Println(buffer.String())

	res.WriteHeader(500)
	re := Result{
		Code: 500,
		Msg:  "server failure",
	}
	json.NewEncoder(res).Encode(re)
}

var logger = log.New(os.Stdout, "", 0)

func ncsaCommonLogFormatLogger() restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		var username = "-"
		if req.Request.URL.User != nil {
			if name := req.Request.URL.User.Username(); name != "" {
				username = name
			}
		}
		chain.ProcessFilter(req, resp)
		logger.Printf("%s - %s [%s] \"%s %s %s\" %d %d",
			strings.Split(req.Request.RemoteAddr, ":")[0],
			username,
			time.Now().Format("02/Jan/2006:15:04:05 -0700"),
			req.Request.Method,
			req.Request.URL.RequestURI(),
			req.Request.Proto,
			resp.StatusCode(),
			resp.ContentLength(),
		)
	}
}
