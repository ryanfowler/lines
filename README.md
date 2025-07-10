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
  -o, --output <FORMAT>    Output format ("table" or "json") [default: table]
  -t, --timing             Show timing information
  -e, --exclude <EXCLUDE>  Exclude regex patterns (can be used multiple times)
  -h, --help               Print help
  -V, --version            Print version
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

### Exclude Patterns

You can exclude files and directories from the count using regex patterns with the `--exclude` flag:

```bash
# Exclude all test files (exact match)
lines --exclude "^test$"

# Exclude files containing "test" anywhere in the name
lines --exclude "test"

# Exclude multiple patterns
lines --exclude "target" --exclude "\.git" --exclude "node_modules"

# Exclude by file extension
lines --exclude "\.log$" --exclude "\.tmp$"

# Exclude files starting with "temp"
lines --exclude "^temp"

# Exclude hidden files (starting with .)
lines --exclude "^\."

# Complex pattern: exclude test files and directories
lines --exclude "(test|spec)" --exclude ".*\.test\."
```

## License

lines is released with the MIT license.
Please see the [LICENSE](./LICENSE) file for more details.
