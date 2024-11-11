#!/bin/env python3
import os
import subprocess
import tomllib


def configload(configpath):
    with open(configpath, "rb") as f:
        tomldata = tomllib.load(f)
    return tomldata

def bemenu(items,launcher):
    p = subprocess.run(['echo "' + items + '" | ' + launcher], shell=True,capture_output=True,text=True)
    launcherOutput = p.stdout.strip()
    return launcherOutput



def main():
    xdg_config_home = os.environ.get('XDG_CONFIG_HOME') or \
         os.path.join(_home, '.config')
    
    configpath = xdg_config_home + "/bookmarks.toml"

    tomldata = configload(configpath)
    #print(tomldata)

    try:
        tomldata_launcher = tomldata["options"]["launcher"]
    except KeyError:
        tomldata_launcher = "bemenu"

    try:
        tomldata_browserlaunch = tomldata["options"]["browser_launch"]
    except KeyError:
        print("browser launch args not set in bookmarks.toml")
        print("exiting...")
        exit()
    
    tomldata_bookmarks = tomldata['bf']
    print(tomldata_bookmarks)
    

    tomldata_namelist = []
    for x in tomldata_bookmarks:
        print(x['name'])
        tomldata_namelist.append(x['name'])
        print(x['urls'])
       # print(x)
       # print(x)
    
    itemsclean = '\n'.join(tomldata_namelist)
    launcherOutput = bemenu(itemsclean,tomldata_launcher)
    print(launcherOutput)
                
    urls_list = []
    urls_list_raw = []
    
    for x in tomldata_bookmarks: # this is a really slow way of doing this
        if launcherOutput == x['name']:
            folder_urls = x['urls']
            print(folder_urls)
            for y in folder_urls[0]:
                print(y)
                print(folder_urls[0][y])
                urls_list.append(folder_urls[0][y])
                urls_list_raw.append(y)
                #urls_list.append(y)
     
    itemsclean = '\n'.join(urls_list)
    launcherOutput = bemenu(itemsclean,tomldata_launcher)
    url_list_raw_id = urls_list.index(launcherOutput)
    print("oppening " + urls_list_raw[url_list_raw_id])
    
    #tomldata_makeitems = tomldata["itemcreate"]
    #tomldata_modifyitems = tomldata["itemmodify"]
    #launcher = tomldata["options"]["launcher"]
    #modifyitems = tomldata_makeitems | tomldata_modifyitems



if __name__ == "__main__":
    main()
