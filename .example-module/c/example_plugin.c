#include <stdio.h>
#include <stdlib.h>

void lpm_entrypoint(char *db_path, unsigned int argc,
                    char **argv) {
  printf("db_path: %s\n", db_path);

  for (int i = 0; i < argc; i++) {
    printf("arg[%i] %s\n", i, argv[i]);
  }
}