# cfgtar
Configuration templates in TAR files.

Usage: `cat template.tar | cfgtag configfile > compiled.tar`

The purpose of cfgtar is to create a set of local configuration files from a global template that is enriched with
local configuration settings.

## Configfile format

`name (type)=value`

Name is a hierarchy supporting "." dot separation (`level1.level2.name`). Arrays are supported as well (`level1.level2[1].name`).

Type is extendable type enforcement, currently supported: string, int.

Value is a string value (to be converted to type), which can optionally be quoted (`="value"`) and supports backslash escaping (`="value\"value"`).

Value can also be a variable name `$level1.level2.name` which refers to a previously defined name=value pair.

Comments start with `#`.

Allows multiline composition:

```
level1.level2.
.name (type)=value1
.name2 (type)=value2
```

is equal to

```
level1.level2.name (type)=value1
level1.level2.name2 (type)=value2
```

## Template format

Each file contained in the input tar stream is considered a go text/template. The configfile will be applied to that
template and the result written to the output stream with all other tar parameters left untouched.