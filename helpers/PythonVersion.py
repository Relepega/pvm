from ast import parse
import re
import subprocess
from typing import TypedDict

from packaging.version import Version, VERSION_PATTERN

class PipVersion(TypedDict):
	version: str
	filename: str
	downloadUrl: str

class PythonVersion:
	def __init__(self, versionNumber: str, releaseDate: str, downloadUrl: str, installerFilename: str) -> None:
		self.versionNumber = versionNumber
		self.versionInfo = Version(versionNumber)
		self.isStable = not self.versionInfo.is_prerelease
		self.releaseDate = releaseDate
		self.downloadUrl = downloadUrl
		self.installerFilename = installerFilename
		self.pipVersion = self.getPipVersion()

	def __str__(self) -> str:
		return f'Python {self.versionInfo.major}, {self.versionNumber}, released in date {self.releaseDate}'
	
	def __repr__(self) -> str:
		return f'{self.versionInfo.major}, {self.versionNumber}, {self.isStable}, {self.releaseDate}, {self.downloadUrl}'

	def getPipVersion(self) -> PipVersion:
		# fetch the right version from here:
		# https://bootstrap.pypa.io/pip/

		if self.versionInfo >= Version('3.7.0'):
			return {
				'version': 'latest',
				'filename': f'get-pip-latest.py',
				'downloadUrl': f'https://bootstrap.pypa.io/pip/get-pip.py'
			}
		else:
			v = f'{self.versionInfo.major}.{self.versionInfo.minor}'
			return {
				'version': v,
				'filename': f'get-pip-v{v}.py',
				'downloadUrl': f'https://bootstrap.pypa.io/pip/{v}/get-pip.py'
			}


class PythonVersions(TypedDict):
	all: list[str]
	stable: list[str]
	unstable: list[str]
	classes: dict[str, PythonVersion]
	creationDate: int


def getVersionInUse() -> str:
	try:
		output = str(subprocess.check_output(['python', '-V']).strip())
		res = output.split(' ')[1][:-1]
		# print(f'"{res}"')
		return res
	except subprocess.CalledProcessError:
		pass

	return ''