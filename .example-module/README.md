add module(requires root privileges):

```sh
lpm --module --add example-module {path to libexample_module.so}
# check if it's added
lpm --module --list
```

then run:
```sh
lpm --module example-module
```
