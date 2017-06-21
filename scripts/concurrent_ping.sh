#!/bin/bash
program=$0
USAGE="""
This is a tiny script to ping ip addresses concurrently.\n\n
[Usage]\n
\t $program  [-n CONCURRENT_NUM] [-h] ADDR_LIST_FILE
"""

# subprocess num to run this job
concurrent=50
while getopts "n:h" opt; do
    case $opt in
    n)
        concurrent=$OPTARG
        ;;
    h)
        echo -e $USAGE
        exit 0
        ;;
    \?)
      echo -e $USAGE
      exit 1
      ;;
    esac
done
shift $(($OPTIND - 1))
# check argument
if [[ $# < 1 ]]; then
    echo "Need file name as argv"
    echo -e $USAGE
    exit 1
fi
iplist_file=$1

rm .tmp -rf && mkdir .tmp

# function for do ping with input_file and output_file as arguments
function doping()
{
    if [[ $@ < 2 ]];then
        echo "Need two args"
        return
    fi
    ifile=$1
    ofile=$2
    echo "process file $ifile"
    ips=`cat $ifile`
    for ip in $ips
    {
        res=`ping $ip -c 1 | grep "time="`
        if [[ "$res" == "" ]]; then
            echo "$ip is bad: $res"
            echo  "$ip" >> "$ofile"_bad
        else
            echo "$ip is ok: $res"
            echo "$ip is ok: $res" >> "$ofile"_good
        fi  
    }
    echo "process file $ifile done"
}

# function for killing all subprocess
function kill_process()
{
    ps -ef | grep "$program" | grep -v "grep" | awk '{print $2}' | xargs -i kill {}
}

# register quit singal
trap "echo 'You want to kill me' && kill_proce
ss && exit 0" TERM
trap "echo 'You want to interept me' && kill_process && exit 0" INT

# splite ip list into several filesm
awk '{mod = NR % '$concurrent'}{print >> ".tmp/file_"mod}{close(".tmp/file_"mod)}' $iplist_file
cd .tmp
inputs=`ls`

for i in $inputs
do
    (
        doping $i "$i"_output
        ) &
done 

# wait for all subprocess' done
wait

# end
cat *_bad > total_bad && mv total_bad ../
cat *_good > total_good && mv total_good ../

echo "Done"