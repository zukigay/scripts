#include <unistd.h>
#include <fcntl.h>
char * zero_str = "0";
int fd;

int main() {
    fd = open("/sys/module/hid_apple/parameters/fnmode",O_WRONLY);
    if (fd == -1) {
        write(1,"a error happend in opening file\n",33);
        return 1;
    }
    if (!write(fd,zero_str,1)) {
        write(1,"a error happend in writting to file\n",35);
        return 1;
    }
}
