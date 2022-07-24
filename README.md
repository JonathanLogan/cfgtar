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
  - string(min=int,max=int,len=int): Must be a string. min length, max length, precise length. 
  - int(min=int,max=int): Must be an int. Minimum/maximum value.
  - float(min=int,max=int): Must be a float. Minimum/maximum value.
  - dir: Directory must exist.
  - file: File must exist.
  - duration(min=intSeconds,max=intSeconds). Duration. min/max value in seconds.
  - hex(min=int,max=int,len=int). Is hexadecimal encoded. min length, max length, precise length.
  - base64(min=int,max=int,len=int). Is hexadecimal encoded. min length, max length, precise length.
  - base58(min=int,max=int,len=int). Is hexadecimal encoded. min length, max length, precise length.
  - ipv4. Is IPv4 address.
  - ipv4net. Is IPv4 network.
  - ipv6. Is IPv6 address.
  - ipv6net. Is IPv6 network.
  - hostname. Value must match local hostname.
  - nic. Network interface name must exist.
  - nic4. Network interface name must exist and have an ipv4 addr.
  - nic6. Network interface name must exist and have an ipv6 addr.
  - lookup4. DNS resolves to ipv4 address.
  - lookup6. DNS resolves to ipv6 address.
  - ...more to come.

Furthermore keys can be marked as required by adding "%required" to the key name or value.

```json
{
  "key%required": [
    "int%required"
  ],
  "key1": {
    "key2": "int(min=1,max=30)%required",
    "key3": "string"
  }
}
```

## Template

Input templates are go text/template. Additional functions are provided:
  - durationAs [<h/m/s>] dur: Format duration as hour/minute/second (default).
  - file (get content)
  - hostname: Returns the hostname of the local machine.
  - ipv4CIDR network: Return the CIDR for network. IPv4 version.
  - ipv4Mask network: Return the netmask for network. IPv4 version.
  - ipv6CIDR network: Return the CIDR for network. IPv6 version.
  - ipv6Mask network: Return the netmask for network. IPv6 version.
  - ipv4NICAddr nic: Returns list of ipv4 addresses of network interface nic. 
  - ipv6NICAddr nic: Returns list of ipv6 addresses of network interface nic.
  - ipv4addr addr/net: Return IP address. IPv4 version.
  - ipv6addr addr/net: Return IP address. IPv6 version.
  - ipv4lookup hostname: Lookup IP addresses of hostname. IPv4 version.
  - ipv6lookup hostname: Lookup IP addresses of hostname. IPv6 version.
  - dnsTXT name: Lookupt TXT records for name.

ToDo:
  - ipv4Addr pos[first,last,_current_] network 
  - ipv6Addr pos[first,last,_current_] network
  