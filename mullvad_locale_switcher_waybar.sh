#!/usr/bin/env bash
tempdir="/tmp/mullvad_locale_waybar"
tempcodelist="$tempdir/mullvadrelaycodes"
tempcodepicked="$tempdir/mullvadrelaypicked"
input1="$1"

if [ ! -f "$tempcodelist" ]
then
mkdir /tmp/mullvad_locale_waybar 
#mullvadrelaylistcodes="$(mullvad relay list | awk '!/^\t\t.*$/' | awk '/./' | sed 's/ @.*//g' | sed 's/^\t/ ∙ /g' | awk '{print (substr($0,length($0)-3,3))}' | sed 's/(//g' )" 

mullvadrelaylistcodes="$(mullvad relay list | awk '!/^\t\t.*$/' | awk '/./' | sed 's/ @.*//g' | sed 's/^\t/ ∙ /g' | grep -v "∙")"


echo "$mullvadrelaylistcodes" > "$tempcodelist"
else
mullvadrelaylistcodes="$(cat "$tempcodelist")"
fi

if [ ! -f "$tempcodepicked" ]
then
echo "1" > "$tempcodepicked" 
fi
NUM="$(cat "$tempcodepicked")"

if [ "$input1" == "-" ]
then
    NUM=$((NUM-1))
    if [ "$NUM" -gt 0 ]
    then
        echo "$NUM" > "$tempcodepicked"
    elif [ "$NUM" -gt -1 ]
    then
        echo "1" > "$tempcodepicked"
    fi
elif [ "$input1" == "set" ] 
then
    location="$(sed "${NUM}q;d" "$tempcodelist" | awk -F'[()]' '{print $2}')" 
    mullvad relay set location "$location"
elif [ "$input1" == "+" ]
then
    NUM=$((NUM+1))
    NUMFILE="$(sed -n '$=' $tempcodelist)" # gets number of lines in file
    if [ "$NUM" -le "$NUMFILE" ]
    then
        echo "$NUM" > "$tempcodepicked" 
    fi
else
echo "$(sed "${NUM}q;d" "$tempcodelist")"
fi



