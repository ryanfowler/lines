# lines
Counts lines of code

## Usage

Run lines with an argument of the directory you'd like counted:

```
$  ./lines-OS-ARCH /path/to/directory
```

will print something like:

```
+-------------+-------+---------+--------+----------+--------+---------+
| Language    | Files | Code    | Mixed  | Comments | Empty  | Total   |
+-------------+-------+---------+--------+----------+--------+---------+
| Go          | 2,320 | 375,812 | 17,395 |   62,496 | 54,571 | 510,274 |
| Javascript  |    95 |  90,663 |     38 |    2,978 |  3,268 |  96,947 |
| C           |    20 |  17,489 |     50 |      265 |  1,019 |  18,823 |
| HTML        |    87 |  11,925 |     37 |      410 |  2,100 |  14,472 |
| XML         |    14 |   5,365 |      1 |       12 |    242 |   5,620 |
| Python      |     2 |   2,501 |     44 |      423 |    392 |   3,360 |
| C++         |     2 |   2,514 |     11 |      176 |    373 |   3,074 |
| CSS         |     9 |   1,514 |      1 |       25 |    196 |   1,736 |
| Objective-C |     2 |     248 |      0 |        0 |     36 |     284 |
| Ruby        |     1 |      32 |     14 |        2 |     13 |      61 |
| Rust        |     1 |      43 |      0 |        0 |      1 |      44 |
+-------------+-------+---------+--------+----------+--------+---------+
| Totals:     | 2,553 | 508,106 | 17,591 |   66,787 | 62,211 | 654,695 |
+-------------+-------+---------+--------+----------+--------+---------+
Time: 505.31865ms
```

* Files: # of individual files
* Code: lines of code only
* Mixed: lines with code and comment(s)
* Comments: lines of comments only
* Empty: all empty lines
* Total: total # of lines

## Flags

`-filter=REGEXP`
* filter all file and directory names with the provided regex

`-exclude=REGEXP`
* exclude all file and directory names with the provided regex

`-filterDir=REGEXP`
* filter all directory names with the provided regex

`-excludeDir=REGEXP`
* exclude all directory names with the provided regex

## TODO

* Add more tests!
* Set directory to the current directory if none provided

## License

lines is released with the MIT license.
Please see the LICENSE file for more details.
