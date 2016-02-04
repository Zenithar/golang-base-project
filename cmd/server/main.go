package main

import (
	"net"
	"net/http"
	"time"

	"github.com/Zenithar/goproject/cmd/server/shared"
	"github.com/Zenithar/goproject/cmd/server/system"
	"github.com/Zenithar/goproject/version"

	"github.com/Sirupsen/logrus"
	"github.com/namsral/flag"
	"github.com/zenazn/goji/graceful"
)

var (
	// Enable namsral/flag logality
	configFlag = flag.String("config", "", "config file to load")

	// General flags
	bindAddress      = flag.String("bind", ":5000", "Network address used to bind")
	logFormatterType = flag.String("log", "text", "Log formatter type. Either \"json\" or \"text\"")
	forceColors      = flag.Bool("force_colors", false, "Force colored prompt?")
	ravenDSN         = flag.String("raven_dsn", "", "Defines the sentry endpoint dsn")

	databaseDriver    = flag.String("db_driver", "mongodb", "Specify the database to use (mongodb, rethinkdb)")
	databaseHost      = flag.String("db_host", "localhost:27017", "Database hosts, split by ',' to add several hosts")
	databaseNamespace = flag.String("db_namespace", "hammurapi", "Select the database")
	databaseUser      = flag.String("db_user", "", "Database user")
	databasePassword  = flag.String("db_password", "", "Database user password")

	memcachedHosts = flag.String("memcached_hosts", "", "Memcached servers for cache (ex: 127.0.0.1:11211)")
	redisHost      = flag.String("redis_host", "", "Redis server for cache")
)

func main() {
	// Parse the flags
	flag.Parse()

	logrus.Infoln("**********************************************************")
	logrus.Infoln("goproject server starting ...")
	logrus.Infof("Version : %s (%s-%s)", version.Version, version.Revision, version.Branch)

	// Set localtime to UTC
	time.Local = time.UTC

	// Put config into the environment package
	shared.Config = &shared.Flags{
		BindAddress:      *bindAddress,
		LogFormatterType: *logFormatterType,
		ForceColors:      *forceColors,
		RavenDSN:         *ravenDSN,

		DatabaseDriver:    *databaseDriver,
		DatabaseHost:      *databaseHost,
		DatabaseNamespace: *databaseNamespace,
		DatabaseUser:      *databaseUser,
		DatabasePassword:  *databasePassword,

		MemcachedHosts: *memcachedHosts,
		RedisHost:      *redisHost,
	}

	// Generate a mux
	mux := system.Setup(shared.Config)

	// Make the mux handle every request
	http.Handle("/", mux)

	// Log that we're starting the server
	shared.Log.WithFields(logrus.Fields{
		"address": shared.Config.BindAddress,
	}).Info("Starting the HTTP server")

	// Initialize the goroutine listening to signals passed to the app
	graceful.HandleSignals()

	// Pre-graceful shutdown event
	graceful.PreHook(func() {
		shared.Log.Info("Received a signal, stopping the application")
	})

	// Post-shutdown event
	graceful.PostHook(func() {
		shared.Log.Info("Stopped the application")
	})

	// Listen to the passed address
	listener, err := net.Listen("tcp", shared.Config.BindAddress)
	if err != nil {
		shared.Log.WithFields(logrus.Fields{
			"error":   err,
			"address": *bindAddress,
		}).Fatal("Cannot set up a TCP listener")
	}

	// Start the listening
	err = graceful.Serve(listener, http.DefaultServeMux)
	if err != nil {
		// Don't use .Fatal! We need the code to shut down properly.
		shared.Log.Error(err)
	}

	// If code reaches this place, it means that it was forcefully closed.

	// Wait until open connections close.
	graceful.Wait()
}
