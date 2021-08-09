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
  // handle the case with no flags or files specified
  if (argc == 1) {
    my_ls(".");
    return 0;
  }

  // find all the flags
  while (--argc > 0)
    if ((*++argv)[0] == '-') {
      while (isalpha((++(*argv))[0])) {
        switch((*argv)[0]) {
          case 'a':
            if (show_hidden_files == 0)
              printf("showing hidden files\n");
            show_hidden_files = 1;
            break;
          case 'l':
            if (long_format == 0)
              printf("using long format\n");
            long_format = 1;
            break;
          default:
            printf("unknown flag: %c\n", (*argv)[0]);
            break;
        }
      }
    }
    else {
      ++argc;
      --argv;
      break;
    }

  if (argc > 1) {
    printf("has multiple files!!\n");
    has_multiple_files = 1;
  }

  while (--argc > 0)
    my_ls(*++argv);

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
        fprintf(stderr, "dirwalk: name %s/%s too long \n", dir, dp->d_name);
      else {
        fprintf(stdout, "%s\t", dp->d_name);
        sprintf(name, "%s/%s", dir, dp->d_name);
      }
    }
    closedir(dfd);
  }
  // printf("%8lld %s\n", stbuf.st_size, name);
}
