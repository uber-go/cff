package basic

import (
	"context"
	"errors"
	"io"
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

	body, err := io.ReadAll(r)
	require.NoError(t, err)

	assert.Equal(t, "hello world", string(body))
}

func TestSerialFailures(t *testing.T) {
	t.Run("first function fails", func(t *testing.T) {
		err := SerialFailableFlow(
			context.Background(),
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
			context.Background(),
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

func TestContextCancelation(t *testing.T) {
	dontCallMe := func(t *testing.T) func() error {
		return func() error {
			t.Fatal("this function must not be called")
			return nil
		}
	}

	t.Run("cancel before first", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		cancel()
		err := SerialFailableFlow(ctx, dontCallMe(t), dontCallMe(t))
		require.Error(t, err)
		assert.Equal(t, ctx.Err(), err)
	})

	t.Run("cancel before second", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		err := SerialFailableFlow(ctx,
			func() error {
				require.NoError(t, ctx.Err(), "context can't be done yet")
				cancel()
				return nil
			},
			dontCallMe(t),
		)
		require.Error(t, err)
		assert.Equal(t, ctx.Err(), err)
	})

	t.Run("cancel before third", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		err := SerialFailableFlow(ctx,
			func() error {
				require.NoError(t, ctx.Err(), "context can't be done yet")
				return nil
			},
			func() error {
				require.NoError(t, ctx.Err(), "context can't be done yet")
				cancel()
				return nil
			},
		)
		require.Error(t, err)
		assert.Equal(t, ctx.Err(), err)
	})
}
