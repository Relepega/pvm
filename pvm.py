import argparse
import os
import re
import shutil

from SystemHandler import Client
from helpers import (
	PYTHON_DOWNLOAD_PATH,
	PYTHON_VERSION_REGEX
)

class PVM:
	def __init__(self) -> None:
		self.client = Client()


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

	
	def uninstallSingleVersion(self, version: str) -> None:
		versionPath = os.path.join(PYTHON_DOWNLOAD_PATH, version)
		if re.match(string=version, pattern=PYTHON_VERSION_REGEX) and os.path.exists(versionPath) and os.path.isdir(versionPath):
			shutil.rmtree(versionPath)


	def uninstallParserHandler(self, version: str) -> None:
		match version.lower():
			case 'all':
				if os.path.exists(PYTHON_DOWNLOAD_PATH):
					shutil.rmtree(PYTHON_DOWNLOAD_PATH)
			case _:
				self.uninstallSingleVersion(version)

	
	def reinstallParserHandler(self, version: str) -> None:
		installed = [ f.path.split('\\')[-1] for f in os.scandir(PYTHON_DOWNLOAD_PATH) if f.is_dir() ]
		reinstall = installed if version.lower() == 'all' else [version]

		if not version in installed and not version.lower() == 'all':
			print(f'"{version}" is not a valid version or is not currently installed.')
			return

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
		installParser.set_defaults(func=self.client.installNewVersion)

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