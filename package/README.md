# aux4/secret-aux4

Local secret provider for `aux4/secret`. Secrets are kept in a single encrypted file on the machine, with no external service, no daemon, and no account.

It is the portable baseline for `secret://` references: it works on every platform aux4 runs on, including headless servers, containers and CI, where an OS keychain or a desktop keyring is unavailable.

## Installation

```bash
aux4 aux4 pkger install aux4/secret-aux4
```

## Quick Start

```bash
# Store a secret
aux4 secret aux4 create --vault Work --item "Billing API" --fields "clientId=my-client,clientSecret=s3cr3t"
# -> secret://aux4/Work/Billing API

# Reference it anywhere aux4 resolves parameters
aux4 curl request https://api.example.com --header "Authorization: Bearer secret://aux4/Work/Billing API/clientSecret"
```

You rarely call `get` yourself. Put the `secret://` reference in a config file or an `.aux4`, and aux4 resolves it when the command runs.

## How References Work

A reference is `secret://aux4/<vault>/<item>/<field>`.

The vault is a **namespace, not a container**. Nothing is created on disk for it and nothing has to exist beforehand, which is what lets a reference committed to a shared config file resolve on a machine that has never seen that vault.

```text
secret://aux4/Work/Billing API/clientSecret
          │    │     │           └─ field
          │    │     └───────────── item
          │    └─────────────────── vault
          └──────────────────────── provider
```

The item may contain `/`; the vault may not, since the first segment always separates them.

## Commands

### `aux4 secret aux4 create`

Create a secret. Prints the reference to paste into a config file.

```bash
aux4 secret aux4 create --vault <name> --item <title> --fields <key=value,...> [--category <type>] [--store <path>] [--encrypter <name>]
```

| Flag | Description | Default |
|------|-------------|---------|
| `--vault` | Vault namespace | required |
| `--item` | Item title | required |
| `--fields` | Comma-separated `key=value` pairs | required |
| `--category` | Item category | `Login` |
| `--store` | Path to the store | `${aux4HomeDir}/secret/aux4/store.json` |
| `--encrypter` | Encrypter protecting the values | `aux4` |

Field values may contain `=` — only the first one separates name from value, so base64 padding and connection strings survive intact.

Creating an item that already exists is an error; use `set` to change a field.

### `aux4 secret aux4 get`

Retrieve fields as a JSON object. Called automatically when aux4 resolves a `secret://` reference.

```bash
aux4 secret aux4 get --ref <vault/item> --fields <field1,field2>
```

```json
{
  "clientId": "my-client",
  "clientSecret": "s3cr3t"
}
```

### `aux4 secret aux4 set`

Update one field of an existing secret.

```bash
aux4 secret aux4 set --ref <vault/item> --field <name> --value <value>
```

### `aux4 secret aux4 list`

List stored secrets as references, ready to copy into a config file.

```bash
aux4 secret aux4 list [--vault <name>] [--withFields <true|false>]
```

```text
secret://aux4/Personal/GitHub
secret://aux4/Work/Billing API
```

With `--withFields true`, one line per field.

### `aux4 secret aux4 search`

Search items by title, case-insensitively.

```bash
aux4 secret aux4 search <query> [--vault <name>] [--withFields <true|false>]
```

### `aux4 secret aux4 remove`

Delete a secret and all of its fields.

```bash
aux4 secret aux4 remove --ref <vault/item>
```

## Encryption

Values are encrypted through `aux4/encrypter`, which this package does not reimplement. `--encrypter` selects the provider, so the same store can be protected by local AES-256-GCM or by a managed key service:

```bash
aux4 secret aux4 create --vault Work --item Deploy --fields "token=abc" --encrypter aws-kms
```

Only the values are encrypted. Vault names, item titles and field names are stored in the clear so the store can be listed and searched without unlocking it — do not put anything sensitive in a name.

## The Key

The key is resolved in this order:

1. `AUX4_SECRET_KEY`, when set
2. a key file beside the store, generated on first write

**Note:** a key file sitting next to the store is not meaningfully stronger than a plaintext file against anything running as your user — it can read both. What it does protect against is broader exposure: backups, a synced home directory, and casual inspection of the file. For protection against local access, supply `AUX4_SECRET_KEY` from somewhere the store cannot reach, such as a CI secret or a hardware-backed agent.

Both the store and the key file are written with `0600` permissions in a `0700` directory. The store is written atomically, so an interrupted write cannot truncate it.

## Environment Variables

| Variable | Description |
|----------|-------------|
| `AUX4_SECRET_KEY` | Key material protecting the store. Overrides the key file |

## Limitations

- **No one-time passwords.** `--otp` is accepted for contract compatibility but this provider stores no TOTP seeds.
- **Names are not encrypted.** See Encryption above.
