#!/bin/env python3
import json
import os
import subprocess
import signal

def getpidstate(pid):
    try:
        with open(f"/proc/{pid}/stat") as status_file:
            #piddata = status_file.read.split()
            return status_file.read().split()[2]
    except FileNotFoundError:
        print("hyprctl reporting nonexistent process?\nExiting...")
        exit()


def main():
    p = subprocess.run(['hyprctl', 'activewindow', '-j'],capture_output=True,text=True)
    activewindowdata = p.stdout.strip()
    activewindowdata_json = json.loads(activewindowdata)
    pid = activewindowdata_json['pid']
    
    pid_state = getpidstate(pid)
    if pid_state == "T":
        os.kill(int(pid), signal.SIGCONT)
        print("Resuming " + str(pid))
    else:
        os.kill(int(pid), signal.SIGSTOP)
        print("Pausing " + str(pid))

main()
