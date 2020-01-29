# Fileagebeat

Welcome to Fileagebeat.

Ensure that this folder is at the following location:
`${GOPATH}/src/github.com/Ardiea/fileagebeat`

## Getting Started with Fileagebeat

### Requirements

* [Golang](https://golang.org/dl/) 1.7

### Init Project
To get running with Fileagebeat and also install the
dependencies, run the following command:

```
make setup
```

It will create a clean git history for each major step. Note that you can always rewrite the history if you wish before pushing your changes.

To push Fileagebeat in the git repository, run the following commands:

```
git remote set-url origin https://github.com/Ardiea/fileagebeat
git push origin master
```

For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).

### Build

To build the binary for Fileagebeat run the command below. This will generate a binary
in the same directory with the name fileagebeat.

```
make
```

### Configure

The primary configuration element is a list of input structures as described below. 

onfiguration Element | Type | Description | Required? | Default Value |
|-----------------------|-----------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-----------|----------------|
| `fields:` | map | A map of fields that will be added under `fields:` to every result and document the beat produces  | No | Nil |
| `inputs:` | List | List of configuration inputs | Yes | Not applicable |
| `  name:` | String | A unique name for the input. | Yes | Not applicable |
| `  paths:` | List of Strings | A list of paths for this input to check. | Yes | Not Applicable |
| `  disable:` | Bool | A boolean for to enable or disable the input. Default is false. | No | `false` |
| `  period:` | Golang Duration | How often to check the paths. Defaults to 60 seconds | No | `60s` |
| `  threshold:` | Golang Duration | A period of time that if the file's age exceeds it it will be considered to be an aging file. | No | `60s` |
| `  whiteslist:` | List of Strings | A list of regular expressions that will be tested against discovered filenames. Any file names found that match one of these regular exprssions will be included for consideration in age testing. Whitelist and blacklist are mutually exclusive. | No | Empty List |
| `  blacklist:` | List of Strings | A list of regular expressions that will be tested against discovered filenames. Any file names found that match one of these regular expressions will be excluded from consideration in age testing. Whitelist and blacklist are mutually exclusive. | No | Empty List |
| `  max_depth:` | Integer | A restriction on how deeply into the directory structure of each path to descend. 0 means no restriction. 1 means do not descend into any sub-directories. To prevent unrestricted recursion, 0 == 128. To go deeper than 128 directories specify a value > 128.  | No | 128 |
| `  attribute:` | String | Specifies which time attribute to use in age testing. Valid options are `mtime`, `ctime`, and `atime`.  | No | `mtime` |
| `  heartbeat:` | Bool | Enables a heartbeat message to be sent to the outputs every `period`. This is useful to know if the monitor is still running. | No | `false` |

### Run

To run Fileagebeat with debugging output enabled, run:

```
./fileagebeat -c fileagebeat.yml -e -d "*"
```


### Test

To test Fileagebeat, run the following command:

```
make testsuite
```

alternatively:
```
make unit-tests
make system-tests
make integration-tests
make coverage-report
```

The test coverage is reported in the folder `./build/coverage/`

### Update

Each beat has a template for the mapping in elasticsearch and a documentation for the fields
which is automatically generated based on `fields.yml` by running the following command.

```
make update
```


### Cleanup

To clean  Fileagebeat source code, run the following command:

```
make fmt
```

To clean up the build directory and generated artifacts, run:

```
make clean
```


### Clone

To clone Fileagebeat from the git repository, run the following commands:

```
mkdir -p ${GOPATH}/src/github.com/Ardiea/fileagebeat
git clone https://github.com/Ardiea/fileagebeat ${GOPATH}/src/github.com/Ardiea/fileagebeat
```


For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).


## Packaging

The beat frameworks provides tools to crosscompile and package your beat for different platforms. This requires [docker](https://www.docker.com/) and vendoring as described above. To build packages of your beat, run the following command:

```
make release
```

This will fetch and create all images required for the build process. The whole process to finish can take several minutes.
