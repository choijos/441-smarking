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
const twilioPhone = process.env.TWILIO_PHONE_NUMBER;

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

const initSms = new Map();

// Creating, storing, and ending parking sessions (for notifications)
var sms = {
  Start: (endTime, phone, parkid, currCar) => {
    let sentPrem = false;
    let now = new Date();
    let endDate = new Date(endTime);
    let total = (endDate.getTime() - now.getTime()) / 1000;
    let short = false;
    if (total <= 350) {
      short = true;
  
    }

    let intervalID = setInterval(() => {
      let endDate = new Date(endTime);
      let now = new Date();
      let secLeft = (endDate.getTime() - now.getTime()) / 1000;

      if (!short && secLeft <= 300 && !sentPrem) {
        let msgBody = "You have 5 minutes remaining in your parking session ("+ currCar.Make + " " + currCar.Model + " - " + currCar.LicensePlate + ") \n[" + parkid + "]";
        client.messages
          .create({
            body: msgBody,
            from: twilioPhone,
            to: phone,
          })
          .then((message) => {
            console.log(message.sid)
            sentPrem = true;

          });

      } else if (secLeft <= 1) {
        let msgBody = "Smarking: Your parking session has ended\n("+ currCar.Make + " " + currCar.Model + " - " + currCar.LicensePlate + ") \n\n Parking ID: [" + parkid + "]";
        client.messages
          .create({
            body: msgBody,
            from: twilioPhone,
            to: phone,
          })
          .then((message) => {
            console.log(message.sid)

            Parking.findByIdAndUpdate(JSON.parse(parkid), { isComplete: true }, function (err, doc) {
              if (err) {
                throw err;

              }

            });

            clearInterval(initSms.get(parkid));

          });

      }
    }, 4000)

    parkid = JSON.stringify(parkid);
    initSms.set(parkid, intervalID);

  },

  Stop: (parkid) => {
    parkid = JSON.stringify(parkid);
    let intervalID = initSms.get(parkid);
    clearInterval(intervalID);
    initSms.delete(parkid);

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

    let userCars = [];

    connection.query(
      "select ID, LicensePlate, Make, Model from cars where UserID = ?", //
      [user.id],
      (err, results, fields) => {
        if (err) throw err;
        for (let i = 0; i < results.length; i++) {
          let currCar = results[i];
          let oneCar = {
            ID: currCar.ID,
            LicensePlate: currCar.LicensePlate,
            Make: currCar.Make,
            Model: currCar.Model

          }

          userCars.push(oneCar);

        }

      }

    );

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
        SchemeAndDbForwarder.uCars = userCars;
        SchemeAndDbForwarder.sms = sms;
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