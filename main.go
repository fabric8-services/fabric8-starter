package main

import (
	"flag"
	"net/http"
	"os"
	"os/user"
	"runtime"

	"github.com/fabric8-services/fabric8-common/configuration"
	"github.com/fabric8-services/fabric8-common/log"
	"github.com/fabric8-services/fabric8-common/sentry"
	"github.com/fabric8-services/fabric8-starter/app"
	"github.com/fabric8-services/fabric8-starter/controller"
	"github.com/goadesign/goa"
	goalogrus "github.com/goadesign/goa/logging/logrus"
	"github.com/goadesign/goa/middleware"
	"github.com/goadesign/goa/middleware/gzip"
	"github.com/google/gops/agent"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	// --------------------------------------------------------------------
	// Parse flags
	// --------------------------------------------------------------------
	var configFilePath string
	var printConfig bool
	flag.StringVar(&configFilePath, "config", "", "Path to the config file to read")
	flag.BoolVar(&printConfig, "printConfig", false, "Prints the config (including merged environment variables) and exits")
	flag.Parse()

	// Override default -config switch with environment variable only if -config switch was
	// not explicitly given via the command line.
	configSwitchIsSet := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "config" {
			configSwitchIsSet = true
		}
	})
	if !configSwitchIsSet {
		if envConfigPath, ok := os.LookupEnv("F8_CONFIG_FILE_PATH"); ok {
			configFilePath = envConfigPath
		}
	}

	config, err := configuration.New(configFilePath)
	if err != nil {
		log.Panic(nil, map[string]interface{}{
			"config_file_path": configFilePath,
			"err":              err,
		}, "failed to setup the configuration")
	}

	if printConfig {
		os.Exit(0)
	}

	// Initialized developer mode flag and log level for the logger
	log.InitializeLogger(config.IsLogJSON(), config.GetLogLevel())

	// Initialize sentry client
	haltSentry, err := sentry.InitializeSentryClient(
		sentry.WithRelease(app.Commit),
		sentry.WithEnvironment(config.GetEnvironment()),
	)
	if err != nil {
		log.Panic(nil, map[string]interface{}{
			"err": err,
		}, "failed to setup the sentry client")
	}
	defer haltSentry()

	printUserInfo()

	// Create service
	service := goa.New("starter")

	// Mount middleware
	service.Use(middleware.RequestID())
	// Use our own log request to inject identity id and modify other properties
	service.Use(gzip.Middleware(9))
	service.Use(app.ErrorHandler(service, true))
	service.Use(middleware.Recover())

	service.WithLogger(goalogrus.New(log.Logger()))

	// service.Use(metric.Recorder())

	// Mount the 'status controller
	statusCtrl := controller.NewStatusController(service)
	app.MountStatusController(service, statusCtrl)

	// Mount the 'label' controller
	labelCtrl := controller.NewLabelController(service)
	app.MountLabelController(service, labelCtrl)

	log.Logger().Infoln("Git Commit SHA: ", app.Commit)
	log.Logger().Infoln("UTC Build Time: ", app.BuildTime)
	log.Logger().Infoln("UTC Start Time: ", app.StartTime)
	log.Logger().Infoln("Dev mode:       ", config.IsPostgresDeveloperModeEnabled())
	log.Logger().Infoln("GOMAXPROCS:     ", runtime.GOMAXPROCS(-1))
	log.Logger().Infoln("NumCPU:         ", runtime.NumCPU())

	http.Handle("/api/", service.Mux)
	http.Handle("/favicon.ico", http.NotFoundHandler())

	if config.GetDiagnoseHTTPAddress() != "" {
		log.Logger().Infoln("Diagnose:       ", config.GetDiagnoseHTTPAddress())
		// Start diagnostic http
		if err := agent.Listen(agent.Options{Addr: config.GetDiagnoseHTTPAddress(), ConfigDir: "/tmp/gops/"}); err != nil {
			log.Error(nil, map[string]interface{}{
				"addr": config.GetDiagnoseHTTPAddress(),
				"err":  err,
			}, "unable to connect to diagnose server")
		}
	}

	// // Start/mount metrics http
	if config.GetHTTPAddress() == config.GetMetricsHTTPAddress() {
		http.Handle("/metrics", promhttp.Handler())
	} else {
		go func(metricAddress string) {
			mx := http.NewServeMux()
			mx.Handle("/metrics", promhttp.Handler())
			if err := http.ListenAndServe(metricAddress, mx); err != nil {
				log.Error(nil, map[string]interface{}{
					"addr": metricAddress,
					"err":  err,
				}, "unable to connect to metrics server")
				service.LogError("startup", "err", err)
			}
		}(config.GetMetricsHTTPAddress())
	}

	// Start http
	if err := http.ListenAndServe(config.GetHTTPAddress(), nil); err != nil {
		log.Error(nil, map[string]interface{}{
			"addr": config.GetHTTPAddress(),
			"err":  err,
		}, "unable to connect to server")
		service.LogError("startup", "err", err)
	}

}

func printUserInfo() {
	u, err := user.Current()
	if err != nil {
		log.Warn(nil, map[string]interface{}{
			"err": err,
		}, "failed to get current user")
	} else {
		log.Info(nil, map[string]interface{}{
			"username": u.Username,
			"uuid":     u.Uid,
		}, "Running as user name '%s' with UID %s.", u.Username, u.Uid)
		g, err := user.LookupGroupId(u.Gid)
		if err != nil {
			log.Warn(nil, map[string]interface{}{
				"err": err,
			}, "failed to lookup group")
		} else {
			log.Info(nil, map[string]interface{}{
				"groupname": g.Name,
				"gid":       g.Gid,
			}, "Running as as group '%s' with GID %s.", g.Name, g.Gid)
		}
	}

}
