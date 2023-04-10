import httpx
import bs4

def fetchHtml(url) -> bs4.BeautifulSoup:
	res = httpx.get(url=url, timeout=None)
	return bs4.BeautifulSoup(res.text, features="lxml")

def fetchJson(client: httpx.Client, url: str) -> dict:
	res = client.get('https://api.nuget.org/v3-flatcontainer/python/index.json')
	return res.json()