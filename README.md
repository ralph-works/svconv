# svconv

Windows file versions, .NET `AssemblyVersion`, and MSI product versions all
use a four-number scheme: `major.minor.patch.revision`, each field a plain
16-bit integer. There's no room in that format for a prerelease tag or
arbitrary build metadata. SemVer 2.0.0, which is what most tooling and
release pipelines actually produce (`1.4.0-rc.1+a1b2c3`), doesn't fit.

Every project that ships both a public SemVer and a Windows build ends up
writing this conversion by hand, usually badly (dropping the prerelease
silently, or crashing on the first version with a `+build` tag). `svconv`
does the conversion explicitly in both directions and refuses to guess when
information would be lost.

## Usage

```
$ go run . 1.4.0
1.4.0.0

$ go run . 1.4.0+7
1.4.0.7

$ go run . 2.0.0.7
2.0.0+7

$ go run . 2.0.0.0
2.0.0

$ go run . 1.4.0-rc.1
svconv: 1.4.0-rc.1 has a prerelease tag ("rc.1"); quad versions can't represent one
```

The tool detects the input format for you: a strict four-part numeric string
is read as a quad version, anything else is parsed as SemVer. Whichever one
it detects, it converts to the other.

## Conversion rules

**SemVer -> quad**: major, minor, and patch map straight across. Build
metadata that is a single plain integer (e.g. `+42`) becomes the revision
field, since that's the closest thing quad has to "extra identifying
number". A prerelease tag, or build metadata that isn't a bare integer, has
no quad equivalent and is a hard error rather than silent data loss.

**Quad -> SemVer**: always lossless. A revision of `0` is dropped (it means
"nothing to say here"); a nonzero revision becomes `+<revision>` build
metadata.

## Building

Standard library only, no dependencies to fetch:

```
go build -o svconv .
```

## Status

Early skeleton. The CLI takes exactly one positional argument; batch/stdin
input isn't wired up yet. See the roadmap in the project notes for what's
planned next.

## License

MIT, see [LICENSE](LICENSE).
