// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"go/build"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	flagAll    = flag.Bool("all", false, "build for all GOOS/GOARCH platforms")
	flagRace   = flag.Bool("race", false, "build with -race")
	flagRemake = flag.Bool("remake", false, "build new toolchain with make.bash instead of go install std cmd")
	flagWork   = flag.Bool("work", false, "build with -work")
)

func usage() {
	fmt.Fprintln(os.Stderr, "usage: toolstash-check [options] [commit]")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if *flagAll {
		if *flagRace {
			log.Fatal("-all and -race are incompatible")
			os.Exit(2)
		}
		if *flagWork {
			log.Fatal("-all and -work are incompatible")
			os.Exit(2)
		}
	}

	spec := "HEAD"
	switch flag.NArg() {
	case 0:
	case 1:
		spec = flag.Arg(0)
	default:
		flag.Usage()
		os.Exit(2)
	}

	goroot := runtime.GOROOT()
	commit, err := revParse(goroot, spec)
	must(err)

	pkg, err := build.Import("golang.org/x/tools/cmd/toolstash", "", build.FindOnly)
	must(err)

	tmpdir, err := ioutil.TempDir("", "toolstash-check-")
	must(err)
	defer os.RemoveAll(tmpdir)

	tmproot := filepath.Join(tmpdir, "go")
	must(command("git", "clone", goroot, tmproot).Run())

	cmd := command("git", "checkout", commit+"^")
	cmd.Dir = tmproot
	must(cmd.Run())

	must(ioutil.WriteFile(filepath.Join(tmproot, "VERSION"), []byte("devel"), 0666))

	cmd = command("./make.bash")
	cmd.Dir = filepath.Join(tmproot, "src")
	must(cmd.Run())

	envPath := os.Getenv("PATH")
	if envPath != "" {
		envPath = string(os.PathListSeparator) + envPath
	}
	must(os.Setenv("PATH", filepath.Join(tmproot, "bin")+envPath))
	must(os.Setenv("GOROOT", tmproot))

	must(command("toolstash", "save").Run())

	cmd = command("git", "checkout", commit)
	cmd.Dir = tmproot
	must(cmd.Run())

	if *flagRemake {
		cmd = command("./make.bash")
		cmd.Dir = filepath.Join(tmproot, "src")
	} else {
		cmd = command("go", "install", "std", "cmd")
	}
	must(cmd.Run())

	if *flagAll {
		must(command(filepath.Join(pkg.Dir, "buildall")).Run())
	} else {
		buildArgs := []string{"build", "-a"}
		if *flagRace {
			buildArgs = append(buildArgs, "-race")
		}
		if *flagWork {
			buildArgs = append(buildArgs, "-work")
		}
		buildArgs = append(buildArgs, "-toolexec", "toolstash -cmp", "std", "cmd")
		must(command("go", buildArgs...).Run())
	}

	fmt.Println("toolstash-check passed for " + commit)
}

// revParse runs "git rev-parse $spec" in $GOROOT to parse a Git
// revision specifier.
func revParse(dir, spec string) (string, error) {
	cmd := exec.Command("git", "rev-parse", spec)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(string(out), "\n"), nil
}

func command(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
