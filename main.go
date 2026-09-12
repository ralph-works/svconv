// Command svconv converts version strings between SemVer 2.0.0 and the
// four-number "quad" scheme used by Windows/.NET assembly versions.
package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <version>\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "converts a semver string to quad form, or a quad string to semver")
		os.Exit(2)
	}

	out, err := convert(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "svconv:", err)
		os.Exit(1)
	}
	fmt.Println(out)
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
