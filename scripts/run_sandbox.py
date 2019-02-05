#!/usr/bin/python

import os
import yaml
import pprint
import docker

cmd = 'docker inspect pingpen_sandbox_db  | grep \'IPAddress"\' | head -n1 | cut -d":" -f2 | cut -d\'"\' -f2'
db_host = os.popen(cmd).read().strip()
#client.containers.get('3d0de1a160').attrs["NetworkSettings"]["IPAddress"]

if db_host == "":
    cmd = 'docker run --name pingpen_sandbox_db -e MYSQL_ROOT_PASSWORD=password -e MYSQL_DATABASE=pingpen -d mysql:5.6'
    os.system(cmd) 
    cmd = 'docker inspect pingpen_sandbox_db  | grep \'IPAddress"\' | head -n1 | cut -d":" -f2 | cut -d\'"\' -f2'
    db_host = os.popen(cmd).read().strip()

cmd = 'cd db && goose mysql "root:password@tcp(%s:3306)/pingpen?parseTime=true" up'%(db_host)
os.system(cmd)

#with open("serverless.yml", 'r') as yamlfile:
#    config = yaml.load(yamlfile)
#pprint.pprint(config)

services = {
    "post_note": {"db_host": db_host, "port": "8080", "bin": "./api/bin/create"},
}

cmd = '''
    export GRID=sandbox; 
    export DBHOST=%(db_host)s;
    export DBUSER=root; 
    export DBPW=password;  
    export GRID=sandbox; 
    export DBNAME=pingpen; 
    export PORT=%(port)s; 
    %(bin)s
'''

for s in services:
    print("spinning up %s"%s)
    #print(cmd%services[s])
    os.system(cmd%services[s])



