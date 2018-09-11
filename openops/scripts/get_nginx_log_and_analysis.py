#!/usr/local/bin/python
# -*-coding=utf8-*-
from elasticsearch import Elasticsearch


es = Elasticsearch("172.16.66.217")
# 最近一小时nginx的日志
#query={"sort": [ { "@timestamp": { "order": "desc", "unmapped_type": "boolean" } } ],"query": {"filtered":{"filter":{ "range":{"@timestamp":{"gt":"now-30m"}} },"query": { "match": {"@log_name": { "query": "nginx.access","type": "phrase" } } }  } }  }
query={"sort": [ { "@timestamp": { "order": "desc", "unmapped_type": "boolean" } } ],"query": {"filtered":{"filter":{ "range":{"@timestamp":{"gt":"now-1h"}} },"query": { "match": {"@log_name": { "query": "nginx.access","type": "phrase" } } }  } }  }
# 获取日志数量
total=es.count(index="logstash-2016.12.24", body=query)['count']
data=es.search(index="logstash-2016.12.24", body=query, size=int(total))

# 获取http 状态
'''
   100-199 用于指定客户端应相应的某些动作。
   200-299 用于表示请求成功。
   300-399 用于已经移动的文件并且常被包含在定位头信息中指定新的地址信息。
   400-499 用于指出客户端的错误。
   500-599 用于支持服务器错误。
   http 常用状态码：200, 206, 301, 302, 404, 499, 500

'''
status_200 = 0
status_206 = 0
status_301 = 0
status_302 = 0
status_404 = 0
status_499 = 0
status_500 = 0
status_other = 0


'''
   UserAgent
       -			 // 平板
       go-request/0.7.0      // 监控
       Java/1.7.0_79	 //
       Mozilla/4.0 (compatible; MSIE 6.0)   // IE 6.0
       Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1; .NET CLR 2.0.50727; .NET CLR 3.0.04506.648; .NET CLR 3.5.21022) // IE 6.0
       Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0; WAC-OFU)   // IE 11.0
       Mozilla/5.0 (Windows NT 6.1; rv:45.0) Gecko/20100101 Firefox/45.0  // Firefox
       Mozilla/5.0 (Windows NT 6.1; Trident/7.0; rv:11.0; JuziBrowser) like Gecko // 橘子浏览器
       Mozilla/5.0 (Windows NT 6.1; Trident/7.0; rv:11.0) like Gecko   // IE
       okhttp/3.0.1
'''
pad_browser = 0
monitor_browser = 0
java_browser = 0
okhttp_browser = 0
ie6_browser = 0
ie9_browser = 0
firefox_browser = 0
juzi_browser = 0
others_browser = 0


allvisitips = []
iplists = []

for i in range(len(data['hits']['hits'])):
    '''
       http status
    '''
    http_status = data['hits']['hits'][i]['_source']['code']
    if http_status == "200":
        status_200 += 1
    elif http_status == "206":
        status_206 += 1
    elif http_status != "301":
        status_301 += 1
    elif http_status != "302":
        status_302 += 1
    elif http_status != "404":
        status_404 += 1
    elif http_status != "499":
        status_499 += 1
    elif http_status != "500":
        status_500 += 1
    else:
        status_other +=1

    '''
       UserAgent
    '''
    if data['hits']['hits'][i]["_source"]['agent'] == str("-"):
        pad_browser += 1
    elif data['hits']['hits'][i]["_source"]['agent'] == str("go-request/0.7.0"):
        monitor_browser += 1
    elif data['hits']['hits'][i]["_source"]['agent'] == str("Java/1.7.0_79"):
        java_browser += 1
    elif data['hits']['hits'][i]["_source"]['agent'] == str("Mozilla/4.0 (compatible; MSIE 6.0)") or data['hits']['hits'][i]["_source"]['agent'] == str("Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1; .NET CLR 2.0.50727; .NET CLR 3.0.04506.648; .NET CLR 3.5.21022)"):
        ie6_browser += 1
    elif data['hits']['hits'][i]["_source"]['agent'] == str("Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0; WAC-OFU)"):
        ie9_browser += 1
    elif data['hits']['hits'][i]["_source"]['agent'] == str("Mozilla/5.0 (Windows NT 6.1; rv:45.0) Gecko/20100101 Firefox/45.0"):
        firefox_browser += 1
    elif data['hits']['hits'][i]["_source"]['agent'] == str("Mozilla/5.0 (Windows NT 6.1; Trident/7.0; rv:11.0; JuziBrowser) like Gecko"):
        junzi_browser += 1
    else:
        others_browser += 1

    '''
       Client
    '''
    client_ip = data['hits']['hits'][i]["_source"]['remote']
    allvisitips.append(client_ip)


sortips = set(allvisitips)  #sortips是另外一个列表，里面的内容是allvisitips里面的无重复 项
ips_count = {}
for item in sortips:
    ips_count["%s" % item] = allvisitips.count(item)

result = {u"ips": ips_count,  u"httpstatus" : {"status_200": status_200 ,"status_206": status_206 ,"status_301": status_301 ,"status_302": status_302 ,"status_404": status_404 ,"status_499": status_499 ,"status_500": status_500 ,"status_other": status_other},  u"useragent":{"pad_browser" : pad_browser,"monitor_browser" : monitor_browser,"java_browser" : java_browser,"okhttp_browser" : okhttp_browser,"ie6_browser" : ie6_browser,"ie9_browser" : ie9_browser,"firefox_browser" : firefox_browser,"juzi_browser" : juzi_browser,"others_browser" : others_browser}}
print result