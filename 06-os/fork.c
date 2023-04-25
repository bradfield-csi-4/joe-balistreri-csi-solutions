#include <stdio.h>
#include <unistd.h>
#include <sys/wait.h>

int main() {
    printf("the pid of parent is %d\n", (int)getpid());
    int pid = fork();

    if (pid == 0) {
        char *argv[3] = {"Command-line", ".", NULL};
        execvp("find", argv);
        printf("after fork, the pid of child is %d\n", (int)getpid());
    } else {
        int duration = 2;
        wait(&duration);
        printf("after fork, the pid of parent is %d, returned %d\n", (int)getpid(), pid);
    }
    return 0;
}