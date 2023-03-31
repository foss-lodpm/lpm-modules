example configuration for lpm:

```json
{
	"modules": [
		{
			"name": "example-module",
			"dylib_path": "{path to libexample_module.so}"

		}
	]
}
```

then run:
```sh
lpm --module example-module
```
