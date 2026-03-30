# TODO

- [ ] Refactor application to use `//go:embed` for the `static/` and `templates/` directories. This will bundle these assets directly into the `main` executable, simplifying distribution (especially for Homebrew) by removing the need to ship these folders individually alongside the binary.
