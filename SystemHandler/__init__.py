import sys

match sys.platform:
	case 'linux','linux2':
		pass
	case 'darwin':
		pass
	case 'win32':
		from SystemHandler.Windows import Client
	case _:
		print('OS Not supported...')
		sys.exit()