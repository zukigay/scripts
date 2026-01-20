#!/usr/bin/tcc -run -lX11
#include <stdio.h>
#include <X11/Xlib.h>
#include <X11/Xatom.h>


#define MAX(A, B)               ((A) > (B) ? (A) : (B))
#define MIN(A, B)               ((A) < (B) ? (A) : (B))
#define CLAMP(a, b, c) ( MIN(MAX(b, c), MAX(MIN(b, c), a)) )
#define LENGTH(X)               (sizeof X / sizeof X[0])
#define END(A)                  ((A) + LENGTH(A))

typedef struct{
    char * wm_name;
    char * wm_class;
} Program;

static Display *dpy;
static Atom utf8;
static Window root;

// Program p[] = {
//     {"Slay the Spire", "Slay the Spire"}
// };



int main() 
{
	if(!(dpy = XOpenDisplay(NULL))) {
		fputs("spirefix: cannot open display\n", stderr);
		return 1;
	}

    root = DefaultRootWindow(dpy);
    Atom property = XInternAtom(dpy, "_NET_ACTIVE_WINDOW", False);
    //return values
    Atom type_return;
    int format_return;
    unsigned long nitems_return;
    unsigned long bytes_left;
    unsigned char *data;
	

    XGetWindowProperty(
        dpy,
        root,
        property,
        0,              //no offset
        1,              //one Window
        False,
        XA_WINDOW,
        &type_return,   //should be XA_WINDOW
        &format_return, //should be 32
        &nitems_return, //should be 1 (zero if there is no such window)
        &bytes_left,    //should be 0 (i'm not sure but should be atomic read)
        &data           //should be non-null
    );
    Window win = ((Window*) data)[0];
    if (win != 0) {
	    property = XInternAtom(dpy, "_NET_WM_NAME", True);
	    unsigned char *p_str = NULL;
	    Atom da;

	    int di;
	    unsigned long dl;
	    utf8 = XInternAtom(dpy, "UTF8_STRING", True);
	    XGetWindowProperty(dpy, win, property, 0, BUFSIZ, False, utf8, &da, &di, &dl, &dl, &p_str);
        printf("_NET_WM_NAME %s\n",p_str);
    }

    // printf("window %li\n",win);

    XFree(data);

	XCloseDisplay(dpy);
    return 0;
}

/*
datatype="class"
targetclass=""
while read type ; do 
    # to add a game to this script simply add to this array first the class name (although only the first word of it with a prefixed "
    # then add the games pid (note games with mutiable pids might not work
    for data in '"Slay' SlayTheSpire '"Mosa' 'Mosa Lina'  ; do
        if [ "$datatype" == "class" ] ; then
            targetclass="$data"
            echo "targetclass = $data" 
            datatype="pid"
            continue
        else
            pidname=$data
            datatype="class"
            echo "targetpid = $data" 
        fi


        pid="$(pgrep "$pidname")"
        if [ "$pid" ] ; then
            if [ "$type" = "_NET_ACTIVE_WINDOW(WINDOW): window id # 0x0" ]
            then
                wpctl set-mute -p "$(pgrep "$pidname")" 1
            else
            	set -- $type
                classcmd=$(xprop -id "$5" WM_CLASS)
                set -- $classcmd
                class=$3
                echo "class = $class"
                if [ "$class" == "$targetclass" ] ; then
                    wpctl set-mute -p "$pid" 0
                else
                    wpctl set-mute -p "$pid" 1
                fi
            fi
        fi
    done
done < <(xprop -spy -root _NET_ACTIVE_WINDOW)
*/
