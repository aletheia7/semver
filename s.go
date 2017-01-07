// Copyright 2016 aletheia7. All rights reserved. Use of this source code is
// governed by a BSD-2-Clause license that can be found in the LICENSE file.

// semver.org version strings only allow [0-9A-Za-z-]. This package allows
// unicode letters in place of A-Z and a-z; i.e. 3.24.3-β+20150115102400 is
// acceptable. The allowance of unicode characters makes this package
// noncompliant with semver.org.

// Package semver compares semver.org version strings.

package semver

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Version represents a parsed version. See http://semver.org/ for
// detailed description of the various components.
type Version struct {
	Major      int      // The major version number.
	Minor      int      // The minor version number.
	Patch      int      // The patch version number.
	Prerelease []string // The pre-release version (dot-separated elements)
	Build      []string // The build version (dot-separated elements)
}

var charClasses = strings.NewReplacer("d", `[\pNd]`, "c", `[\-\pNd\pL]`)

const pattern = `^(d{1,9})\.(d{1,9})\.(d{1,9})(-c+(\.c+)*)?(\+c+(\.c+)*)?$`

var versionPat = regexp.MustCompile(charClasses.Replace(pattern))

// Parse parses the version, which is of one of the following forms:
//     1.2.3
//     1.2.3-prerelease
//     1.2.3+build
//     1.2.3-prerelease+build
func Parse(s string) (*Version, error) {
	m := versionPat.FindStringSubmatch(s)
	if m == nil {
		return nil, fmt.Errorf("invalid version %q", s)
	}
	v := new(Version)
	v.Major = atoi(m[1])
	v.Minor = atoi(m[2])
	v.Patch = atoi(m[3])
	if m[4] != "" {
		v.Prerelease = strings.Split(m[4][1:], ".")
	}
	if m[6] != "" {
		v.Build = strings.Split(m[6][1:], ".")
	}
	return v, nil
}

// atoi is the same as strconv.Atoi but assumes that
// the string has been verified to be a valid integer.
func atoi(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return n
}

func (v Version) String() string {
	var pre, build string
	if v.Prerelease != nil {
		pre = "-" + strings.Join(v.Prerelease, ".")
	}
	if v.Build != nil {
		build = "+" + strings.Join(v.Build, ".")
	}
	return fmt.Sprintf("%d.%d.%d%s%s", v.Major, v.Minor, v.Patch, pre, build)
}

func allDigits(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// lessIds returns whether the slice of identifiers a is less than b,
// as specified in semver.org,
func lessIds(a, b []string) (v bool) {
	i := 0
	for ; i < len(a) && i < len(b); i++ {
		if c := cmp(a[i], b[i]); c != 0 {
			return c < 0
		}
	}
	return i < len(b)
}

// eqIds returns whether the slice of identifiers a is equal to b,
// as specified in semver.org,
func eqIds(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, s := range a {
		if cmp(s, b[i]) != 0 {
			return false
		}
	}
	return true
}

// cmp implements comparison of identifiers as specified at semver.org.
// It returns 1, -1 or 0 if a is greater-than, less-than or equal to b,
// respectively.
//
// Identifiers consisting of only digits are compared numerically and
// identifiers with letters or dashes are compared lexically in ASCII
// sort order.  Numeric identifiers always have lower precedence than
// non-numeric identifiers.
func cmp(a, b string) int {
	numa, numb := allDigits(a), allDigits(b)
	switch {
	case numa && numb:
		return numCmp(a, b)
	case numa:
		return -1
	case numb:
		return 1
	case a < b:
		return -1
	case a > b:
		return 1
	}
	return 0
}

// numCmp 1, -1 or 0 depending on whether the known-to-be-all-digits
// strings a and b are numerically greater than, less than or equal to
// each other.  Avoiding the conversion means we can work correctly with
// very long version numbers.
func numCmp(a, b string) int {
	a = strings.TrimLeft(a, "0")
	b = strings.TrimLeft(b, "0")
	switch {
	case len(a) < len(b):
		return -1
	case len(a) > len(b):
		return 1
	case a < b:
		return -1
	case a > b:
		return 1
	}
	return 0
}

// Less returns whether v is semantically earlier in the
// version sequence than w.
func (v *Version) Less(w *Version) bool {
	switch {
	case v.Major != w.Major:
		return v.Major < w.Major
	case v.Minor != w.Minor:
		return v.Minor < w.Minor
	case v.Patch != w.Patch:
		return v.Patch < w.Patch
	case !eqIds(v.Prerelease, w.Prerelease):
		if v.Prerelease == nil || w.Prerelease == nil {
			return v.Prerelease != nil
		}
		return lessIds(v.Prerelease, w.Prerelease)
	case !eqIds(v.Build, w.Build):
		return lessIds(v.Build, w.Build)
	}
	return false
}

// Equal returns whether v is semantically equal with w
func (v *Version) Equal(w *Version) bool {
	if v.Major == w.Major &&
		v.Minor == w.Minor &&
		v.Patch == w.Patch {
		if len(v.Prerelease) == len(w.Prerelease) {
			for i := range v.Prerelease {
				if v.Prerelease[i] != w.Prerelease[i] {
					goto not
				}
			}
			return true
		}
	not:
	}
	return false
}
