# secret aux4 get

## reading fields

```beforeAll
aux4 secret aux4 create --vault Work --item "Billing API" --fields "clientId=my-client,clientSecret=s3cr3t" --store .secret/store.json > /dev/null 2>&1
```

```afterAll
rm -rf .secret
```

### should return the requested fields as json

```execute
aux4 secret aux4 get --ref "Work/Billing API" --fields clientId,clientSecret --store .secret/store.json
```

```expect:json
{
  "clientId": "my-client",
  "clientSecret": "s3cr3t"
}
```

### should return a single field

```execute
aux4 secret aux4 get --ref "Work/Billing API" --fields clientSecret --store .secret/store.json
```

```expect:json
{
  "clientSecret": "s3cr3t"
}
```

### should fail on an unknown item

A mistyped reference has to fail loudly; resolving to nothing would hand the
caller an empty credential and fail somewhere far less obvious.

```execute
aux4 secret aux4 get --ref Work/Nope --fields token --store .secret/store.json
```

```error:partial
no secret found at secret://aux4/Work/Nope
```

### should fail on an unknown field

```execute
aux4 secret aux4 get --ref "Work/Billing API" --fields nosuchfield --store .secret/store.json
```

```error:partial
field "nosuchfield" not found in secret://aux4/Work/Billing API
```

### should reject a reference without an item

```execute
aux4 secret aux4 get --ref Work --fields token --store .secret/store.json
```

```error:partial
must be in the form <vault>/<item>
```
