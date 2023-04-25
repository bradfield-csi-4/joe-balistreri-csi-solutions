#include <stdio.h>
#include <unistd.h>
#include <string.h>
#include <stdlib.h>

int main(int argc, char *argv[]) {
    if (argc != 2) {
        return 1;
    }

    int i, num;
    num = 0;

    for (i=0; i< strlen(argv[1]); i++) {
        num *= 10;
        num += (argv[1][i] - '0');
        printf("%c\n", argv[1][i]);
    }

    int *array = malloc(sizeof(int) * num * 1000000);

    for (;;) {
        for (i = 0; i < 1000000 * num; i++)
        {
            array[i] = i;
            printf("%d\n", i);
        }
    }
    

    return 0;
}