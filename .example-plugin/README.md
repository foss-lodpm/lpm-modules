example configuration for lpm:

```json
{
	"plugins": [
		{
			"name": "example-plugin",
			"dylib_path": "{path to libexample_plugin.so}"

		}
	]
}
```

and then run:
```sh
lpm --plugin example-plugin
```