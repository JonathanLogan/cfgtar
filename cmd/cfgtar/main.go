package main

import (
	"encoding/json"
	"fmt"
	"github.com/JonathanLogan/cfgtar/pkg/jsonschema"
	"github.com/JonathanLogan/cfgtar/pkg/schemareg"
	"github.com/JonathanLogan/cfgtar/pkg/tarpipe"
	"io/ioutil"
	"os"
)

// cat template.tar | cfgtar config | tar -x -C /

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
	var err error
	var config, schema interface{}
	if len(os.Args) < 2 || len(os.Args) > 3 {
		_, _ = fmt.Fprintf(os.Stderr, "%s [<schema.json>] <config.json>\n", os.Args[0])
		os.Exit(1)
	}
	if config, err = parseJsonFile(os.Args[1]); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[1], err)
		os.Exit(2)
	}
	if len(os.Args) == 3 {
		schema = config
		if config, err = parseJsonFile(os.Args[2]); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[1], err)
			os.Exit(2)
		}
		errPath, _, err := jsonschema.Validate(schema, config)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Schema validation: %v %s\n", errPath, err)
			os.Exit(3)
		}
	}

	if err := tarpipe.TarPipe(os.Stdin, os.Stdout, schemareg.New(config)); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(4)
	}
	_ = os.Stdout.Sync()
	_ = os.Stdout.Close()
}
