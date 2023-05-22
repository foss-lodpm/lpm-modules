add module(requires root privileges):

```sh
lpm --add-module example-module {path to libexample_module.so}
# check if it's added
lpm --modules
```

then run:
```sh
lpm --module example-module
```
