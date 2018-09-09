#!/usr/bin/env python
# -*- coding:utf-8 -*-
# data collector
# 数据收集服务端
# 时间：2014-10-21
#

import multiprocessing
import socket
import pymongo
import datetime
import json

def insert_into_mongodb(data, address, dbconfig):
    dbaddress = dbconfig[0]['dbaddress']
    dbport = dbconfig[1]['dbport']
    dbuser = dbconfig[2]['dbuser']
    dbpass = dbconfig[3]['dbpass']
    dbname = dbconfig[4]['dbname']
    # print dbaddress, dbport, dbuser, dbpass, dbname
    connection=pymongo.Connection('%s' % dbaddress, int(dbport))

    # 数据库名称 未解决
    db = connection.mokylin
    posts = db.posts
    #data = json.dumps(data)
    print "-------------"
    print data
    b=eval(data)
    disk_info = b[0]
    cpu_info= b[1]['load_status']
    mem_info = b[2]['mem_use_rate']
    print type(b)
    print "-------------"
#   disk_info = json.dumps(data.split('+')[0])
#   cpu_info = json.dumps(data.split('+')[1])
#   mem_info = json.dumps(data.split('+')[2])

    print type(mem_info)
   #print mem_info[0]['mem_use_rate']
   #print cpu_info[0]['load_status']
    post = {"address": address,
            "diskinfo":disk_info,
            "mem":mem_info,
            "cpu":cpu_info,
            "date": datetime.datetime.utcnow()}
    print post
    posts.insert(post)


def handle(connection, address, dbconfig):
    import logging
    logging.basicConfig(level=logging.DEBUG)
    logger = logging.getLogger("process-%r" % (address,))
    try:
        #logger.debug("Connected %r at %r", connection, address)
        logger.debug("Connected %r",address)
        while True:
            data = connection.recv(1024)
            if data == "":
                logger.debug("Socket closed remotely")
                break
            #数据处理
            address=address[0]
                   
            insert_into_mongodb(data, address, dbconfig)
            # logger.debug("Received data %r", data)
            #发送数据到客户端，暂时无用，注释
           #connection.sendall(data)
           #logger.debug("Sent data")
    except:
        logger.exception("Problem handling request")
    finally:
        logger.debug("Closing socket")
        connection.close()

class Server(object):
    def __init__(self, hostname, port, dbconfig):
        import logging
        logging.basicConfig(level=logging.DEBUG)
        self.logger = logging.getLogger("server")
        self.hostname = hostname
        self.port = port
        self.dbconfig = dbconfig

    def start(self):
        self.logger.debug("listening")
        self.socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.socket.bind((self.hostname, self.port))
        self.socket.listen(1)

        while True:
            conn, address = self.socket.accept()
            self.logger.debug("Got connection")
            process = multiprocessing.Process(target=handle, args=(conn, address, self.dbconfig))
            process.daemon = True
            process.start()
            self.logger.debug("Started process %r", process)

#if __name__ == "__main__":
def socket_start(listen_ip, port, dbconfig):
    import logging
    #loglevel='INFO'
    loglevel='DEBUG'
    logging.basicConfig(level=str(loglevel),
                    format='[%(asctime)s] %(levelname)-4s %(message)s',
                    datefmt='%Y-%m-%d %H:%M:%S',
                    filename='/opt/mokylin/datacollection/logs/mokylin_server.log',
                    filemode='w')
    #server = Server("0.0.0.0", 9000)
    server = Server("%s" % listen_ip, int(port), dbconfig)
    try:
        logging.info("Listening")
        server.start()
    except:
        logging.exception("Unexpected exception")
    finally:
        logging.info("Shutting down")
        for process in multiprocessing.active_children():
            logging.info("Shutting down process %r", process)
            process.terminate()
            process.join()
    logging.info("All done")
