package main

import (
	"crypto/x509"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"path"
	"regexp"
	"strings"
	"time"
)

const (
	usageCredentials = "use to set the info/credentials on Imgur"
	usageConfig      = "use to set the config file parameter with info/credentials on Imgur"
	usageHttpPort    = "use to set HTTP @ Port No."

	//status
	StatusInProgress = "in-progress"
	StatusSuccess    = "success"
	StatusFailed     = "failed"
	StatusPending    = "pending"
	StatusDownloaded = "downloaded"
	StatusComplete   = "complete"
)

var (
	pLogDir = "."
	//loggers
	infoLog  *log.Logger
	warnLog  *log.Logger
	errorLog *log.Logger

	//signal flag
	pStillRunning = true

	pBuildTime = "0"
	pVersion   = "0.1.0" + "-" + pBuildTime
	//console
	pShowConsole = true
	//envt
	pEnvVars = map[string]*string{
		"MONGERS_IMGUR_API": &pLogDir,
	}

	//ssl certs
	pool *x509.CertPool

	pDumpTmp          = "/tmp/imgur-tmf"
	pHttpPort         = "7777"
	pUserCredentials  = ""
	pConfCredentials  = ""
	pParamConfig      *UserCredential
	pParamConfigOk    bool
	userCodeChan      chan string
	pShowDebug        bool
	pImageUniqHash    map[string]string
	pImageHistory     map[string]*UploadRecords
	pUploaderChan     chan *UploadRecords
	userBearerChan    chan string
	pImageHistoryChan chan *UploadRecords
)

func init() {
	//uniqueness
	rand.Seed(time.Now().UnixNano())

	//defaults here
	userCodeChan = make(chan string, 100)
	userBearerChan = make(chan string, 100)
	pShowDebug = true

	//for testing purposes, dont set limit :-)
	pImageUniqHash = make(map[string]string)
	pImageHistory = make(map[string]*UploadRecords)
	pUploaderChan = make(chan *UploadRecords, 100)
	pImageHistoryChan = make(chan *UploadRecords, 100)

	//init certs
	pool = x509.NewCertPool()
	pool.AppendCertsFromPEM(pemCerts)
}

//initRecov is for dumpIng segv in
func initRecov() {
	//might help u
	defer func() {
		recvr := recover()
		if recvr != nil {
			fmt.Println("MAIN-RECOV-INIT: ", recvr)
		}
	}()
}

//os.Stdout, os.Stdout, os.Stderr
func initLogger(i, w, e io.Writer) {
	//just in case
	if !pShowConsole {
		infoLog = makeLogger(i, pLogDir, "imgur", "INFO: ")
		warnLog = makeLogger(w, pLogDir, "imgur", "WARN: ")
		errorLog = makeLogger(e, pLogDir, "imgur", "ERROR: ")
	} else {
		infoLog = log.New(i,
			"INFO: ",
			log.Ldate|log.Ltime|log.Lmicroseconds)
		warnLog = log.New(w,
			"WARN: ",
			log.Ldate|log.Ltime|log.Lshortfile)
		errorLog = log.New(e,
			"ERROR: ",
			log.Ldate|log.Ltime|log.Lshortfile)
	}
}

//initEnvParams enable all OS envt vars to reload internally
func initEnvParams() {
	//just in-case, over-write from ENV
	for k, v := range pEnvVars {
		if os.Getenv(k) != "" {
			*v = os.Getenv(k)
		}
	}
	//get options
	flag.StringVar(&pUserCredentials, "credentials", pUserCredentials, usageCredentials)
	flag.StringVar(&pConfCredentials, "config", pConfCredentials, usageConfig)

	flag.StringVar(&pHttpPort, "port", pHttpPort, usageHttpPort)
	flag.Parse()

}

//initConfig set defaults for initial reqmts
func initConfig() {

	//either 1 should be present
	if pConfCredentials == "" && pUserCredentials == "" {
		showUsage()
		return
	}

	//check user credentials
	pParamConfigOk, pParamConfig = formatConfig(pConfCredentials, pUserCredentials)
	if pParamConfig == nil || !pParamConfigOk {
		log.Println("USER CREDENTIALS PARAMETER IS EMPTY/INVALID_FORMAT")
		os.Exit(0)
	}

	//dont wait
	isReady := make(chan bool, 1)
	go manageImageLocalDownloader(isReady)
	<-isReady
	isHistReady := make(chan bool, 1)
	go manageImageHistory(isHistReady)
	<-isHistReady

	//show permission URL
	showMacroMsg(initAuthBearerToken())
}

//initDownloadDir try to init all filehandles for logs
func initDownloadDir() {
	if _, err := os.Stat(pDumpTmp); os.IsNotExist(err) {
		//mkdir -p
		os.MkdirAll(pDumpTmp, os.ModePerm)
	}
}

//formatLogger try to init all filehandles for logs
func formatLogger(fdir, fname, pfx string) string {
	t := time.Now()
	r := regexp.MustCompile("[^a-zA-Z0-9]")
	p := t.Format("2006-01-02") + "-" + r.ReplaceAllString(strings.ToLower(pfx), "")
	s := path.Join(pLogDir, fdir)
	if _, err := os.Stat(s); os.IsNotExist(err) {
		//mkdir -p
		os.MkdirAll(s, os.ModePerm)
	}
	return path.Join(s, p+"-"+fname+".log")
}

//makeLogger initialize the logger either via file or console
func makeLogger(w io.Writer, ldir, fname, pfx string) *log.Logger {
	logFile := w
	if !pShowConsole {
		var err error
		logFile, err = os.OpenFile(formatLogger(ldir, fname, pfx), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
		if err != nil {
			log.Println(err)
		}
	}
	//give it
	return log.New(logFile,
		pfx,
		log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

}

//showUsage
func showUsage() {
	fmt.Println("Version:", pVersion)
	flag.PrintDefaults()
	os.Exit(0)
}
