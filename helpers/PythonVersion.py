import re
import subprocess
from typing import TypedDict

PYTHON_VERSION_REGEX = r'^\d+\.\d+\.\d+(\.\d+)?$'
PYTHON_VERSION_REGEX_SUBPROCESS_VER = r'\d+\.\d+\.\d+(\.\d+)?'

class PythonVersion:
	def __init__(self, version: str, majorRelease: int, releaseDate: str, downloadUrl: str) -> None:
		self.version = version
		self.stable = True if re.match(PYTHON_VERSION_REGEX, version) else False
		self.majorRelease = majorRelease
		self.releaseDate = releaseDate
		self.downloadUrl = downloadUrl

	def __str__(self) -> str:
		return f'Python {self.majorRelease}, {self.version}, released in date {self.releaseDate}'
	
	def __repr__(self) -> str:
		return f'{self.majorRelease}, {self.version}, {self.stable}, {self.releaseDate}, {self.downloadUrl}'


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