package main

import "flag"

type cmdOptions struct {
	Filename string
}

func parseCommandLineFlags() *cmdOptions {
	var filename string
	flag.StringVar(&filename, "config", "./config.json", "Path to config file")

	flag.Parse()

	return &cmdOptions{
		Filename: filename,
	}
}
