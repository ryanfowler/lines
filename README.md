# lines

Counts lines of code

## Usage

```
$ lines -h
lines 0.1.0
Count lines of code.

USAGE:
    lines [FLAGS] [OPTIONS] <path>

FLAGS:
    -h, --help       Prints help information
    -t, --timing     Show timing
    -V, --version    Prints version information

OPTIONS:
    -o, --output <format>    Output format [default: table]

ARGS:
    <path>    Directory or file to scan
```

Using `lines` in this repo outputs:

```
 Language | Files | Lines 
----------+-------+-------
 Rust     |     4 |   526 
 Markdown |     1 |    43 
 TOML     |     1 |    24 
---------- ------- -------
 Total    |     6 |   593 
```

## License

lines is released with the MIT license.
Please see the LICENSE file for more details.
