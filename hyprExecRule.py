#!/bin/env python3
import json
import os
import subprocess
import re
import time

def main():
    while 1 == 1:
        gsettheme()

def gsettheme():
    theme = 'prefer-dark'
    sleeptime = 1
    p = subprocess.run(['hyprctl', 'clients', '-j'],capture_output=True,text=True)
    clientsData = p.stdout.strip()
    clientsData_json = json.loads(clientsData)
    for client in clientsData_json:

        awClass = client['class']
        awTitle = client['title']
        
        if awClass == "firefox" and len(re.findall(".*Google Sheets — Mozilla Firefox", awTitle)) != 0:
            theme = 'prefer-light'
            sleeptime = 0.2
            print(awClass, awTitle)
        #elif awClass == "your app here":
    
    subprocess.run(['gsettings', 'set', 'org.gnome.desktop.interface', 'color-scheme', f"'{theme}'"])
    print('gsettings', 'set', 'org.gnome.desktop.interface', 'color-scheme', f"'{theme}'")
    time.sleep(sleeptime)

    
main()
