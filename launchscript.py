#!/bin/env python3
import os
import subprocess
import tomllib


def configload(configpath):
    with open(configpath, "rb") as f:
        tomldata = tomllib.load(f)
    return tomldata

def getunixpathitems():
    UNIXPATH = os.environ.get('PATH').rsplit(":")
    UNIXPATHITEMS = []
    for x in UNIXPATH:
        if os.path.exists(x):
            listdir = os.listdir(x)
            for y in listdir:
                file = (x + "/" + y)
                if os.access(file, os.X_OK) and not os.path.isdir(file):
                    #print(file + " is executable")
                    UNIXPATHITEMS.append(y)

    return UNIXPATHITEMS

def main():
    xdg_config_home = os.environ.get('XDG_CONFIG_HOME') or \
         os.path.join(_home, '.config')
    
    configpath = xdg_config_home + "/launchscript.toml"
    tomldata = configload(configpath)
    tomldata_makeitems = tomldata["itemcreate"]
    tomldata_modifyitems = tomldata["itemmodify"]
    launcher = tomldata["options"]["launcher"]
    modifyitems = tomldata_makeitems | tomldata_modifyitems
    #print(modifyitems)

    items = getunixpathitems()
    rawpathitems = items 
    for z in tomldata_makeitems:
        items.append(z)
    
    items.sort()
    
    itemclean = '\n'.join(items) # probably a cleaner way to do this
    
    p = subprocess.run(['echo "' + itemclean + '" | ' + launcher], shell=True,capture_output=True,text=True)
    launcherOutput = p.stdout.strip()
    LAUNCH = ""
    #if launcherOutput in modifyitems:
    if launcherOutput in modifyitems and launcherOutput in rawpathitems:
        LAUNCH = modifyitems[launcherOutput]
    else:
        LAUNCH = launcherOutput

    subprocess.run([LAUNCH], shell=True) #without shell=True it will not take into account spaces

if __name__ == "__main__":
    main()
