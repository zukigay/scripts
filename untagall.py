#!/bin/env python3
import json
import os
import subprocess
 
def main():
    p = subprocess.run(['hyprctl', 'clients', '-j'],capture_output=True,text=True)
    hyprctlClientData = p.stdout.strip()
    hyprctlClientDataJson = json.loads(hyprctlClientData)
    hyprCommand = ''
    for client in hyprctlClientDataJson:
        clientAddress = client["address"]
        clientTags = client["tags"]
        for tag in clientTags:
            if tag == "overlay": # add support for using arg1 as tag input
                hyprCommand = hyprCommand + f"dispatch tagwindow -overlay address:{clientAddress} ; "

    if hyprCommand != "":
        print(f'hyprctl --batch "{hyprCommand}"') 
        subprocess.run(['hyprctl', '--batch', f'{hyprCommand}']) 

if __name__ == "__main__":
    main()
