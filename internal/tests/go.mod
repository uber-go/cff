module go.uber.org/cff/internal/tests

go 1.19

require (
	github.com/gofrs/uuid v4.3.0+incompatible
	github.com/golang/mock v1.6.0
	github.com/stretchr/testify v1.8.0
	github.com/uber-go/tally v3.5.0+incompatible
	go.uber.org/atomic v1.10.0
	go.uber.org/cff v0.1.0
	go.uber.org/multierr v1.8.0
	go.uber.org/zap v1.23.0
	golang.org/x/exp v0.0.0-20221012211006-4de253d81b95
)

require (
	github.com/benbjohnson/clock v1.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/twmb/murmur3 v1.1.6 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace go.uber.org/cff => ../../
