# aux4/secret-aux4 0.0.1

First release.

## Local encrypted secret provider

A `aux4/secret` provider that keeps secrets in a single encrypted file on the machine — no external service, no daemon, no account.

It is the portable baseline for `secret://` references. An OS keychain or a desktop keyring is stronger where it exists, but neither is available on a headless server, in a container, or on a CI runner. This one works everywhere aux4 does.

```bash
aux4 secret aux4 create --vault Work --item "Billing API" --fields "clientId=my-client,clientSecret=s3cr3t"
# -> secret://aux4/Work/Billing API
```

## Commands

`create`, `get`, `set`, `list`, `search`, `remove` — the full `aux4/secret` provider contract plus `remove`, which a local store needs in order to forget anything.

## Encryption

Values are encrypted through `aux4/encrypter` rather than reimplemented here, so `--encrypter aws-kms` protects the same store with a managed key instead of the local one.

Only values are encrypted. Vault names, item titles and field names stay in the clear so the store can be listed and searched without unlocking it.

The key comes from `AUX4_SECRET_KEY` when set, otherwise from a key file beside the store, generated on first write. A key file next to the store is not meaningfully stronger than a plaintext file against anything running as the same user; it protects the store in backups, in synced home directories, and from casual inspection. Set `AUX4_SECRET_KEY` from elsewhere for more.

Neither the key nor the secret is ever passed as a command-line argument — both travel by stdin or environment, so neither appears in `ps`.

## Vaults are namespaces

Nothing is created on disk for a vault and nothing has to exist beforehand. That is what lets a `secret://aux4/<vault>/<item>/<field>` reference committed to a shared config file resolve on a machine that has never seen that vault.
