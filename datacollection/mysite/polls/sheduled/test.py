#!/usr/bin/python

import os
def testSchedule():
        path=os.path.abspath('.')
	shell = 'echo aaa >> ~/crontab'.format(path)
	os.system(shell)
