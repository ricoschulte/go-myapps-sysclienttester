package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"strings"

	"time"

	"github.com/ricoschulte/go-myapps-sysclienttester/server"
	"github.com/ricoschulte/go-myapps/sysclient"
	log "github.com/sirupsen/logrus"
)

//go:embed static/*
var StaticDevicesFolderFS embed.FS

var secretKey = flag.String("secretkey", "", "secretkey used to encrypt local session files")
var idFile = flag.String("idfile", "identities.json", "path to a JSON file with device identities")
var logLevel = flag.String("loglevel", "info", "log level to use.")
var sysclientUrl = flag.String("sysclient", "", "the sysclient url of the devices app to connect to (wss://apps.company.com/company.com/devices/sysclients)")
var insecureSkipVerify = flag.Bool("skipverifytls", false, "disabled verification of TLS certificates")
var staticDir = flag.String("staticdir", "", "path to a local folder to serve static files under <device-passthrough>/static")
var sessionDir = flag.String("sessiondir", "", "path to a local folder to save the sessions")

func main() {
	flag.Parse()

	if *secretKey == "" {
		fmt.Println("secretkey cant be empty. see -h")
		os.Exit(1)
	}
	if *sysclientUrl == "" {
		fmt.Println("sysclientUrl cant be empty, no host to connect to. see -h")
		os.Exit(1)
	}

	if *sessionDir != "" && !strings.HasSuffix(*sessionDir, "/") {
		fmt.Printf("the path to a sessionsDir must end with a slash /")
		os.Exit(1)
	} else if *sessionDir != "" {
		staticDirInfo, err := os.Stat(*sessionDir)
		if err != nil {
			fmt.Printf("error stat dir '%s': %v\n", *sessionDir, err)
			os.Exit(1)
		}
		if !staticDirInfo.IsDir() {
			fmt.Printf("error stat dir '%s' is not a directory\n", *sessionDir)
			os.Exit(1)
		}
	}

	if *staticDir != "" {
		staticDirInfo, err := os.Stat(*staticDir)
		if err != nil {
			fmt.Printf("error stat dir '%s': %v\n", *staticDir, err)
			os.Exit(1)
		}
		if !staticDirInfo.IsDir() {
			fmt.Printf("error stat dir '%s' is not a directory\n", *staticDir)
			os.Exit(1)
		}
	}

	fb, err := os.ReadFile(*idFile)
	if err != nil {
		fmt.Printf("error reading identity file '%s': %v\n", *idFile, err)
		os.Exit(1)
	}

	err_init_logging := initLogging()
	if err_init_logging != nil {
		fmt.Printf("error init logging: %v\n", err_init_logging)
		os.Exit(1)
	}

	var identities Identities
	err_f := json.Unmarshal(fb, &identities)
	if err_f != nil {
		fmt.Printf("d: %v\n", err_f)
	}

	for _, identity := range identities.List {
		sc := startClient([]byte(*secretKey), &identity)
		go sc.Connect()
	}
	select {}
}

func startClient(secretkey []byte, identity *sysclient.Identity) *sysclient.Sysclient {
	url, _ := url.Parse(*sysclientUrl)
	sc := sysclient.Sysclient{
		Identity:           *identity,
		Url:                *sysclientUrl,
		Timeout:            time.Duration(2 * time.Second),
		InsecureSkipVerify: *insecureSkipVerify,
		Tunnels:            map[int32]*sysclient.SysclientTunnel{},

		FileSysclientPassword:      fmt.Sprintf("%ssysclient_%s_%s_password.bin", *sessionDir, url.Hostname(), identity.Id),
		FileAdministrativePassword: fmt.Sprintf("%ssysclient_%s_%s_administrativepassword.bin", *sessionDir, url.Hostname(), identity.Id),
		SecretKey:                  secretkey,
	}

	if *staticDir != "" {
		sc.ServeMux = server.GetServerMux(&sc, http.Dir(*staticDir))
	} else {
		fsRoot, _ := fs.Sub(StaticDevicesFolderFS, "static")
		sc.ServeMux = server.GetServerMux(&sc, http.FS(fsRoot))
	}
	return &sc
}

type Identities struct {
	List []sysclient.Identity `json:"ids"`
}

func initLogging() error {
	// Log as JSON instead of the default ASCII formatter.
	//log.SetFormatter(&log.JSONFormatter{})

	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
		PadLevelText:  true,
	})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	//log.SetLevel(log.InfoLevel)
	level, err := log.ParseLevel(*logLevel)
	if err != nil {
		return err
	}
	log.SetLevel(level)
	return nil
}
