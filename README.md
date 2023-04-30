# Naming

Rename video and subtitle files to match so players can load subtitle files without any effort on your part.

Since naming stuff is extremely inconsistent, this uses regular expressions to match.

Extremely alpha. Make a backup of your data files before you use this.

## Usage

```
    go run main.go -dir "/Users/user/Anime/Shirokuma Cafe" -dryrun=false
```

### Flags

```
-dryrun, defaults to true
-dir, defaults to empty
```