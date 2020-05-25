# clips
A simple tool to generate snippets on the command line using templates. For example, generating bash one liners straight to the clipboard in situations where functions and aliases in config files aren't practical.

See the [Tips](#tips) section for an example use case.

## Installation
With `go` installed and a properly configured `$GOPATH`, simply run
```
go get -u github.com/danielthatcher/clips
```

## Usage
clips uses a template system based around JSON files. To view the available templates, use
```
clips list
```

To generate a snippet using one of the available templates, run
```
clips [template]
```
By default, this will ask you to fill in any missing variables for the template. To use defaults, add `-q`. To copy the result to the clipboard, add `-c`.

For a list of all options, use
```
clips -h
```

### Templates
No templates come by default. To add a new template and open it in an editor (`-e`), run
```
clips new -e [template]
```
Templates are JSON files, stored under `$HOME/.config/clips/[profile]/[template].json`. An example of a template can be seen [here](https://github.com/danielthatcher/dotfiles2/blob/master/clips/ctflinux/shell.json).

The `line` field of a template should provide the base of the template. The `variables` field provides a dictionary of substitutions to make on the `line` field, with the keys being strings to find, and the values being the names of the variables that should be used to replace them. The `defaults` sections is a dictionary of default values for variables in the case that they are unset.

### Variables
Variables can be set globally using
```
clips set [variable] [value]
```
The value of this variable will then be used to fill in templates.

Variables can also be set when generating a template by using the `-s` flag. For example,
```
clips [template] -s MYVAR=zzz
```

### Profiles
Profiles provide a way to switch between sets of active templates. You can view your current profile with
```
clips profile
```
You can view all profiles with
```
clips profile list
```
The active profile can be changed with
```
clips profile use [profile]
```
Profiles can be added using
```
clips profile add [profile]
```
and removed using
```
clips profile remove [profile]
```

## Tips
This is particularly useful when combined with a program such as [rofi](https://github.com/davatorium/rofi) or [dmenu](https://tools.suckless.org/dmenu/). For example, I use the [profiles in my dotfiles](https://github.com/danielthatcher/dotfiles2/tree/master/clips) for Boot2Root style CTFs with the following command bound to a hotkey:
```
clips -q -c $(clips list | rofi -dmenu)
```
When I (typically) start a VPN connection and focus on a target machine, I set the `LHOST` and `RHOST` variables using
```
clips set LHOST <vpn-ip>
clips set RHOST <target-ip>
```
I also keep a public SSH key stored in a variable.

I can then press my assigned hotkey, and make a selection from a number of useful commands to run on the target machine (as well as some simpler templates to generate the URL of my local HTTP server, or copy the `LHOST` or `RHOST` variables to the clipboard), which will then be copied to the clipboard, ready to be used. This saves remembering and repeatedly typing out repetitive commands.
