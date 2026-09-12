package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Quad is the four-number version scheme used by Windows file/product
// versions and .NET AssemblyVersion attributes: major.minor.patch.revision,
// each field a plain uint16 with no room for prerelease or build tags.
type Quad struct {
	Major, Minor, Patch, Revision uint16
}

// ParseQuad parses a strict "a.b.c.d" numeric version, the form Windows
// resource files and MSBuild expect.
func ParseQuad(s string) (Quad, error) {
	parts := strings.Split(strings.TrimSpace(s), ".")
	if len(parts) != 4 {
		return Quad{}, fmt.Errorf("%q does not have four dot-separated parts", s)
	}

	nums := make([]uint16, 4)
	for i, p := range parts {
		n, err := strconv.ParseUint(p, 10, 16)
		if err != nil {
			return Quad{}, fmt.Errorf("part %d (%q) is not a valid uint16: %w", i+1, p, err)
		}
		nums[i] = uint16(n)
	}

	return Quad{Major: nums[0], Minor: nums[1], Patch: nums[2], Revision: nums[3]}, nil
}

// String renders the quad back into dotted form.
func (q Quad) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", q.Major, q.Minor, q.Patch, q.Revision)
}

// SemverToQuad converts a semver into a quad version.
//
// Quad has no concept of prerelease or arbitrary build metadata, so this
// only succeeds for versions that fit the format losslessly: no prerelease,
// and build metadata that is either absent or a single plain integer (which
// becomes the revision field). Anything else is rejected rather than
// silently discarding information.
func SemverToQuad(v Semver) (Quad, error) {
	if v.Major > 0xFFFF || v.Minor > 0xFFFF || v.Patch > 0xFFFF {
		return Quad{}, fmt.Errorf("%s has a component too large for a uint16 quad version", v)
	}
	if v.Prerelease != "" {
		return Quad{}, fmt.Errorf("%s has a prerelease tag (%q); quad versions can't represent one", v, v.Prerelease)
	}

	var revision uint64
	if v.Build != "" {
		n, err := strconv.ParseUint(v.Build, 10, 16)
		if err != nil {
			return Quad{}, fmt.Errorf("%s has build metadata %q; only a plain integer can map to a quad revision", v, v.Build)
		}
		revision = n
	}

	return Quad{
		Major:    uint16(v.Major),
		Minor:    uint16(v.Minor),
		Patch:    uint16(v.Patch),
		Revision: uint16(revision),
	}, nil
}

// QuadToSemver converts a quad version into a semver. This direction is
// always lossless: the revision field becomes build metadata when nonzero,
// and is dropped when zero since "+0" and "no build metadata" both mean
// "nothing to say here".
func QuadToSemver(q Quad) Semver {
	v := Semver{Major: uint64(q.Major), Minor: uint64(q.Minor), Patch: uint64(q.Patch)}
	if q.Revision != 0 {
		v.Build = strconv.FormatUint(uint64(q.Revision), 10)
	}
	return v
}
