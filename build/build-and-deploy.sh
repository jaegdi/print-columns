#!/usr/bin/env bash
set -eo pipefail

script="$(basename "$0")"
scriptdir="$(dirname "$0")"
dir=$(dirname "$scriptdir")
echo "script: $script, scriptdir: $scriptdir, Dir: $dir"
cd "$dir"
quayurl="registry-quay-quay.apps.pro-scp1.sf-rz.de"
if podman login -u "$USER" -p "$(eval "$LDAPPASSWORDPROVIDER")" "$quayurl"; then
    echo "quayurl: $quayurl"
else
    echo "ERROR: podman login failed" 1>&2
    exit 1
fi

# Set defaults and evaluate commandline parameters
CLUSTER=cid-scp0
tagversion="$(get-git-tag.sh)"
# per default make it all
build='false'

hilfe() {
    if [[ -n $1 ]]; then
        echo
        echo '***'  "$*"  '***'
        echo
    fi
    cat <<-EOH

    SYNOPSIS
        $script [-b|--build] [-h|--help]

    OPTIONS
        -b <name> | --build[=]<name>
            Enable to build the pc executable and deploy to artifactory
        -t <tag> | --tag[=]<tag>
            The tag version of the pc to build
        -h | --help
            Print this help message

    DESCRIPTION
        $script builds and deploys the pc to the specified cluster.

    EXAMPLES
        $script -b name -t v1.0.0
            Build the pc with tag v1.0.0 and deploy it
        $script -b name -t latest
            Build the pc with the latest tag
        $script -b name
            Build the pc with the newest tag that is defined in git
EOH
}


optspec=":vhb:t:-:"
while getopts "$optspec" optchar; do
    case "${optchar}" in
        -)  # Evaluate long options
            case "${OPTARG}" in
                build)
                    val="${!OPTIND}"; OPTIND=$(( $OPTIND + 1 ))
                    echo "Parsing option: '--${OPTARG}', value: '${val}'" >&2;
                    build="$val"
                    ;;
                build=*)
                    val=${OPTARG#*=}
                    opt=${OPTARG%=$val}
                    echo "Parsing option: '--${opt}', value: '${val}'" >&2
                    build="$val"
                    ;;
                tag)
                    val="${!OPTIND}"; OPTIND=$(( $OPTIND + 1 ))
                    echo "Parsing option: '--${OPTARG}', value: '${val}'" >&2;
                    tagversion="$val"
                    ;;
                tag=*)
                    val=${OPTARG#*=}
                    opt=${OPTARG%=$val}
                    echo "Parsing option: '--${opt}', value: '${val}'" >&2
                    tagversion="$val"
                    ;;
                help)
                    hilfe ''
                    exit
                    ;;
                *)
                    if [ "$OPTERR" = 1 ] && [ "${optspec:0:1}" != ":" ]; then
                        hilfe "Unknown option --${OPTARG}"
                        exit
                    fi
                    ;;
            esac;;
        # Evaluate short options
        b)
            val="${OPTARG#*=}";
            echo "Parsing option: '-${optchar}', value: '${val}'" >&2
            build="$val"
            ;;
        t)
            val="${OPTARG#*=}";
            echo "Parsing option: '-${optchar}', value: '${val}'" >&2
            tagversion="$val"
            ;;
        h)
            val="${!OPTIND}";
            hilfe ''
            exit 0
            ;;
        *)
            if [ "$OPTERR" != 1 ] || [ "${optspec:0:1}" = ":" ]; then
                echo "Non-option argument: '-${OPTARG}'" >&2
            fi
            ;;
    esac
done

remember-current-cluster
# login in the build cluster, which is cid-scp0, where the scp-build namespace is located
if ocl cid-scp0 -d > /dev/null 2>&1 &&
    ocl > /dev/null 2>&1; then  # login in the current cluster
    echo "working on CLUSTER: $CLUSTER"
else
    echo "ERROR: ocl login failed" 1>&2
    exit 1
fi

echo "# B U I L D   L O C A L   A N D   D E P L O Y   T O   A R T I F A C T O R Y"
if [ "$build" != 'false' ]; then
    if [ "$tagversion" != 'latest' ]; then
        # Ensure git checkout master is executed on script exit
        trap 'git checkout master' EXIT
        echo "git checkout $tagversion"
        git checkout "$tagversion"
    fi
    echo "Build $build local and deploy to artifactory"
    "$scriptdir"/_build-and-deploy-to-artifactory.sh $build

fi


switch-back-to-current-cluster
