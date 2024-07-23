# lines

Counts lines of code, fast.

## Installation

Using `cargo`:

```sh
cargo install lines-cli
```

## Usage

```
❯ lines -h
Counts lines of code, fast.

Usage: lines [OPTIONS] [PATH]

Arguments:
  [PATH]  Directory or file to scan [default: .]

Options:
  -o, --output <FORMAT>  Output format ("table" or "json") [default: table]
  -t, --timing           Show timing information
  -h, --help             Print help
  -V, --version          Print version
```

Using `lines` in this repo outputs:

```
 Language | Files | Lines
----------+-------+-------
 Rust     |     4 |   600
 Markdown |     1 |    47
 TOML     |     1 |    36
---------- ------- -------
 Total    |     6 |   683
```

## License

lines is released with the MIT license.
Please see the [LICENSE](./LICENSE) file for more details.
