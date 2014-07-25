leejo
=====

**leejo** is an open source board to add and display skillsets
of individuals in a programmer group(s).


Dependencies
------------
**leejo** dependes on these software:

  1. Go (a.k.a. golang) version 1.2+
  2. PostgreSQL version 9+


Compile
-------

To compile the project, you need to install golang (v1.2+)
first. Then in the directory where this file locates, type:

    $ make

The server program binary will be compiled to
`./bin/leejo_server`


Install
-------

There is no way to configure / install the binary yet .You
may manually copy `./bin/leejo_server` to any folder in 
your PATH.

You should create the database with the schema file provided
with the source. Run this:

    $ psql [your_db_name [your_user]] \
      < ./data/leejo_schema.sql

Remember to replace your database name and user. If your 
database is password protected, it will prompt you for it.


Run
---

To run the server, simply use this command:

    $ leejo_server -config "path/to/config.json"

You may reference `data/config.example.json` for the config file format. Remember to replace your database name, user and password

If success, the server should be running at
http://localhost:8080


Integration Test
----------------

When the server is running, you can test the installation
by running our integration test. To run, simply run this
command within the folder this file locates:

    $ make check

The make script should build and run the integration test
for you. Please note that the test fails if the leejo
server process is not running or is not listening to
localhost:8080.


Bug Reports
-----------

To report issue, please visit our
[issue tracker](https://github.com/opensourcehk/leejo/issues).
