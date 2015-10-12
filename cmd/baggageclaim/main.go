package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/pivotal-golang/clock"
	"github.com/pivotal-golang/lager"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/grouper"
	"github.com/tedsuo/ifrit/http_server"
	"github.com/tedsuo/ifrit/sigmon"
	"github.com/xoebus/zest"

	"github.com/concourse/baggageclaim/api"
	"github.com/concourse/baggageclaim/reaper"
	"github.com/concourse/baggageclaim/uidjunk"
	"github.com/concourse/baggageclaim/volume"
	"github.com/concourse/baggageclaim/volume/driver"
)

var listenAddress = flag.String(
	"listenAddress",
	"0.0.0.0",
	"address to listen on",
)

var listenPort = flag.Int(
	"listenPort",
	7788,
	"port for the server to listen on",
)

var volumeDir = flag.String(
	"volumeDir",
	"",
	"directory where volumes and metadata will be stored",
)

var driverType = flag.String(
	"driverType",
	"",
	"the backend driver to use for filesystems",
)

var reapInterval = flag.Duration(
	"reapInterval",
	10*time.Second,
	"interval on which to reap expired containers",
)

var yellerAPIKey = flag.String(
	"yellerAPIKey",
	"",
	"API token to output error logs to Yeller",
)
var yellerEnvironment = flag.String(
	"yellerEnvironment",
	"development",
	"environment label for Yeller",
)

func main() {
	flag.Parse()
	if *volumeDir == "" {
		fmt.Fprintln(os.Stderr, "-volumeDir must be specified")
		os.Exit(1)
	}

	logger := lager.NewLogger("baggageclaim")
	sink := lager.NewReconfigurableSink(lager.NewWriterSink(os.Stdout, lager.DEBUG), lager.INFO)
	logger.RegisterSink(sink)

	if *yellerAPIKey != "" {
		yellerSink := zest.NewYellerSink(*yellerAPIKey, *yellerEnvironment)
		logger.RegisterSink(yellerSink)
	}

	listenAddr := fmt.Sprintf("%s:%d", *listenAddress, *listenPort)

	var volumeDriver volume.Driver

	if *driverType == "btrfs" {
		volumeDriver = driver.NewBtrFSDriver(logger.Session("driver"))
	} else {
		volumeDriver = &driver.NaiveDriver{}
	}

	var namespacer uidjunk.Namespacer

	maxUID, maxUIDErr := uidjunk.DefaultUIDMap.MaxValid()
	maxGID, maxGIDErr := uidjunk.DefaultGIDMap.MaxValid()

	if runtime.GOOS == "linux" && maxUIDErr == nil && maxGIDErr == nil {
		maxId := uidjunk.Min(maxUID, maxGID)

		mappingList := uidjunk.MappingList{
			{
				FromID: 0,
				ToID:   maxId,
				Size:   1,
			},
			{
				FromID: 1,
				ToID:   1,
				Size:   maxId - 1,
			},
		}

		uidTranslator := uidjunk.NewUidTranslator(
			mappingList,
			mappingList,
		)

		namespacer = &uidjunk.UidNamespacer{
			Translator: uidTranslator,
			Logger:     logger.Session("uid-namespacer"),
		}
	} else {
		namespacer = uidjunk.NoopNamespacer{}
	}

	locker := volume.NewLockManager()

	filesystem, err := volume.NewFilesystem(volumeDriver, *volumeDir)
	if err != nil {
		logger.Fatal("failed-to-initialize-filesystem", err)
	}

	volumeRepo := volume.NewRepository(
		logger.Session("repository"),
		filesystem,
		locker,
	)

	strategerizer := volume.NewStrategerizer(namespacer, locker)

	apiHandler, err := api.NewHandler(
		logger.Session("api"),
		strategerizer,
		volumeRepo,
	)
	if err != nil {
		logger.Fatal("failed-to-create-handler", err)
	}

	clock := clock.NewClock()

	morbidReality := reaper.NewReaper(clock, volumeRepo)

	memberGrouper := []grouper.Member{
		{"api", http_server.New(listenAddr, apiHandler)},
		{"reaper", reaper.NewRunner(logger, clock, *reapInterval, morbidReality.Reap)},
	}

	group := grouper.NewParallel(os.Interrupt, memberGrouper)
	running := ifrit.Invoke(sigmon.New(group))

	logger.Info("listening", lager.Data{
		"addr": listenAddr,
	})

	err = <-running.Wait()
	if err != nil {
		logger.Error("exited-with-failure", err)
		os.Exit(1)
	}
}
