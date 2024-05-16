# itectl - yet another tool for managing ITE 8291r3 keyboard backlight controller

## Description

This utility manages the backlight of the ITE 8291r3 keyboard. It
allows you to adjust the backlight brightness, set the keyboard
backlight controller to one of the built-in modes and retrieve its
status. It also provides the ability to customize mode speed,
color(s), direction, brightness, whether the mode should respond to
user input or be static, etc.

Furthermore, it can be configured using combined system
(`/etc/xdg/itectl.yml`) and user (`~/.config/itectl.yml`) yaml
configuration files and environment variables. Configuration also
allows you to name rgb colors and use those names as values for
corresponding command options and configuration properties.

`itectl` makes several attempts to discover the ITE 8291 device (with
a configurable timeout/interval) when the device cannot be found
immediately. This feature is useful when executing the command
directly after modprobing ITE 8291 module.

In addition, the project includes

- default configuration file;
- `udev` rules to make the ITE 8291 device accessible from user-space;
- `initcpio` hooks to configure keyboard backlighting at boot time
  (useful, for example, when entering a passphrase to decrypt hard drive(s))
- `bash`, `zsh` and `fish` shells completions.

The following ITE 8291r3 devices are supported

- Vendor-ID: **0x048d**, Product-ID: **0x6004**;
- Vendor-ID: **0x048d**, Product-ID: **0x6006**;
- Vendor-ID: **0x048d**, Product-ID: **0xce00**.

`itectl` recognizes the following keyboard backlight modes:
**aurora**, **breath**, **fireworks**, **marquee**, **rainbow**,
**raindrop**, **random**, **ripple**, **single-color**, **wave**.

### Dependencies

**This project does not contain the ITE 8291r3 driver and requires it
to be installed (e.g., from
[tuxedo-drivers](https://github.com/tuxedocomputers/tuxedo-drivers))
to work properly.**

### License: [MIT](./LICENSE)

## Why another ITE 8291 controller ?

This project is highly inspired
by [ite8291r3-ctl](https://github.com/pobrn/ite8291r3-ctl). Under
words _"highly inspired"_, I mean that `itectl` supports practically the
same functionality as
[ite8291r3-ctl](https://github.com/pobrn/ite8291r3-ctl). Of course,
some [ite8291r3-ctl](https://github.com/pobrn/ite8291r3-ctl) commands
are
[absent](#features-missing-from-itectl-and-present-in-ite8291r3-ctl)
in `itectl`, but some [new
features](#features-present-in-itectl-and-missing-from-ite8291r3-ctl)
have been added.

So why another project?
[ite8291r3-ctl](https://github.com/pobrn/ite8291r3-ctl) is a great
tool, but I needed the following features that that
[ite8291r3-ctl](https://github.com/pobrn/ite8291r3-ctl) does not have.

- Ability to configure keyboard backlight mode properties (e.g.,
  brightness, speed, color(s), etc.) through combined system and user
  configuration files and/or environment variables.

- Ability to retry discovery of ITE 8291 device if it cannot be
  detected immediately. This feature is useful when executing the
  utility directly after modprobing the ITE 8291r3 module (for
  example, during boot).

- Ability to launch the utility and configure the keyboard backlight
  during boot. This feature helps to enter information (such as a
  password to decrypt hard drive(s)) during boot. A Python script is
  not the best choice in this case. Small executable included in
  `initramfs` image will work better.

- It would also be nice to have shell completions, as well as `udev`
  rules and `initcpio` hooks.

These missing features led to the creation of this project. I wrote it
for myself. If you need or like these new additional features, you can
give this project a chance.

If you don’t, [ite8291r3-ctl](https://github.com/pobrn/ite8291r3-ctl)
would be a better choice.

## Differences between itectl and [ite8291r3-ctl](https://github.com/pobrn/ite8291r3-ctl)

### Minor differences in command names and functionality

- The family of `ite8291r3-ctl effect` commands are implemented as
  corresponding _"-mode"_ commands. For example `ite8291r3-ctl effect
wave` is `itectl wave-mode`.

- The following commands were renamed:

  - `ite8291r3-ctl monocolor` -> `itectl single-color-mode`
  - `ite8291r3-ctl off` -> `itectl off-mode`
  - `ite8291r3-ctl brightness` -> `itectl set-brightness`
  - `ite8291r3-ctl query --fw-version` -> `itectl firmware-version`
  - `ite8291r3-ctl query --brightness` -> `itectl brightness`
  - `ite8291r3-ctl query --state` -> `itectl state`
  - `ite8291r3-ctl palette --set-color` -> `itectl set-color`

- `ite8291r3-ctl palette --restore` command is implemented as
  `--reset` flag of every _"-mode"_ command.
- The values of the `s|--speed` option are reversed: a value of **0**
  indicates the slowest mode, **10** indicates the fastest.

All these changes are purely cosmetic. The functionality has not
changed significantly.

### Features missing from itectl and present in [ite8291r3-ctl](https://github.com/pobrn/ite8291r3-ctl)

The following commands were not implemented: `ite8291r3-ctl
test-pattern`, `ite8291r3-ctl freeze`, `ite8291r3-ctl palette
--random`, `ite8291r3-ctl query --devices`, `ite8291r3-ctl mode
--screen` and `ite8291r3-ctl anim`.

### Features present in itectl and missing from [ite8291r3-ctl](https://github.com/pobrn/ite8291r3-ctl)

- #### Tool configuration

  The system configuration file can be used to configure the `itectl`
  common and mode options, default mode, predefined and named colors,
  etc. The system configuration can be expanded/overridden using the
  user configuration file and/or environment variables.

- #### Named colors

  Any rgb color can be given a name that can be used as the
  corresponding value of the option or configuration property.

- #### Device polling

  In cases where the ITE 8291 device cannot be detected immediately,
  `itectl` will retry at the specified time intervals and stop after
  the specified timeout. This feature can be disabled.

- #### Additional configuration files

  Furthermore, the project includes files that will help you configure
  the utility.

  - default `itectl` system configuration file.
  - `udev` rules to access ITE 8291 device from user-space.
  - `initcpio` hooks to configure keyboard backlight during boot time.
  - shells completions.

## Installation

**This project does not contain the ITE 8291r3 driver and requires it
to be installed (e.g., from
[tuxedo-drivers](https://github.com/tuxedocomputers/tuxedo-drivers))
to work properly.**

### Installation from AUR repo (for arch systems only)

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

### Manual installation

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
sudo install -Dm 0755 -o root -g root ./itectl /usr/local/bin/
```

---

Default documented [configuration file](./config/etc/xdg/itectl.yml)
provided by the project can be used as the system configuration file

```
sudo install -Dm 0644 -o root -g root ./config/etc/xdg/itectl.yml /etc/xdg/
```

and/or the user configuration file

```
cp ./config/etc/xdg/itectl.yml ~/.config/
```

---

The `udev` [rules](./config/usr/lib/udev/rules.d/10-ite8291r3.rules)
provided by the project allow access to the discovered ITE 8291r3
device for users of the systemd `input` group. For the rules to take
effect, they must be copied to one of the locations recognized by
`udev` on your system. For example

```
sudo install -Dm 0644 -o root -g root ./config/usr/lib/udev/rules.d/10-ite8291r3.rules /etc/udev/rules.d/

```

---

The `initcpio` hook provided by the project is based on the
`ite_8291` module from
[tuxedo-drivers](https://github.com/tuxedocomputers/tuxedo-drivers). If
another ITE 8291 module is installed, `ite_8291` must be replaced with
this module name in
[config/usr/lib/initcpio/hooks/itectl](./config/usr/lib/initcpio/hooks/itectl).

This hook calls `itectl` without any options and uses the system
configuration to set the default mode, device polling, named colors,
etc. If you prefer a different solution, you can edit the hook
[config/usr/lib/initcpio/hooks/itectl](./config/usr/lib/initcpio/hooks/itectl)
and add additional options and/or environment variables.

The hook files must be copied to locations recognized by your
system. For example, for arch systems

```
sudo install -Dm 0644 -o root -g root ./config/usr/lib/initcpio/hooks/itectl /usr/lib/initcpio/hooks
sudo install -Dm 0644 -o root -g root ./config/usr/lib/initcpio/install/itectl /usr/lib/initcpio/install

```

To work properly, the `itectl` hook must be included in the
`/etc/mkinitcpio.conf` file after the `udev` hook. For example

```
HOOKS=(base udev autodetect itectl microcode modconf kms keyboard keymap consolefont block filesystems fsck)
```

These hooks rely on ITE 8291 module `ite_8291`. If another ITE 8291
module is installed, `ite_8291` must be replaced with this module
name in [config/usr/lib/initcpio/hooks/itectl](./config/usr/lib/initcpio/hooks/itectl).

---

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
- predefined colors set to

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

- **mode** - default mode of the keyboard backlight to set if `itectl`
  is called with no sub-command specified.<br/> Environment variable:
  `ITECTL_MODE`.
- **brightness** - brightness of the keyboard backlight.<br/>Minimum
  value: **0**. Maximum value: **50**. Default value: **25**.<br/>
  Environment variable: `ITECTL_BRIGHTNESS`.<br/> Command line
  option(s): `-b`, `--brightness`.
- **speed** - speed of the keyboard backlight effect.<br/>Slowest
  effect: **0**. Fastest effect: **10**. Default value:
  **5**.<br/>Environment variable: `ITECTL_SPEED`.<br/> command line
  option(s): `-s`, `--speed`.
- **direction** - direction of the keyboard backlight effect.<br/>
  Allowed values: **none**, **left**, **right**, **up**,
  **down**. Default value: **right**.<br/>Environment variable:
  `ITECTL_DIRECTION`.<br/>Command line option(s): `-d`, `--direction`.
- **colorNum** - number of the predefined color of the keyboard
  backlight controller to use by effect.<br/>No color: **0**. Random
  color: **8**. Customizable colors: **1**-**7**. Default value:
  **8**.<br/> Environment variable: `ITECTL_COLORNUM`.<br/>Command
  line option(s): `-c`, `--color-num`.
- **reactive** - determines whether the keyboard backlight effect
  should react to keypresses.<br/>Default value:
  **false**.<br/>Environment variable: `ITECTL_REACTIVE`.<br/>Command
  line option: `--reactive`.
- **reset** - specifies whether the customizable predefined colors
  should be reset to their corresponding configured/default values
  before setting an effect. Used by all _"mode"_ sub-commands.<br/>
  Default value: **false**.<br/>Environment variable:
  `ITECTL_RESET`.<br/>Command line option: `--reset`.
- **save** - indicates whether ITE 8291 controller should save its
  state.<br/>Default value: **false**.<br/>Environment variable:
  `ITECTL_SAVE`.<br/>Command line option: `--save`.
- **singleModeColor** - color of the keyboard backlight to use by
  single color mode. The option value can be a name of one of
  configured named colors or an rgb value in one of the following
  formats: **0xHHHHHH**, **#xHHHHHH**, **#HHHHHH**, **HHHHHH**,
  **#HHH**, **HHH**.<br/>Default value: **#FFFFFF**.<br/>Environment
  variable: `ITECTL_SINGLEMODECOLOR`.<br/>Command line option(s):
  `--color-name` or `--rgb` or (`--red` and/or `--green` and/or
  `--blue`)
- **poll** - device probing related properties.

  - **interval** - time interval to wait between device detection
    attempts. The value is ignored if **timeout** is set to
    **0**.<br/>Default value: **200ms**.<br/>Environment variable:
    `ITECTL_POLL_INTERVAL`.<br/>Command line option:
    `--poll-interval`.

  - **timeout** - maximum duration of time to wait for an ITE 8291
    device to be available. If set to **0**, only one attempt to
    discover ITE 8291 device is made and if it is not detected,
    `itectl` returns immediately with non zero exit code.<br/>Default
    value: **0**.<br/>Environment variable:
    `ITECTL_POLL_TIMEOUT`.<br/>Command line option: `--poll-timeout`.

  For instance

  ```

  poll:
    interval: "200ms"
    timeout: "500ms"

  ```

- **device** - device address related properties.

  - **bus** - bus number of the ITE 8291 device to use. If it's set to
    **0**, the option is ignored.<br/>Default value:
    **0**.<br/>Environment variable: `ITECTL_DEVICE_BUS`.<br/>Command
    line option: `--device-bus`.

  - **address** - device number of the ITE 8291 to use. If it's set to
    **0**, the option is ignored.<br/>Default value:
    **0**.<br/>Environment variable:
    `ITECTL_DEVICE_ADDRESS`.<br/>Command line option:
    `--device-address`.

  Both values must be either both positive or non-positive
  (i.e., ignored). For instance

  ```

  device:
    bus: 1
    address: 2

  ```

- **predefinedColors** - predefined customizable color values.<br/>The
  keys have the format **color*<N\>***, where _<N\>_ is one of the
  color numbers (**1**-**7**).<br/>The color value can be either the
  name of one of the configured named colors, or an RGB value in one
  of the following formats: **0xHHHHHH**, **#xHHHHHH**, **#HHHHHH**,
  **HHHHHH**, **#HHH**, **HHH**.<br/>Default values:

  1. `#FFFFFF`
  1. `#FF0000`
  1. `#FFFF00`
  1. `#00FF00`
  1. `#0000FF`
  1. `#00FFFF`
  1. `#FF00FF`

  Example

  ```

  predefinedColors:
    color1: lemon
    color2: "#FFBF00"
    color3: coral
    color4: "#008000"
    color5: opal
    color6: "#B53389"
    color7: cyan

  ```

- **namedColors** - color name -> RGB color mapping.<br/> The color
  name can be an arbitrary string.<br/> The color value can be an RGB
  value in one of the following formats: **0xHHHHHH**, **#xHHHHHH**,
  **#HHHHHH**, **HHHHHH**, **#HHH**, **HHH**. For instance

  ```

  namedColors:
    aero: "#7CB9E8"
    alloy_orange: "#C46210"
    azure: "#007FFF"

  ```

  The [default configuration file](./config/etc/xdg/itectl.yml)
  provided by the project contains settings for all colors from the
  Wikipedia [Lists of
  colors](https://en.wikipedia.org/wiki/Lists_of_colors). Unfortunately,
  the color representation of the ITE 8291 keyboard backlight is not
  entirely accurate. Some colors look too greenish, some too
  bluish. This configuration property allows you to reconfigure any of
  these colors to your liking.

## Usage

### Common options

The following options are supported by every command

- `--config` - path to the configuration file. If specified, system
  and user configuration files and environment variables are ignored.
- `--poll-interval` - timeout interval between attempts to detect a
  device. The value is ignored if `--poll-timeout` is set to `0`. It
  defaults to the configured value or 200ms if no value is configured.
- `--poll-timeout` - maximum duration of time to wait for an ITE 8291
  device to become available. If set to `0`, only one attempt is made
  to discover an ITE 8291 device, and if it is not found, itectl
  immediately returns with a non-zero exit code. It defaults to the
  configured value or `0` if no value is configured.
- `--device-bus` - bus number of the ITE 8291. If set to `0`, the
  option is ignored. The default is the configured value or `0` if the
  value is not configured.
- `--device-address` - address number of the ITE 8291. If it is set to
  `0`, the option is ignored. The default is the configured value or
  `0` if the value is not configured.
- `--help` - prints a command's help.

### Mode options

The following options are supported by _-mode_ commands

- `-b`, `--brightness` - brightness of the keyboard backlight; minimum
  value: `0`; maximum value: `50`. The default is the configured value
  or `25` if the value is not configured.

- `-s`, `--speed` - speed of the keyboard backlight effect; slowest
  effect: `0`; fastest effect: `10`. The default is the configured value
  or `5` if the value is not configured.

- `-d`, `--direction` - direction of the keyboard backlight
  effect. The allowed values are `none`, `left`, `right`, `up`,
  `down`. The default is the configured value or `right` if the value
  is not configured.

- `--color-num` - number of the predefined color of the keyboard
  backlight effect; no color: `0`; random color: `8`; customizable
  colors: `1`-`7`. The default is the configured value or `8` if the
  value is not configured.

- `--reactive` - if specified, the keyboard backlight effect will
  respond to keypresses. Defaults to the configured value, or `false`
  if no value is configured.

- `--reset` - if specified, custom predefined keyboard backlight
  colors will be reset to the corresponding configured values or
  default values if not configured. Default custom predefined colors
  are

  1. `#FFFFFF`
  1. `#FF0000`
  1. `#FFFF00`
  1. `#00FF00`
  1. `#0000FF`
  1. `#00FFFF`
  1. `#FF00FF`

- `--save` - if specified, the keyboard backlight controller will
  retain its state. Defaults to the configured value, or `false` if no
  value is configured.

- `--color-name` - name of the configured color.

- `--rgb` - rgb color in one of the following formats: **0xHHHHHH**,
  **#xHHHHHH**, **#HHHHHH**, **HHHHHH**, **#HHH**, **HHH**.

- `--red`, `--green`, `--blue` - the corresponding red, green and blue
  parts of the color.

### Commands

- `aurora-mode` - sets the keyboard backlight to _aurora_ mode.
- `breath-mode` - sets the keyboard backlight to _breathing_ mode.
- `fireworks-mode` - sets the keyboard backlight to _fireworks_ mode.
- `brightness` - prints out brightness of the keyboard backlight.
- `firmware-version` - prints out firmware version of the keyboard
  backlight controller.
- `marquee-mode` - sets the keyboard backlight to _marquee_ mode.
- `off-mode` - turns off the keyboard backlight.
- `rainbow-mode` - sets the keyboard backlight to _rainbow_ mode.
- `raindrop-mode` - sets the keyboard backlight to _raindrop_ mode.
- `random-mode` - sets the keyboard backlight to _random_ mode.
- `ripple-mode` - sets the keyboard backlight to _ripple_ mode.
- `set-brightness` - sets the keyboard backlight brightness to the
  specified value.
- `set-color` - sets the keyboard backlight custom predefined color
  sepcified via `-c`|`--color-num` option to the color specified by
  either `--color-name` or `--rgb` or (`--red` and/or `--green` and/or
  `--blue`) options.
- `single-color-mode` - sets the keyboard backlight to _single-color_
  mode. In this mode, all keys are assigned a single color, specified
  by either `--color-name` or `--rgb` or (`--red` and/or `--green`
  and/or `--blue`) options.
- `state` - prints out `Off` if the keyboard backlight is turned off
  by `off-mode` command. Otherwise it prints `On` (even if brightness
  is set to `0`).
- `wave-mode` - sets the keyboard backlight to _wave_ mode.

## TODO

- Implement `-v|--verbose` option.
- Implement mode's specific configuration (i.e. different
  configuration properties for different modes).

## Acknowledgements

- [ite8291r3-ctl](https://github.com/pobrn/ite8291r3-ctl) is a great
  tool. If I didn't need a small executable instead of a Python
  script, I wouldn't create this project.
- [ite-backlight](https://github.com/hexagonal-sun/ite-backlight)
- [tuxedo-drivers](https://github.com/tuxedocomputers/tuxedo-drivers)

## Similar projects

- [ite8291r3-ctl](https://github.com/pobrn/ite8291r3-ctl)
- [ite-backlight](https://github.com/hexagonal-sun/ite-backlight)
- [ite8291r2ctl](https://github.com/leewis101/ite8291r2ctl)

_Keep coding and have fun._ :metal:
