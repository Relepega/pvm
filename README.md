# Python Version Manager

Python Version Manager (PVM for short) is a project that has been inspired by [nvm for windows](https://github.com/coreybutler/nvm-windows) and aims to easily manage multiple python enviroments on your windows system.

![](media/pvm.png)

# Table of contents

- [Python Version Manager](#python-version-manager)
- [Table of contents](#table-of-contents)
- [How to use](#how-to-use)
	- [Install](#install)
	- [Uninstall](#uninstall)
- [Developing](#developing)
	- [Get started](#get-started)
- [FAQ](#faq)

# How to use

This is a CLI application, so you need to open a terminal and use it from there.

## Install

At the moment it is available only as a portable version and the installation process is not too complex:

1. Download the latest version [from here](https://github.com/Relepega/PythonVersionManager/releases).
2. Unzip the app where you won't move it never again until you uninstall it.
3. Inside the app folder, go to the subfolder `scripts`, right-click on `install.bat` and then click on `Run as administrator`. If the UAC kicks in, please click on `yes`.
4. Wait until you don't see any flying terminal and then close the folder.
5. Restart your terminal.
6. Type `pvm -h` and press enter. If the installation was successful then you will see an output like in [this image](#python-version-manager).
7. Profit 🎉!

Please keep in mind that you must install the app for each user you want to use it with. 

## Uninstall

It's basically the same process as the installation one:

1. Open the app folder.
2. Go to the `scripts` subfolder, right-click on `uninstall.bat` and then click on `Run as administrator`. If the UAC kicks in, please click on `yes`.
3. Wait until you don't see any flying terminal and then close the folder.
4. Restart your terminal.
5. Type `pvm -h` and press enter. If pvm was uninstalled successfully now in the terminal you should see an error.
6. Profit 🎉!
7. (optional) Please let me know why you decided to uninstall pvm: fill [this form](https://github.com/Relepega/PythonVersionManager/issues/new) and i'll be here to read your struggles with pvm.

Please keep in mind that you must uninstall the app for each user you want to use it with. 

# Developing

Any type of contribution is well accepted, just create a PR and i'll review it as soon as possible!

## Get started

1. Make sure to have these tools installed on your system

- Go 1.21.0 (minimum required version)
- Git

2. When you're sure that you have installed them correctly, proceed by cloning the repository

`$ git clone https://github.com/Relepega/PythonVersionManager.git`

3. Hop into the project directory

`$ cd pvm`

4. Export the `GO111MODULE` Enviroment Variable to ensure all modules to be installed correctly

- On Linux/MacOS: `$ export GO111MODULE=on`
- On Windows: `> set GO111MODULE=on`

5. Install the dependencies

`$ go mod download`

6. You're now ready to go! If you want to build from source because you don't trust some random guy on the internet, run the build script:

`build.bat`

# FAQ

Q: Why reinvent the wheel when [pyenv-win](https://github.com/pyenv-win/pyenv-win) exists and does the same thing?

A: Forgive my ignorance, but i didn't know of its existence until the creation of this project.
