alter user root identified with mysql_native_password by 'thisbetterwork';
flush privileges;

USE userinfo;

-- creating the user table
CREATE TABLE IF NOT EXISTS users (
    id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
    email varchar(254) NOT NULL UNIQUE,
    first_name varchar(64) NOT NULL,
    last_name varchar(128) NOT NULL,
    username varchar(255) NOT NULL UNIQUE,
    passhash varchar(60) NOT NULL,
    photourl varchar(128) NOT NULL
);

-- creating table to track user sign-ins
CREATE TABLE IF NOT EXISTS usersignin (
    id int NOT NULL REFERENCES users(id),
    whensignin time NOT NULL,
    clientip varchar(255) NOT NULL
);

-- creating table consisting of car information
CREATE TABLE IF NOT EXISTS cars (
  ID int NOT NULL AUTO_INCREMENT PRIMARY KEY,
  LicensePlate varchar(15) NOT NULL,
  UserID int NOT NULL REFERENCES users(id),
  Make varchar(128),
  Model varchar(128),
  ModelYear varchar(128),
  Color varchar(128)

);