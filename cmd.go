// Copyright 2016 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Command git-checkout-pr fetches a pull request from GitHub and checks it out into a local branch.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

var (
	verbose = flag.Bool("v", false, "verbose")
	branch  = flag.String("branch", "", "branch name. default is `pullN`, where N is the pull ID.")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage of git-checkout-pr:
	git-checkout-pr [url-to-pull-request]

Flags:
`)
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	u, err := url.Parse(flag.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not parse URL %q: %v\n", flag.Arg(0), err)
		os.Exit(1)
	}

	parts := strings.Split(u.EscapedPath(), "/")
	if len(parts) < 5 || parts[3] != "pull" || parts[4] == "" {
		fmt.Fprintf(os.Stderr, "URL malformed: %q must have a path of `/{owner}/{repo}/pull/{pullid}`\n", flag.Arg(0))
		os.Exit(1)
	}

	org := parts[1]
	repo := parts[2]
	id := parts[4]

	remote := fmt.Sprintf("https://github.com/%s/%s.git", org, repo)

	if *branch == "" {
		*branch = "pull" + id
	}

	if *verbose {
		log.Print("Running:", "git", "fetch", remote, "pull/"+id+"/head:"+*branch)
	}
	cmd := exec.Command("git", "fetch", remote, "pull/"+id+"/head:"+*branch)
	out, err := cmd.CombinedOutput()
	if err != nil {
		os.Stderr.Write(out)
		os.Exit(1)
	}

	if *verbose {
		log.Print("Running:", "git", "checkout", *branch)
	}
	cmd = exec.Command("git", "checkout", *branch)
	out, err = cmd.CombinedOutput()
	if err != nil {
		os.Stderr.Write(out)
		os.Exit(1)
	}

	if *verbose {
		log.Print("Done")
	}
}
