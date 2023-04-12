from ast import parse
import re
import subprocess
from typing import TypedDict

from packaging.version import parse as version_parse

PYTHON_VERSION_REGEX = r'^\d+\.\d+\.\d+(\.\d+)?$'
PYTHON_VERSION_REGEX_SUBPROCESS_VER = r'\d+\.\d+\.\d+(\.\d+)?'

class PipVersion(TypedDict):
	version: str
	filename: str
	downloadUrl: str

class PythonVersion:
	def __init__(self, versionNumber: str, majorRelease: int, releaseDate: str, downloadUrl: str) -> None:
		self.versionNumber = versionNumber
		self.stable = True if re.match(PYTHON_VERSION_REGEX, versionNumber) else False
		self.majorRelease = majorRelease
		self.releaseDate = releaseDate
		self.downloadUrl = downloadUrl
		self.filename = downloadUrl.split('/')[-1] + '.zip'
		self.versionInfo = version_parse(versionNumber)
		self.pipVersion = self.getPipVersion()

	def __str__(self) -> str:
		return f'Python {self.majorRelease}, {self.versionNumber}, released in date {self.releaseDate}'
	
	def __repr__(self) -> str:
		return f'{self.majorRelease}, {self.versionNumber}, {self.stable}, {self.releaseDate}, {self.downloadUrl}'

	def getPipVersion(self) -> PipVersion:
		major = self.versionInfo.major
		minor = self.versionInfo.minor

		# fetch the right version from here:
		# https://bootstrap.pypa.io/pip/
		v = '2.7'

		if major >= 3:
			if minor <= 2:
				v = '3.2'
			elif minor <= 3:
				v = '3.3'
			elif minor <= 4:
				v = '3.4'
			elif minor <= 5:
				v = '3.5'
			elif minor <= 6:
				v = '3.6'
			else:
				return {
					'version': 'latest',
					'filename': f'get-pip-latest.py',
					'downloadUrl': f'https://bootstrap.pypa.io/pip/get-pip.py'
				}
		
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
		output = subprocess.check_output(['python', '-V'])
		v = re.search(PYTHON_VERSION_REGEX_SUBPROCESS_VER, output.decode())
		res = v.group() # type: ignore
		# print(res)
		return res
	except subprocess.CalledProcessError:
		pass

	return ''