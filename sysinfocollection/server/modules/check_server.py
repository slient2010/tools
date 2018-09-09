#!/usr/bin/env python
# -*- coding: utf-8 -*-
#
# Author : Jun
# Desc: 获取游戏进程数
# Date: 2014-10-16 
#
import os
import commands
import httplib

class checkService:
    #构造函数
   # def  __init__():
   #     pass
    
    #检查服务, TODO
    def check_service(self, service):
        self.service = service
        get_service_status = commands.getstatusoutput("ps aux |grep '%s'|grep -v grep |wc -l" % self.service)
        return int(get_service_status[1])
        #if int(get_service_status[1]) == 0:
           #os.system("echo '[%s] starting service %s'" %(self.start_time, self.service))
           #os.system("%s"%self.service)
            # 启动完成后，发邮件通知
            #self.sendmail(self.service)
    
    
    # 异常后重启服务，发送邮件通知 TODO
    def sendmail(self,service_msg):
       self.service_msg = service_msg
       dt = time.strftime('%Y-%m-%d %H:%M:%S',time.localtime(time.time()))
       self.msg = "%s %s start" % (dt, self.service_msg)
       httpClient = None
       try:
           httpClient = httplib.HTTPConnection(self.server, int(self.server_port), 30)
           httpClient.request('GET', '/sendmail/sendmail.php?msg='+self.msg)
           response = httpClient.getresponse()
       except Exception, e:
           print e
       #finally:
       if httpClient:
           httpClient.close()
       #获取HTTP对象
       h = httplib2.Http()
       # to do
       resp, content = h.request("http://%s:%d/sendmail/sendmail.php?msg="% (self.server, int(self.server_port))+self.msg)

    def kill_service():
        self.service = service
        get_service_pid = commands.getstatusoutput("ps aux |grep '%s'|grep -v grep |awk '{print $2}'" % self.service)
        get_service_pid = commands.getstatusoutput("kill %s" % get_service_pid)
        
