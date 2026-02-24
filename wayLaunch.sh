export MOZ_ENABLE_WAYLAND=1
export NVD_BACKEND=direct
export ELECTRON_OZONE_PLATFORM_HINT=x11
export XDG_SESSION_TYPE=wayland
# export XDG_CURRENT_DESKTOP=river
export XDG_CURRENT_DESKTOP=$1
# export WLR_NO_HARDWARE_CURSORS=1
# export WLR_DRM_NO_MODIFIERS=1
# export WLR_SCENE_DISABLE_DIRECT_SCANOUT=1

# can't live without this xkb option
export XKB_DEFAULT_OPTIONS="lv5:caps_switch"

#nvidia
if [ "$(which nvidia-smi)" ] ; then
    echo "nvidia detected"
    export LIBVA_DRIVER_NAME=nvidia
    export __GLX_VENDOR_LIBRARY_NAME=nvidia
    # export WLR_DRM_NO_MODIFIERS=1 # might be required for wayfire
    export GBM_BACKEND=nvidia-drm
fi

exec dbus-launch --exit-with-session $@
