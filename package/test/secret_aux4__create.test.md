# secret aux4 create

## storing a secret

```afterAll
rm -rf .secret
```

### should print the reference to the new secret

```execute
aux4 secret aux4 create --vault Work --item "Billing API" --fields "clientId=my-client,clientSecret=s3cr3t" --store .secret/store.json
```

```expect
secret://aux4/Work/Billing API
```

### should never write a value in the clear

The whole point of the store is that its contents are not readable, so the
plaintext must not appear anywhere in the file.

```execute
aux4 secret aux4 create --vault Work --item Deploy --fields "token=plaintext-canary" --store .secret/store.json > /dev/null
grep -c "plaintext-canary" .secret/store.json || true
```

```expect
0
```

### should keep a value that contains an equals sign intact

Only the first `=` separates the name from the value, so base64 padding and
connection strings survive.

```execute
aux4 secret aux4 create --vault Work --item Jira --fields "token=abc=def==" --store .secret/store.json > /dev/null
aux4 secret aux4 get --ref Work/Jira --fields token --store .secret/store.json
```

```expect:json
{
  "token": "abc=def=="
}
```

### should store several fields on one item

```execute
aux4 secret aux4 create --vault Personal --item GitHub --fields "username=john,password=hunter2" --store .secret/store.json > /dev/null
aux4 secret aux4 get --ref Personal/GitHub --fields username,password --store .secret/store.json
```

```expect:json
{
  "password": "hunter2",
  "username": "john"
}
```

### should refuse to overwrite an existing item

```execute
aux4 secret aux4 create --vault Work --item Once --fields "a=b" --store .secret/store.json > /dev/null
aux4 secret aux4 create --vault Work --item Once --fields "a=c" --store .secret/store.json
```

```error:partial
secret://aux4/Work/Once already exists, use set to change a field
```

### should reject a vault containing a slash

```execute
aux4 secret aux4 create --vault "Work/Sub" --item Thing --fields "a=b" --store .secret/store.json
```

```error:partial
must not contain '/'
```

### should require at least one field

```execute
aux4 secret aux4 create --vault Work --item Empty --fields "" --store .secret/store.json
```

```error:partial
at least one field is required
```

### should reject a malformed field assignment

```execute
aux4 secret aux4 create --vault Work --item Bad --fields "justaname" --store .secret/store.json
```

```error:partial
must be in the form name=value
```
