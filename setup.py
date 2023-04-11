import sys
from cx_Freeze import setup, Executable # type: ignore

version = '0.1.1'

win_exe = Executable(
	script='pvm.py',
	icon=None
)

setup(
	name="pvm",
	version=version,
	description="Python Version Manager",
	author="Relepega",
	url='https://github.com/Relepega/PythonVersionManager',
	download_url='https://github.com/Relepega/PythonVersionManager/releases',
	license='GNU GPLv3',
	license_files='.\\LICENSE',
	project_urls='https://github.com/Relepega/PythonVersionManager',
	options={
		"build_exe": {
			"optimize": 2,
			'build_exe': f'.\\build\\pvm-win-v{version}',
			'include_files': ['scripts', 'LICENSE'],
			'include_msvcr': True,
			'silent_level': 1
		}
	},
	executables=[win_exe],
)