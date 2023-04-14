import os
import platform
import shutil
import subprocess
import sys
from typing import Callable
import zipfile
import httpx
import pickle
from datetime import datetime
import time
from packaging.version import Version
from rich import print
from helpers import (
	PythonVersion,
	PythonVersions,
	downloadFile,
	getAppRootPath,
	getVersionInUse,
	fetchJson,
	rmPath,
	PYTHON_DOWNLOAD_PATH
)

SYMLINK_DEST = f"{os.getenv('LOCALAPPDATA')}\\Python"

class Client:
	def __init__(self, checkIfValidPythonVersion: Callable) -> None:
		self.appRoot = getAppRootPath()
		self.__checkIfValidPythonVersion = checkIfValidPythonVersion

		self.arch = 'amd64' if platform.architecture()[0] == '64bit' else 'win32'
		self.httpClient = httpx.Client(
			timeout=httpx.Timeout(10.0, connect=60.0), # A client with a 60s timeout for connecting, and a 10s timeout elsewhere.
			follow_redirects=True
		)
		self.pythonVersions: PythonVersions = {
			'all': [],
			'stable': [],
			'unstable': [],
			'classes': {},
			'creationDate': 0
		}

		self.cachedDataExists = False


	def __parsePythonPackages(self, items: list[dict]) -> None:
		for version in items:
			versionNumber = version['catalogEntry']['version']
			packageContent = version['catalogEntry']['packageContent']
			dt = datetime.fromisoformat(version['catalogEntry']['published']) # ISO 8601 date parsing
			
			pv: Version = self.__checkIfValidPythonVersion(versionNumber)

			self.pythonVersions['all'].insert(0, versionNumber)

			if pv.is_prerelease:
				self.pythonVersions['unstable'].insert(0, versionNumber)
			else:
				self.pythonVersions['stable'].insert(0, versionNumber)

			vnNoDash = versionNumber.replace('-', '')
			installerFilename = packageContent.split('/')[-1] + '.zip' if pv.major >= 3 else f"python-{vnNoDash}{'.amd64' if self.arch == '64bit' else ''}.msi"
			downloadUrl = packageContent if pv.major >= 3 else f"https://www.python.org/ftp/python/{pv.base_version}/{installerFilename}"

			self.pythonVersions['classes'][versionNumber] = PythonVersion(
				versionNumber=versionNumber,
				releaseDate=f'{dt.month}/{dt.day}/{dt.year}',
				downloadUrl=downloadUrl,
				installerFilename=installerFilename
			)

	
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
		nugetPackages = ['python2', 'python'] if self.arch == 'amd64' else ['python2x86', 'pythonx86']

		for packageID in nugetPackages:
			url = f'https://api.nuget.org/v3/registration5-semver1/{packageID}/index.json'
			microsoftInconsistentPaginationElements: list[dict] = fetchJson(client=self.httpClient, url=url)['items']

			if 'python2' in packageID:
				self.__parsePythonPackages(items=microsoftInconsistentPaginationElements[0]['items']) # type: ignore
			else:
				for item in microsoftInconsistentPaginationElements:
					paginationElements = fetchJson(self.httpClient, item['@id'])
					self.__parsePythonPackages(items=paginationElements['items'])

		self.pythonVersions['creationDate'] = int(time.time_ns() / 1000)

		# dump data
		with open(file=cacheFile, mode='wb') as f:
			pickle.dump(obj=self.pythonVersions, file=f)

		self.cachedDataExists = True


	def listLatest(self) -> None:
		self.fetchAllAvailableVersions()

		print('Latest python versions:')
		print('(first 5 for each major version):')

		pv = [p for p in self.pythonVersions['classes'].values() if p.isStable]
		pv.sort(
			key=lambda p: datetime.strptime(p.releaseDate, '%m/%d/%Y'),
			reverse=True
		)

		python3 = [p for p in pv if p.versionInfo.major == 3][: 5]
		python2 = [p for p in pv if p.versionInfo.major == 2][: 5]

		for p in python3:
			print(str(p))

		for p in python2:
			print(str(p))

		print('\nThis is a partial list. For a complete list, visit https://www.python.org/downloads/')

	
	def listAll(self) -> None:
		self.fetchAllAvailableVersions()

		print('All python versions:')
		print('(First 20 of the list)\n')

		# last20 = list(self.pythonVersions['classes'].values())[-20:]
		last20 = list(self.pythonVersions['classes'].values())
		last20.reverse()

		for pv in last20:
			print(str(pv))

		# print('\nThis is a partial list. For a complete list, visit https://www.python.org/downloads/')


	def listInstalled(self) -> None:
		print('Installed python versions:\n')

		if not os.path.exists(PYTHON_DOWNLOAD_PATH):
			print("No installed version found.")
			return

		installed = [ f.path.split('\\')[-1] for f in os.scandir(PYTHON_DOWNLOAD_PATH) if f.is_dir() ]
		inUse = getVersionInUse()

		if len(installed) == 0:
			print("No installed version found.")
			return

		for v in installed:
			print(f"{v}{' (in use)' if v in inUse else ''}")


	def installNewVersion(self, version: str) -> None:
		self.fetchAllAvailableVersions()

		if not version in self.pythonVersions['all'] and version != 'latest':
			print(f'"{version}" is not a valid python version.')
			return
		
		pythonVersion = self.pythonVersions['classes'][self.pythonVersions['stable'][0]] if version == 'latest' else self.pythonVersions['classes'][version]

		# set system paths
		unpackedPythonPath = os.path.join(PYTHON_DOWNLOAD_PATH, pythonVersion.versionNumber)
		offlineFilePath = os.path.join(self.appRoot, pythonVersion.installerFilename)

		if os.path.exists(unpackedPythonPath):
			rmPath(unpackedPythonPath)

		# download python zip/installer file
		print(f'Downloading "{pythonVersion.installerFilename}" ...')
		if not downloadFile(url=pythonVersion.downloadUrl, absoluteFilePath=offlineFilePath, client=self.httpClient):
			print('File not downloaded.')
			return
		
		if pythonVersion.versionInfo.major == 2:
			# unpacking installer using msiexec
			print("Unpacking installer data...")

			try:
				p = subprocess.check_output(
					f'msiexec /n /a "{offlineFilePath}" /qn TARGETDIR={unpackedPythonPath}',
					shell=True,
					stderr=subprocess.STDOUT # "subprocess.DEVNULL" for no output
				)
			except Exception:
				print("Couldn't unpack the requested data. Aborting...")

			print("Done...")

			# move all the 'DLLs' files to '{unpackedPythonPath}'
			[shutil.move(f.path, os.path.join(unpackedPythonPath, f.name)) for f in os.scandir(os.path.join(unpackedPythonPath, 'DLLs'))]

			# delete 'DLLs' directory
			os.removedirs(os.path.join(unpackedPythonPath, 'DLLs'))

			# delete installer leftover
			os.remove(os.path.join(unpackedPythonPath, pythonVersion.installerFilename))

			# install pip into fresh python download
			print('Installing "pip" package ...')
			pythonExe = os.path.join(unpackedPythonPath, "python.exe")

			try:
				p = subprocess.check_output(
					f'{pythonExe} -m ensurepip --default-pip && {pythonExe} -m pip install --upgrade pip',
					shell=True,
					stderr=subprocess.STDOUT # "subprocess.DEVNULL" for no output
				)
			except Exception:
				pass

			print('Done!')
		else:
			self.python3Install(
				pythonVersion=pythonVersion,
				unpackedPythonPath=unpackedPythonPath,
				offlineZipPath=offlineFilePath
			)

		# remove python zip
		os.remove(offlineFilePath)

		# set symlink
		if self.symlinkDownloadedVersion(pythonVersion.versionNumber):
			print(f'Python {pythonVersion.versionNumber} installed successfully!')
		else:
			print(f"Error creating the symlink. Python {pythonVersion.versionNumber} wasn't set as active version.")

		print('-----------------------------------------------')

	
	def python3Install(self, pythonVersion: PythonVersion, unpackedPythonPath: str, offlineZipPath: str) -> None:
		# set system paths
		getPipScriptRootPath = os.path.join(unpackedPythonPath, 'Tools')
		getPipScriptPath = os.path.join(getPipScriptRootPath, pythonVersion.pipVersion['filename'])

		print("Hacking python's folder :) ...")

		# unpack python
		shutil.unpack_archive(offlineZipPath, unpackedPythonPath)

		# remove all folders apart 'tools'
		[rmPath(f.path) for f in os.scandir(unpackedPythonPath) if not 'tools' in f.path] # type: ignore

		# rename 'tools' folder to avoid conflicts in next step
		shutil.move(os.path.join(unpackedPythonPath, 'tools'), os.path.join(unpackedPythonPath, 'pythonContainer'))

		# move all the 'pythonContainer' subfolders to '{unpackedPythonPath}'
		[shutil.move(f.path, os.path.join(unpackedPythonPath, f.name)) for f in os.scandir(os.path.join(unpackedPythonPath, 'pythonContainer'))]

		# delete 'pythonContainer'
		os.removedirs(os.path.join(unpackedPythonPath, 'pythonContainer'))

		# zip all 'Lib' content apart 'site-packages' to 'pythonXXX.zip'
		zipFilename = 'python' + ''.join(pythonVersion.versionNumber.split('.')[:2]) + '.zip'

		with zipfile.ZipFile(os.path.join(unpackedPythonPath, zipFilename), "w") as zf:
			libRoot = os.path.join(unpackedPythonPath, 'Lib')

			for dirname, subdirs, files in os.walk(libRoot):
				if 'site-packages' in subdirs:
					subdirs.remove('site-packages')
				
				relativeDirname = dirname.replace(libRoot, '')

				if not relativeDirname == '':
					zf.write(dirname, arcname=relativeDirname)
				
				for filename in files:
					absPathFilename = os.path.join(dirname, filename)
					zf.write(absPathFilename, arcname=absPathFilename.replace(libRoot, ''))

		# remove all directories from 'Lib' apart 'site-packages'
		[rmPath(f.path) for f in os.scandir(os.path.join(unpackedPythonPath, 'Lib')) if not 'site-packages' in f.path] # type: ignore

		# move all the 'DLLs' files to '{unpackedPythonPath}'
		[shutil.move(f.path, os.path.join(unpackedPythonPath, f.name)) for f in os.scandir(os.path.join(unpackedPythonPath, 'DLLs'))]

		# delete 'DLLs' directory
		os.removedirs(os.path.join(unpackedPythonPath, 'DLLs'))

		# fix site-packages (https://stackoverflow.com/a/68891090)
		with open(file=os.path.join(unpackedPythonPath, f"python{''.join(pythonVersion.versionNumber.split('.')[:2])}._pth"), mode='a', encoding='utf-8') as f:
			f.writelines([
				zipFilename + '\n',
				'.\n',
				'\n',
				'# Uncomment to run site.main() automatically\n',
				'#import site\n',
				'\n',
				r"Lib\site-packages"
			])

		print('Done!')

		# download "get-pip.py" if not already downloaded
		if not os.path.exists(getPipScriptPath):
			print(f'Downloading "get-pip.py" from "{pythonVersion.pipVersion["downloadUrl"]}" ...')
			print('Done!')
			downloadFile(url=pythonVersion.pipVersion["downloadUrl"], absoluteFilePath=getPipScriptPath, client=self.httpClient)

		# install pip into fresh python download
		print('Installing "pip" package ...')

		try:
			pythonExe = os.path.join(unpackedPythonPath, "python.exe")
			p = subprocess.check_output(
				f"{pythonExe} {getPipScriptPath} && {pythonExe} -m pip install --upgrade pip",
				shell=True,
				stderr=subprocess.STDOUT # "subprocess.DEVNULL" for no output
			)
		except Exception:
			pass

		print('Done!')


	def symlinkDownloadedVersion(self, version: str) -> bool:
		versionPath = os.path.join(PYTHON_DOWNLOAD_PATH, version)
		if not os.path.exists(versionPath) or not os.path.isdir(versionPath):
			print(f'Python "{version}" is not installed. Try installing it first...')
			return False

		command = f"New-Item -Force -ItemType SymbolicLink -Path '{SYMLINK_DEST}' -Target '{versionPath}'"

		print('Making symlink ...')
		
		try:
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

			print('Done!')

			return True if p.returncode == 0 else False
		except Exception:
			print("Couldn't create the symlink, exiting ...")
			return False
