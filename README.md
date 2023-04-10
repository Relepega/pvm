# Python Version Manager

PythonVersionManager (PVM for short) is a project that has been inspired by [nvm for windows](https://github.com/coreybutler/nvm-windows) and aims to easily manage multiple python enviroments on your windows system.

![](media/pvm.png)

# Table of contents

- [Python Version Manager](#python-version-manager)
- [Table of contents](#table-of-contents)
- [How to use](#how-to-use)
	- [Install](#install)
	- [Uninstall](#uninstall)
	- [Commands](#commands)
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

## Install

At the moment is available only a portable version of the app and the installation process is nothing too complex:

1. Download the latest version of the app [from here](https://github.com/Relepega/PythonVersionManager/releases).
2. Unzip the app in a place where you won't move the app folder never again until you uninstall it.
3. Inside the app folder go to the subfolder `scripts`, right-click on the `install.bat` and then click on `Run as administrator`. If the UAC kicks in, please click on `yes`.
4. Wait until you don't see any flying terminal and then close the folder.
5. Restart your terminal.
6. type `pvm -h` and press enter. If the installation was successful then you will see an output like in [this image](#python-version-manager).
7. Profit 🎉!

## Uninstall

It's basically the same process as the installation one:

1. Open the app folder.
2. go to the `scripts` subdolder, right-click on the `uninstall.bat` and then click on `Run as administrator`. If the UAC kicks in, please click on `yes`.
3. Wait until you don't see any flying terminal and then close the folder.
4. Restart your terminal.
5. type `pvm -h` and press enter. If pvm was uninstalled successfully now in the terminal you should see an error.
6. Profit 🎉!
7. (optional) Please let me know why you decided to uninstall pvm: fill [this form](https://github.com/Relepega/PythonVersionManager/issues/new) and i'll be here to read your struggles with pvm.

## Commands

The options you can use are: `list`, `install`, `uninstall`, `reinstall`, `use`, `--help`

Here's a description of each command:

### `list` (or the alias `l`)

Displays the user the requested info.

`$ pvm list latest` lists all the latest python releases.

`$ pvm list all` lists all the available python versions.

`$ pvm list installed` lists all the currently installed versions.

### `install` (or the alias `i`)

Installs the requested version. If you typed the wrong version nothing will happen.

`$ pvm install latest` Installs the latest stable python version (now is 3.11.3).

`$ pvm install 3.11.0` Installs python 3.11.0.

### `uninstall` (or the alias `u`)

Uninstalls the specified version. If you typed the wrong version nothing will happen.

`$ pvm uninstall all` Uninstalls ALL the installed python versions. Be careful with this one...

`$ pvm uninstall 3.11.0` Unistalls python 3.11.0.

### `reinstall` (or the alias `r`)

Uninstalls and then installs again the specified version. If you typed the wrong version nothing will happen.

`$ pvm reinstall all` Reinstalls ALL the installed python versions.

`$ pvm reinstall 3.11.0` Reinstalls python 3.11.0.

### `use`

No alias for this one.

If already downloaded, makes active the requested version.

`$ pvm use 3.11.0`

### `--help` (or the alias `-h`)

Basically shows in the console what is written here. You can see it in action in the image [here](#python-version-manager).

# Developing

Any type of contribution is well accepted, just create a PR and i'll review it as soon as possible!

## Get started

First of all, check if you have all the tools needed to run the software installed:

- Python 3.8.x or higher (actually this project has beed written and tested with python 3.11.2)
- Pip3
- Git
- Visual C++ Redist

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
