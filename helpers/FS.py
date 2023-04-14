import shutil
import sys
import os

def getAppRootPath() -> str:
    # https://www.codegrepper.com/code-examples/python/pyinstaller+onefile+current+working+directory
    
	# determine if the application is a frozen `.exe` (e.g. pyinstaller --onefile) 
	if getattr(sys, 'frozen', False):
		return os.path.dirname(sys.executable)
	# or a script file (e.g. `.py` / `.pyw`)
	# elif __file__:
	else:
		return os.path.abspath(os.path.join(os.path.dirname(__file__), '..'))
	
PYTHON_DOWNLOAD_PATH = os.path.join(getAppRootPath(), 'Python')

def rmPath(path: str) -> None:
	try:
		if os.path.isdir(path):
			shutil.rmtree(path)
		else:
			os.remove(path)
	except Exception:
		print("Couldn't delete the file/folder. Exiting...")
		sys.exit()