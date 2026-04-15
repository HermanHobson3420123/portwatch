package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"portwatch/internal/alert"
	"portwatch/internal/config"
	"portwatch/internal/monitor"
	"portwatch/internal/scanner"
)

func main() {
	configPath := flag.String("config", "", "path to config file (optional)")
	flag.Parse()

	cfg := config.DefaultConfig()
	if *configPath != "" {
		loaded, err := config.Load(*configPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error loading config: %v\n", err)
			os.Exit(1)
		}
		cfg = loaded
	}

	fmt.Printf("portwatch starting — scanning ports %d-%d every %s\n",
		cfg.PortRangeStart, cfg.PortRangeEnd, cfg.Interval)

	sc := scanner.New(cfg.PortRangeStart, cfg.PortRangeEnd, cfg.Timeout)
	al := alert.New(os.Stdout)
	mon := monitor.New(sc, al, cfg.Interval)

	if err := mon.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "error starting monitor: %v\n", err)
		os.Exit(1)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\nshutting down portwatch...")
	mon.Stop()
	_ = time.Second // ensure time import used if needed
}
