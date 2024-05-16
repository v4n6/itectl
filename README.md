# itectl - yet another tool for managing ITE 8291 keyboard backlight controller

## Description

This utility controls the backlight of the ITE 8291 keyboard. It
allows you to set the backlight brightness, set the keyboard backlight
controller to one of the built-in modes and retrieve its status. It
also provides the functionality to customize mode speed, color(s),
direction, brightness, and whether the mode should respond to user
input or be static. Furthermore, it can be configured using merged
system (`/etc/xdg/itectl.yml`) and user (`~/.config/itectl.yml`) yaml
configuration files and environment variables. The configuration also
allows you to name `rgb` colors and use those names as values for
corresponding command options and configuration properties.`itectl`
makes several attempts to discover the ITE 8291 device (with a
configurable timeout/interval) when the device cannot be found
immediately (useful when running the command directly after modprobing
an ITE 8291 module). The project also provides

- default configuration file;
- `udev` rules to make the ITE 8291 device accessible from user-space;
- `initcpio` hooks to configure keyboard backlighting at boot time
  (useful when entering a passphrase to decrypt hard drive(s));
- `bash`, `zsh` and `fish` shells completions.

The following ITE 8291 devices are supported

- Vendor-ID: **0x048d**, Product-ID: **0x6004**;
- Vendor-ID: **0x048d**, Product-ID: **0x6006**;
- Vendor-ID: **0x048d**, Product-ID: **0xce00**.

`itectl` recognizes the following keyboard backlight modes:
**aurora**, **breath**, **fireworks**, **marquee**, **rainbow**,
**raindrop**, **random**, **ripple**, **single-color**, **wave**.

### Dependencies

This project does not contain the ITE 8291 driver and requires it to
be installed (e.g., from
[tuxedocomputers/tuxedo-drivers](https://github.com/tuxedocomputers/tuxedo-drivers))
to work properly.

### License: [MIT](./LICENSE)

## Why another ITE 8291 controller ?

This project is highly inspired by
[pobrn/ite8291r3-ctl](https://github.com/pobrn/ite8291r3-ctl). Under
words _"highly inspired"_, I mean that `itectl` supports almost the
same functionality as
[pobrn/ite8291r3-ctl](https://github.com/pobrn/ite8291r3-ctl). Of
course, some
[pobrn/ite8291r3-ctl](https://github.com/pobrn/ite8291r3-ctl) commands
are
[absent](#features-missing-in-v4n6itectl-and-present-in-pobrnite8291r3-ctl)
in `itectl`, and new features were
[added](#features-present-in-v4n6itectl-and-missing-in-pobrnite8291r3-ctl).

So why another project?
[pobrn/ite8291r3-ctl](https://github.com/pobrn/ite8291r3-ctl) is a
great utility, but I needed the following features that are absent in
[pobrn/ite8291r3-ctl](https://github.com/pobrn/ite8291r3-ctl)

- Possibility to configure keyboard backlight mode properties (e.g.,
  brightness, speed, color(s), etc.) via merged system and user
  configuration files and/or environment variables.

- Repeated attempts to detect the ITE 8291 device, if it cannot be
  found immediately. This function is useful when executing the
  utility immediately after modprobing an ITE 8291 module (e.g., at
  boot time).

- The ability to execute the utility and configure the keyboard
  backlight during the boot time. This function helps to enter
  information (e.g., a password for deciphering a hard drive) at
  boot. A Python script is not a good choice in this case. I needed
  a small executable to include in the `initramfs` image.

- Shells completions, as well as the `udev` rules and `initcpio`
  hooks, would also be nice to have.

These missing functions led to a new project. I wrote it for
myself. If you need/like these new functions, you can give this
project a chance.

If you do not,
[pobrn/ite8291r3-ctl](https://github.com/pobrn/ite8291r3-ctl) would be
a better choice. It is more mature and is used more widely.

## Differences between [v4n6/itectl](https://github.com/v4n6/itectl) and [pobrn/ite8291r3-ctl](https://github.com/pobrn/ite8291r3-ctl)

### Minor differences in commands naming and functionality

- The family of `ite8291r3-ctl effect` commands are implemented as
  corresponding _"-mode"_ commands. For example<br/>
  `ite8291r3-ctl effect wave` is `itectl wave-mode`.
- The following commands were renamed:

  - `ite8291r3-ctl monocolor` -> `itectl single-color-mode`
  - `ite8291r3-ctl off` -> `itectl off-mode`
  - `ite8291r3-ctl brightnes` -> `itectl set-brightness`
  - `ite8291r3-ctl query --fw-version` -> `itectl firmware-version`
  - `ite8291r3-ctl query --brightness` -> `itectl brightness`
  - `ite8291r3-ctl query --state` -> `itectl state`
  - `ite8291r3-ctl palette --set-color` -> `itectl set-color`

- `ite8291r3-ctl palette --restore` command is implemented as
  `--reset` flag of every _"-mode"_ command.
- `-s|--speed` option values are reversed: **0** value indicates the slowest mode,
  **10** - the fastest.

All these changes are purely cosmetic. The functionality has not been
significantly changed.

### Features missing in [v4n6/itectl](https://github.com/v4n6/itectl) and present in [pobrn/ite8291r3-ctl](https://github.com/pobrn/ite8291r3-ctl)

The following commands were not implemented: `ite8291r3-ctl
test-pattern`, `ite8291r3-ctl freeze`, `ite8291r3-ctl palette
--random`, `ite8291r3-ctl query --devices`, `ite8291r3-ctl mode
--screen` and `ite8291r3-ctl anim`.

### Features present in [v4n6/itectl](https://github.com/v4n6/itectl) and missing in [pobrn/ite8291r3-ctl](https://github.com/pobrn/ite8291r3-ctl)

- #### Configuration

  The system configuration file can be used to configure the itectl
  options, default mode, predefined and named colors. The `itectl`
  system configuration can be expanded/overridden using the user
  configuration file and/or environment variables.

- #### Named colors

  Any rgb color can be assigned a name that can be used as the
  corresponding value of the option or configuration property.

- #### Device polling

  In cases where the ITE 8291 device cannot be detected immediately,
  the itectl will repeat the attempt to discover it with the specified
  time interval and will stop after the specified timeout. This
  feature can be disabled.

- #### Additional configuration files

  Additionally, there are files included in the project that can help
  configure the utility.

  - default `itectl` system configuration file.
  - `udev` rules to access ITE 8291 device from user-space.
  - `initcpio` hooks to configure keyboard backlight during boot time.
  - shells completions.

## Installation

This project does not contain the ITE 8291 driver and requires it to
be installed separately (e.g., from
[tuxedocomputers/tuxedo-drivers](https://github.com/tuxedocomputers/tuxedo-drivers))
to work properly.

### Automatically (for arch systems only)

You can install `itectl` through the dedicated
[itectl-aur](https://github.com/v4n6/itectl-aur) AUR repository. To do
this, you need to clone the repository and install the package using
`makepkg`. For example

```
git clone https://github.com/v4n6/itectl-aur
cd ./itectl-aur
makepkg -fi
```

</code>

### Manualy

You can test the project using the `go` tool

```
go test ./...
```

You can build `itectl` by executing

```
go build
```

and copy it to a location in `PATH`. For example

```
sudo install -Dm 0755 -o root -g root ./itectl /usr/bin/
```

Default documented [configuration file](./config/etc/xdg/itectl.yml)
provided by the project can be used as a system configuration file

```
sudo install -Dm 0644 -o root -g root ./config/etc/xdg/itectl.yml /etc/xdg/
```

or a user configuration file

```
cp ./config/etc/xdg/itectl.yml ~/.config/
```

The `udev` [rules](./config/usr/lib/udev/rules.d/10-ite8291r3.rules)
provided by the project allow access to the discovered ITE 8291 device
for users of the systemd `input` group. For the rules to take effect,
they must be copied to one of the locations recognized by `udev` on
your system. For example

```
sudo install -Dm 0644 -o root -g root ./config/usr/lib/udev/rules.d/10-ite8291r3.rules /etc/udev/rules.d/

```

The provided `initcpio` hooks can be used by copying them to locations
recognized by your system. For example, for arch systems

```
sudo install -Dm 0644 -o root -g root ./config/usr/lib/initcpio/hooks/itectl /usr/lib/initcpio/hooks
sudo install -Dm 0644 -o root -g root ./config/usr/lib/initcpio/install/itectl /usr/lib/initcpio/install

```

These hooks rely on ITE 8291 module `ite_8291`. If another ITE 8291
module is installed, `ite_8291` must be replaced with this module
name in [config/usr/lib/initcpio/hooks/itectl](./config/usr/lib/initcpio/hooks/itectl).

To work properly, the `itectl` hook must be included in the
`/etc/mkinitcpio.conf` file after the `udev` hook. For example

```
HOOKS=(base udev autodetect itectl microcode modconf kms keyboard keymap consolefont block filesystems fsck)
```

Shells completions can be created using the `itectl completion` command. For example

```
# bash shell completions
itectl completion bash >>~/.bash_completion

# zsh shell completions
sudo zsh -c 'itectl completion zsh >/usr/share/zsh/site-functions/_itectl'
sudo zsh -c "chmod 'u+rw-x,go+r-wx' /usr/share/zsh/site-functions/_itectl"
sudo zsh -c 'chown root:root /usr/share/zsh/site-functions/_itectl'

# fish shell completions
mkdir -p ~/.local/share/fish/vendor_completions.d
itectl completion fish > ~/.local/share/fish/vendor_completions.d/itectl.fish

```

## Configuration

It's possible to configure `itectl` via

1. system configuration file: `/etc/xdg/itectl.yml`;
1. user configuration file: `~/.config/itectl.yml`;
1. environment variables (starting with `ITECTL_` prefix).

Both configuration files are merged. It's possible to add or override
properties in system configuration file using the user configuration
file and/or environment variables.

For instance

system configuration file `/etc/xdg/itectl.yml`

```
mode: rainbow
brightness: 25
reset: true
save: true
singleModeColor: "#FFFFFF"
predefinedColors:
  color1: "#FFFFFF"
  color2: "#FF0000"
  color3: "#FFFF00"
  color4: "#00FF00"
  color5: "#0000FF"
  color6: "#00FFFF"
  color7: "#FF00FF"
poll:
  interval: "100ms"
  timeout: "1s"
namedColors:
  barn_red: "#7C0A02"
  alloy_orange: "#C46210"
  carmine: "#960018"
  goldenrod: "#DAA520"
  melon: "#FEBAAD"

```

and user configuration file `~/.config/itectl.yml`;

```
mode: fireworks
brightness: 40
save: false
reactive: true
singleModeColor: melon
predefinedColors:
  color2: "#D1E231"
  color7: carmine
poll:
  timeout: 0
namedColors:
  purple: "#6A0DAD"
  steel_blue: "#4682B4"
  carmine: "0xC60010"
```

will result in the following default values

- default mode set to _fireworks_
- default brightness - 40
- default reset - true
- default reactive - true
- default save - false
- default single-mode color - _melon_ (`#FEBAAD`)
- device polling - disabled (timeout = 0)
- predfined colors set to

  1. `#FFFFFF`
  1. `#D1E231`
  1. `#FFFF00`
  1. `#00FF00`
  1. `#0000FF`
  1. `#00FFFF`
  1. carmine (`0xC60010`)

- named colors set to

  - _barn_red_ = `#7C0A02`
  - _alloy_orange_ = `#C46210`
  - _goldenrod_ = `#DAA520`
  - _melon_ = `#FEBAAD`
  - _purple_ = `#6A0DAD1`
  - _steel_blue_ = `#4682B4`
  - _carmine_ = `0xC60010`

### Configuration Properties

- **mode** - default mode of the keyboard backlight to set if
  `itectl` was called without any sub-command is specified.
  Corresponding environment variable `ITECTL_MODE`.
- **brightness** - brightness of the keyboard backlight; minimum value:
  **0**; maximum value: **50**. Default value: **25**. Corresponding
  environment variable: `ITECTL_BRIGHTNESS`. Corresponding command
  line option(s): `-b`, `--brightness`.
- **speed** - speed of the keyboard backlight effect; slowest effect:
  **0**; fastest effect: **10**. Default value: **5**. Corresponding
  environment variable: `ITECTL_SPEED`. Corresponding command line
  option(s): `-s`, `--speed`.
- **direction** - direction of the keyboard backlight effect. Allowed
  values: **none**, **left**, **right**, **up**, **down**. Default
  value: **right**. Corresponding environment variable:
  `ITECTL_DIRECTION`. Corresponding command line option(s): `-d`,
  `--direction`.
- **colorNum** - number of the predfined color of keyboard backlight
  controller to use by effect; **0** means no color; **8** - random
  color; **1**-**7** - a customizable color. Default value: **8**.
  Corresponding environment variable:
  `ITECTL_COLOR_NUM`. Corresponding command line option(s): `-c`,
  `--color-num`.
- **reactive** - specifies whether keyboard backlight effect should
  react to key press. Default value: **false**. Corresponding
  environment variable: `ITECTL_REACTIVE`. Corresponding command
  line option: `--reactive`.
- **reset** - specifies whether the customizable predefined colors
  should be reset to their corresponding configured/default values
  before setting the effect. Used by all _"mode"_ sub-commands.
  Default value: **false**. Corresponding environment variable:
  `ITECTL_RESET`. Corresponding command line option: `--reset`.
- **save** - specifies whether ITE 8291 controller should save its
  state. Default value: **false**. Corresponding environment variable:
  `ITECTL_SAVE`. Corresponding command line option: `--save`.
- **singleModeColor** - color of the keyboard backlight to use by
  single color mode. The value of the option can be a name of one of
  the configured named colors or an RGB value in one of the following
  formats: **0xHHHHHH**, **#xHHHHHH**, **#HHHHHH**, **HHHHHH**,
  **#HHH**, **HHH**. Default value: **#FFFFFF**. Corresponding
  environment variable: `ITECTL_SINGLE_MODE_COLOR`. Corresponding
  command line option(s): `--color-name` or `--rgb` or (`--red` and/or
  `--green` and/or `--blue`)

- **poll** - dictionary specifying poll interval and poll timeout.

  - **interval** - time interval to wait between controller polls. The
    value is ignored if timeout is set to **0**. Default value:
    **200ms**. Corresponding environment variable:
    `ITECTL_POLL_INTERVAL`. Corresponding command line option: `--poll-interval`.

  - **timeout** - maximum time to wait for ITE 8291 device to be
    available. If set to **0**, only one attempt to discover ITE 8291
    device is made. If it's not available, `itectl` returns with non zero exit
    code. Default value: **0**. Corresponding environment variable:
    `ITECTL_POLL_TIMEOUT`. Corresponding command line option: `--poll-timeout`.

  For instance:

  ```

  poll:
    interval: "200ms"
    timeout: "500ms"

  ```

- **device** - dictionary specifying bus and address of ITE 8291
  device to use by `itectl`.

  - **bus** - bus of ITE 8291 device to use. If it's set to **0**, the
    option is ignored. Default value: **0**. Corresponding environment
    variable: `ITECTL_DEVICE_BUS`. Corresponding command line option:
    `--device-bus`.

  - **address** - address of ITE 8291 device to use. If it's set to
    **0**, the option is ignored. Default value: **0**. Corresponding
    environment variable: `ITECTL_DEVICE_ADDRESS`. Corresponding
    command line option: `--device-address`.

    Both values should be set either to zero or to non zero values.
    For instance:

  ```

  device:
    bus: 1
    address: 2

  ```

- **predefinedColors** - dictionary associating customizable color
  number to a color. The key is in the form **color*<N\>***, where
  _<N\>_ is one of the customizable predefined color numbers
  (**1**-**7**). The value of the color can be a name of one of the
  configured named colors or an RGB value in one of the following
  formats: **0xHHHHHH**, **#xHHHHHH**, **#HHHHHH**, **HHHHHH**,
  **#HHH**, **HHH**. For instance:

  ```

  predefinedColors:
    color1: "#FFFFFF"
    color2: "#FF0000"
    color3: "#FFFF00"
    color4: "#00FF00"
    color5: "#0000FF"
    color6: "#00FFFF"
    color7: "#FF00FF"

  ```

- **namedColors** - dictionary associating a string with a color. It
  allows to use this string as value of `--color-name` property or a
  predefined color. The value of the color can be an RGB value in one
  of the following formats: **0xHHHHHH**, **#xHHHHHH**, **#HHHHHH**,
  **HHHHHH**, **#HHH**, **HHH**. For instance:

```

namedColors:
aero: "#7CB9E8"
alloy_orange: "#C46210"
azure: "#007FFF"

```

The [supplied default configuration
file](./config/etc/xdg/itectl.yml) contains configuration of all
colors from Wikipedia [Lists of
colors](https://en.wikipedia.org/wiki/Lists_of_colors). Color
representation of ITE 8291 keyboard backlight is not quite
precise. Some colors look too greenish some too bluish. This
property allows to reconfigure any of these colors at your
convenience.

## Usage

### General options

Following options are supported by every command:

- `--poll-interval` - time interval to wait between ITE 8291 device
  polls. The value is ignored if `--poll-timeout` is set to **0**. It
  defaults to a configured value or **200ms** if the value is not
  configured.
- `--poll-timeout` - maximum time to wait for ITE 8291 device to be
  available. If set to **0**, only one attempt to discover ITE 8291
  device is made. If the device is not available, `itectl` returns
  with non zero exit code. It defaults to a configured value or
  **0** if the value is not configured.
- `--device-bus` - bus of ITE 8291 device to use by `itectl`. If it's
  set to **0**, the option is ignored. It defaults to a configured value or
  **0** if the value is not configured.
- `--device-address` - address of ITE 8291 device to use by `itectl`.
  If it's set to **0**, the option is ignored. It defaults to a
  configured value or **0** if the value is not configured.
- `--config` - path to a configuration file to use. If it's specified,
  system and user xdg configuration files are ignored.
- `--help` - print command's help.

### Mode options

Following options are supported by some mode commands:

- `-b`, `--brightness` - brightness of the keyboard backlight; minimum value:
  **0**; maximum value: **50**. It defaults to a configured value or
  **25** if the value is not configured.

- `-s`, `--speed` - speed of the keyboard backlight effect; slowest effect:
  **0**; fastest effect: **10**. It defaults to a configured value or
  **5** if the value is not configured.

- `-d`, `--direction` - direction of the keyboard backlight
  effect. Allowed values: **none**, **left**, **right**, **up**,
  **down**. It defaults to a configured value or **right** if the
  value is not configured.

- `--color-num` - number of the predfined color of keyboard backlight
  controller to use by effect; **0** means no color; **8** - random
  color; **1**-**7** - a customizable color. It defaults to a
  configured value or **8** if the value is not configured.

- `--reactive` - if specified keyboard backlight effect will react to
  a key press. It defaults to a configured value or **false** if the
  value is not configured.
- `--reset` - if specified the customizable predefined colors
  will be reset to their corresponding configured values or

```

predefinedColors:
color1: "#FFFFFF"
color2: "#FF0000"
color3: "#FFFF00"
color4: "#00FF00"
color5: "#0000FF"
color6: "#00FFFF"
color7: "#FF00FF"

```

if not configured.

`--reset` defaults to a configured value or **false** if not configured.

- `--save` - if specified ITE 8291 controller will save its
  state. It defaults to a configured value or **false** if the
  value is not configured.

- `--color-name` - name of the configured color to use.

- `--rgb` - a color to use. It's value is an RGB color in one of the
  following formats: **0xHHHHHH**, **#xHHHHHH**, **#HHHHHH**,
  **HHHHHH**, **#HHH**, **HHH**.

- `--red`, `--green`, `--blue` - a color to use. Values of these
  properties are numbers in **0**-**255** range.

### Commands

- **aurora-mode** - set ITE 8291 device to a _aurora_ mode.
- **breath-mode** - set ITE 8291 device to a _breathing_ mode.
- **fireworks-mode** - set ITE 8291 device to a _fireworks_ mode.
- **brightnes** - prints out brightness of keyboard backlight.
- **firmware-version** - prints out firmware version of keyboard backlight controller.
- **marquee-mode** - set ITE 8291 device to a _marquee_ mode.
- **off-mode** - set keyboard backlight off.
- **rainbow-mode** - set ITE 8291 device to a _rainbow_ mode.
- **raindrop-mode** - set ITE 8291 device to a _raindrop_ mode.
- **random-mode** - set ITE 8291 device to a _random_ mode.
- **ripple-mode** - set ITE 8291 device to a _ripple_ mode.
- **set-brightness** - set keyboard backlight brightness.
- **set-color** - set keyboard backlight predefined color
  (**1**-**7**) sepcified via `-c`|`--color-num` option to color
  specified by either `--color-name` or `--rgb` or (`--red` and/or
  `--green` and/or `--blue`) options.
- **single-color-mode** - set ITE 8291 to a single color mode. All
  keys are set to a single color specified by either `--color-name` or
  `--rgb` or (`--red` and/or `--green` and/or `--blue`) options.
- **state** - prints out **Off** if the keyboard backlight is switched
  off via **off-mode** command. It prints **On** otherwise (even if
  brightness is set to **0**).
- **wave-mode** - set ITE 8291 device to a _wave_ mode.

## TODO

- Implement properties like brightness, speed, direction, etc. configurable per mode.
- Implement test-pattern command.

## Aknowledges

## Similar projects

If you like this tool give it a star.

_Keep coding and have fun._ :metal:
