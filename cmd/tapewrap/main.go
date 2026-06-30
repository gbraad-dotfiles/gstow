package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ducttape-infra/gstow/pkg/gstow"
	"github.com/spf13/cobra"
)

var (
	verbose bool
	dryRun  bool
	force   bool
	target  string
)

var version = "dev"

func main() {
	cmdName := filepath.Base(os.Args[0])

	switch cmdName {
	case "stow", "stow.exe":
		os.Exit(runStowCompatible())
	default:
		os.Exit(runTapewrap())
	}
}

// runStowCompatible implements GNU Stow compatible mode.
// Activated when the binary is invoked as "stow".
func runStowCompatible() int {
	var delete, restow bool
	var packages []string
	for i := 1; i < len(os.Args); i++ {
		a := os.Args[i]
		switch a {
		case "-D", "--delete":
			delete = true
		case "-R", "--restow":
			restow = true
		case "-v", "--verbose":
			verbose = true
		case "-n", "--dry-run":
			dryRun = true
		case "-f", "--force":
			force = true
		case "-t", "--target":
			i++
			if i >= len(os.Args) {
				fmt.Fprintln(os.Stderr, "error: -t requires a value")
				return 1
			}
			target = os.Args[i]
		default:
			if strings.HasPrefix(a, "-") {
				fmt.Fprintf(os.Stderr, "unknown flag: %s\n", a)
				return 1
			}
			packages = append(packages, a)
		}
	}
	if delete && restow {
		fmt.Fprintln(os.Stderr, "error: -D and -R are mutually exclusive")
		return 1
	}
	if len(packages) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: stow [-D] [-R] [-t TARGET] [-v] [-n] [-f] PACKAGE...\n")
		return 1
	}
	cfg := buildConfig()
	for _, pkg := range packages {
		var err error
		switch {
		case restow:
			err = gstow.Restow(cfg, pkg)
		case delete:
			err = gstow.Unstow(cfg, pkg)
		default:
			err = gstow.Stow(cfg, pkg)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s: %v\n", pkg, err)
			return 1
		}
	}
	return 0
}

// runTapewrap is the advanced tapewrap mode.
// Activated when the binary is invoked as "tapewrap".
func runTapewrap() int {
	rootCmd := &cobra.Command{
		Use:          "tapewrap [packages...]",
		Short:        "Manage symlink farms for dotfiles",
		SilenceUsage: true,
		Version:      version,
		Args:         cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}
			cfg := buildConfig()
			for _, pkg := range args {
				if err := gstow.Stow(cfg, pkg); err != nil {
					return err
				}
			}
			return nil
		},
	}

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "n", false, "Dry run")
	rootCmd.PersistentFlags().BoolVarP(&force, "force", "f", false, "Force")
	rootCmd.PersistentFlags().StringVarP(&target, "target", "t", "", "Target directory")

	rootCmd.AddCommand(&cobra.Command{
		Use:   "stow [packages...]",
		Short: "Stow packages",
		Args:  cobra.MinimumNArgs(1),
		RunE:  runAction(gstow.Stow),
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:   "unstow [packages...]",
		Short: "Unstow packages",
		Args:  cobra.MinimumNArgs(1),
		RunE:  runAction(gstow.Unstow),
	})
	rootCmd.AddCommand(&cobra.Command{
		Use:   "wrap [packages...]",
		Short: "Stow packages (alias for stow)",
		Args:  cobra.MinimumNArgs(1),
		RunE:  runAction(gstow.Stow),
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "restow [packages...]",
		Short: "Restow packages (unstow then stow)",
		Args:  cobra.MinimumNArgs(1),
		RunE:  runAction(gstow.Restow),
	})

	if err := rootCmd.Execute(); err != nil {
		return 1
	}
	return 0
}

func runAction(fn func(*gstow.Config, string) error) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		cfg := buildConfig()
		for _, pkg := range args {
			if err := fn(cfg, pkg); err != nil {
				return err
			}
		}
		return nil
	}
}

func buildConfig() *gstow.Config {
	stowDir, _ := filepath.Abs(".")
	tgt := target
	if tgt == "" {
		tgt = filepath.Dir(stowDir)
	}
	targetDir, _ := filepath.Abs(tgt)
	return &gstow.Config{
		StowDir:   stowDir,
		TargetDir: targetDir,
		Verbose:   verbose,
		DryRun:    dryRun,
		Force:     force,
	}
}
