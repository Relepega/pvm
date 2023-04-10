import httpx
import bs4

def fetchHtml(url) -> bs4.BeautifulSoup:
	res = httpx.get(url=url, timeout=None)
	return bs4.BeautifulSoup(res.text, features="lxml")