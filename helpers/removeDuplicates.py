def removeDuplicatesFromList(myList):
	res = []
	[res.append(x) for x in myList if x not in res]
	return res