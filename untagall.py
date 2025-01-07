#!/bin/env python3
import json
import os
import subprocess
import argparse

def main():
    args = parseArgs()
    if args.tag == None:
        print("no tag given, check help for proper usage, exiting...")
        exit(1)
    clientHandle(args.tag, args.tagAddMode)

def parseArgs():
    parser = argparse.ArgumentParser(
                    prog='untagall',
                    description='script to tag/untag all windows, using a user defined tag',
                    epilog='specifying a tag with -t is requared!')
    parser.add_argument('-t', '--tag', type=str)
    parser.add_argument('-ta', '--tagAddMode',action='store_true',default=False)
    args = parser.parse_args()
    return args

def clientHandle(targetTag, tagAddMode):
    p = subprocess.run(['hyprctl', 'clients', '-j'],capture_output=True,text=True)
    hyprctlClientData = p.stdout.strip()
    hyprctlClientDataJson = json.loads(hyprctlClientData)
    hyprCommand = ''
    for client in hyprctlClientDataJson:
        clientAddress = client["address"]
        clientTags = client["tags"]
        for clientTag in clientTags:
            if tagAddMode == False:
                if clientTag == targetTag: # add support for using arg1 as tag input
                    hyprCommand = hyprCommand + f"dispatch tagwindow -{targetTag} address:{clientAddress} ; "
            else:
                if clientTag != targetTag:
                    hyprCommand = hyprCommand + f"dispatch tagwindow +{targetTag} address:{clientAddress} ; "


    if hyprCommand != "":
        print(f'hyprctl --batch "{hyprCommand}"') 
        subprocess.run(['hyprctl', '--batch', f'{hyprCommand}']) 

if __name__ == "__main__":
    main()
