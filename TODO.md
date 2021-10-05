# TODO

Some areas that need improvement 

## kubectl check endoflife

- [X] `endoflife` command should exit 1 or 0 based on the date, if past due exit 1
- [X] `endoflife` command should have an expiration threshold where it will exit 1 if say within 30 days of expiration
- [ ] `endoflife` command should allow version override so it doesn't require access to a cluster
- [ ] `endoflife` command should support JSON output
- [X] write simple examples of how this could be used in scripting
- [X] duplicated code in `pkg/endoflife` could be simplified

## kubectl check versions

- [ ] `versions` command should allow overriding namespaces that need to be checked
- [ ] `versions` command should allow overriding embedded config, embedded config should be sane defaults
- [ ] `versions` command should support JSON output
- [ ] `versions` command has lots of logic that could be moved into `pkg` directory