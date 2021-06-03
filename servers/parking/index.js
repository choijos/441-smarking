"use strict";

const mongoose = require("mongoose");
const express = require("express");
const morgan = require("morgan");
var mysql = require("mysql");

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

const { parkingSchema } = require("./models/schemas");
const {
  postParkingHandler,
  getParkingHandler,
  getSpecParkingHandler,
  patchSpecParkingHandler,
  deleteSpecParkingHandler,
  invalidMethod,
} = require("./handlers/handlers");

const mongoEndpoint = "mongodb://customMongoContainer:27017/test"; // test is name of database

const addr = process.env.PARKINGADDR || ":80";
const [host, port] = addr.split(":");

// Getting twilio info
const accountSid = process.env.TWILIO_ACCOUNT_SID;
const authToken = process.env.TWILIO_AUTH_TOKEN;

const client = require("twilio")(accountSid, authToken);

const Parking = mongoose.model("Parking", parkingSchema);

const app = express();

// Add middleware
app.use(express.json());
app.use(morgan("dev"));
app.use((err, req, res, next) => {
  console.error(err);
  console.error(err.stack);

  res.set("Content-Type", "text/plain");
  res.status(500).send("Server experienced an error");
});

const smsNotif = async (sec, phone) => {
  timer = setTimeout(() => {
    // So this will have to be some sort of callback function for event handlers - when the client presses the start button or something, this function should start so maybe shouldn't be it's own endpoint/microservice?
    let msgBody = secs + " seconds have elapsed";
    client.messages
      .create({
        body: msgBody,
        from: "+12512734782", // set as environment variable later?
        to: phone, // get from db
      })
      .then((message) => console.log(message.sid));
  }, sec * 1000);

  return stop;

  function stop() {
    if (timer) {
      clearTimeout(timer);
      timer = 0;
    }
  }
};

// Request Wrapper
const RequestWrapper = (handler, SchemeAndDbForwarder) => {
  return (req, res) => {
    let user = JSON.parse(req.get("X-User"));

    if (user == null || user.id == null) {
      res.status(401).send("You need to be signed in to do that");
      return;
    }

    connection.query(
      "select id, email, phonenumber from users where id = ?",
      [user.id],
      (err, results, fields) => {
        if (err) throw err;
        if (results == null || results.length == 0) {
          res.status(404).send("User not found");
          return;
        }

        let insertUser = {
          _id: results[0].id,
          email: results[0].email,
          phonenumber: results[0].phonenumber,
        };

        SchemeAndDbForwarder.user = insertUser;
        SchemeAndDbForwarder.smsNotif = smsNotif;
        handler(req, res, SchemeAndDbForwarder);
      }
    );
  };
};

// Requests
//  Parking
app.post("/v1/usersparking", RequestWrapper(postParkingHandler, { Parking }));
app.get("/v1/usersparking", RequestWrapper(getParkingHandler, { Parking }));
app.patch("/v1/usersparking", RequestWrapper(invalidMethod, {}));
app.delete("/v1/usersparking", RequestWrapper(invalidMethod, {}));

// Spec. Parking
app.get(
  "/v1/parking/:parkid",
  RequestWrapper(getSpecParkingHandler, { Parking })
);
app.patch(
  "/v1/parking/:parkid",
  RequestWrapper(patchSpecParkingHandler, { Parking })
);
app.delete(
  "/v1/parking/:parkid",
  RequestWrapper(deleteSpecParkingHandler, { Parking })
);
app.post("/v1/parking/:parkid", RequestWrapper(invalidMethod, {}));

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
