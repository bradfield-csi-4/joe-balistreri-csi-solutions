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

void print_target(char *);
struct stat get_stat(char *);

void my_ls(char *name)
{
  struct stat stbuf = get_stat(name);
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
      if (dp->d_name[0] == '.' && !show_hidden_files)
        continue;
      if (strlen(dir)+strlen(dp->d_name)+2 > sizeof(name))
        fprintf(stderr, "my_ls: name %s/%s too long \n", dir, dp->d_name);
      else {
        // TODO: Fix the issue where we can't access things outside the current directory
        sprintf(name, "%s/%s", dir, dp->d_name);
        print_target(name);
        print_target(dp->d_name);
      }
    }
    closedir(dfd);
  } else if ((stbuf.st_mode & S_IFMT) == S_IFREG) {
    fprintf(stdout, "%s\t", name);
  }
  // printf("%8lld %s\n", stbuf.st_size, name);
}

struct stat get_stat(char *name) {
  struct stat stbuf;
  printf("getting stat for %s\n", name);

  if (stat(name, &stbuf) == -1) {
    fprintf(stderr, "my_ls: can't access %s\n", name);
    return stbuf;
  }
  return stbuf;
}

void red() {
  printf("\033[0;31m");
}

void blue() {
  printf("\033[0;36m");
}

void reset() {
  printf("\033[0m");
}

void print_target(char *name) {
  struct stat stbuf = get_stat(name);
  if ((stbuf.st_mode & S_IFMT) == S_IFDIR) {
    blue();
  }
  // TODO: handle if executable
  fprintf(stdout, "%s\t", name);
  reset();
}
