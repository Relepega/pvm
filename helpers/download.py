import os
import sys
import httpx
import rich.progress
from sys import exit


def downloadFile(url: str, absoluteFilePath: str, client: httpx.Client) -> bool:
	try:
		downloadFile = open(file=absoluteFilePath, mode='wb')

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

	except (httpx.HTTPError, httpx.StreamError):
		print('File not downloaded. Try again later.')
		sys.exit()

	return True