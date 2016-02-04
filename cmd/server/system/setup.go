package system

import (
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/getsentry/raven-go"
	"github.com/zenazn/goji/web"
	"zenithar.org/go/common/cache"
	"zenithar.org/go/common/eventbus"
	"zenithar.org/go/common/logging/logrus/hooks"
	"zenithar.org/go/common/registrar"

	"github.com/Zenithar/goproject/cmd/server/shared"
	"github.com/Zenithar/goproject/version"
)

// Setup the server
func Setup(flags *shared.Flags) *web.Mux {
	// Set up a new logger
	log := logrus.New()

	// Set the formatter depending on the passed flag's value
	if flags.LogFormatterType == "text" {
		log.Formatter = &logrus.TextFormatter{
			ForceColors: flags.ForceColors,
		}
	} else if flags.LogFormatterType == "json" {
		log.Formatter = &logrus.JSONFormatter{}
	}

	// Pass it to the environment package
	shared.Log = log

	// Connect to raven
	var rc *raven.Client
	if len(strings.TrimSpace(flags.RavenDSN)) > 0 {
		shared.Log.Infoln("**********************************************************")
		shared.Log.Infoln("Initializing Sentry client")
		shared.Log.Infof(" DSN : %s", flags.RavenDSN)

		h, err := os.Hostname()
		if err != nil {
			log.Fatal(err)
		}

		rc, err = raven.NewClient(flags.RavenDSN, map[string]string{
			"hostname": h,
			"app":      "goproject",
			"version":  version.Version,
			"revision": version.Revision,
			"branch":   version.Branch,
		})
		if err != nil {
			log.Fatal(err)
		}

		// Sentry hook
		sentryHook, err := hooks.NewWithClientSentryHook(rc, []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
		})
		if err == nil {
			log.Hooks.Add(sentryHook)
		}
	}

	shared.Raven = rc

	// Initialize the registrar
	shared.Registrar = registrar.Registry()

	shared.Log.Infoln("**********************************************************")
	// Initialize event bus
	log.Infoln("Initializing EvenBus : local mode")
	bus := eventbus.NewLocal()
	shared.Registrar.Register("eventBus", bus)

	shared.Log.Infoln("**********************************************************")
	// Initialize event bus
	if len(strings.TrimSpace(flags.MemcachedHosts)) > 0 {
		log.Infoln("Initializing CacheManager : memcached")
		cacheMgr := cache.NewMemcachedStore(strings.Split(flags.MemcachedHosts, ","), cache.DEFAULT)
		shared.Registrar.Register("cacheManager", cacheMgr)
	} else if len(strings.TrimSpace(flags.RedisHost)) > 0 {
		log.Infoln("Initializing CacheManager : redis")
		cacheMgr := cache.NewRedisCache(flags.RedisHost, "", cache.DEFAULT)
		shared.Registrar.Register("cacheManager", cacheMgr)
	} else {
		log.Infoln("Initializing CacheManager : inMemory")
		cacheMgr := cache.NewInMemoryStore(cache.DEFAULT)
		shared.Registrar.Register("cacheManager", cacheMgr)
	}

	shared.Log.Infoln("**********************************************************")
	// Initialize database
	log.Infoln("Initializing database connection")

	log.Infof("Database type : %s", flags.DatabaseDriver)
	switch flags.DatabaseDriver {
	case "mongodb":
		setupMongoDB(flags)
		break
	case "rethinkdb":
		setupRethinkDB(flags)
		break
	default:
		log.Fatalf("Unrecognised database driver %s, please select between 'mongodb' or 'rethinkdb'", flags.DatabaseDriver)
	}

	// Initialize services
	shared.Log.Infoln("**********************************************************")

	// Create a new goji mux
	mux := web.New()

	// Compile the routes
	mux.Compile()

	return mux
}

// Initialize MongoDB connection and dao
func setupMongoDB(flags *shared.Flags) {

}

// Initialize RethinkDB connection and dao
func setupRethinkDB(flags *shared.Flags) {

}
