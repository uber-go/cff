package builtincallexpr

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlow(t *testing.T) {
	buf := bytes.Buffer{}
	Flow("0", &buf)
	assert.Nil(t, buf.Bytes())
}

func TestFlowError(t *testing.T) {
	buf := bytes.Buffer{}
	Flow("notanumber", &buf)
	assert.NotNil(t, buf.Bytes())
}
