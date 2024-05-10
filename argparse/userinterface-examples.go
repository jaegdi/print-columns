package pc

import (
	"flag"
	"fmt"
	"os"
)

// cmdExamples print examples of usage
func cmdExamples() {
	usageText := `

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
        host-wrk-v08.my-domain.de | app-monitoring        | prometheus-prometheus-0                            | 11       | 1d
        host-wrk-v10.my-domain.de | br-test               | rsync-container-1-trkwt                            | 1        | 27d
        host-inf-v01.my-domain.de | cluster-tasks         | ldapgroupsync-1583331300-86bg8                     | 0        | 22d
        host-inf-v01.my-domain.de | cluster-tasks         | ldapgroupsync-1583334900-fsh48                     | 0        | 22d
        host-inf-v01.my-domain.de | cluster-tasks         | prune-builds-1585239000-lrncj                      | 0        | 1h
        host-inf-v01.my-domain.de | cluster-tasks         | prune-deployments-1585242300-vr22s                 | 0        | 24m
        host-inf-v01.my-domain.de | cluster-tasks         | registry-image-pruning-1585235220-prbj7            | 0        | 2h
        host-inf-v03.my-domain.de | default               | docker-registry-5-bxk5x                            | 0        | 27d
        host-mst-v00.my-domain.de | default               | registry-console-7-sj72f                           | 0        | 8d

    - Filter the output of a command and convert to json
        cmd:  oc get pod -o wide --all-namespaces |pc -json --filter='wrk-v01'   8 1 2 5 6
        result:
        {
            "data": [
                [
                    "host-wrk-v01.my-domain.de",
                    "fpc-fa2",
                    "datenkopie-zulieferung-46-46dhb",
                    "1",
                    "8h"
                ],
                [
                    "host-wrk-v01.my-domain.de",
                    "fpc-int1",
                    "datenkopie-zulieferung-64-pdp5r",
                    "1",
                    "8h"
                ],
                [
                    "host-wrk-v01.my-domain.de",
                    "openshift-logging",
                    "logging-fluentd-6bg5h",
                    "3",
                    "23d"
                ]
            ]
        }

        If you define the filter, to also get the header and set the flag -ts (titleseperator), then the json output respects the header info:
        cmd:  oc get pod -o wide --all-namespaces |pc -ts -json --filter='NAME|wrk-v01'   8 1 2 5 6
        result:
        {
            "NODES": [
              {
                "NODE": "cid-scp0-wrk-v01.sf-rz.de",
                "data": {
                  "NAMESPACE": "b2c-fpc-int1",
                  "NAME": "fpc-request-history-store-8ff8d8794-b8rf4",
                  "RESTARTS": "0",
                  "AGE": "15d"
                }
              },
              {
                "NODE": "cid-scp0-wrk-v01.sf-rz.de",
                "data": {
                  "NAMESPACE": "ibs-bil-app",
                  "NAME": "billing-controlling-56d5fb5fb4-bl475",
                  "RESTARTS": "0",
                  "AGE": "15d"
                }
              },
              {
                "NODE": "cid-scp0-wrk-v01.sf-rz.de",
                "data": {
                  "NAMESPACE": "openshift-dns",
                  "NAME": "dns-default-8d698",
                  "RESTARTS": "2",
                  "AGE": "16d"
                }
              }
            ]
        }
    - format and filter output of oc get nodes

        At first get overview of columns
        oc get nodes -o wide | pc -filter="NAME|-(mst|inf)-v" -mb -gcol=3 -sortcol=3 -pp -num

        Then select the columns that you want to display
        oc get nodes -o wide | pc -filter="NAME|-(mst|inf)-v" -mb -gcol=3 -sortcol=3 -pp 1:3 6:7 9

        | ------------------------- | ------ | ------ | ------------ | ----------- | --------------------------- |
        | NAME                      | STATUS | ROLES  | INTERNAL-IP  | EXTERNAL-IP | KERNEL-VERSION              |
        | ------------------------- | ------ | ------ | ------------ | ----------- | --------------------------- |
        | ------------------------- | ------ | ------ | ------------ | ----------- | --------------------------- |
        | host-inf-v05.my-domain.de | Ready  | cicd   | 192.68.42.42 | <none>      | 3.10.0-1160.42.2.el7.x86_64 |
        | host-inf-v06.my-domain.de | Ready  | ''     | 192.68.42.47 | <none>      | 3.10.0-1160.42.2.el7.x86_64 |
        | host-inf-v07.my-domain.de | Ready  | ''     | 192.68.42.43 | <none>      | 3.10.0-1160.42.2.el7.x86_64 |
        | ------------------------- | ------ | ------ | ------------ | ----------- | --------------------------- |
        | host-inf-v00.my-domain.de | Ready  | infra  | 192.68.42.91 | <none>      | 3.10.0-1160.42.2.el7.x86_64 |
        | host-inf-v01.my-domain.de | Ready  | ''     | 192.68.42.40 | <none>      | 3.10.0-1160.42.2.el7.x86_64 |
        | host-inf-v02.my-domain.de | Ready  | ''     | 192.68.42.45 | <none>      | 3.10.0-1160.42.2.el7.x86_64 |
        | host-inf-v03.my-domain.de | Ready  | ''     | 192.68.42.41 | <none>      | 3.10.0-1160.42.2.el7.x86_64 |
        | host-inf-v04.my-domain.de | Ready  | ''     | 192.68.42.46 | <none>      | 3.10.0-1160.42.2.el7.x86_64 |
        | ------------------------- | ------ | ------ | ------------ | ----------- | --------------------------- |
        | host-mst-v00.my-domain.de | Ready  | master | 192.68.42.90 | <none>      | 3.10.0-1160.42.2.el7.x86_64 |
        | host-mst-v01.my-domain.de | Ready  | ''     | 192.68.42.20 | <none>      | 3.10.0-1160.42.2.el7.x86_64 |
        | host-mst-v02.my-domain.de | Ready  | ''     | 192.68.42.25 | <none>      | 3.10.0-1160.42.2.el7.x86_64 |
        | ------------------------- | ------ | ------ | ------------ | ----------- | --------------------------- |

        `
	fmt.Printf("Usage: %s [OPTIONS] argument ...\n", os.Args[0])
	fmt.Println(usageText)
	flag.PrintDefaults()
}
