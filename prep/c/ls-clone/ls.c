#define NAME_MAX 14

#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include <fcntl.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <stdlib.h>
#include <ctype.h>
#include <dirent.h>

static int has_multiple_files = 0;
static int show_hidden_files = 0;
static int long_format = 0;

void my_ls(char *);

int main(int argc, char **argv)
{
  // find all the flags
  while (--argc > 0)
    if ((*++argv)[0] == '-') {
      while (isalpha((++(*argv))[0])) {
        switch((*argv)[0]) {
          case 'a':
            show_hidden_files = 1;
            break;
          case 'l':
            long_format = 1;
            break;
          default:
            printf("unknown flag: %c\n", (*argv)[0]);
            break;
        }
      }
    }
    else {
      --argv;
      break;
    }
  ++argc;

  // print information for each directory
  if (argc == 1) {
    my_ls(".");
  } else {
    if (argc > 2) {
      has_multiple_files = 1;
    }
    while (--argc > 0) {
      ++argv;
      if (has_multiple_files) {
        printf("%s:\n", *argv);
        my_ls(*argv);
        if (argc > 1) {
          printf("\n\n");
        }
      } else {
        my_ls(*argv);
      }
    }
  }
  return 0;
}

#define MAX_PATH 1024

void my_ls(char *name)
{
  struct stat stbuf;

  if (stat(name, &stbuf) == -1) {
    fprintf(stderr, "my_ls: can't access %s", name);
    return;
  }
  // handle if it's a directory
  if ((stbuf.st_mode & S_IFMT) == S_IFDIR) {
    char *dir = name;
    char name[MAX_PATH];
    struct dirent *dp;
    DIR *dfd;

    if ((dfd = opendir(dir)) == NULL) {
      fprintf(stderr, "my_ls: can't open %s\n", dir);
      return;
    }
    while ((dp = readdir(dfd)) != NULL) {
      if (strcmp(dp->d_name, ".") == 0 || strcmp(dp->d_name, "..") == 0)
        continue;
      if (strlen(dir)+strlen(dp->d_name)+2 > sizeof(name))
        fprintf(stderr, "my_ls: name %s/%s too long \n", dir, dp->d_name);
      else {
        fprintf(stdout, "%s\t", dp->d_name);
        sprintf(name, "%s/%s", dir, dp->d_name);
      }
    }
    closedir(dfd);
  } else if ((stbuf.st_mode & S_IFMT) == S_IFREG) {
    fprintf(stdout, "%s\t", name);
  }
  // printf("%8lld %s\n", stbuf.st_size, name);
}
