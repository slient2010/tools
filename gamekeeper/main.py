#!/usr/bin/python
# -*-coding:utf-8-*-
import sys, os, time, atexit, string
from signal import SIGTERM 
import time
import httplib2
import commands
import ConfigParser
import codecs

'''
    程序说明：
        守护指定程序
    注意：
    1.如果你的守护进程是由inetd启动的，不要这样做！inetd完成了
    所有需要做的事情，包括重定向标准文件描述符，需要做的事情只有
    chdir() 和 umask()了
    2.程序是单线程，执行速度1秒检查一个游戏服务进程
    
'''

class Daemon:
    #def __init__(self, pidfile, stdin=stdin, stdout=stdout, stderr=stderr):
    def __init__(self, pidfile, stdin, stdout, stderr):
        self.pidfile= pidfile 
        self.stdin = stdin
        self.stdout = stdout
        self.stderr = stderr

    def _daemonize(self):
        try:
            pid = os.fork()
            # 退出主进程
            if pid > 0:
                sys.exit(0)
        except OSError, e:
            sys.stderr.wirte('fork #1 failed: %d (%s)\n' % (e.errno, e.strerror))
            sys.exit(1)
        os.chdir("/")
        os.umask(0)
        os.setsid()

        try:
            pid = os.fork()
            if pid > 0:
                sys.exit(0)
        except OSError, e:
            sys.stderr.wirte('fork #2 failed: %d (%s)\n' % (e.errno, e.strerror))
            sys.exit(1)
        #进程已经是守护进程，重定向标准文件描述符
        for f in sys.stdout, sys.stderr: f.flush()
        si = file(self.stdin, 'r')
        so = file(self.stdout,'a+')
        se = file(self.stderr,'a+',0)
        os.dup2(si.fileno(), sys.stdin.fileno())
        os.dup2(so.fileno(), sys.stdout.fileno())
        os.dup2(se.fileno(), sys.stderr.fileno())
            
        #创建processid文件
        atexit.register(self.delpid)
        pid = str(os.getpid())
        file(self.pidfile,'w+').write('%s\n' % pid)

    def delpid(self):
        os.remove(self.pidfile)

    #检查服务，若未起来，启动服务
    def check_service(self, service):
        self.service = service
        get_service_status = commands.getstatusoutput("ps aux |grep '%s'|grep -v grep |wc -l" % self.service)
        if int(get_service_status[1]) == 0:
            os.system("echo 'starting service %s'" % self.service)
            # 启动完成后，发邮件通知
            self.sendmail(self.service)
             
    def _run(self):
        while True:
            #服务所在路径，若不存在，退出程序
            if not os.path.exists(self.game_path):
                sys.exit(0)
            
            get_domain_dir = commands.getstatusoutput('ls %s |egrep "%s$"' % (self.game_path, self.domain))
            if int(get_domain_dir[0]) == 0:
                domain_data = get_domain_dir[1].split("\n")
                for i in range(0, len(domain_data)):
                    # 得到每个服所在路径 
                    game_domain_path = self.game_path+"/"+domain_data[i]+"/socket"
                    get_sub_domain_dir = commands.getstatusoutput('ls %s |egrep "rxtl_s[0-9]+"' % game_domain_path)
                    if int(get_sub_domain_dir[0]) == 0:
                        sub_domain_dir = get_sub_domain_dir[1].split("\n")
                        for j in range(0, len(sub_domain_dir)):
                            # 逐一检查服务
                            # 选出进程数
                            #game_process_num = commands.getstatusoutput("ps -ef | egrep %s | grep -v 'egrep' |wc -l " % sub_domain_dir[j]).read().strip()
                            game_process_num = commands.getstatusoutput("ps -ef | egrep %s | grep -v 'egrep' |wc -l " % sub_domain_dir[j])
                            # 当前进程数为7
                            if game_process_num != self.services_num:
                                game_base_dir = game_domain_path + "/" + sub_domain_dir[j] + "/servers/"
                                loginServer = game_base_dir + "loginServer" 
                                dbServer = game_base_dir + "dbServer" 
                                worldServer = game_base_dir + "worldServer" 
                                chatServer = game_base_dir + "chatServer" 
                                dipServer = game_base_dir + "dipServer" 
                                gameServer_1 = game_base_dir + "gameServer 1" 
                                gameServer_2 = game_base_dir + "gameServer 2" 

                                self.check_service(loginServer)
                                self.check_service(dbServer)
                                self.check_service(worldServer)
                                self.check_service(chatServer)
                                self.check_service(gameServer_1)
                                self.check_service(gameServer_2)
                                self.check_service(dipServer)
               
                        # 休眠5秒
                        time.sleep(5)
    # 异常后，发送邮件
    def sendmail(self,service_msg):
        self.service_msg = service_msg
        dt = time.strftime('%Y-%m-%d %H:%M:%S',time.localtime(time.time()))
        self.msg = "%s %s start" % (dt, self.service_msg)
        print self.msg
        #获取HTTP对象
        h = httplib2.Http()
        # to do
        resp, content = h.request("http://%s:%d/?" % (self.ip, int(self.port))+self.msg)

    def readConfig(self, config_file):
        self.config_file = config_file
        cp =  ConfigParser.ConfigParser()
        try:
            cp.readfp(codecs.open(config_file, 'r', 'utf-8-sig'))
        except IOError as e:
            print "Open %s failed, no such file." % (config_file)
            #if raw_input('Press any key to continue'):
            sys.exit()

        self.domain = cp.get('main', 'domain')
        self.game_path = cp.get('main', 'game_path')
        self.role = cp.get('main', 'roles')
        self.ip = cp.get('main', 'ip')
        self.port = cp.get('main', 'port')
        self.services_num = cp.get('servers', 'running')
    
    def start(self, config):

        self.config = config
        #检查pid文件是否存在以探测是否存在进程
        try:
            pf = file(self.pidfile,'r')
            pid = int(pf.read().strip())
            pf.close()
        except IOError:
            pid = None
        
        if pid:
            message = 'pidfile %s already exist. Daemon already running?\n'
            sys.stderr.write(message % self.pidfile)
            sys.exit(1)
        
        #启动监控
        self.readConfig(self.config)
        self._daemonize()
        self._run()

    def stop(self):
        #读取pid
        try:
            pf = file(self.pidfile, 'r')
            pid = int(pf.read().strip())
            pf.close()
        except IOError:
            pid = None
     
        if not pid:
            message = 'pidfile %s does not exist. Daemon not running?\n'
            sys.stderr.write(message % self.pidfile)
            return #重启不报错

        #杀进程
        try:
            while 1:
              os.kill(pid, SIGTERM)
              time.sleep(0.1)
              # os.system('/tmp/test.sh')
        except OSError, err:
            err = str(err)
            if err.find('No such process') > 0:
                if os.path.exists(self.pidfile):
                    os.remove(self.pidfile)
            else:
                print str(err)
                sys.exit(1)
    
    def restart(self):
        self.stop()
        self.start()


if __name__ == '__main__':
    daemon = Daemon('/tmp/watch_process.pid', '/dev/null', '/tmp/testdeamon.log', '/tmp/testdeamon.log')
    if len(sys.argv) == 2:
        if 'start' == sys.argv[1]:
            # 读取配置文件
            config_pre = os.getcwd()
            config= config_pre + '/' + 'config.ini'
            if os.path.exists(config):
                daemon.start('%s' % config)
        elif 'stop' == sys.argv[1]:
            daemon.stop()
        elif 'restart' == sys.argv[1]:
            daemon.restart()
        else:
            print 'Unknown command'
            sys.exit(2)
        sys.exit(0)
    else:
      print 'usage: %s (start|stop|restart)' % sys.argv[0]
      sys.exit(2)
