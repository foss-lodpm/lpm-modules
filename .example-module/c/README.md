example build command:
```sh
cc -D_POSIX_THREAD_SAFE_FUNCTIONS -c ./example_module.c -o example_module.o
cc -shared ./example_module.o -o libexample_module.so
```