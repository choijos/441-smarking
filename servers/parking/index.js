"use strict";

const mongoose = require("mongoose");
const express = require("express");
var mysql = require("mysql");

var connection = mysql.createConnection({
  host: "db",
  user: "root",
  password: "password",
  database: "db",
  insecureAuth: true,
});

connection.connect(function (err) {
  if (err) {
    console.error("error connecting: " + err.stack);
    return;
  }
  console.log("connected as id " + connection.threadId);
});

const { parkingSchema } = require("./models/schemas");
const {
  postParkingHandler,
  getParkingHandler,
  getSpecParkingHandler,
  patchSpecParkingHandler,
  deleteSpecParkingHandler,
} = require("./handlers/handlers");

const mongoEndpoint = "mongodb://customMongoContainer:27017/test"; // test is name of database

const addr = process.env.MESSAGESADDR || ":80";
const [host, port] = addr.split(":");

const Parking = mongoose.model("Parking", parkingSchema);

const app = express();

// Add middleware
app.use(express.json());
app.use((err, req, res, next) => {
  console.error(err);
  console.error(err.stack);

  res.set("Content-Type", "text/plain");
  res.status(500).send("Server experienced an error");
});

// Request Wrapper
const RequestWrapper = (handler, SchemeAndDbForwarder) => {
  return (req, res) => {
    let user = JSON.parse(req.get("X-User"));

    if (user == null || user.id == null) {
      res.status(401).send("You need to be signed in to do that");
      return;
    }

    connection.query(
      "SELECT ID, Email FROM `Users` WHERE `ID` = ?",
      [user.id],
      (err, results, fields) => {
        if (err) throw err;
        if (results == null) {
          res.status(404).send("User not found");
          return;
        }

        let insertUser = { _id: results[0].ID, email: results[0].Email };

        SchemeAndDbForwarder.user = insertUser;
        handler(req, res, SchemeAndDbForwarder);
      }
    );
  };
};

// Requests
//  Parking
app.post(
  "/v1/users/:id/parking",
  RequestWrapper(postParkingHandler, { Parking })
);
app.get(
  "/v1/users/:id/parking",
  RequestWrapper(getParkingHandler, { Parking })
);

// Spec. Parking
app.get(
  "/v1/users/:id/parking/:parkid",
  RequestWrapper(getSpecParkingHandler, { Parking })
);
app.patch(
  "/v1/users/:id/parking/:parkid",
  RequestWrapper(patchSpecParkingHandler, { Parking })
);
app.delete(
  "/v1/users/:id/parking/:parkid",
  RequestWrapper(deleteSpecParkingHandler, { Parking })
);

const connect = () => {
  mongoose.connect(mongoEndpoint);
};

// Start Server
connect();
mongoose.connection
  .on("error", console.error)
  .on("disconnected", connect)
  .once("open", main);

async function main() {
  app.listen(port, host, () => {
    console.log(`server is listening at http://${host}:${port}`);
  });
}
