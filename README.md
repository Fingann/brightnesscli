# BrightnessCLI

BrightnessCLI is a command-line interface (CLI) tool for controlling your computer's brightness. It's written in Go and uses the Cobra library for CLI interactions. The tool requires root privileges to run, unless you install the provided udev rules.

## Features

- Increase brightness by 5% (`up` command)
- Decrease brightness by 5% (`down` command)
- Set brightness to a specific value (`set` command)
- Get the current brightness (`get` command)
- Install brightness rules (`install-rules` command)

## Installation

The preferred way to install BrightnessCLI is by downloading the latest binary from the [Releases](https://github.com/fingann/brightnesscli/releases) page.

1. Download the latest release.
2. Extract the binary from the downloaded file.
3. Move the binary to a directory in your `PATH`, for example `/usr/local/bin`.

For example, if you downloaded the binary to your `Downloads` directory and it's named `brightnesscli`, you can install it using the following commands:

```bash
chmod +x ~/Downloads/brightnesscli
sudo mv ~/Downloads/brightnesscli /usr/local/bin
```
## Installing Rules to Run Without Sudo

BrightnessCLI requires root privileges to control the system's brightness. However, you can install udev rules to allow BrightnessCLI to run without sudo. This can be done using the `install-rules` command. This command installs the necessary udev rules and adds the specified user (or the current user if no user is specified) to the 'video' group.

```bash
sudo brightnesscli install-rules
```

Please note that allowing a program to control system settings without root privileges can be a security risk. In this case, any process that the user runs can change the system's brightness, which might not be desirable in a multi-user environment or if the user runs untrusted code.

## Usage

After installing BrightnessCLI, you can run it from the command line:

```bash
brightnesscli [command]
```

Available commands:

- `up`: Increase brightness by 5%
- `down`: Decrease brightness by 5%
- `set [value]`: Set brightness to a specific value
- `get`: Get the current brightness
- `install-rules [--user username]`: Install brightness rules and add the specified user to the 'video' group. If no user is specified, the current user is used.

Example:

```bash
sudo brightnesscli install-rules
```

After running the `install-rules` command, you'll be added to the 'video' group. You need to log out and back in, or start a new shell with 'newgrp video' or 'su - $USER', to use this group.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is open source and available under the [MIT License](LICENSE).
