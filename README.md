# SFS (Stupid File Storage)

## Send chunk request

Format:
```
*<filename_size><filename><id><size><data>
```

Where:
- `*` means "send"
- `filename_size` is a little-endian uint64
- `id` is a little-endian uint64
- `size` is a little-endian uint64
- `filename` is []byte with len of `filename_size`, containing the name
  of uploading file
- `data` is []byte with len of `size` containing the `id`'th chunk of file

## Receive file request

Request all chunks of the file that present in node.

Format:
```
/<filename_size><filename>
```
Where:
- `/` means "receive"
- `filename_size` is a little-endian uint64
- `filename` is []byte with len of `filename_size`, containing the name
  of uploading file

## Receive specific chunks request

Request specific chunks of the file.

Format:
```
%<filename_size><filename><count><...ids>
```
Where:
- `%` means "receive specific chunks"
- `filename_size` and `count` are little-endian uint64
- `...ids` is sequens (with len of `count`) containing little-endian uint64 ids of requesting chunks
- `filename` is []byte with len of `filename_size`, containing the name
  of uploading file
