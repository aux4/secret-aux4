#### Description

The `set` command updates a single field of an existing secret, encrypting the new value before writing it. Other fields on the item are untouched.

The item must already exist. Creating one here would let a mistyped reference produce a half-populated item that fails much later, at the point something tries to resolve it — use `create` first.

#### Usage

```bash
aux4 secret aux4 set --ref <vault/item> --field <name> --value <value> [--store <path>] [--encrypter <aux4|aws-kms>]
```

--ref         The secret reference, without the provider or field (required)
--field       The field name to update (required)
--value       The new value (required)
--store       Path to the encrypted store (default: `${aux4HomeDir}/secret/aux4/store.json`)
--encrypter   Encrypter protecting the values (default: `aux4`)

#### Example

```bash
aux4 secret aux4 set --ref "Work/Billing API" --field clientSecret --value new-s3cr3t
```

```text
secret://aux4/Work/Billing API/clientSecret updated
```
