# dirsync

A small Go utility that copies files from a source directory to a destination when size or modification time differs.

## Usage

```bash
go run . -src ./from -dst ./to -dry-run
```

## Flags

- `-src` / `-dst` — directories
- `-dry-run` — only print actions
- `-verbose` — print every file
## License

MIT

## Notes

Contributions and small patches are welcome.
