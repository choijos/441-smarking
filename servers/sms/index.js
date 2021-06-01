"use strict";

const express = require("express");
const morgan = require("morgan");
const mysql = require("mysql");
const { startTimer } = require("./message");

var connection = mysql.createConnection({
  host: "finaldb",
  user: "root",
  password: "thisbetterwork",
  database: "userinfo",
  insecureAuth: true,

});

connection.connect(function (err) {
  if (err) {
    console.error("error connecting: " + err.stack);
    return;

  }

  console.log("connected as id " + connection.threadId);

});

const addr = process.env.NOTIFICATIONADDR || ":80"; // might have to add notificationsaddr to main.go
const [host, port] = addr.split(":");

app.use(express.json());
app.use((err, req, res, next) => {
  console.error(err);
  console.error(err.stack);

  res.set("Content-Type", "text/plain");
  res.status(500).send("Server experienced an error");

});

app.post("somepath", (req, res, next) => {
  let currUser = req.header("X-User");
  
  if (!currUser || currUser == "{}") {
    res.status(401).send("User not authenticated");
    return;

  }

  currUser = JSON.parse(currUser);
  
  let q = 'SELECT phonenumber FROM users WHERE id = ?';

  connection.query(q, [currUser.id], function (err, rows, fields) {
    if (err) {
      throw err;

    }

    if (rows.length == 0) {
      res.status(400).send("Cannot find you in store");
      return;

    }

    // let creator = {"id": rows[0].id, "email": rows[0].email};
    // handler(req, res, SchemeAndDbForwarder, creator);

    // user inputted seconds
    startTimer(10, rows[0].phonenumber);


  });

})