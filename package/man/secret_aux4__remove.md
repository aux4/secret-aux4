#### Description

The `remove` command deletes a secret and every field on it. The vault disappears from listings once its last item is removed, since a vault is only a namespace and has no existence of its own.

Removing a secret does not revoke it wherever it is used — it only forgets it locally.

#### Usage

```bash
aux4 secret aux4 remove --ref <vault/item> [--store <path>]
```

--ref     The secret reference, without the provider or field (required)
--store   Path to the encrypted store (default: `${aux4HomeDir}/secret/aux4/store.json`)

#### Example

```bash
aux4 secret aux4 remove --ref "Work/Billing API"
```

```text
secret://aux4/Work/Billing API removed
```
