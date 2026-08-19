#### Description

The `get` command retrieves one or more fields from a stored secret and prints them as a JSON object.

You rarely call this directly. It is the command the aux4 core invokes when it resolves a `secret://aux4/...` reference in a parameter — put the reference in a config file or an `.aux4` and it is resolved when the command runs.

A missing item or a missing field is an error rather than an empty result, so a mistyped reference fails loudly instead of silently supplying an empty credential.

`--otp` is accepted for compatibility with the secret provider contract, but this provider stores no one-time password seeds.

#### Usage

```bash
aux4 secret aux4 get --ref <vault/item> --fields <field1,field2> [--otp <true|false>] [--store <path>] [--encrypter <aux4|aws-kms>]
```

--ref         The secret reference, without the provider or field (required)
--fields      Comma-separated field names to retrieve (required)
--otp         Include a one-time password (not supported by this provider)
--store       Path to the encrypted store (default: `${aux4HomeDir}/secret/aux4/store.json`)
--encrypter   Encrypter protecting the values (default: `aux4`)

#### Example

```bash
aux4 secret aux4 get --ref "Work/Billing API" --fields clientId,clientSecret
```

```json
{
  "clientId": "my-client",
  "clientSecret": "s3cr3t"
}
```
