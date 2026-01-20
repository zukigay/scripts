#include <stdio.h>
#include <stdbool.h>
#include <string.h>
#include <X11/Xlib.h>
#include <X11/Xatom.h>
#include <X11/Xutil.h>
#include <unistd.h>
#include <stdlib.h>


#define MAX(A, B)               ((A) > (B) ? (A) : (B))
#define MIN(A, B)               ((A) < (B) ? (A) : (B))
#define CLAMP(a, b, c) ( MIN(MAX(b, c), MAX(MIN(b, c), a)) )
#define LENGTH(X)               (sizeof X / sizeof X[0])
#define END(A)                  ((A) + LENGTH(A))

typedef struct{
    const char * wm_name;
    const char * wm_class;
    const char * process_id;
    bool running;
} Program;

static Display *dpy;
static Atom utf8;
static Window root;

Program prog_list[] = {
    {"Slay the Spire", "Slay the Spire","SlayTheSpire",false},
};

unsigned char * get_active() {
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
    return data;
}
char * get_prop(Window win,char * prop) {
    Atom property = XInternAtom(dpy, prop, True);
    unsigned char *p_str = NULL;
    Atom da;
    /* dl and di are unused */
    unsigned long dl;
    int di;
    XGetWindowProperty(dpy, win, property, 0, BUFSIZ, False, utf8, &da, &di, &dl, &dl, &p_str);
    return p_str;
}
int get_classname(Window win,unsigned char **str) {
    XClassHint classhint;
    Status ret = XGetClassHint(dpy,win,&classhint);
    if (ret) {
        *str = (unsigned char*)classhint.res_class;
    } else {
        *str = NULL;
    }
    return ret;
}
/* todo remove all exec usage and switch to using native api's */
void mute_prog(const char * prog_name,char * arg) {
    char pgrep_arg[50];
    snprintf(pgrep_arg,sizeof pgrep_arg, "pgrep %s",prog_name);
    FILE * fp = popen(pgrep_arg,"r");
    char pid[10];
    fgets(pid,10,fp);
    pclose(fp);
    char * mute_args[] = {"wpctl","set-mute","-p",pid,arg,NULL};
    if (fork() == 0) {
        execvp(mute_args[0],mute_args);
        exit(1);
    }
}



int main() 
{
	if(!(dpy = XOpenDisplay(NULL))) {
		fputs("spirefix: cannot open display\n", stderr);
		return 1;
	}

    root = DefaultRootWindow(dpy);
	utf8 = XInternAtom(dpy, "UTF8_STRING", True);
    
    if (root == 0) {fputs("No root window\n",stderr); return 1;}
	

	XSelectInput(dpy, root, PropertyChangeMask | StructureNotifyMask);
	XEvent event;
    for (;;) {
	    fflush(stdout);
	    XNextEvent(dpy, &event);
 	    if (event.type == DestroyNotify)
 		break;
 	    if (event.type != PropertyNotify)
 		continue;
        unsigned char * d = get_active();
        Window win = ((Window*) d)[0];
        if (win != 0) {
            unsigned char * wm_class;
            if (get_classname(win,&wm_class)) {
                char * wm_name = get_prop(win,"_NET_WM_NAME");
                // printf("wm_name %s wm_class %s\n",wm_name,wm_class);
                Program * p;
                for (p = prog_list; p < END(prog_list); p++) {
                    if (strcmp(wm_name,p->wm_name) == 0 && strcmp(wm_class,p->wm_class) == 0) {
                        // printf("unmute %s\n",p->process_id);
                        mute_prog(p->process_id,"0");
                        p->running = true;
                    } else if (p->running == true) {
                        /* mute program */
                        mute_prog(p->process_id,"1");
                        // printf("mute %s\n",p->process_id);
                    }
                }
            }
        } else { /* if there is no active window mute all running */
            Program * p;
            for (p = prog_list; p < END(prog_list); p++) {
                if (p->running == true) {
                    /* mute program */
                    mute_prog(p->process_id,"1");
                    // printf("mute %s\n",p->process_id);
                }
            }
        }
        XFree(d);
    }
	XCloseDisplay(dpy);
    return 0;
}
