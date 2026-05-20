package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/YangYuS8/phanes/internal/builder"
	"github.com/YangYuS8/phanes/internal/config"
	"github.com/YangYuS8/phanes/internal/contracts"
	runtimehttp "github.com/YangYuS8/phanes/internal/runtime"
)

var (
	version = "dev"
	commit  = "unknown"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "phanes: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		return usage()
	}

	switch args[0] {
	case "version", "--version", "-v":
		fmt.Printf("phanes %s (%s)\n", version, commit)
		return nil
	case "runtime":
		return runtimeCommand(args[1:])
	case "builder":
		return builderCommand(args[1:])
	case "cache":
		return cacheCommand(args[1:])
	case "save":
		return saveCommand(args[1:])
	case "launcher":
		return launcherCommand(args[1:])
	case "help", "--help", "-h":
		return usage()
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func runtimeCommand(args []string) error {
	if len(args) == 0 {
		return errors.New("runtime subcommand is required")
	}
	switch args[0] {
	case "start":
		fs := flag.NewFlagSet("runtime start", flag.ContinueOnError)
		cfgPath := fs.String("config", "", "config path")
		bind := fs.String("bind", contracts.DefaultBindHost, "bind host")
		port := fs.Uint("port", 0, "HTTP port")
		cacheDir := fs.String("cache-dir", "", "cache directory")
		saveDB := fs.String("save-db", "", "save database path")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		cfg, err := configFromFlags(*cfgPath, *bind, uint32(*port), *cacheDir, *saveDB)
		if err != nil {
			return err
		}
		server, err := runtimehttp.NewServer(*cfg, runtimehttp.Options{OnReady: func(server *runtimehttp.Server) {
			fmt.Fprintf(os.Stderr, "runtime ready: %s\n", server.URL())
		}})
		if err != nil {
			return err
		}
		return server.Run(context.Background())
	case "status", "stop":
		fs := flag.NewFlagSet("runtime "+args[0], flag.ContinueOnError)
		url := fs.String("url", "", "runtime status URL")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		if *url == "" {
			return errors.New("--url is required")
		}
		return notImplemented("runtime " + args[0] + " validated")
	default:
		return fmt.Errorf("unknown runtime subcommand %q", args[0])
	}
}

func builderCommand(args []string) error {
	if len(args) == 0 {
		return errors.New("builder subcommand is required")
	}
	switch args[0] {
	case "verify":
		fs := flag.NewFlagSet("builder verify", flag.ContinueOnError)
		cacheDir := fs.String("cache-dir", "", "cache directory")
		deep := fs.Bool("deep", false, "deep hash check")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		if *cacheDir == "" {
			return errors.New("--cache-dir is required")
		}
		verification, err := contracts.VerifyCacheRoot(*cacheDir, contracts.CacheVerifyOptions{DeepHashCheck: *deep})
		if err != nil {
			return err
		}
		fmt.Printf("cache verification passed: %s\n", verification.Manifest.GetCacheId())
		return nil
	case "build":
		fs := flag.NewFlagSet("builder build", flag.ContinueOnError)
		output := fs.String("output", "", "output cache directory")
		mode := fs.String("mode", "", "build mode")
		_ = fs.String("input", "", "local input path")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		if *output == "" {
			return errors.New("--output is required")
		}
		if *mode != "embedded-minimal" {
			return notImplemented("builder build currently supports only --mode embedded-minimal")
		}
		manifest, err := builder.BuildEmbeddedMinimal(*output)
		if err != nil {
			return err
		}
		fmt.Printf("embedded-minimal cache built: %s\n", manifest.GetCacheId())
		return nil
	case "detect":
		return notImplemented("builder detect command is declared but not implemented")
	default:
		return fmt.Errorf("unknown builder subcommand %q", args[0])
	}
}

func cacheCommand(args []string) error {
	if len(args) == 0 {
		return errors.New("cache subcommand is required")
	}
	switch args[0] {
	case "inspect":
		fs := flag.NewFlagSet("cache inspect", flag.ContinueOnError)
		cacheDir := fs.String("cache-dir", "", "cache directory")
		_ = fs.Bool("json", false, "JSON output")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		if *cacheDir == "" {
			return errors.New("--cache-dir is required")
		}
		if err := contracts.ValidateCacheRoot(*cacheDir); err != nil {
			return err
		}
		return notImplemented("cache inspect validated")
	case "clean":
		fs := flag.NewFlagSet("cache clean", flag.ContinueOnError)
		cacheDir := fs.String("cache-dir", "", "cache directory")
		saveDB := fs.String("save-db", "", "save database path")
		if err := fs.Parse(args[1:]); err != nil {
			return err
		}
		if *cacheDir == "" || *saveDB == "" {
			return errors.New("--cache-dir and --save-db are required")
		}
		if err := contracts.CacheCleanAllowed(*cacheDir, *saveDB); err != nil {
			return err
		}
		return notImplemented("cache clean validated")
	default:
		return fmt.Errorf("unknown cache subcommand %q", args[0])
	}
}

func saveCommand(args []string) error {
	if len(args) == 0 {
		return errors.New("save subcommand is required")
	}
	switch args[0] {
	case "list-profiles", "export", "import":
		return notImplemented("save " + args[0] + " command is declared but not implemented")
	default:
		return fmt.Errorf("unknown save subcommand %q", args[0])
	}
}

func launcherCommand(args []string) error {
	if len(args) == 0 {
		return errors.New("launcher subcommand is required")
	}
	if args[0] != "dev" {
		return fmt.Errorf("unknown launcher subcommand %q", args[0])
	}
	return notImplemented("launcher dev command is declared but not implemented")
}

func configFromFlags(configPath string, bind string, port uint32, cacheDir string, saveDB string) (*config.Config, error) {
	if configPath != "" {
		cfg, err := config.Load(configPath)
		if err != nil {
			return nil, err
		}
		if bind != "" {
			cfg.Runtime.BindHost = bind
		}
		if port != 0 {
			cfg.Runtime.HTTPPort = port
		}
		if cacheDir != "" {
			cfg.Cache.Dir = cacheDir
		}
		if saveDB != "" {
			cfg.Save.DBPath = saveDB
		}
		return cfg, cfg.Validate()
	}
	cfg := &config.Config{
		Runtime: config.RuntimeConfig{BindHost: bind, HTTPPort: port},
		Cache:   config.CacheConfig{Dir: cacheDir},
		Save:    config.SaveConfig{DBPath: saveDB},
	}
	return cfg, cfg.Validate()
}

func usage() error {
	fmt.Println("phanes <version|runtime|builder|cache|save|launcher> ...")
	return nil
}

func notImplemented(message string) error {
	return fmt.Errorf("%s; implementation pending", message)
}
