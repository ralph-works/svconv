package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Semver holds the parsed pieces of a SemVer 2.0.0 version string.
type Semver struct {
	Major, Minor, Patch uint64
	Prerelease          string
	Build               string
}

// semverPattern is a trimmed-down version of the official SemVer 2.0.0 regex.
// It is stricter than most "loose" parsers on purpose: this tool is meant to
// round-trip versions, not guess at malformed ones.
var semverPattern = regexp.MustCompile(
	`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)` +
		`(?:-([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?` +
		`(?:\+([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?$`,
)

// ParseSemver parses a version string into its numeric core plus the
// optional prerelease and build metadata fields.
func ParseSemver(s string) (Semver, error) {
	m := semverPattern.FindStringSubmatch(strings.TrimSpace(s))
	if m == nil {
		return Semver{}, fmt.Errorf("%q is not a valid semver string", s)
	}

	major, err := strconv.ParseUint(m[1], 10, 64)
	if err != nil {
		return Semver{}, fmt.Errorf("major version out of range: %w", err)
	}
	minor, err := strconv.ParseUint(m[2], 10, 64)
	if err != nil {
		return Semver{}, fmt.Errorf("minor version out of range: %w", err)
	}
	patch, err := strconv.ParseUint(m[3], 10, 64)
	if err != nil {
		return Semver{}, fmt.Errorf("patch version out of range: %w", err)
	}

	return Semver{
		Major:      major,
		Minor:      minor,
		Patch:      patch,
		Prerelease: m[4],
		Build:      m[5],
	}, nil
}

// String renders the version back into standard SemVer form.
func (v Semver) String() string {
	s := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.Prerelease != "" {
		s += "-" + v.Prerelease
	}
	if v.Build != "" {
		s += "+" + v.Build
	}
	return s
}
