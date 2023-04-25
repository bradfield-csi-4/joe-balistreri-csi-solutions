#include <sys/stat.h>
#include <sys/types.h>
#include <stdio.h>
#include <fcntl.h>
#include <unistd.h>
#include <stdbool.h>
#include <string.h>
#include <dirent.h>


int main(int argc, char *argv[]) {
    char *cwd;
    bool lFlagEnabled;

    if (argc > 3) {
        printf("invalid args - you may specify a directory and the -l flag\n");
        return 1;
    }

    if (argc == 1) {
        cwd = getcwd(NULL, 0);
    } else if (argc == 3) {
        if (strcmp(argv[1], "-l") == 0) {
            lFlagEnabled = true;
        } else {
            printf("invalid flag: %s\n", argv[1]);
            return 1;
        }
        cwd = argv[2];
    } else {
        if (strcmp(argv[1], "-l") == 0) {
            lFlagEnabled = true;
            cwd = getcwd(NULL, 0);
        } else {
          cwd = argv[1];
        }
    }

    // printf("using cwd: %s\n", cwd);
    // printf("lflag enabled: %d\n", lFlagEnabled);

    DIR *d = opendir(cwd);
    if (d == NULL) {
        char output[500];
        sprintf(output, "failed to open dir %s", cwd);
        perror(output);
        return 1;
    }

    struct dirent *dr = readdir(d);
    struct stat s;

    for (;dr != NULL;) {
        if (lFlagEnabled) {
          char fullName[500];
          sprintf(fullName, "%s/%s", cwd, dr->d_name);

          char *username = uid_to_username(sb.st_uid);
          if (username == NULL)
          {
              return 1;
          }

          printf("The file %s is owned by %s\n", argv[1], username);
          
          free(username);
          if (stat(fullName, &s)) {
            perror("stat failed");
          }
          printf("%s%s%s%s%s%s%s%s%s - %s\n",
                 s.st_mode & S_IRUSR ? "r" : "-",
                 s.st_mode & S_IWUSR ? "w" : "-",
                 s.st_mode & S_IXUSR ? "x" : "-",
                 s.st_mode & S_IRGRP ? "r" : "-",
                 s.st_mode & S_IWGRP ? "w" : "-",
                 s.st_mode & S_IXGRP ? "x" : "-",
                 s.st_mode & S_IROTH ? "r" : "-",
                 s.st_mode & S_IWOTH ? "w" : "-",
                 s.st_mode & S_IXOTH ? "x" : "-",
                 dr->d_name);
        } else {
            printf("%s\n", dr->d_name);
        }
        // TODO: add stat syscall here if flag enabled
        dr = readdir(d);
    }

    return 0;
}

// Print out file size, number of blocks allocated, 
// reference(link) count, and so forth.What is the 
// link count of a directory, as the number of entries 
// in the directory changes ?