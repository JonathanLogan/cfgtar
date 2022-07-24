# cfgtar

Read templates from tarfile, apply configuration data, output new tar file.
  - Supports schema validation
  - Supports lookup functions
  - Dry & Pre-validation runs (optional)

For generating configuration file trees from global templates and local configuration settings.

Usage: `cat template.tar | cfgtag config.json > compiled.tar`\
Apply config.json to template.tar. Uses embedded schema only.

Usage: `cat template.tar | cfgtag schema.json config.json > compiled.tar`\
Apply additional schema.json.

Usage: `cat template.tar | cfgtag -d schema.json config.json`\
Dry run - do not produce output except for errors.

Usage: `cfgtag -v -i template.tar schema.json config.json > compiled.tar`\
Pre-validation run - first check for errors, then produce output. Additional schema is optional.

## Schema

Unless a schema.json is given on the commandline, only embedded ._config-schema.json files are considered for
schema validation. Embedded schema files are applied to the directory they are contained in and to subdirectories, unless
replaced by another embedded schema file.

Schema files are json files that define the structure and types of valid config.json files. Instead of data they contain
validation parameters.

Validation happens on the value level, currently supported are:
  - string: Must be a string.
  - int: Must be an int.
  - float: Must be a float.
  - ...more to come.

Furthermore keys can be marked as required by adding "%required" to the key name or value.

```json
{
  "key%required": [
    "int%required"
  ],
  "key1": {
    "key2": "int%required",
    "key3": "string"
  }
}
```