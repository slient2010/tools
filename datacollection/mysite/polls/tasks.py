#!/usr/bin/python
# -*- coding:utf8 -*-

from sheduled import client
from celery import task
from polls.models import DataSave

@task()
def addsum(x=1, y=3):
    return x+y

@task()
def savedata(log_type="error"):
    testdata = client.getdata(log_type)
    if testdata == 0:
        data = '' 
    else:
        data = eval(testdata)[0]['hits']['hits']
        for i in range(len(data)):
            lid = data[i]['_id']
            logname = data[i]['_index']
            game = data[i]['_source']['@log_name']
            thread = data[i]['_source']['thread']
            loglevel = data[i]['_source']['level']
            logtime = data[i]['_source']['@timestamp']
            loginfo = data[i]['_source']['message']
            need_insert = DataSave(lid=r'%s'%lid, logname=r'%s'%logname, game=r'%s'%game, thread=r'%s'%thread, loglevel=r'%s'%loglevel, logtime=r'%s'%logtime, loginfo=r'%s'%loginfo)
            need_insert.save()
