"use strict";

// https://shockoe.com/ideas/development/creating-a-push-notification-server-with-node-js/
const mongoose = require("mongoose");
const express = require("express");
const morgan = require("morgan");
const restify = require('restify');
const apns = require('apns');
const gcm = require('node-gcm');
const device = require("./schemas");
const parkingSchema = require("../parking/models/schemas"); // grabbing the parking model

const app = express(); // not sure if we need this
const addr = process.env.ADDR || ":80"
const[host, port] = addr.split(":")
app.use(express.json()); 
app.use(morgan("dev"));

// const db = mongoose.connection;

// db.on('error', console.error.bind(console, 'connection error:'));

// db.once('open', function() { // maybe could be an arrow function
//   console.log('db open');

// });


// mongoose.connect('mongodb://localhost/pushserver');
const dev = mongoose.createConnection("mongodb://choijos.me/dev"); // have to make this container or something
const park = mongoose.createConnection("mongodb://customMongoContainer:27017/test"); // the stuff we end up doing in parking

// Connecting
const devConn = () => {
  dev.connect((err) => {
    if (err) {
      throw err;

    }

    console.log("Connected to device mongo db for notifcations");

  });

};

const parkConn = () => {
  park.connect((err) => {
    if (err) {
      throw err;

    }

    console.log("Connected to parking mongo db for notifcations");

  });

};

// // connecting to mongodb
// const connect = () => {
//   mongoose.connect(mongoEndpoint);

// }

// mongoose.connection.on("error", console.error) // Not sure if we need one of these statements for each of our db connections
//   .on("disconnected", connect)
//   .once("open", main);

// async function main() {
//   app.listen(port, "", () => { // check console.log
//     console.log(`server is listening at port ${port}`);

//   });

// }

const nextConn = park.connection.on("error", console.error).on("disconnected", parkConn).once("open", main); // error because mongodb containers aren't running currently

dev.connection.on("error", console.error)
  .on("disconnected", devConn)
  .once("open", nextConn); // idk if this acutally will work, checking for errors in the parking after getting device connection right

async function main() {
  app.listen(port, host, () => {
    console.log(`server is listening at http://${host}:${port}`);

  });

}

// const DeviceSchema = mongoose.model('Device', Device);
const Device = dev.model("Device", device);
const Parking = park.model("Parking", parkingSchema);

const server = restify.createServer({
  name : 'pushServer'

});


const options = {
  keyFile  : 'key.pem',
  certFile : 'cert.pem',
  debug    : true,
  gateway  : 'gateway.sandbox.push.apple.com',
  errorCallback : function(num, err) {
    console.error(err);

  }

};


function sendIos(deviceId) {
  let connection = new apns.Connection(options);

  let notification = new apns.Notification();
  notification.device = new apns.Device(deviceId);
  notification.alert = 'Hello World !';

  connection.sendNotification(notification);

}


function sendAndroid(devices) {
  let message = new gcm.Message({
    notification : {
      title : 'Hello, World!'
      
    }

  });

  let sender = new gcm.Sender('AAAABfuSCZs:APA91bFEwkkWhsah5lOq8BRhrE_P2NWYTLFRIs3_ChFYCBHaOSdyRTTWbfTlJJupA4dniRYAgCsXhxR0icMOfRtHVqzOXUHMYILVkDeZvLnqAudIr2wNhs_DWk885YaJVQmmv0SX_YhG'); // was lowercase sender but fcm says uppercase?

  sender.send(message, {
    registrationTokens : devices
  }, function(err, response) {
    if (err) {
      console.error(err);

    } else {
      console.log(response);

    }

  });

}


const RequestWrapper = (handler, SchemeAndDbForwarder) => {
  return (req, res) => {
    let currUser = req.header("X-User");
  
    if (!currUser || currUser == "{}") {
      res.status(401).send("User not authenticated");
      return;
  
    }
  
    currUser = JSON.parse(currUser);
  
    let q = 'SELECT id, email FROM users WHERE id = ?';
  
    sConnect.query(q, [currUser.id], function (err, rows, fields) {
      if (err) {
        throw err;
  
      }
  
      if (rows.length == 0) {
        res.status(400).send("Cannot find you in store");
        return;
  
      }
  
      let creator = {"id": rows[0].id, "email": rows[0].email};
      handler(req, res, SchemeAndDbForwarder, creator);

    });

  }

}


server.post('/register', (req, res, next) => { // might or might not need this handler, depending on where we do the registration stuff
  let body = JSON.parse(req.body); // when registering, the request body should have the user id so we can look up the parking stuff

  if (body) {
    let newDevice = new DeviceSchema(body);
    newDevice.save((err) => {
      if (!err) {
        res.send(200);

      } else {
        res.send(500);

      }

    });

  }

});


server.get('/send', (req, res) => {
  DeviceSchema.find( (err, devices) => {
    if (!err && devices) {
      let androidDevices = [];
      devices.forEach((device) => {
        if (device.platform === 'ios') {
          sendIos(device.deviceId);

        } else if (device.platform === 'android') {
          androidDevices.push(device.deviceId);

        }

      });

      sendAndroid(androidDevices);
      res.send(200);

    } else {
      res.send(500);

    }

  });

});