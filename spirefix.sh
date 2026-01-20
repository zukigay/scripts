#!/bin/bash
datatype="class"
targetclass=""
while read type ; do 
    # to add a game to this script simply add to this array first the class name (although only the first word of it with a prefixed "
    # then add the games pid (note games with mutiable pids might not work
    for data in '"Slay' SlayTheSpire '"Mosa' 'Mosa Lina'  ; do
        if [ "$datatype" == "class" ] ; then
            targetclass="$data"
            echo "targetclass = $data" 
            datatype="pid"
            continue
        else
            pidname=$data
            datatype="class"
            echo "targetpid = $data" 
        fi


        pid="$(pgrep "$pidname")"
        if [ "$pid" ] ; then
            if [ "$type" = "_NET_ACTIVE_WINDOW(WINDOW): window id # 0x0" ]
            then
                wpctl set-mute -p "$(pgrep "$pidname")" 1
            else
            	set -- $type
                classcmd=$(xprop -id "$5" WM_CLASS)
                set -- $classcmd
                class=$3
                echo "class = $class"
                if [ "$class" == "$targetclass" ] ; then
                    wpctl set-mute -p "$pid" 0
                else
                    wpctl set-mute -p "$pid" 1
                fi
            fi
        fi
    done
done < <(xprop -spy -root _NET_ACTIVE_WINDOW)
