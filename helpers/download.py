import os
import re
import shutil
import httpx
import rich.progress
from sys import exit

from .FS import (
	getAppRootPath
)

def downloadFile(url: str, absoluteFilePath: str) -> bool:
	try:
		downloadFile = open(file=absoluteFilePath, mode='wb')

		client = httpx.Client(timeout=None, follow_redirects=True)
	
		with client.stream("GET", url) as response:
			if response.status_code == 404:
				downloadFile.close()
				os.remove(absoluteFilePath)
				return False

			total = int(response.headers["Content-Length"])

			with rich.progress.Progress(
				"[progress.percentage]{task.percentage:>3.2f}%",
				rich.progress.BarColumn(bar_width=None),
				rich.progress.DownloadColumn(),
				rich.progress.TransferSpeedColumn(),
			) as progress:
				downloadTask = progress.add_task("Download", total=total)
				for chunk in response.iter_bytes():
					downloadFile.write(chunk)
					progress.update(downloadTask, completed=response.num_bytes_downloaded)

		downloadFile.close()

	except httpx.ConnectTimeout:
		print('File not downloaded. Try again later.')
		return False
	except httpx.ConnectError:
		print('This version of python does not exist.')
		return False

	return True