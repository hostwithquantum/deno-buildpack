# deno-buildpack

A deno buildpack (for [runway](https://runway.planetary-quantum.com/)).

> Deno is an alternative JavaScript runtime for the server. This buildpacks generates a Docker/OCI image from your application code and includes the correct version of deno to run the server with.

## Configuration

- `BP_RUNWAY_DENO_VERSION=v1.25.1`
- `BP_RUNWAY_DENO_FILE_VERSION=runtime.txt`
- `BP_RUNWAY_DENO_PERM_ENV=PORT`
- `BP_RUNWAY_DENO_PERM_HRTIME=false`
- `BP_RUNWAY_DENO_PERM_NET=true`
- `BP_RUNWAY_DENO_PERM_FFI=false`
- `BP_RUNWAY_DENO_PERM_READ=true`
- `BP_RUNWAY_DENO_PERM_RUN=false`
- `BP_RUNWAY_DENO_PERM_WRITE=false`
- `BP_RUNWAY_DENO_PERM_ALL=false`
- `BP_RUNWAY_DENO_MAIN=server.ts`

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

##### BP_RUNWAY_DENO_FILE_VERSION

`runtime.txt` should contain a version, such as `vA.B.C`.

> (last) The buildpack also supports a `.dvmrc` file.

## Development

Run `make setup` to configure the default builder and trust it.

Run `make test` to build an (app) image from `./samples/deno` with one entrypoints:

- web: `docker run --rm -p 8080:8080 test-deno-app`
