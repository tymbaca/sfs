# SFS (Stupid File Storage)

Chunk format:

```
$<filename_size><filename><id><size><data>
```
Where:
- `$` is just a dollar simbol as delimiter (why not?)
- `filename_size`, `id` and `size` are little-endian uint64
- `filename` is []byte with len of `filename_size`, containing the name
  of uploading file
- `data` is []byte with len of `size` containing the `id`'th chunk of file
