
# PrintColumns Documentation

## NAME

`pc` - PrintColumns - tidy and filter input and generate formatted output

## SYNOPSIS

`pc [options] [column-numbers]`

**Options:**

`-file=input-filename`  Read the text from this file. If there is also data from STDIN, this is added together.<br>
`-header='col1-header col2-header ...'` Define header line if the text has no headers.  Headers are left-adjusted; if they start with a dash (`-`), they are right-adjusted.<br>
`-sep=' '` Define the separator used to split the data (default is a space).<br>
`-mb` MoreBlanks: Assumes columns are separated by more than one blank (default: false). Useful when input is a preformatted ASCII table.<br>
`-w=1` Number of blanks between column separator and column content (default is 1).<br>
`-colsep='|'` Define the character to separate columns (default is '|').<br>
`-filter='regex'` Process only lines where the string or regex matches.<br>
`-sortcol=colnum:` Sort by the specified column number (output column number). Only one column can be defined for sorting.<br>
`-gcol=colnum:` Group lines when the value in this column differs from the previous line.  In the grouped column, the second and following lines of a group get the value `""`. This behavior can <br>be disabled with `-gcolval`.
`-gcolval` Do not replace values in the group column with `""`.<br>
`-nf` No format: Don't format columns for common column width.<br>
`-nn` No numerical: Don't format numerical content right-adjusted.<br>
`-nhl` No headline: The data contains no headline.<br>
`-ts` TitleSeparator: Draws a separator line between the first and second lines of output.<br>
`-fs` FooterSeparator: Draws a separator line between the last and second-to-last lines of output.<br>
`-cs` ColumnSeparator: Draws a separator (default '|') between columns of output.<br>
`-pp` PrettyPrint: Draw cell borders and all separators.<br>
`-rh` RemoveHeader: Removes the first line.<br>
`-num` Numbering: Insert column numbers in the first line.<br>
`-csv` CSV: Write output in CSV format.<br>
`-json` JSON: Write output in JSON format.<br>
`-jtc` TitleColumn: Relevant for JSON with defined headers; use the first column as the main key and put all other columns as sub-keys.<br>
`-help` Print help and exit.<br>
`-man` Print help and manual, then exit.<br>
`-v` Verify: Print parameter verify info.<br>

**Column-numbers:** `[1 2 n:m]`, a range of columns can be given with `n:m`.  All options must be set before the column numbers.

## DESCRIPTION

This command formats text columns from stdin or a file and prints them as an ASCII table, CSV, or JSON.  The columns for output can be selected by single numbers separated by spaces or one or more ranges separated by colons. The input should have the same number of columns in each line. The formatted result is printed to stdout.  Without column-number parameters, `pc` prints all columns from the input in formatted form.

## INSTALLATION

To install this tool, follow these steps:

1.  Clone the repository: `git clone <repository_url>`
2.  Navigate to the project directory: `cd print-columns`
3.  Build the project: `go build`

## USAGE EXAMPLES

**Example 1: Formatting a data file**

```sh
pc --file=test/data/data2.txt 
# or
cat test/data/data2.txt | pc
```

**Example 2: Formatting with a header**

```sh
pc -file=test/data/data2.txt -header='col-1 col-2 col-3 col-4 col-5 col-6'
# or
cat test/data/data2.txt | pc --header='col-1 col-2 col-3 -col-4 col-5 -col-6'
```

**Example 3: Formatting with a header and column separator**

```sh
pc --file=test/data/data2.txt --header='col-1 col-2 col-3 col-4 col-5 col-6' -cs
# or
cat test/data/data2.txt | pc --header='col-1 col-2 col-3 col-4 col-5 col-6' -cs
```

**Example 4: Formatting with a header and pretty print**

```sh
pc --file=test/data/data2.txt --header='col-1 col-2 col-3 col-4 col-5 col-6' -pp
# or
cat test/data/data2.txt | pc --header='col-1 col-2 col-3 col-4 col-5 col-6' -pp
```

**Example 5: Selecting specific columns**

```sh
pc --file=test/data/data2.txt --header='col-2 col-5' 5 2
# or
cat test/data/data2.txt | pc --header='col-2 col-5' 5 2
```

**Example 6: Formatting and filtering command output**

```sh
oc get pod -o wide --all-namespaces |head -n15| pc -ts -cs  8 1 2 5 6
```

**Example 7: Filtering and converting to JSON**

```sh
oc get pod -o wide --all-namespaces |pc -json --filter='wrk-v01'   8 1 2 5 6
```

## PROJECT STRUCTURE

The project is structured as follows:

*   `internal/argparse`: Contains code for argument parsing.
*   `internal/dataformat`: Contains code for data formatting.
*   `internal/loaddata`: Contains code for loading data.
*   `test`: Contains test files.
*   `docs`: Contains documentation files.

## CONTRIBUTING

Contributions are welcome! Please open an issue or submit a pull request.

## AUTHOR

Dirk Jäger (dirk.jaeger.dj@gmail.com)

## COPYRIGHT

Copyright © 2020 Free Software Foundation, Inc. License GPLv3+: GNU GPL version 3 or later <https://gnu.org/licenses/gpl.html>. This is free software: you are free to change and redistribute it. There is NO WARRANTY, to the extent permitted by law.

