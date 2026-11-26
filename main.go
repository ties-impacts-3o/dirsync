package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type options struct {
	src, dst string
	dryRun   bool
	verbose  bool
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

func walk(opts options) error {
	return filepath.Walk(opts.src, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		rel, err := filepath.Rel(opts.src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(opts.dst, rel)
		existing, statErr := os.Stat(target)
		if statErr == nil && existing.Size() == info.Size() && !existing.ModTime().Before(info.ModTime()) {
			if opts.verbose {
				fmt.Printf("skip %s\n", rel)
			}
			return nil
		}
		if opts.dryRun {
			fmt.Printf("would copy %s (%d bytes)\n", rel, info.Size())
			return nil
		}
		if err := copyFile(path, target); err != nil {
			return err
		}
		fmt.Printf("copied %s\n", rel)
		return nil
	})
}

func main() {
	var opts options
	flag.StringVar(&opts.src, "src", ".", "source directory")
	flag.StringVar(&opts.dst, "dst", "", "destination directory")
	flag.BoolVar(&opts.dryRun, "dry-run", false, "print actions only")
	flag.BoolVar(&opts.verbose, "verbose", false, "print skipped files")
	flag.Parse()
	if opts.dst == "" {
		fmt.Fprintln(os.Stderr, "missing -dst")
		os.Exit(2)
	}
	if err := walk(opts); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
