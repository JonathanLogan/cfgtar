package main

import (
	"fmt"
	"github.com/JonathanLogan/cfgtar/pkg/lineconfig"
	"github.com/JonathanLogan/cfgtar/pkg/tarpipe"
	"os"
)

// cat template.tar | cfgtar config | tar -x -C /

func main() {
	if len(os.Args) != 2 {
		_, _ = fmt.Fprintf(os.Stderr, "%s: Requires config file argument\n", os.Args[0])
		os.Exit(1)
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Cannot open config file (%s): %s\n", os.Args[1], err)
		os.Exit(2)
	}
	defer func() { _ = f.Close() }()
	config, err := lineconfig.ParseConfig(f)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(3)
	}
	if err := tarpipe.TarPipe(os.Stdin, os.Stdout, config); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(4)
	}
	_ = os.Stdout.Sync()
	_ = os.Stdout.Close()
}
