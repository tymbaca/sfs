package chunkio

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriter(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		f, err := os.Create("testdata/writer-test.txt")
		require.NoError(t, err)

		w := NewWriter(f, 5, 10)

		n, err := w.Write([]byte("hello"))
		require.Equal(t, 5, n)
		require.NoError(t, err)

		// assert file
		result, err := io.ReadAll(f)
		require.NoError(t, err)
		require.Equal(t, []byte(string([]byte{0, 0, 0, 0, 0})+"hello"), result)
	})

	t.Run("ok multiple writes", func(t *testing.T) {
		f, err := os.Create("testdata/writer-test.txt")
		require.NoError(t, err)

		w := NewWriter(f, 5, 10)

		n, err := w.Write([]byte("hel"))
		require.Equal(t, 3, n)
		require.NoError(t, err)

		n, err = w.Write([]byte("lo"))
		require.Equal(t, 2, n)
		require.NoError(t, err)

		// assert file
		result, err := io.ReadAll(f)
		require.NoError(t, err)
		require.Equal(t, []byte(string([]byte{0, 0, 0, 0, 0})+"hello"), result)
	})

	t.Run("jump over limit", func(t *testing.T) {
		f, err := os.Create("testdata/writer-test.txt")
		require.NoError(t, err)

		w := NewWriter(f, 5, 10)

		n, err := w.Write([]byte("helloo"))
		require.Equal(t, 5, n)
		require.ErrorIs(t, err, io.ErrShortWrite)

		// assert file
		result, err := io.ReadAll(f)
		require.NoError(t, err)
		require.Equal(t, []byte(string([]byte{0, 0, 0, 0, 0})+"hello"), result)
	})

	t.Run("jump over limit", func(t *testing.T) {
		f, err := os.Create("testdata/writer-test.txt")
		require.NoError(t, err)

		w := NewWriter(f, 5, 10)

		n, err := w.Write([]byte("hel"))
		require.Equal(t, 3, n)
		require.NoError(t, err)

		n, err = w.Write([]byte("loo"))
		require.Equal(t, 2, n)
		require.ErrorIs(t, err, io.ErrShortWrite)

		// assert file
		result, err := io.ReadAll(f)
		require.NoError(t, err)
		require.Equal(t, []byte(string([]byte{0, 0, 0, 0, 0})+"hello"), result)
	})
}
