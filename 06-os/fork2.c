#include <stdio.h>
#include <unistd.h>
#include <sys/wait.h>
#include <signal.h>

int main() {
    int childPid = fork();
    int childPid2 = fork();

    int pipefds[2] = {childPid, childPid2};
    pipe(pipefds);

    char *file = "/bin/ls";
    // char *arg1 = "ls";
    // char *const env[] = { "ENV1=ls", NULL };
    char *const args[] = {"/bin/ls", "-al", ".", NULL};

    if (childPid == 0)
    {
        printf("poopie\n");
    } else if childP {
        printf("peepee\n");
        int time = 1;
        waitpid(childPid, &time, 0);
        printf("dingdong\n");
    }
    return 0;
}