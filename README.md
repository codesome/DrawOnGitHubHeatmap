# DrawOnGitHubHeatmap
Draw on your GitHub heatmap

Live demo can be found here: [https://github.com/codesomeExperimental](https://github.com/codesomeExperimental) (will be update occasionally)

**WARNING: Use this only on a dummy repo, all the commits will be deleted when this is run.**

## How this works

* Deletes the master branch and creates another with no commits.
* Considering 52 weeks (width) in heatmap and the height to be 7 days a week (height), make commits on appropriate dates using `--date` flag in `git`.
* Highly customisable, not just limited to letters. You can map letter to some drawing and use this; you will know how when you see the config example.

## How to use

### Download

```bash
$ git clone https://github.com/codesome/DrawOnGitHubHeatmap.git
$ cd DrawOnGitHubHeatmap
```

### Build

```bash
# Executable will be created with the name DrawOnGitHubHeatmap.
$ make build # Tested with go1.11.2
```

### Help

```bash
$ ./DrawOnGitHubHeatmap -h
```

### Preparing config

It is a YAML config, and follows this structure:

```yaml
# The starting point

width: <int>

characters:
  [ - <character> ]
```

```yaml
# <character>

- char: <character> # Case sensitive.
  layout: # There should be exactly 7 strings (height), and length of the string should match the width above.
    [ - <string> ]
```

The string inside the layout can consist only `.` or `*`
* `.` means it's an empty space
* `*` means it's filled.

You can get a good idea about the config from the prefilled configs [here](https://github.com/codesome/DrawOnGitHubHeatmap/tree/master/prefilledConfigs) and below example.

Example:

```yaml
# myPixelConfig.yaml

width: 5 # If you set this to 52, you can map a letter to a drawing on entire heatmap!

characters:
  - char: H
    layout:
      - "*...*"
      - "*...*"
      - "*...*"
      - "*****"
      - "*...*"
      - "*...*"
      - "*...*"
  - char: I
    layout:
      - "*****"
      - "..*.."
      - "..*.."
      - "..*.."
      - "..*.."
      - "..*.."
      - "*****"
```

### Run

Create an empty (dummy) repository in GitHub (preferably private, to avoid distrubing public stats). Let's assume the repo is `foo/bar`. (If you want to try this on existing **dummy* project, see the Gotchas section.)

`cd` into the repo and run the following command

```bash
$ /path/to/DrawOnGitHubHeatmap --pixel-layout-config=<path_to_config> --commits-per-day=<int> --text=<string>
```

* `--pixel-layout-config`: Path to the config file to use.
* `--commits-per-day`: Number of commits on each day on heatmap. This should be `>=` the number of commits on darkest spot of your heatmap.
* `--text`: The text to write on heatmap. This should not be too long. Do the math of letter width and 52 weeks (the width of heatmap), and choose accordingly.

Example
```bash
$ /path/to/DrawOnGitHubHeatmap --pixel-layout-config=/path/to/myPixelConfig.yaml --commits-per-day=5 --text="HI"
```

## Gotchas

* If you want to re-run this on your dummy repo, or run this on an existing repo, you need to delete the repo and and create it again before using.
    * Because, even if you force push by deleting old commits, they don't dissapear from the heatmap.

## TODO

* Prefilled config file for all alphabets and numbers.
* Cleanup code.
* Comments in the code.