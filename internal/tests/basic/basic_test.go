package basic

import (
	"context"
	"errors"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSimpleFlow(t *testing.T) {
	msg, err := SimpleFlow()
	assert.NoError(t, err)
	assert.Equal(t, "hello world", msg)
}

func TestFlowWithoutParameters(t *testing.T) {
	r, err := NoParamsFlow(context.Background())
	require.NoError(t, err)

	body, err := ioutil.ReadAll(r)
	require.NoError(t, err)

	assert.Equal(t, "hello world", string(body))
}

func TestSerialFailures(t *testing.T) {
	t.Run("first function fails", func(t *testing.T) {
		err := SerialFailableFlow(
			func() error {
				return errors.New("great sadness")
			},
			func() error {
				t.Fatal("this function must not be called")
				return nil
			},
		)
		assert.Equal(t, errors.New("great sadness"), err)
	})

	t.Run("second function fails", func(t *testing.T) {
		err := SerialFailableFlow(
			func() error { return nil },
			func() error {
				return errors.New("failure")
			},
		)
		assert.Equal(t, errors.New("failure"), err)
	})
}

func TestProduceMultiple(t *testing.T) {
	require.NoError(t, ProduceMultiple())
}
