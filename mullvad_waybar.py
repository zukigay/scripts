#!/usr/bin/env python
import re
#import os
import subprocess
import json
import sys

mullcli = "mullvad"
preface = "mull "

def mullcall(mullvadcommand):
    g = subprocess.run(mullvadcommand, capture_output=True,text=True).stdout.strip()
    return g

#def mull_get_state():


def localemode(): # w.i.p 
    relay_list = mullcall([mullcli, "relay", "list"])
    county_lines = [line for line in relay_list.splitlines() if not re.search(r'	', line)]
    relay_dict = {}
    relay_list_codes = []
    relay_list_names = []
    for line in county_lines:
        if not line == "":
            line_split = line.split("(")
            relay_list_codes.append(line_split[1])
            relay_list_names.append(line_split[0])
            relay_dict.update({line_split[0]: line_split[1]})
    
    locale_wanted = 20
    print(relay_list_names[locale_wanted] + "(" + relay_list_codes[locale_wanted])
    

def main():
    mull_status_raw = mullcall([mullcli, "status", "-j"])
    mull_status = json.loads(mull_status_raw)
    mull_state = mull_status['state']
    mull_locale = "UK"
    match mull_state:
        case "connected":
            mull_locale = mull_status['details']['location']['country']
            print(preface + mull_locale + " ")
        case "disconnected":
            print(preface  +  " ")
        case "connecting":
            print("Connecting...")
        case "disconnecting":
            print("Disconnecting...")
        case _:
            print("nothing")



if len(sys.argv) > 1:
    if sys.argv[1] == "l":
        localemode()
    else:
        main()
else:
    main()

