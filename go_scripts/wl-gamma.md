# script to auto set bluelight upon day and night

## extra notes:
clean up output such that it could be used as a replacement for the wl-gamma-rs waybar module

## usage
`wl-gamma -nt 3300 -dt 6500 -mh 10 -nh 21

explaination:
`wl-gamma`: name of script
`-nt` command line option to set night temp (defualt 3300)
`-dt` command line option to set day temp (defualt 6500)
`-nh` command line options to set night hour (defualt 21)
`-mh` command line options to set morning hour (defualt 10)


## implementiation notes/choices

### can settings be changed at runtime?
The settings can only be set via hard coding or command flags.

### what if temp is user set?
only change setting if night is lower temp then current

### how is user sleep times set?
launch option sets morning and then night temp like so "10,21" with comma speration
