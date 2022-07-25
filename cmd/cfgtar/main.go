package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/JonathanLogan/cfgtar/pkg/jsonschema"
	"github.com/JonathanLogan/cfgtar/pkg/schemareg"
	"github.com/JonathanLogan/cfgtar/pkg/tarpipe"
	"io"
	"io/ioutil"
	"os"
)

// cat template.tar | cfgtar config | tar -x -C /

var (
	flagDryRun      bool
	flagValidateRun bool
	inputFile       string
	inputFd         *os.File
	configFile      string
	schemaFile      string
	delim           string
	delimLeft       string
	delimRight      string
	configData      interface{}
	schemaData      interface{}
)

func init() {
	flag.BoolVar(&flagDryRun, "d", false, "dry run (no output)")
	flag.BoolVar(&flagValidateRun, "v", false, "validate before generating output, requires input file")
	flag.StringVar(&inputFile, "i", "", "Input tarfile")
	flag.StringVar(&delim, "D", "{{}}", "Left|Right delimiter")
}

func params() {
	var err error
	flag.Parse()
	if flagValidateRun && len(inputFile) == 0 {
		printError(1, "%s: -v implies -i", os.Args[0])
	}
	args := flag.Args()
	switch len(args) {
	case 1:
		configFile = args[0]
	case 2:
		configFile = args[1]
		schemaFile = args[0]
	default:
		printError(1, "%s [<schema.json>] <config.json>", os.Args[0])
	}
	if configData, err = parseJsonFile(configFile); err != nil {
		printError(2, "%s: %s\n", configFile, err)
	}
	if schemaFile != "" {
		if schemaData, err = parseJsonFile(schemaFile); err != nil {
			printError(3, "%s: %s\n", schemaFile, err)
		}
		errPath, _, err := jsonschema.Validate(schemaData, configData)
		if err != nil {
			printError(4, "Schema validation: %v %s\n", errPath, err)
		}
	}
	if inputFile != "" {
		if inputFd, err = os.Open(inputFile); err != nil {
			printError(5, "%s: %s\n", inputFile, err)
		}
	} else {
		inputFd = os.Stdin
	}
	if len(delim) > 0 && len(delim)%2 == 0 {
		delimLeft = delim[:len(delim)/2]
		delimRight = delim[len(delim)/2:]
		fmt.Println(delim[:len(delim)/2], delim[len(delim)/2:])
	} else {
		printError(5, "'%s': delimiter must have even length >0.\n", delim)
	}
}

func printError(exitCode int, format string, v ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", v...)
	os.Exit(exitCode)
}

func parseJsonFile(filename string) (interface{}, error) {
	var ret interface{}
	d, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(d, &ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func main() {
	params()

	if flagDryRun || flagValidateRun {
		if err := tarpipe.TarPipe(inputFd, nil, schemareg.New(configData), delimLeft, delimRight); err != nil {
			printError(6, "%s\n", err)
		}
		if flagValidateRun {
			if _, err := inputFd.Seek(io.SeekStart, 0); err != nil {
				printError(7, "%s\n", err)
			}
		}
	}
	if !flagDryRun || flagValidateRun {
		if err := tarpipe.TarPipe(inputFd, os.Stdout, schemareg.New(configData), delimLeft, delimRight); err != nil {
			printError(20, "%s\n", err)
		}
		_ = os.Stdout.Sync()
		_ = os.Stdout.Close()
	}
}
