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
lines 0.3.1
Count lines of code

USAGE:
    lines [OPTIONS] [PATH]

ARGS:
    <PATH>    Directory or file to scan [default: .]

OPTIONS:
    -h, --help               Print help information
    -o, --output <FORMAT>    Output format [default: table]
    -t, --timing             Show timing information
    -V, --version            Print version information
```

Using `lines` in this repo outputs:

```
 Language | Files | Lines 
----------+-------+-------
 Rust     |     4 |   611 
 Markdown |     1 |    49 
 TOML     |     1 |    35 
---------- ------- -------
 Total    |     6 |   695 
```

## License

lines is released with the MIT license.
Please see the [LICENSE](./LICENSE) file for more details.
