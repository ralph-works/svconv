// Command svconv converts version strings between SemVer 2.0.0 and the
// four-number "quad" scheme used by Windows/.NET assembly versions.
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "-h" || os.Args[1] == "-help" || os.Args[1] == "--help") {
		usage(os.Stdout)
		return
	}
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

func usage(w io.Writer) {
	fmt.Fprintln(w, "usage: svconv [version ...]")
	fmt.Fprintln(w, "converts a semver string to quad form, or a quad string to semver")
	fmt.Fprintln(w, "with no arguments, or with \"-\", versions are read from stdin, one per line")
}

// run converts every version it is given and returns the process exit code.
// One bad version doesn't stop the batch: the failure is reported on stderr
// and the rest are still converted, so a single typo in a long list doesn't
// hide the other results. The exit code is 1 if anything failed.
func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	failed := false

	emit := func(input string) {
		out, err := convert(input)
		if err != nil {
			fmt.Fprintln(stderr, "svconv:", err)
			failed = true
			return
		}
		fmt.Fprintln(stdout, out)
	}

	if len(args) == 0 || (len(args) == 1 && args[0] == "-") {
		sc := bufio.NewScanner(stdin)
		for sc.Scan() {
			line := strings.TrimSpace(sc.Text())
			if line == "" {
				continue
			}
			emit(line)
		}
		if err := sc.Err(); err != nil {
			fmt.Fprintln(stderr, "svconv: reading stdin:", err)
			return 1
		}
	} else {
		for _, a := range args {
			emit(a)
		}
	}

	if failed {
		return 1
	}
	return 0
}

// convert detects which of the two formats the input is in and produces the
// other one. A strict four-part numeric string is treated as quad; anything
// else is parsed as semver.
func convert(input string) (string, error) {
	if q, err := ParseQuad(input); err == nil {
		return QuadToSemver(q).String(), nil
	}

	v, err := ParseSemver(input)
	if err != nil {
		return "", fmt.Errorf("%q is neither a valid quad version nor a valid semver", input)
	}

	q, err := SemverToQuad(v)
	if err != nil {
		return "", err
	}
	return q.String(), nil
}
