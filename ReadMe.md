
NAME
    pc - PrintColumns - tidy and filter input and generate formated output

SYNOPSYS
    pc [options]  [column-numbers]

    options:        [-file=input-filename] [-header='col1-header col2-header coln-header'] [-colsep='|'] [-filter='string]
                    [-csv] [-json] [-ts] [-cs] [-rh] [-pp] [-num] [-h,-help] [-man]
    column-numbers: [ 1 2 n:m ], a range of columns can be given with [n:m]

    !! All the options must be set before the column-numbers!

DESCRIPTION
    Text that is in unformatted columns can be formatted and filtered with this command.

    pc - formats text columns from stdin or file and print them as a ASCII table or CSV or JSON.
         With the optional parameters the look and content of the output can be modified.
         The columns for output can be selected by single numbers seperated by space or one or more ranges seperated by colon.
         The input should have the same number of columns in each line.
         The formated result is printed to stdout.

         All named parameters must be defined before the column numbers.

    Without column-numbers parameters 'pc' print all columns from input in formated form

    The parameters:

        -file=filename                     read the text from this file,
                                           if there is also data from STDIN, this is added together


        -header='...'    Headerline,       if the text has no headers, you can define headers.
                                           They must be defined in the original order of the incoming text.
                                           Headers are left adjusted,
                                             if they start with a dash (-), then they right adjusted.
        -sep=' '                           define the seperator, which is used to split the data, default is blank ' '
        -mb                                MoreBlanks, assumes, that columes are separated by more than one blank,
                                           default=false
                                           This can be used, if some fields contains blanks, but only works correct,
                                           if all columns consequently delimited by more than one blank.
                                           Typically useful when input is a preformatted ASCII table with blanks as delimiters.
                                           If a headerline is defined with -header=..., then the fields must also delimited by
                                           two blanks ore more.
        -w=1                               no of blanks between colums seperator and column content, default is 1.
        -colsep='|'       ColumnSeparator  define the character to separate the columns, default='|'.
        -filter='regex'   Filter lines,    process only lines where 'string' or 'regex' matches.
        -sortcol=colnum:  SortColumn       number of column, to sort for. Only one column can be defined for sort.
                                           Number refers to the number of the output column.
        -gcol=colnum:     GroupCol         write a separator when the value in this column is different
                                           to the value in the previous line to group the values in this column.
                                           Number refers to the number of the output column.
        -nf               no format        don't format the colums for common column width.
        -nhl              no headline      The data contains no headline.
        -ts               TitleSeparator   draws a separator line between first and second line of output.
        -fs               FooterSeparator  draws a separator line between last and second last line of output.
        -cs               ColumnSeparator  draws a separator (default=|) between columns of output.
        -pp               PrettyPrint      draw cell borders and all separators.
        -rh               RemoveHeader     removes the first line.
        -mb               MoreBlanks       more than one blank to split columns.
        -num              Num-bering       insert col numbers in the first line.
        -csv              CSV              write output in CSV format.
        -json             JSON             write output in JSON format.
        -help             Help             print help and exit.
        -man              Manual           print help and manual, then exit.
        -v                Verify           print parameter verirfy info.

        1 2 m n      ColumnNumbers,
        m:n          ColumnNumber-ranges   The number or ranges of the columns from the incoming text,
                                           that should printed out. To rearrange the columns
                                           the columns can given in the wanted order.
                                           This parameters must be defined after all other parameters.

        -h -help          Help,            print help and exit
        -man              Manual,          print help and manual, then exit

AUTHOR
    written by Dirk Jäger (dirk.jaeger.dj@gmail.com)

COPYRIGHT
    Copyright © 2020 Free Software Foundation, Inc.  License GPLv3+: GNU GPL version 3 or later <https://gnu.org/licenses/gpl.html>.
    This is free software: you are free to change and redistribute it.  There is NO WARRANTY, to the extent permitted by law.
    
