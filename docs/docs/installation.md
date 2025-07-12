---
hide:
  - toc
---

mkrgen provides a convenient CLI tool to effortlessly set up your Go projects. Follow the steps below to install the tool on your system.

## Binary Installation

To install the mkrgen CLI tool as a binary, run the following command:

```sh
go install github.com/noxymon-mekari/mkrgen@latest
```

This command installs the mkrgen binary, automatically binding it to your `$GOPATH`.

> If you’re using Zsh, you’ll need to add it manually to `~/.zshrc`.

> After running the installation command, you need to update your `PATH` environment variable. To do this, you need to find out the correct `GOPATH` for your system. You can do this by running the following command:
> Check your `GOPATH`
>
> ```
> go env GOPATH
> ```
>
> Then, add the following line to your `~/.zshrc` file:
>
> ```
> GOPATH=$HOME/go  PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
> ```
>
> Save the changes to your `~/.zshrc` file by running the following command:
>
> ```
> source ~/.zshrc
> ```

## Building and Installing from Source

If you prefer to build and install mkrgen directly from the source code, you can follow these steps:

Clone the mkrgen repository from GitHub:

```sh
git clone https://github.com/noxymon-mekari/mkrgen
```

Build the mkrgen binary:

```sh
go build
```

Install in your `$PATH` to make it accessible system-wide:

```sh
go install
```

Verify the installation by running:

```sh
mkrgen version
```

This should display the version information of the installed mkrgen.

Now you have successfully built and installed mkrgen from the source code.
