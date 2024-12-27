# script to auto set bluelight upon day and night

## NOTE script is heavy W.I.P and is currently using HARD CODED values This document is partly for planing out the script, while also writing docs for the tool

### usage
`wl-gamma -nt "3300" -dt "6500" -n "10,21"

explaination:
`wl-gamma`: name of script
`-nt` command line option to set night temp (defualt 3300)
`-dt` command line option to set day temp (defualt 6500)
`-n` command line options to set night hours (defualt 10,21)

The way the `-n` option works is that blue light temp is trigged IF ether

1. first number (In this case 10) is higher then the current hour in 24hour time.
2. second number (In this case 21) is lower then the current hour in 24hour time.



## implementiation notes/choices

### can settings be changed at runtime?
The settings can only be set via hard coding or command flags.

### what if temp is user set?
only change setting if night is lower temp then current

### how is user sleep times set?
launch option sets morning and then night temp like so "10,21" with comma speration
