#!/usr/bin/env python2.7
#coding=utf-8
import json
import urllib2
import datetime
import time

# based url and required header
url = "http://ljops.monitor.ue.cn/api_jsonrpc.php"
header = {"Content-Type":"application/json"}
# request json

now_time = datetime.datetime.now()
current_timestamp =int(time.mktime(now_time.timetuple()))
last_twelve_hour_timestamp =int(time.mktime(now_time.timetuple())) - 28800
# yesterday_timestamp =int(time.mktime(yes_time.timetuple())) + 43200


def dataCenterTraffic():
    itemids = (28538, 28539)

    history_time = []
    data_in = []
    data_out = []

    for itemid in itemids:
        data = json.dumps(
        {
            "jsonrpc": "2.0",
            "method": "history.get",
            "params": {
                "output": "extend",
                "history": 3,
                "itemids":itemid,
                "time_from":last_twelve_hour_timestamp,
                "time_till":current_timestamp,
                "sortfield": "clock",
                "sortorder": "ASC",
                "limit": 500,
            },
            "auth": "1df4b3f4bec2fc4ed6c6555cfb2042e7",
            "id": 1
        })


        # create request object
        request = urllib2.Request(url,data)
        for key in header:
            request.add_header(key,header[key])
        try:
            result = urllib2.urlopen(request)
        except URLError as e:
            if hasattr(e, 'reason'):
               print 'We failed to reach a server.'
               print 'Reason: ', e.reason
            elif hasattr(e, 'code'):
               print 'The server could not fulfill the request.'
               print 'Error code: ', e.code
        else:
           response = json.loads(result.read())
           result.close()
           for host in response['result']:
              if itemid == 28538:
                   # history_time = history_time.append("%s" % time.strftime('%Y-%m-%d %H:%M:%S',time.localtime(int(host['clock']))))
                   time_string = str(time.strftime('%Y-%m-%d %H:%M:%S',time.localtime(int(host['clock']))))
                   history_time.append(time_string)
                   # print history_time
                   # print time.strftime('%Y-%m-%d %H:%M:%S',time.localtime(int(host['clock']))), float(host['value'])/ 1024.00 /1024.00
                   data_in.append("%.2f" % float(float(host['value'])/ 1024.00 /1024.00))
              else:
                   data_out.append("%.2f" % float(float(host['value'])/ 1024.00 /1024.00))
           # print response

    return history_time, data_in, data_out
    # return json.dumps(history_time), json.dumps(data_in), json.dumps(data_out)




data_time, data_in, data_out = dataCenterTraffic()
print str(data_time) + "+" + str(data_in) + "+" + str(data_out)
# print data_time + "+" + data_in + "+" + data_out
