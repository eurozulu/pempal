module github.com/eurozulu/pempal

go 1.24.0

toolchain go1.24.2

require (
	github.com/eurozulu/colout v0.0.2
	github.com/eurozulu/spud v0.6.4
	gopkg.in/yaml.v2 v2.4.0
	golang.org/x/text v0.22.0 // indirect
)

replace (
	github.com/eurozulu/spud v0.6.4 => /Users/rpgi/GolandProjects/spud
)