# README

### Run Forgelogs

See manual:

```console
forgelogs --help
```

Sample commands:

1. Failed outcomes on MyLoginJourney tree
```console
forgelogs -s am -b 2023-10-01T00:00:00Z -d 2d -t MyLoginJourney --failed-only
```

2. Failed outcomes on MyLoginJourney tree with detailed reason
```console
forgelogs -s am -b 2023-10-01T00:00:00Z -d 2d -t MyLoginJourney --failed-only --detailed
```

3. Failed outcomes on all trees
```console
forgelogs -s am -b 2023-10-01T00:00:00Z -d 2d --all-trees --failed-only
```

4. Failed outcomes on specific trees listed in `config.json`
```console
forgelogs -s am -b 2023-10-01T00:00:00Z -d 2d --filter-trees --failed-only
```

5. All outcomes on all trees
```console
forgelogs -s am -b 2023-10-01T00:00:00Z -d 2d --all-trees
```

6. No report, just raw logs
```console
forgelogs -s am -b 2023-10-01T00:00:00Z -d 2d
```

### Try it out

Download in the **Releases** section, or click [here](https://github.com/constantiux/forgelogs/releases).

### Build this project

**Note**: This is optional.

```console
go build
```

