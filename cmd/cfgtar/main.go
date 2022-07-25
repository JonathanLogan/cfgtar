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
	"path"
	"strings"
)

// cat template.tar | cfgtar config | tar -x -C /
const (
	SchemaFileName = "._config-schema.json"
)

// SELECTOR SUPPORT: -s KEY. Iterate over config[KEY] and produce output to KEY.tar

var (
	flagDryRun      bool
	flagValidateRun bool
	inputFile       string
	inputFd         *os.File
	outputFd        *os.File
	configFile      string
	schemaFile      string
	delim           string
	delimLeft       string
	delimRight      string
	schemaFileName  string
	configData      interface{}
	schemaData      interface{}
	selector        string
	target          string
	selectorData    []string
)

func init() {
	flag.BoolVar(&flagDryRun, "d", false, "dry run (no output)")
	flag.BoolVar(&flagValidateRun, "v", false, "validate before generating output, requires input file")
	flag.StringVar(&inputFile, "i", "", "Input tarfile")
	flag.StringVar(&delim, "D", "{{.}}", "Left|Right delimiter")
	flag.StringVar(&schemaFileName, "S", SchemaFileName, "Name of embedded schema file")
	flag.StringVar(&selector, "s", "", "Selector: Iterate over config.selector and write to selector.tar(s)")
	flag.StringVar(&target, "t", "", "Target directory for selector runs")
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
	outputFd = os.Stdout
	if len(delim) > 0 {
		d := strings.Split(delim, ".")
		if len(d) != 2 || len(d[0]) == 0 || len(d[1]) == 0 {
			printError(5, "Invalid delimiter '%s'", delim)
		}
		delimLeft = d[0]
		delimRight = d[1]
	}
	findSelector()
}

func findSelector() {
	if selector != "" {
		if inputFile == "" {
			printError(4, "Selector requires input file.\n", nil)
		}
		if m, ok := configData.(map[string]interface{}); ok {
			if k, ok := m[selector]; ok {
				if a, ok := k.([]interface{}); ok {
					selectorData = make([]string, 0, len(a))
					for _, v := range a {
						if s, ok := v.(string); ok {
							selectorData = append(selectorData, s)
						} else {
							printError(4, "Selector data not a string: %s\n", selector)
						}
					}
				}
			}
		}
		if selectorData == nil {
			printError(4, "Selector not found: %s\n", selector)
		}
	}
}

func setSelector(pos int, s string) {
	configData.(map[string]interface{})[selector] = struct {
		Pos   int
		Value string
	}{
		Pos:   pos,
		Value: s,
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

func dryRun() {
	if err := tarpipe.TarPipe(inputFd, nil,
		schemareg.New(configData),
		delimLeft, delimRight,
		schemaFileName); err != nil {
		printError(6, "%s\n", err)
	}
	if flagValidateRun {
		inputReset()
	}
}

func selectorDryRun() {
	for k, v := range selectorData {
		setSelector(k, v)
		dryRun()
	}
}

func selectorRun() {
	var err error
	for k, v := range selectorData {
		fn := path.Join(target, v) + ".tar"
		outputFd, err = os.Create(fn)
		if err != nil {
			printError(6, "Cannot create target: %s %s\n", fn, err)
		}
		setSelector(k, v)
		run()
	}
}

func inputReset() {
	if _, err := inputFd.Seek(io.SeekStart, 0); err != nil {
		printError(7, "%s\n", err)
	}
}

func run() {
	if err := tarpipe.TarPipe(inputFd, outputFd,
		schemareg.New(configData),
		delimLeft, delimRight,
		schemaFileName); err != nil {
		printError(20, "%s\n", err)
	}
	_ = outputFd.Sync()
	_ = outputFd.Close()
}

func main() {
	params()
	if flagDryRun || flagValidateRun {
		if len(selector) > 0 {
			selectorDryRun()
		} else {
			dryRun()
		}
	}
	if !flagDryRun || flagValidateRun {
		if len(selector) > 0 {
			selectorRun()
		} else {
			run()
		}
	}
}
