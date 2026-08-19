#### Description

The `create` command stores a new secret and prints the `secret://` reference for it, ready to paste into a config file or an `.aux4`.

Every field value is encrypted before it reaches disk. Field values may contain `=` — only the first one separates the name from the value, so base64 padding and connection strings are stored intact.

Creating an item that already exists is an error rather than an overwrite, since silently replacing a secret would discard fields that are not mentioned. Use `set` to change a field of an existing item.

The vault must not contain `/`, which always separates the vault from the item in a reference. The item may contain `/`.

#### Usage

```bash
aux4 secret aux4 create --vault <name> --item <title> --fields <key=value,...> [--category <type>] [--store <path>] [--encrypter <aux4|aws-kms>]
```

--vault       The vault namespace to create the secret in (required)
--item        The item title (required)
--fields      Comma-separated `key=value` assignments (required)
--category    The item category (default: `Login`)
--store       Path to the encrypted store (default: `${aux4HomeDir}/secret/aux4/store.json`)
--encrypter   Encrypter protecting the values (default: `aux4`)

#### Example

```bash
aux4 secret aux4 create --vault Work --item "Billing API" --fields "clientId=my-client,clientSecret=s3cr3t"
```

```text
secret://aux4/Work/Billing API
```

Protecting the values with a managed key service instead of the local key:

```bash
aux4 secret aux4 create --vault Work --item Deploy --fields "token=abc123" --encrypter aws-kms
```
