# lines

Counts lines of code, fast.

## Installation

Using `cargo`:

```sh
cargo install lines-cli
```

## Usage

```
$ lines -h
lines-cli 0.3.3
Counts lines of code, fast.

USAGE:
    lines [OPTIONS] [PATH]

ARGS:
    <PATH>    Directory or file to scan [default: .]

OPTIONS:
    -h, --help               Print help information
    -o, --output <FORMAT>    Output format ("table" or "json") [default: table]
    -t, --timing             Show timing information
    -V, --version            Print version information
```

Using `lines` in this repo outputs:

```
 Language | Files | Lines 
----------+-------+-------
 Rust     |     4 |   611 
 Markdown |     1 |    49 
 TOML     |     1 |    36 
---------- ------- -------
 Total    |     6 |   696 
```

## License

lines is released with the MIT license.
Please see the [LICENSE](./LICENSE) file for more details.
