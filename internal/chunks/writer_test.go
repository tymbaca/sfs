package chunks

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tymbaca/sfs/internal/files"
)

func TestWriteChunksTo(t *testing.T) {
	t.Parallel()

	t.Run("normal", func(t *testing.T) {
		filename := "testdata/" + t.Name()

		chk0 := newChunk(0, "0----")
		chk1 := newChunk(1, "1----")
		chk2 := newChunk(2, "2----")
		expected := "0----1----2----"

		f, err := files.Create(filename)
		require.NoError(t, err)

		err = WriteTo(f, chk0, chk1, chk2)
		require.NoError(t, err)

		f.Close()

		// ASSert
		f, err = os.Open(filename)
		require.NoError(t, err)

		data, err := io.ReadAll(f)
		require.NoError(t, err)
		actual := string(data)

		require.Equal(t, expected, actual)
		require.NoError(t, err)
	})

	t.Run("smaller last chunk", func(t *testing.T) {
		filename := "testdata/" + t.Name()

		chk0 := newChunk(0, "0----")
		chk1 := newChunk(1, "1----")
		chk2 := newChunk(2, "2--")
		expected := "0----1----2--"

		f, err := files.Create(filename)
		require.NoError(t, err)

		err = WriteTo(f, chk0, chk1, chk2)
		require.NoError(t, err)

		f.Close()

		// ASSert
		f, err = os.Open(filename)
		require.NoError(t, err)

		data, err := io.ReadAll(f)
		require.NoError(t, err)
		actual := string(data)

		require.Equal(t, expected, actual)
		require.NoError(t, err)
	})

	t.Run("different sizes", func(t *testing.T) {
		filename := "testdata/" + t.Name()

		chk0 := newChunk(0, "0-")
		chk1 := newChunk(1, "1---")
		chk2 := newChunk(2, "2--")
		expected := "0-1---2--"

		f, err := files.Create(filename)
		require.NoError(t, err)

		err = WriteTo(f, chk0, chk1, chk2)
		require.NoError(t, err)

		f.Close()

		// ASSert
		f, err = os.Open(filename)
		require.NoError(t, err)

		data, err := io.ReadAll(f)
		require.NoError(t, err)
		actual := string(data)

		require.Equal(t, expected, actual)
		require.NoError(t, err)
	})
	t.Run("empty chunk", func(t *testing.T) {
		filename := "testdata/" + t.Name()

		chk0 := newChunk(0, "0----")
		chk1 := newChunk(1, "")
		chk2 := newChunk(2, "2----")
		expected := "0----2----"

		f, err := files.Create(filename)
		require.NoError(t, err)

		err = WriteTo(f, chk0, chk1, chk2)
		require.NoError(t, err)

		f.Close()

		// ASSert
		f, err = os.Open(filename)
		require.NoError(t, err)

		data, err := io.ReadAll(f)
		require.NoError(t, err)
		actual := string(data)

		require.Equal(t, expected, actual)
		require.NoError(t, err)
	})
	t.Run("all chunks are empty", func(t *testing.T) {
		filename := "testdata/" + t.Name()

		chk0 := newChunk(0, "")
		chk1 := newChunk(1, "")
		chk2 := newChunk(2, "")
		expected := ""

		f, err := files.Create(filename)
		require.NoError(t, err)

		err = WriteTo(f, chk0, chk1, chk2)
		require.NoError(t, err)

		f.Close()

		// ASSert
		f, err = os.Open(filename)
		require.NoError(t, err)

		data, err := io.ReadAll(f)
		require.NoError(t, err)
		actual := string(data)

		require.Equal(t, expected, actual)
		require.NoError(t, err)
	})
	t.Run("no chunks", func(t *testing.T) {
		filename := "testdata/" + t.Name()

		expected := ""

		f, err := files.Create(filename)
		require.NoError(t, err)

		err = WriteTo(f)
		require.NoError(t, err)

		f.Close()

		// ASSert
		f, err = os.Open(filename)
		require.NoError(t, err)

		data, err := io.ReadAll(f)
		require.NoError(t, err)
		actual := string(data)

		require.Equal(t, expected, actual)
		require.NoError(t, err)
	})
}

func newChunk(id uint64, data string) Chunk {
	return Chunk{
		ID:       id,
		Filename: "test",
		Size:     uint64(len(data)),
		Body:     strings.NewReader(data),
	}
}
