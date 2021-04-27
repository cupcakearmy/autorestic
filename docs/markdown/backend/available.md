# Available Backends

In theory [all the restic backends](https://restic.readthedocs.io/en/stable/030_preparing_a_new_repo.html) are supported.

Those tested are the following:

## Local

```yaml
backends:
  name-of-backend:
    type: local
    path: /data/my/backups
```

## Backblaze

```yaml
backends:
  name-of-backend:
    type: b2
    path: 'backblaze_bucketID'
    env:
      B2_ACCOUNT_ID: 'backblaze_keyID'
      B2_ACCOUNT_KEY: 'backblaze_applicationKey'
```

#### API Keys gotcha

If you use a _File name prefix_ when making the application key, do not include a leading slash. Make sure to include this prefix in the path (e.g. `path: 'backblaze_bucketID:my/path'`).

## S3 / Minio

```yaml
backends:
  name-of-backend:
    type: s3
    path: s3.amazonaws.com/bucket_name
    # Minio
    # path: http://localhost:9000/bucket_name
    env:
      AWS_ACCESS_KEY_ID: my_key
      AWS_SECRET_ACCESS_KEY: my_secret
```

## SFTP

For SFTP to work you need to use configure your host inside of ~/.ssh/config as password prompt is not supported. For more information on this topic please see the [official docs](https://restic.readthedocs.io/en/stable/030_preparing_a_new_repo.html#sftp) on the matter.

```yaml
backends:
  name-of-backend:
    type: sftp
    path: my-host:/remote/path/on/the/server
```

## Rest Server

See [here](https://github.com/restic/rest-server) for how to install a rest server backend and [here](https://restic.readthedocs.io/en/latest/030_preparing_a_new_repo.html#rest-server) for further documentation.

```yaml
backends:
  name-of-backend:
    type: rest
    path: http://localhost:8000/repo_name
    # Or authenticated
    path: https://user:pass@host:6969/path
```

Optionally you can set user and password separately

```yaml
backends:
  rest:
    type: rest
    path: http://localhost:6969/path
    key: ...
    rest:
      user: user
      password: pass
```

> :ToCPrevNext
