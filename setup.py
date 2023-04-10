import sys
from cx_Freeze import setup, Executable # type: ignore

win_exe = Executable(
	script='pvm.py',
	icon=None
)

setup(
	name="pvm",
	version="0.1",
	description="Python Version Manager",
	author="Relepega",
	url='',
	download_url='',
	license='',
	license_files='',
	project_urls='',
	options={
		"build_exe": {
			"optimize": 2,
			'build_exe': '.\\build\\pvm-win',
			'include_files': ['scripts'],
			'include_msvcr': True,
			'silent_level': 1
		}
	},
	executables=[win_exe],
)