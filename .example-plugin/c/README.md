example build command:
```sh
cc -D_POSIX_THREAD_SAFE_FUNCTIONS -c ./example_plugin.c -o example_plugin.o
cc -shared ./example_plugin.o -o example_plugin.so
```