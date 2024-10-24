#!/usr/bin/env bash
signal="$1"
mode="$2"
input_command="$3"
debug=0
p="" #this var is printed at the end of the script if not running in command mode

note_debug(){
if [ "$debug" == "1" ]
then
    notify-send "$@"
fi
}


main(){
mullvad_json="$(mullvad status -j)"
mullvad_status="$(echo "$mullvad_json" | jq '.state')"
mullvad_locale="$(echo "$mullvad_json" | jq -r '.details.location.country' )"


if [ "$mullvad_status" == '"connected"' ]
then 
    p="$mullvad_locale "
    #echo ""
    #echo "🔒"
elif [ "$mullvad_status" == '"disconnected"' ]
then
    p=""
    #echo "🔓"
elif [ "$mullvad_status" == '"connecting"' ]
then
    p="Connecting..."
    sleep 0.1
    pkill -SIGRTMIN+$signal waybar
elif [ "$mullvad_status" == '"disconnecting"' ]
then
    p="Disconnecting..."
    sleep 0.1
    pkill -SIGRTMIN+$signal waybar
else
    notify-send "$mullvad_status not captured"
    p="$mullvad_status"
fi

#json_output=$p}
#echo '{"text": "uk connected", "tooltip": "hello"}' | jq --unbuffered --compact-output
echo "$p"
}
mull_till_new(){ 
# this gets mullvad's state and runs a program of choie then loops till the state changes and then it updates the waybar module
mullvad_status_old="$(mullvad status -j | jq '.state')"
mullvad_status="$mullvad_status_old"
$1
note_debug "command used '$1'"
while [ "$mullvad_status_old" == "$mullvad_status" ]
do
    mullvad_status="$(mullvad status -j | jq '.state')"
    note_debug "updating $mullvad_status_old $mullvad_status"
done
pkill -SIGRTMIN+$signal waybar
}


if [ "$mode" == "command" ]
then
    mull_till_new "$input_command"
else
main
fi
