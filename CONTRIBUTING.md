# Contributing

All contributions are welcome. Features, docs, bugs, ....

We currently lack a proper COC, so the least we ask is: we all respect one another.

Questions, comments, etc. — we are happy to help. :-)

## Development

For features, bugfixes, docs etc. make a branch and open a pull-request.

For local testing of your changes:

- run `make setup` to configure the default builder and trust it.
- run `make test` to build an (app) image from `./samples/deno` with one entrypoints:
- test it: `docker run --rm -p 8080:8080 test-deno-app`

## Release

Make a release branch to adjust the version in `buildpack.toml` and tag it:

```sh
git checkout -b release-x.y.z
vi buildpack.toml # to adjust the version
git commit -a -m 'Chore: update version'
git push origin release-x.y.z
git tag -a vX.Y.Z
git push --tags
```

