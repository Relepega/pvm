import copy
import os
import platform
import re
import shutil
import subprocess
import sys
import httpx
import bs4
import pickle
from datetime import datetime
import time
from rich import print
from helpers import (
	PythonVersion,
	PythonVersions,
	downloadFile,
	getAppRootPath,
	getVersionInUse,
	PYTHON_VERSION_REGEX,
	PYTHON_DOWNLOAD_PATH
)

SYMLINK_DEST = f"{os.getenv('LOCALAPPDATA')}\\Python"

class Client:
	def __init__(self):
		self.appRoot = getAppRootPath()

		self.arch = 'amd64' if platform.architecture()[0] == '64bit' else 'win32'
		self.client = httpx.Client(timeout=None, follow_redirects=True)
		self.pythonVersions: PythonVersions = { # type: ignore
			'names': [],
			'latest': [],
			'versions': [],
			'creationDate': 0
		}

		self.cachedDataExists = False


	def clientInfo(self) -> str:
		return f'Windows client ({self.arch})'

	
	def fetchAllAvailableVersions(self) -> None:
		cacheFile = os.path.join(self.appRoot, f'cache_{self.arch}')
		
		# cache already formed
		if self.cachedDataExists:
			return
		
		# if cache file exists, load it
		if os.path.exists(cacheFile) and os.path.isfile(cacheFile):
			cacheData: PythonVersions

			with open(file=cacheFile, mode='rb') as f:
				cacheData = pickle.load(f)
			
			# if cache is too old, rebuild it
			expiry = 2*60*60*1000 #2 hours
			if not self.pythonVersions['creationDate'] + expiry >= int(time.time_ns() / 1000):
				self.pythonVersions = cacheData
				return

		# fetch data
		nugetPackageId = ['python', 'python2'] if self.arch == 'amd64' else ['pythonx86', 'python2x86']

		for packageID in nugetPackageId:
			soup = bs4.BeautifulSoup(
				self.client.get(f'https://www.nuget.org/packages/{packageID}#versions-body-tab').text,
				features="lxml"
			)

			tableRows = soup.select('#version-history > table > tbody > tr')

			for tr in tableRows:
				td = tr.select('td')
				versionName: str = td[0].find('a')['title'] # type: ignore
				versionReleaseDate: str = td[2].find('span').text # type: ignore

				if not re.search(PYTHON_VERSION_REGEX, versionName):
					continue

				self.pythonVersions['names'].append(versionName)
				self.pythonVersions['versions'].append(
					PythonVersion(
						version=versionName,
						majorRelease=int(versionName.split('.')[0]),
						releaseDate=versionReleaseDate
					)
				)

		self.pythonVersions['creationDate'] = int(time.time_ns() / 1000)

		# dump data
		with open(file=cacheFile, mode='wb') as f:
			pickle.dump(obj=self.pythonVersions, file=f)
		
		self.cachedDataExists = True


	def listLatest(self) -> None:
		self.fetchAllAvailableVersions()

		print('Latest python versions:')
		print('(first 5 for each major version):')

		pv = copy.copy(self.pythonVersions['versions'])
		pv.sort(
			key=lambda p: datetime.timestamp(datetime.strptime(p.releaseDate,'%m/%d/%Y')),
			reverse=True
		)

		python3 = [p for p in pv if p.majorRelease == 3][: 5]
		python2 = [p for p in pv if p.majorRelease == 2][: 5]

		for p in python3:
			print(str(p))

		for p in python2:
			print(str(p))

		print('\nThis is a partial list. For a complete list, visit https://www.python.org/downloads/')

	
	def listAll(self) -> None:
		self.fetchAllAvailableVersions()

		print('All python versions:')
		print('(First 20 of the list)\n')

		for pv in self.pythonVersions['versions'][:20]:
			print(str(pv))

		print('\nThis is a partial list. For a complete list, visit https://www.python.org/downloads/')


	def listInstalled(self) -> None:
		installed = [ f.path.split('\\')[-1] for f in os.scandir(PYTHON_DOWNLOAD_PATH) if f.is_dir() ]		
		inUse = getVersionInUse()
		
		print('Installed python versions:\n')

		for v in installed:
			print(f"{v}{' (in use)' if v == inUse else ''}")


	def installNewVersion(self, version: str) -> None:
		self.fetchAllAvailableVersions()

		if not version in self.pythonVersions['names'] and version != 'latest':
			print(f'"{version}" is not a valid python version.')
			return
		
		ver = self.pythonVersions['names'][0] if version == 'latest' else version

		# set python file naming and architecture
		arch = 'amd64' if platform.architecture()[0] == '64bit' else 'win32'
		filename = f'python-{ver}-embed-{arch}.zip'
		downloadUrl = f'https://www.python.org/ftp/python/{ver}/{filename}'

		# set files path
		appRootPath = getAppRootPath()
		offlineZipPath = os.path.join(appRootPath, filename)
		unpackedPythonPath = os.path.join(PYTHON_DOWNLOAD_PATH, ver)
		getPipScriptPath = os.path.join(appRootPath, 'get-pip.py')

		print(f'Downloading "{filename}" ...')
		# download python zip file
		if not downloadFile(url=downloadUrl, absoluteFilePath=offlineZipPath):
			print('File not downloaded.')
			return
		
		# unpack python
		shutil.unpack_archive(offlineZipPath, unpackedPythonPath)

		# download "get-pip.py" if not already downloaded
		if not os.path.exists(getPipScriptPath):
			print('Downloading newest "get-pip.py" ...')
			downloadFile(url='https://bootstrap.pypa.io/get-pip.py', absoluteFilePath=getPipScriptPath)
			print('Done!')

		# fix site-packages (https://stackoverflow.com/a/68891090)
		with open(file=os.path.join(unpackedPythonPath, f"python{''.join(ver.split('.')[:2])}._pth"), mode='a', encoding='utf-8') as f:
			f.write(r"Lib\site-packages")

		# install pip into fresh python download
		print('Installing "pip" package...')

		p = subprocess.check_output(
			[os.path.join(unpackedPythonPath, "python.exe"), getPipScriptPath],
			shell=True,
			stderr=subprocess.DEVNULL
		)

		# remove python zip
		os.remove(filename)

		# set symlink
		self.symlinkDownloadedVersion(ver)

		print(f'Python {ver} installed successfully!')
		print('-----------------------------------------------')


	def symlinkDownloadedVersion(self, version: str) -> None:
		versionPath = os.path.join(PYTHON_DOWNLOAD_PATH, version)
		if not re.match(string=version, pattern=PYTHON_VERSION_REGEX) or not os.path.exists(versionPath) or not os.path.isdir(versionPath):
			print(f'Python "{version}" is not installed. Try installing it first...')
			return

		command = f"New-Item -Force -ItemType SymbolicLink -Path '{SYMLINK_DEST}' -Target '{versionPath}'"
		p = subprocess.Popen(
			[
				"powershell.exe", 
				"-noprofile", "-c",
				f"""
				Start-Process -WindowStyle hidden -Verb RunAs -Wait powershell.exe -Args "{command}"
				"""
			],
			stdout=sys.stdout
		)
		p.communicate()