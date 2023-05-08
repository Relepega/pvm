import argparse
import os
import re
from sys import exit
import shutil
from typing import Union
from packaging.version import Version, VERSION_PATTERN
from rich import print

from SystemHandler import Client
from helpers import PYTHON_DOWNLOAD_PATH

class PVM:
	def __init__(self) -> None:
		self.client = Client(self.__checkIfValidPythonVersion)


	def __checkIfValidPythonVersion(self, version: str) -> Union[Version, None]:
		try:
			return Version(version)
		except:
			print(f'"{version}" is not a valid version.')
			exit(0)


	def listParserHandler(self, mode: str) -> None:
		match mode:
			case 'latest':
				self.client.listLatest()
			case 'all':
				self.client.listAll()
			case 'installed':
				self.client.listInstalled()
			case _:
				self.client.listLatest()


	def installNewVersion(self, version: str) -> None:
		self.__checkIfValidPythonVersion(version)

		versionPath = os.path.join(PYTHON_DOWNLOAD_PATH, version)
		if os.path.exists(versionPath) and os.path.isdir(versionPath):
			print(f'Python {version} is already installed. Please use the command "reinstall" instead.')
		else:
			self.client.installNewVersion(version)

	
	def uninstallSingleVersion(self, version: str) -> None:
		versionPath = os.path.join(PYTHON_DOWNLOAD_PATH, version)
		if os.path.exists(versionPath) and os.path.isdir(versionPath):
			shutil.rmtree(versionPath)


	def uninstallParserHandler(self, version: str) -> None:
		self.__checkIfValidPythonVersion(version)

		match version.lower():
			case 'all':
				if os.path.exists(PYTHON_DOWNLOAD_PATH):
					shutil.rmtree(PYTHON_DOWNLOAD_PATH)
			case _:
				self.uninstallSingleVersion(version)


	def reinstallParserHandler(self, version: str) -> None:
		if not version.lower() == 'all':
			self.__checkIfValidPythonVersion(version)

		installed = [ f.path.split('\\')[-1] for f in os.scandir(PYTHON_DOWNLOAD_PATH) if f.is_dir() ]

		if not version in installed and not version.lower() == 'all':
			print(f'"{version}" is not a valid version or is not currently installed.')
			return

		reinstall = installed if version.lower() == 'all' else [version]

		for v in reinstall:
			self.uninstallSingleVersion(v)
			self.client.installNewVersion(v)


	def start(self) -> None:
		parser = argparse.ArgumentParser(
			prog='pvm',
			description='(Yet another) Python Version Manager',
			usage='%(prog)s'
		)
		subparsers = parser.add_subparsers(dest='command')

		listParser = subparsers.add_parser(
			name='list',
			aliases='l',
			help='Displays the requested info to the user.',
			usage='%(prog)s [all, latest, installed]'
		)
		listParser.add_argument('mode')
		listParser.set_defaults(func=self.listParserHandler)

		installParser = subparsers.add_parser(
			name='install',
			aliases='i',
			help='Installs the requested version. "latest" parses to the latest stable version available.',
			usage='%(prog)s [latest, version_number]'
		)
		installParser.add_argument('version')
		installParser.set_defaults(func=self.installNewVersion)

		uninstallParser = subparsers.add_parser(
			name='uninstall',
			aliases='u',
			help='Uninstalls the specified version.',
			usage='%(prog)s [all, version_number]'
		)
		uninstallParser.add_argument('version')
		uninstallParser.set_defaults(func=self.uninstallParserHandler)

		reinstallParser = subparsers.add_parser(
			name='reinstall',
			aliases='r',
			help='Uninstalls and then installs again the specified version.',
			usage='%(prog)s version_number'
		)
		reinstallParser.add_argument('version')
		reinstallParser.set_defaults(func=self.reinstallParserHandler)
		
		useParser = subparsers.add_parser(
			name='use',
			help='If already downloaded, activates the requested version.',
			usage='%(prog)s version_number'
		)
		useParser.add_argument('version')
		useParser.set_defaults(func=self.client.symlinkDownloadedVersion)

		args = parser.parse_args()
		args_ = vars(args).copy()

		args_.pop('command', None)
		args_.pop('func', None)

		try:
			args.func(**args_)
		except:
			pass


if __name__ == "__main__":
	pvm = PVM()
	pvm.start()