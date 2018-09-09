#!/usr/bin/python
#coding=utf-8
import re
import sys
import MySQLdb

def usage():
    print sys.argv[0] + " ip dbname(xx_xx_xx1)"
    exit(1)

def check_game_server_status(ip, platform, game_name):
    pass

def truncate_tables(ip, user, password, dbname):
    conn= MySQLdb.connect(
            host = ip,
            port = 3306,
            user = user,
            passwd = password,
            db = dbname,
            charset='utf8',
            )
    cur = conn.cursor()
    
    #创建数据表
    #cur.execute("create table student(id int ,name varchar(20),class varchar(30),age varchar(10))")
    
    #插入数据
    cur.execute("insert into student values(2,'Tom','3 year 2 class','9')")
    
    #修改查询条件的数据
    #cur.execute("update student set class='3 year 1 class' where name = 'Tom'")
    
    #查询数据，返回记录条数
    try:
        count = cur.execute("select * from student")
        
        if count < 2:
            print "less than 2 records"
        else:
            print "more than 2 records, less do something"
            cur.execute("show tables")
            tables = cur.fetchall()  
            print tables
        
            for (table_name, ) in cur:
                try:
                    print "trying to truncate table %s" % table_name
                    cur.execute("truncate table %s" % table_name)
                except MySQLdb.Error, e:
                    print e
                
    except MySQLdb.Error, e:
        print e
        
    #删除查询条件的数据
    #cur.execute("delete from student where age='9'")
        
    cur.close()
    conn.commit()
    conn.close()

if __name__ == "__main__":
    if len(sys.argv) != 3:
        usage()
    ip = sys.argv[1]
    dbname = sys.argv[2]
   #check_db = re.split('_', dbname)
   #if len(check_db) != 3:
   #    print "wrong dbname"
   #    exit(3)
   #game = dbname.split("_")[0]
   #platform = dbname.split("_")[1]
   #server_name = dbname.split("_")[2]
   #game_name = game + "_" + server_name
   #if game == "" or platform == "" or server_name == "" :
   #    usage() 
   #print  platform, game_name
#    check_game_server_status(ip)
    truncate_tables(ip, "root", "password", dbname)
