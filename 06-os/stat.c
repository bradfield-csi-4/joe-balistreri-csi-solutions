#include <sys/stat.h>
#include <stdio.h>
#include <fcntl.h>
#include <unistd.h>


int main(int argc, char *argv[]) {
    if (argc != 2) {
        printf("need an argument!\n");
        return 1;
    }

    struct stat s;

    if (stat(argv[1], &s) != 0) {
        printf("failed to call stat");
        return 1;
    }

    printf("size: %lld\n", s.st_size);
    printf("device no: %d\n", s.st_dev);
    printf("block count: %lld\n", s.st_blocks);
    printf("link count: %d\n", s.st_nlink);

    char buf[4096];

    int fd = open(argv[1], O_RDONLY);
    printf("got fd: %d\n", fd);

    ssize_t n = read(fd, &buf, 4096);
    if (n == -1) {
        perror("fudged up");
        return 1;
    }

    printf("read %zd bytes\n", n);
    write(1, &buf, n);
    
    return 0;
}

// Print out file size, number of blocks allocated, 
// reference(link) count, and so forth.What is the 
// link count of a directory, as the number of entries 
// in the directory changes ?