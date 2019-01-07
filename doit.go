package main

import "os"

func doIt() {
	//recovery
	initRecov()
	//evt
	initEnvParams()
	//loggers
	initLogger(os.Stdout, os.Stdout, os.Stderr)
	//init
	httpInit()
	//cfg
	initConfig()
	//app entry
	var app ApiHandler
	handleIt(app)

}
