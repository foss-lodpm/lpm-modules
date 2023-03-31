#include <iostream>

extern "C" void lpm_entrypoint(char *config_path, char *db_path,
                               unsigned int argc, char **argv) {
  std::cout << "config_path: " << config_path << std::endl;
  std::cout << "db_path: " << db_path << std::endl;

  for (int i = 0; i < argc; i++) {
    std::cout << "arg[" << i << "] " << argv[i] << std::endl;
  }
}