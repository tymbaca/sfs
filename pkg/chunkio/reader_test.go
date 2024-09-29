package chunkio

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReader(t *testing.T) {
	sr := strings.NewReader("1----2----3----")
	r1 := NewReader(sr, 0, 5)
	r2 := NewReader(sr, 5, 10)
	r3 := NewReader(sr, 10, 15)

	t.Run("1st chunk", func(t *testing.T) {
		buf := make([]byte, 3)
		n, err := r1.Read(buf)

		require.Equal(t, 3, n)
		require.Equal(t, buf[:n], []byte("1--"))
		require.NoError(t, err)

		n, err = r1.Read(buf)

		require.Equal(t, 2, n)
		require.Equal(t, buf[:n], []byte("--"))
		require.NoError(t, err)

		n, err = r1.Read(buf)

		require.Equal(t, 0, n)
		require.Equal(t, buf[:n], []byte(""))
		require.ErrorIs(t, err, io.EOF)
	})

	t.Run("2nd chunk", func(t *testing.T) {
		buf := make([]byte, 3)
		n, err := r2.Read(buf)

		require.Equal(t, 3, n)
		require.Equal(t, buf[:n], []byte("2--"))
		require.NoError(t, err)

		n, err = r2.Read(buf)

		require.Equal(t, 2, n)
		require.Equal(t, buf[:n], []byte("--"))
		require.NoError(t, err)

		n, err = r2.Read(buf)

		require.Equal(t, 0, n)
		require.Equal(t, buf[:n], []byte(""))
		require.ErrorIs(t, err, io.EOF)
	})

	t.Run("3rd chunk", func(t *testing.T) {
		buf := make([]byte, 3)
		n, err := r3.Read(buf)

		require.Equal(t, 3, n)
		require.Equal(t, buf[:n], []byte("3--"))
		require.NoError(t, err)

		n, err = r3.Read(buf)

		require.Equal(t, 2, n)
		require.Equal(t, buf[:n], []byte("--"))
		require.NoError(t, err)

		n, err = r3.Read(buf)

		require.Equal(t, 0, n)
		require.Equal(t, buf[:n], []byte(""))
		require.ErrorIs(t, err, io.EOF)
	})
}
