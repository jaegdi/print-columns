
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
        -colsep='|'       ColumnSeparator  define the character to separate the columns, default='|'.
        -filter='string'  Filter lines,    process only lines where 'string' is found.
        -sortcol=colnum:  SortColumn       number of column, to sort for
        -gcol=colnum:     GroupCol         write a separator when the value in this column is different
                                           to the value in the previous line to group the values in this column.
                                           Number refers to the number of the output column.
        -ts               TitleSeparator,  draws a separator line between first and second line of output.
        -fs               FooterSeparator, draws a separator line between last and second last line of output.
        -cs               ColumnSeparator, draws a separator (default=|) between columns of output.
        -pp               PrettyPrint,     draw cell borders and all separators.
        -rh               RemoveHeader,    removes the first line.
        -num              Num-bering,      insert col numbers in the first line.
        -csv              CSV,             write output in CSV format.
        -json             JSON,            write output in JSON format.
        -v                                 verify, show all given parameters

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
    
Usage: pc [OPTIONS] argument ...


EXAMPLES

    - A data-file 'data2.txt' has the following content:

        aaaaa bbbbbbbbbbbb cccccccc dd eeeeeee fffffffffff
        aaaaaaaaaaaaaaa bbbbbbbbbbbbb cc dddddddddddd eeeeaaaaaeeee ff
        aa bbbbb cc dd ee ffffffffff

    - To print the file with formated columns
        cmd: pc --file=data2.txt
             or    cat data2.txt | pc
        result:
        aaaaa             bbbbbbbbbbbb    cccccccc   dd             eeeeeee         fffffffffff
        aaaaaaaaaaaaaaa   bbbbbbbbbbbbb   cc         dddddddddddd   eeeeaaaaaeeee   ff
        aa                bbbbb           cc         dd             ee              ffffffffff

    - To print that file formated with additional header:
        cmd: pc -file=data2.txt -header='col-1 col-2 col-3 col-4 col-5 col-6'
             or   cat data2.txt | pc --header='col-1 col-2 col-3 -col-4 col-5 -col-6'
        result:
        col-1             col-2           col-3             col-4   col-5                 col-6
        ---------------   -------------   --------   ------------   -------------   -----------
        aaaaa             bbbbbbbbbbbb    cccccccc   dd             eeeeeee         fffffffffff
        aaaaaaaaaaaaaaa   bbbbbbbbbbbbb   cc         dddddddddddd   eeeeaaaaaeeee   ff
        aa                bbbbb           cc         dd             ee              ffffffffff

    - To print that file formated with additional header and columnseparator:
        cmd: pc --file=data2.txt --header='col-1 col-2 col-3 col-4 col-5 col-6'
             or   cat data2.txt | pc --header='col-1 col-2 col-3 col-4 col-5 col-6' -cs
        result:
        col-1           | col-2         | col-3    | col-4        | col-5         | col-6
        --------------- | ------------- | -------- | ------------ | ------------- | -----------
        aaaaa           | bbbbbbbbbbbb  | cccccccc | dd           | eeeeeee       | fffffffffff
        aaaaaaaaaaaaaaa | bbbbbbbbbbbbb | cc       | dddddddddddd | eeeeaaaaaeeee | ff
        aa              | bbbbb         | cc       | dd           | ee            | ffffffffff

    - To print that file formated with additional header and prettyprint:
        cmd: pc --file=data2.txt --header='col-1 col-2 col-3 col-4 col-5 col-6'
             or   cat data2.txt | pc --header='col-1 col-2 col-3 col-4 col-5 col-6' -pp
        result:
        | --------------- | ------------- | -------- | ------------ | ------------- | ----------- |
        | col-1           | col-2         | col-3    | col-4        | col-5         | col-6       |
        | --------------- | ------------- | -------- | ------------ | ------------- | ----------- |
        | aaaaa           | bbbbbbbbbbbb  | cccccccc | dd           | eeeeeee       | fffffffffff |
        | aaaaaaaaaaaaaaa | bbbbbbbbbbbbb | cc       | dddddddddddd | eeeeaaaaaeeee | ff          |
        | aa              | bbbbb         | cc       | dd           | ee            | ffffffffff  |
        | --------------- | ------------- | -------- | ------------ | ------------- | ----------- |

    - To print col2 and col5 with additional headers in reverse order
        cmd: pc --file=data2.txt --header='col-2 col-5 ' 5 2
             or cat data2.txt | pc --header='col-2 col-5 ' 5 2
        result:
        col-5           col-2
        -------------   -------------
        eeeeeee         bbbbbbbbbbbb
        eeeeaaaaaeeee   bbbbbbbbbbbbb
        ee              bbbbb

    - Format the output of a command
        cmd: oc get pod -o wide --all-namespaces |head -n15| pc -ts -cs  8 1 2 5 6
        result:
        NODE                      | NAMESPACE             | NAME                                               | RESTARTS | AGE
        ------------------------- | --------------------- | -------------------------------------------------- | -------- | ---
        int-apc0-wrk-v08.sf-rz.de | app-monitoring        | prometheus-prometheus-0                            | 11       | 1d
        int-apc0-wrk-v10.sf-rz.de | br-test               | rsync-container-1-trkwt                            | 1        | 27d
        int-apc0-inf-v01.sf-rz.de | cluster-tasks         | ldapgroupsync-1583331300-86bg8                     | 0        | 22d
        int-apc0-inf-v01.sf-rz.de | cluster-tasks         | ldapgroupsync-1583334900-fsh48                     | 0        | 22d
        int-apc0-inf-v01.sf-rz.de | cluster-tasks         | prune-builds-1585239000-lrncj                      | 0        | 1h
        int-apc0-inf-v01.sf-rz.de | cluster-tasks         | prune-deployments-1585242300-vr22s                 | 0        | 24m
        int-apc0-inf-v01.sf-rz.de | cluster-tasks         | registry-image-pruning-1585235220-prbj7            | 0        | 2h
        int-apc0-inf-v03.sf-rz.de | default               | docker-registry-5-bxk5x                            | 0        | 27d
        int-apc0-mst-v00.sf-rz.de | default               | registry-console-7-sj72f                           | 0        | 8d

    - Filter the output of a command and convert to json
        cmd:  oc get pod -o wide --all-namespaces |pc -ts -cs -json --filter='wrk-v01'   8 1 2 5 6
        result:
        {
            "data": [
                [
                    "int-apc0-wrk-v01.sf-rz.de",
                    "fpc-fa2",
                    "datenkopie-zulieferung-46-46dhb",
                    "1",
                    "8h"
                ],
                [
                    "int-apc0-wrk-v01.sf-rz.de",
                    "fpc-int1",
                    "datenkopie-zulieferung-64-pdp5r",
                    "1",
                    "8h"
                ],
                [
                    "int-apc0-wrk-v01.sf-rz.de",
                    "openshift-logging",
                    "logging-fluentd-6bg5h",
                    "3",
                    "23d"
                ]
            ]
        }
    
