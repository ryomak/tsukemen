DROP DATABASE IF EXISTS test;
CREATE DATABASE test;

DROP USER IF EXISTS 'app' ;
CREATE USER 'app' IDENTIFIED BY 'password' ;
GRANT ALL ON test.* TO app@'localhost';   
