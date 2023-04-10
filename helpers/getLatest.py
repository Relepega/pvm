from typing import Union
from .removeDuplicates import removeDuplicatesFromList

def getLatestMinor(versions: list[str], major: int) -> dict:
	i = 0
	latestMinor = 0
	allMinors: list[int] = []

	while(True):
		for version in versions:
			if version.startswith(f'{major}.{i}'):
				latestMinor = i
				allMinors.append(i)

		if latestMinor != i:
			break

		i += 1

	return {'latestMinor': latestMinor, 'allMinors': removeDuplicatesFromList(allMinors)}


def getLatestBugfix(versions: list[str], major: int, minor: int) -> dict:
	i = 0
	latestBugfix = 0
	allBugfixes: list[int] = []
	while(True):
		for version in versions:
			if version == (f'{major}.{minor}.{i}'):
				latestBugfix = i

		if latestBugfix != i:
			break

		i += 1

	return {'latestBugfix': latestBugfix, 'allBugfixes': removeDuplicatesFromList(allBugfixes)}