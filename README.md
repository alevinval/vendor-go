# The Vending CLI

This tool allows users to install, manage and update dependencies. See the
[example](example/) to see how it works.

## Features

* Local caching of repositories: speeds up the installation and update commands
  since a full clone is not required every time.

* Highly customizable tool: develop your custom presets, adapt, and standardize the
  vendoring process to your needs.

* Parallelized vendoring process: dependencies are vendored in parallel during install
  or update operations, this further enhances speed due to the I/O required to clone or
  fetch the upstreams.

## Usage

* `vending init` initializes a `.vendor.yml` file in the working directory
* `vending add` adds a dependency in the `.vendor.yml` file
* `vending install` downloads and vendors the vendor the specified dependencies
   * The first time this command is executed, it will generate a `.vendor-lock.yml`
     which keeps track of the locked reference that has been vendored (eg. a specific commit)
   * Once the lock file already exists, it vendors dependencies at the
     specified locked reference.
* `vending update` ignores the `vendor-lock.yml` and fetches newest dependencies
   according to the refname that is specified in the `.vendor.yml` file
