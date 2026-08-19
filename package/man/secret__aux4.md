#### Description

The `aux4` secret provider keeps secrets in a single encrypted file on the local machine. There is no external service, no daemon and no account.

It is the portable baseline for `secret://` references. An OS keychain or a desktop keyring is stronger where it exists, but neither is available on a headless server, in a container or on a CI runner — this provider works everywhere aux4 does.

Available subcommands:

- `create` — Store a new secret
- `get` — Retrieve fields (called automatically when aux4 resolves a reference)
- `set` — Update one field
- `list` — List stored secrets as references
- `search` — Find secrets by title
- `remove` — Delete a secret

A reference is `secret://aux4/<vault>/<item>/<field>`. The vault is a namespace rather than a container: nothing is created for it and nothing needs to exist beforehand, which is what allows a reference committed to a shared config file to resolve on a machine that has never seen that vault.

#### Usage

```bash
aux4 secret aux4 <subcommand> [options]
```

#### Example

```bash
aux4 secret aux4 create --vault Work --item "Billing API" --fields "clientId=my-client,clientSecret=s3cr3t"
```

```text
secret://aux4/Work/Billing API
```
