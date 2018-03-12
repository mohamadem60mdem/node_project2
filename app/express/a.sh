#!/bin/bash


#   A simple copy script
    #read a1 
    #echo $a1 ;  
    #echo ;echo ;echo ;echo ;echo \n;echo aa;printf "hello\nworld\n\n\n";
    #echo -n "--" ;
    # echo -e "..\e[32m..\e[33m..... \e[31m ........."
#



printf "\n \e[32m START......\e[39m................................................. \n\n\n";


    a2=r ;a3=4; a4=s 
    # read count;

    # if [ $count == "s" ]
    #     then
    #     echo "Count is 100"
    #     echo "d" >> 1.txt
    # fi
    #  echo "qq" > 1.txt

     mkdir -pv x/a1/a2
    # rmdir -pv x/a1/a2




    # ls -R






printf "\n\n \e[32m END......\e[39m................................................. \n\n\n";









exit


# Link filedescriptor 10 with stdin
exec 10<&0
# stdin replaced with a file supplied as a first argument
exec < $1
# remember the name of the input file
in=$1

# init
file="current_line.txt"
let count=0

# this while loop iterates over all lines of the file
while read LINE
do
    # increase line counter 
    ((count++))
    # write current line to a tmp file with name $file (not needed for counting)
    echo $LINE > $file
    # this checks the return code of echo (not needed for writing; just for demo)
    if [ $? -ne 0 ] 
     then echo "Error in writing to file ${file}; check its permissions!"
    fi
done




echo "Number of lines: $count"
echo "The last line of the file is: `cat ${file}`"

# Note: You can achieve the same by just using the tool wc like this
echo "Expected number of lines: `wc -l $in`"

# restore stdin from filedescriptor 10
# and close filedescriptor 10
exec 0<&10 10<&-



