# Python Version Manager

PythonVersionManager (PVM for short) is a project that took inspiration from [nvm for windows](https://github.com/coreybutler/nvm-windows) and aims to easily manage multiple python enviroments on your windows system.

![](media/pvm.png)

# Table of contents

- [Python Version Manager](#python-version-manager)
- [Table of contents](#table-of-contents)
- [How to use](#how-to-use)
	- [`list` (or the alias `l`)](#list-or-the-alias-l)
	- [`install` (or the alias `i`)](#install-or-the-alias-i)
	- [`uninstall` (or the alias `u`)](#uninstall-or-the-alias-u)
	- [`reinstall` (or the alias `r`)](#reinstall-or-the-alias-r)
	- [`use`](#use)
	- [`--help` (or the alias `-h`)](#--help-or-the-alias--h)
- [Developing](#developing)
	- [Get started](#get-started)
- [FAQ](#faq)


# How to use

This is a CLI application, so you need to open a terminal and use it that way.

The options you can use are: `list`, `install`, `uninstall`, `reinstall`, `use`, `--help`

Here's a description of each command:

## `list` (or the alias `l`)

Displays the user the requested info.

`$ pvm list latest` lists all the latest python releases.

`$ pvm list all` lists all the available python versions.

`$ pvm list installed` lists all the currently installed versions.

## `install` (or the alias `i`)

Installs the requested version. If you typed the wrong version nothing will happen.

`$ pvm install latest` Installs the latest stable python version (now is 3.11.3).

`$ pvm install 3.11.0` Installs python 3.11.0.

## `uninstall` (or the alias `u`)

Uninstalls the specified version. If you typed the wrong version nothing will happen.

`$ pvm uninstall all` Uninstalls ALL the installed python versions. Be careful with this one...

`$ pvm uninstall 3.11.0` Unistalls python 3.11.0.

## `reinstall` (or the alias `r`)

Uninstalls and then installs again the specified version. If you typed the wrong version nothing will happen.

`$ pvm reinstall all` Reinstalls ALL the installed python versions.

`$ pvm reinstall 3.11.0` Reinstalls python 3.11.0.

## `use`

No alias for this one.

If already downloaded, makes active the requested version.

`$ pvm use 3.11.0`

## `--help` (or the alias `-h`)

Basically shows in the console what is written here. You can see it in action in the image [here)](#pythonversionmanager).

# Developing

Any type of contribution is well accepted, just create a PR and i'll review it as soon as possible!

## Get started

First of all, check if you have all the tools needed to run the software installed:

-   Python 3.8.x or higher (actually this project has beed written and tested with python 3.11.2)
-   Pip3
-   Git

(Surprisingly no C++ BuildTool is required)

When you're sure that you have installed them correctly, proceed by cloning the repository:

`git clone https://github.com/Relepega/PythonVersionManager.git`

Once you've cloned the repo, cd into the directory and install the required python packages:

`$ cd PythonVersionManager && pip3 install -r requirements.txt`

You're now ready to go!
If you want to build from source because you don't trust some random guy on the internet, run the build script:
`build.bat`

or do it manually:

`$ python setup.py build`

# FAQ

Q: Why reinvent the wheel when [pyenv-win](https://github.com/pyenv-win/pyenv-win) exists and does the same thing?
A: Excuse my ignorance, but i didn't know of its existence until the creation of this project.

Q: The application says that the specified version of python was not downloaded, why is it?
A: The project is still in early developement, so don't expect everything to work as it should 😉.
