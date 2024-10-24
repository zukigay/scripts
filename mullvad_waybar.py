#!/usr/bin/env python
#import os
import subprocess
import json

mullcli = "mullvad"
preface = "mull "

def mullcall(mullvadcommand):
    g = subprocess.run(mullvadcommand, capture_output=True,text=True)
    return g

#def mull_get_state():
    

def main():
    mull_status_raw = mullcall([mullcli, "status", "-j"]).stdout.strip()
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




main()
