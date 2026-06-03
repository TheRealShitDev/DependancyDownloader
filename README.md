# Dependancy Downloader

A small Go command-line tool for downloading Maven dependency JARs from a
`pom.xml` file.

The tool reads dependencies declared in a Maven POM, builds Maven repository
download URLs, and saves the downloaded JAR files into a `dependencies/`
folder beside the selected `pom.xml`.

## Features

- Reads direct Maven dependencies from `pom.xml`
- Downloads JAR files from Maven-compatible repositories
- Falls back to repository URLs defined in the selected `pom.xml`
- Creates a local `dependencies/` output folder automatically
- Simple command-line workflow

## Requirements

- Go 1.26.4 or newer
- Internet access for downloading artifacts from configured repositories

## Usage

Run the project with a path to a Maven `pom.xml` file:

```bash
go run . path/to/pom.xml
```

By default, the tool searches repositories in this order:

1. Maven Central
2. Repositories defined in the `pom.xml`

If Maven Central does not have an artifact, the tool tries each repository
declared under the POM's `<repositories>` section.

Or build a binary first:

```bash
go build -o downloader
./downloader path/to/pom.xml
```

Downloaded files are saved here:

```text
path/to/dependencies/
```

## Example

```bash
go run . ./example-project/pom.xml
```

For each dependency, the tool prints the Maven coordinates and downloads the
matching JAR file:

```text
groupId artifactId version
```

## Current Limitations

- Maven `pom.xml` files are supported.
- Gradle files are detected, but not supported yet.
- Only dependencies with explicit `groupId`, `artifactId`, and `version` values are handled.
