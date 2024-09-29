package chunkio

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSplit(t *testing.T) {
	t.Run("normal order", func(t *testing.T) {
		sr := strings.NewReader("1----2----3----")
		t.Parallel()

		rs := Split(sr, sr.Size(), 5)

		require.Len(t, rs, 3)

		r1, r2, r3 := rs[0], rs[1], rs[2]

		t.Run("1st chunk", func(t *testing.T) {
			buf := make([]byte, 3)
			n, err := r1.Read(buf)

			require.Equal(t, 3, n)
			require.Equal(t, []byte("1--"), buf[:n])
			require.NoError(t, err)

			n, err = r1.Read(buf)

			require.Equal(t, 2, n)
			require.Equal(t, []byte("--"), buf[:n])
			require.NoError(t, err)

			n, err = r1.Read(buf)

			require.Equal(t, 0, n)
			require.Equal(t, []byte(""), buf[:n])
			require.ErrorIs(t, err, io.EOF)
		})

		t.Run("2nd chunk", func(t *testing.T) {
			buf := make([]byte, 3)
			n, err := r2.Read(buf)

			require.Equal(t, 3, n)
			require.Equal(t, []byte("2--"), buf[:n])
			require.NoError(t, err)

			n, err = r2.Read(buf)

			require.Equal(t, 2, n)
			require.Equal(t, []byte("--"), buf[:n])
			require.NoError(t, err)

			n, err = r2.Read(buf)

			require.Equal(t, 0, n)
			require.Equal(t, []byte(""), buf[:n])
			require.ErrorIs(t, err, io.EOF)
		})

		t.Run("3rd chunk", func(t *testing.T) {
			buf := make([]byte, 3)
			n, err := r3.Read(buf)

			require.Equal(t, 3, n)
			require.Equal(t, []byte("3--"), buf[:n])
			require.NoError(t, err)

			n, err = r3.Read(buf)

			require.Equal(t, 2, n)
			require.Equal(t, []byte("--"), buf[:n])
			require.NoError(t, err)

			n, err = r3.Read(buf)

			require.Equal(t, 0, n)
			require.Equal(t, []byte(""), buf[:n])
			require.ErrorIs(t, err, io.EOF)
		})
	})
	t.Run("last chunk not full", func(t *testing.T) {
		sr := strings.NewReader("1----2----3--")
		t.Parallel()

		rs := Split(sr, sr.Size(), 5)

		require.Len(t, rs, 3)

		_, _, r3 := rs[0], rs[1], rs[2]

		buf := make([]byte, 100)
		n, err := r3.Read(buf)
		require.NoError(t, err)
		require.Equal(t, 3, int(n)) // last readers must provide 3 bytes
		require.Equal(t, []byte("3--"), buf[:n])
	})
}

func TestReader(t *testing.T) {
	t.Run("normal order", func(t *testing.T) {
		sr := strings.NewReader("1----2----3----")
		t.Parallel()

		r1 := NewReader(sr, 0, 5)
		r2 := NewReader(sr, 5, 10)
		r3 := NewReader(sr, 10, 15)

		t.Run("1st chunk", func(t *testing.T) {
			buf := make([]byte, 3)
			n, err := r1.Read(buf)

			require.Equal(t, 3, n)
			require.Equal(t, []byte("1--"), buf[:n])
			require.NoError(t, err)

			n, err = r1.Read(buf)

			require.Equal(t, 2, n)
			require.Equal(t, []byte("--"), buf[:n])
			require.NoError(t, err)

			n, err = r1.Read(buf)

			require.Equal(t, 0, n)
			require.Equal(t, []byte(""), buf[:n])
			require.ErrorIs(t, err, io.EOF)
		})

		t.Run("2nd chunk", func(t *testing.T) {
			buf := make([]byte, 3)
			n, err := r2.Read(buf)

			require.Equal(t, 3, n)
			require.Equal(t, []byte("2--"), buf[:n])
			require.NoError(t, err)

			n, err = r2.Read(buf)

			require.Equal(t, 2, n)
			require.Equal(t, []byte("--"), buf[:n])
			require.NoError(t, err)

			n, err = r2.Read(buf)

			require.Equal(t, 0, n)
			require.Equal(t, []byte(""), buf[:n])
			require.ErrorIs(t, err, io.EOF)
		})

		t.Run("3rd chunk", func(t *testing.T) {
			buf := make([]byte, 3)
			n, err := r3.Read(buf)

			require.Equal(t, 3, n)
			require.Equal(t, []byte("3--"), buf[:n])
			require.NoError(t, err)

			n, err = r3.Read(buf)

			require.Equal(t, 2, n)
			require.Equal(t, []byte("--"), buf[:n])
			require.NoError(t, err)

			n, err = r3.Read(buf)

			require.Equal(t, 0, n)
			require.Equal(t, []byte(""), buf[:n])
			require.ErrorIs(t, err, io.EOF)
		})
	})

	t.Run("concurrent order", func(t *testing.T) {
		sr := strings.NewReader("1----2----3----")
		t.Parallel()

		r1 := NewReader(sr, 0, 5)
		r2 := NewReader(sr, 5, 10)
		r3 := NewReader(sr, 10, 15)

		buf1 := make([]byte, 3)
		buf2 := make([]byte, 3)
		buf3 := make([]byte, 3)

		n, err := r1.Read(buf1)
		require.Equal(t, 3, n)
		require.Equal(t, []byte("1--"), buf1[:n])
		require.NoError(t, err)

		n, err = r2.Read(buf2)
		require.Equal(t, 3, n)
		require.Equal(t, []byte("2--"), buf2[:n])
		require.NoError(t, err)

		n, err = r3.Read(buf3)
		require.Equal(t, 3, n)
		require.Equal(t, []byte("3--"), buf3[:n])
		require.NoError(t, err)

		n, err = r1.Read(buf1)
		require.Equal(t, 2, n)
		require.Equal(t, []byte("--"), buf1[:n])
		require.NoError(t, err)

		n, err = r2.Read(buf2)
		require.Equal(t, 2, n)
		require.Equal(t, []byte("--"), buf2[:n])
		require.NoError(t, err)

		n, err = r3.Read(buf3)
		require.Equal(t, 2, n)
		require.Equal(t, []byte("--"), buf3[:n])
		require.NoError(t, err)

		n, err = r1.Read(buf1)
		require.Equal(t, 0, n)
		require.Equal(t, []byte(""), buf1[:n])
		require.ErrorIs(t, err, io.EOF)

		n, err = r2.Read(buf2)
		require.Equal(t, 0, n)
		require.Equal(t, []byte(""), buf2[:n])
		require.ErrorIs(t, err, io.EOF)

		n, err = r3.Read(buf3)
		require.Equal(t, 0, n)
		require.Equal(t, []byte(""), buf3[:n])
		require.ErrorIs(t, err, io.EOF)
	})
}
