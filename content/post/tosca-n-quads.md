---
categories:
- category
date: 2017-04-30T21:16:38+02:00
description: ""
draft: true
images:
- http://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.1/csprd02/TOSCA-Simple-Profile-YAML-v1.1-csprd02_files/image004.png
tags:
- tag1
- tag2
title: tosca n quads
---

# Tosca

![http://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.1/csprd02/TOSCA-Simple-Profile-YAML-v1.1-csprd02_files/image004.png](http://docs.oasis-open.org/tosca/TOSCA-Simple-Profile-YAML/v1.1/csprd02/TOSCA-Simple-Profile-YAML-v1.1-csprd02_files/image004.png)

{{< highlight yaml >}}
tosca_definitions_version: tosca_simple_yaml_1_0
 
description: Template for deploying MySQL and database content.
 
topology_template:
  inputs:
    # omitted here for brevity
 
  node_templates:
    my_db:
      type: tosca.nodes.Database.MySQL
      properties:
        name: { get_input: database_name }
        user: { get_input: database_user }
        password: { get_input: database_password }
        port: { get_input: database_port }
      artifacts:
        db_content:
          file: files/my_db_content.txt
          type: tosca.artifacts.File
      requirements:
        - host: mysql
      interfaces:
        Standard:
          create:
            implementation: db_create.sh
            inputs:
              # Copy DB file artifact to server's staging area
              db_data: { get_artifact: [ SELF, db_content ] }
 
    mysql:
      type: tosca.nodes.DBMS.MySQL
      properties:
        root_password: { get_input: mysql_rootpw }
        port: { get_input: mysql_port }
      requirements:
        - host: db_server
 
    db_server:
      type: tosca.nodes.Compute
      capabilities:
        # omitted here for brevity
{{</ highlight >}}


# n-quads 


{{< highlight yaml >}}
"topology_template" "input" "database_name" .
"topology_template" "input" "database_user" .
"topology_template" "input" "database_password" .
"topology_template" "input" "database_port" .
"topology_template" "input" "mysql_rootpw" .
"topology_template" "input" "mysql_port" .
"my_db" "type" "tosca.nodes.Database.MySQL" .
"my_db" "property" "name" .
"name" "value" "database_name" .
"my_db" "property" "user" .
"user" "value" "database_user" .
"my_db" "property" "password" .
"password" "value" "database_password" .
"my_db" "property" "dbport" .
"dbport" "value" "database_port" .
"my_db" "artifact" "db_content" .
"db_content" "file" "files/my_db_content.txt" .
"db_content" "type" "tosca.artifacts.File" .
"my_db" "interface" "Standard" .
"my_db" "interface_create" "create" .
"my_db" "requires" "mysql" .
"create" "implementation" "db_create.sh" .
"create" "input" "db_data" .
"db_data" "value" "db_content" .
"mysql" "type" "tosca.nodes.DBMS.MySQL" .
"mysql" "property" "root_password" .
"root_password" "value" "mysql_rootpw" .
"mysql" "property" "port" .
"port" "value" "mysql_port" .
"mysql" "requires" "db_server" .
"db_server" "type" "tosca.nodes.Compute" .
{{</ highlight >}}

Let's visualize that:

{{< highlight bash >}}
#!/usr/bin/zsh
(
  echo "digraph G {"
  cat $1 | while read subject predicate object trash 
  do 
      echo "$subject -> $object [ label = "$predicate" ];"
  done
  echo "}"
) | dot -Tsvg > output.svg
{{</ highlight >}}

![Output](/assets/images/tosca-n-quads.svg)
