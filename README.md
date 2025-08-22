# deno-buildpack

A deno buildpack (for [runway](https://runway.planetary-quantum.com/)).

> Deno is an alternative JavaScript runtime for the server. This buildpacks generates a Docker/OCI image from your application code and includes the correct version of deno to run the server with.

## Configuration

- `BP_RUNWAY_DENO_VERSION=`
- `BP_RUNWAY_DENO_FILE_VERSION=runtime.txt`
- `BP_RUNWAY_DENO_PERM_ENV=PORT`
- `BP_RUNWAY_DENO_PERM_HRTIME=false`
- `BP_RUNWAY_DENO_PERM_NET=true`
- `BP_RUNWAY_DENO_PERM_FFI=false`
- `BP_RUNWAY_DENO_PERM_READ=true`
- `BP_RUNWAY_DENO_PERM_RUN=false`
- `BP_RUNWAY_DENO_PERM_WRITE=false`
- `BP_RUNWAY_DENO_PERM_ALL=false`
- `BP_RUNWAY_DENO_MAIN=main.ts,server.ts`
- `BP_RUNWAY_DENO_SERVE=false`

Supported permissions:

> --allow-env=<allow-env>
> --allow-hrtime
> --allow-net=<allow-net>
> --allow-ffi
> --allow-read=<allow-read>
> --allow-run=<allow-run>
> --allow-write=<allow-write>
> -A, --allow-all

### Environment variables

Configuration is done through environment variables.

Permissions can be generally enabled with a `true` value. So for example:

```sh
export BP_RUNWAY_DENO_PERM_NET=true
```

The above allows all net access, but it could be more granular with:

```sh
export BP_RUNWAY_DENO_PERM_NET=github.com:443
```

#### Deno version

Order of priority:

##### BP_RUNWAY_DENO_VERSION

Supersedes `BP_RUNWAY_DENO_FILE_VERSION`.

Contains a version such as `vA.B.C`.

> [!IMPORTANT]
> Tested with Deno 1.x and 2.x

##### BP_RUNWAY_DENO_FILE_VERSION

`runtime.txt` should contain a version, such as `vA.B.C`.

> (last) The buildpack also supports a `.dvmrc` file.

## Contributions?

See [CONTRIBUTING.md](CONTRIBUTING.md) for everything (local setup, testing, releasing).
