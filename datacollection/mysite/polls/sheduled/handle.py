#!/usr/bin/python
# -*- coding:utf8 -*-

#from client import 
import client
import sys
import MySQLdb
 

def handle():
    log_type = sys.argv[1]
    testdata = client.getdata(log_type)
    if testdata == 0:
        data = '' 
    else:
        try:
            conn=MySQLdb.connect(host='localhost',user='root',passwd='',db='test',port=3306)
            conn.set_character_set('utf8')
            cur=conn.cursor()
            cur.execute('SET NAMES utf8;')
            cur.execute('SET CHARACTER SET utf8;')
            cur.execute('SET character_set_connection=utf8;')
            data = eval(testdata)[0]['hits']['hits']
            values = []
            for i in range(len(data)):
                #print data[i]
                lid = data[i]['_id']
                logname = data[i]['_index']
                game = data[i]['_source']['@log_name']
                thread = data[i]['_source']['thread']
                loglevel = data[i]['_source']['level']
                logtime = data[i]['_source']['@timestamp']
                loginfo = data[i]['_source']['message']
                need_insert = [r'%s'%lid, r'%s'%logname, r'%s'%game, r'%s'%thread, r'%s'%loglevel, r'%s'%logtime, r'%s'%loginfo]
                check = cur.execute("select * from test where lid = '%s'"%lid)
                if check == 0:
                    if need_insert not in values:
                        values.append(need_insert)
            cur.executemany('insert into test values(%s, %s, %s, %s, %s, %s, %s)',values)
            conn.commit()
            cur.close()
            conn.close()
        except MySQLdb.Error,e:
            print "Mysql Error %d: %s" % (e.args[0], e.args[1])


handle()
