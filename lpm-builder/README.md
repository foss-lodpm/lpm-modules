## lpm-builder

lpm-builder is a dynamic module for the lod package manager that allows users to create LOD packages by packing software using a set of template files. This documentation explains how to use the lpm-builder module and provides details on the supported build template files.


## Building lpm-builder from Source

```sh
git clone https://github.com/foss-lodpm/lpm-modules.git
cd lpm-modules/lpm-builder
make build
```

## Update LPM configuration

To use the module on lpm, in `/etc/lpm/conf` file, append the following value in `modules` array
```json
{
    "name": "builder",
    "dylib_path": "$path_to_liblpm_builder.so"
}
```


## Template structure


### `stage0` dir

The `stage0` directory contains scripts that run on the builder module level. The following files should be included in the `stage0` directory:

- `init`: This file is used to download and verify the necessary environments and files before the building phase.
- `build`: This file is used to build the program.
- `install_files`: This file is used to install the program files into the lod package with the same paths that lpm should install them on the system.
- `post_install_files`: This file is optional and is used for any additional work that needs to be done after installing the files into the lod package.


### `stage1` dir

The stage1 directory is optional and contains scripts that run on the lpm core level. The following files can be included in the stage1 directory:

- `pre_install`: This file is an optional step that can be used if additional work needs to be done before installing the package.
- `post_install`: This file is an optional step that can be used if additional work needs to be done after installing the package.
- `pre_delete`: This file is an optional step that can be used if additional work needs to be done before deleting the package.
- `post_delete`: This file is an optional step that can be used if additional work needs to be done after deleting the package.
- `pre_downgrade`: This file is an optional step that can be used if additional work needs to be done before downgrading the package.
- `post_downgrade`: This file is an optional step that can be used if additional work needs to be done after downgrading the package.
- `pre_upgrade`: This file is an optional step that can be used if additional work needs to be done before upgrading the package.
- `post_upgrade`: This file is an optional step that can be used if additional work needs to be done after upgrading the package.


### `template` file

`JSON` file that contains almost all the informations about the package.

- `name`: The name of the package.
- `description`: A description of what the package does.
- `maintainer`: Contact information for the maintainer of the package.
- `repository`(optional): The URL of the source code repository for the package.
- `pkg_repository`: The URL of the package repository where the package can be downloaded from.
- `homepage`(optional): The URL of the package's homepage.
- `arch`: The architecture of the package.
- `kind`: The type of package.
- `file_checksum_algo`: The algorithm used to compute the checksum of the package files.
- `tags`: A list of tags or keywords that describe the package.
- `version`: An object that contains version information for the package.
- `license`(optional): The license under which the package is distributed.
- `mandatory_dependencies.runtime`(optional): A list of the package's mandatory dependencies for runtime.
- `mandatory_dependencies.build`(optional): A list of the package's mandatory dependencies for building it.
- `suggested_dependencies.runtime`(optional): A list of the package's suggested dependencies for runtime.
- `suggested_dependencies.build`(optional): A list of the package's suggested dependencies for building it.


## `stage0` built-in functions

The lpm-builder module provides the following built-in functions that can be used in `stage0` scripts:

- `validate_checksum(file_path, sha256_checksum)`: This function takes two arguments: the path of the file to validate, and the SHA256 checksum of the file. It validates the checksum of the file and throws an error if it doesn't match the provided checksum.
- `install_to_package(source_file_path, dest_path)`: This function puts files into the lod package. It takes two arguments: the path of the source file to be added to the package, and the path where lpm should install the file when the package is installed.


## Real Example of Building a Simple Package

In this example, we'll build a package for the [sbs](https://github.com/ozkanonur/sbs) tool, which is a simple background setter.

Prepare the required files for building [sbs](https://github.com/ozkanonur/sbs).

```sh
mkdir sbs_build_template
cd sbs_build_template

mkdir stage0
touch stage0/init
touch stage0/build
touch stage0/install_files

touch template
```

Copy the following content into `stage0/init`

```sh
curl -L https://github.com/ozkanonur/sbs/archive/refs/tags/v1.0.0.tar.gz > sbs.tar.gz
validate_checksum "sbs.tar.gz" "aa4da5b2315046fc2059599b19c530f08bb870e63ed17111a55991b1ae911367"
tar -xvzf sbs.tar.gz --strip 1 -C $SRC
```

Copy the following content into `stage0/build`

```sh
make sbs
```

Copy the following content into `stage0/install_files`

```sh
install_to_package sbs /usr/bin/sbs
```

Copy the following content into `template`

```sh
{
    "name": "sbs",
    "description": "Simple background setter",
    "maintainer": "Lpm Core Maintainer <contact@onurozkan.dev>",
    "repository": "https://github.com/ozkanonur/sbs",
    "pkg_repository": "https://repository.amd64.lodpm.com",
    "homepage": "https://github.com/ozkanonur/sbs",
    "arch": "amd64",
    "kind": "util",
    "file_checksum_algo": "sha256",
    "tags": [
        "x11",
        "background-setter"
    ],
    "version": {
        "readable_format": "1.0.0",
        "major": 1,
        "minor": 0,
        "patch": 0
    },
    "license": "MIT",
    "mandatory_dependencies": {
        "build": [],
        "runtime": []
    },
    "suggested_dependencies": {
        "build": [],
        "runtime": []
    }
}
```

Now, you can run `lpm --module builder --build .` which will generate `sbs.lod` package for you.
