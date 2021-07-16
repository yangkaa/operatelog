package main

import (
	"github.com/facebookgo/inject"
	"github.com/sirupsen/logrus"
	"goodrain.com/operatelog/cmd/operatelog/config"
	"goodrain.com/operatelog/pkg/controller"
	"goodrain.com/operatelog/pkg/models"
	"goodrain.com/operatelog/pkg/usecase"
	"goodrain.com/operatelog/pkg/utils/db"
	restfulutil "goodrain.com/operatelog/pkg/utils/restful"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	config.Parse()
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)
	switch config.C.LogLevel {
	case "trace":
		logrus.SetLevel(logrus.TraceLevel)
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}
}

func main() {
	service := restfulutil.NewService("api", "v1", config.C.ListenHost, config.C.ListenPort, restfulutil.APIForUI, config.C.EnableAccessLog)
	instance := db.GetDbInst()

	var g inject.Graph
	g.Logger = logrus.StandardLogger()

	mkyAuditLogCtrl:= controller.NewMkyAuditLogController()
	objects := []*inject.Object{
		{Value: instance.DB()},
		{Value: models.NewMkyAuditLogRepo()},
		{Value: usecase.NewMkyAuditLogUcase()},
		{Value: mkyAuditLogCtrl},
	}

	if err := g.Provide(objects...); err != nil {
		logrus.Fatalf("provide objects to the Graph: %v", err)
	}
	if err := g.Populate(); err != nil {
		logrus.Fatalf("populate the incomplete Objects: %v", err)
	}
	service.Add(mkyAuditLogCtrl)

	errChan := make(chan error, 1)
	service.Run(errChan)

	// step finally: listen Signal
	term := make(chan os.Signal)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	select {
	case err := <-errChan:
		logrus.Errorf("Received collapse error %s, exiting gracefully...", err.Error())
	case <-term:
		logrus.Warn("Received SIGTERM, exiting gracefully...")
	}
	logrus.Info("See you next time!")
}
