package chunks

import (
	"fmt"
	"io"
	"slices"

	"github.com/tymbaca/sfs/pkg/chunkio"
	"golang.org/x/sync/errgroup"
)

// If w can't grow automatically, user must make sure to pregrow it before calling
// to match the size of all chunks Size summed up
func WriteTo(w io.WriterAt, chunks ...Chunk) error {
	slices.SortFunc(chunks, func(a, b Chunk) int {
		if a.ID > b.ID {
			return 1
		}
		if a.ID < b.ID {
			return -1
		}
		return 0
	})

	writers := make([]io.Writer, len(chunks))
	offset := int64(0)
	// Splitting writer to windows, each corresponding to his chunk
	for i, chk := range chunks {
		writers[i] = chunkio.NewWriter(w, offset, offset+int64(chk.Size))
		offset += int64(chk.Size)
	}

	// Writing each chunk to his window writer
	var wg errgroup.Group
	for i := range chunks {
		i := i
		wg.Go(func() error {
			ww := writers[i]
			chk := chunks[i]
			_, err := io.Copy(ww, chk.Body)
			return err
		})
	}

	if err := wg.Wait(); err != nil {
		return fmt.Errorf("can't write chunks: %w", err)
	}

	return nil
}
