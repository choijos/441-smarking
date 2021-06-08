// postParkingHandler starts a parking session for specific user
const postParkingHandler = async (req, res, { Parking, user, uCars, sms }) => {
  if (req.get("Content-Type") != "application/json") {
    res.status(415).send("Must be JSON");
  }

  let newParking = req.body;
  if (!newParking) {
    res.status(400).send("Must provide a new parking");
    return;
  }

  if (newParking.carID == null) {
    res.status(400).send("Car must be specified");
    return;
  }

  if (newParking.endTime == null) {
    res.status(400).send("Please specify an end time for this parking session");
    return;

  }

  let convCar = Number(req.body.carID);
  if (convCar == NaN) {
    res.status(400).send("Valid Car must be specified");
    return;

  }

  // checking if this car belongs to this user
  let currCar = null;
  for (let i = 0; i < uCars.length; i++) {
    let car = uCars[i];
    if (car.ID == convCar) {
      currCar = car;
      break;

    }

  }

  if (!currCar) {
    res.status(400).send("Car provided is not a car you have registered");
    return;

  }
  
  // checking if there is currently a parking session that hasn't been completed
  let carPark = await Parking.find({ carID: convCar, owner: user._id }, function (err, docs) {
    if (err) {
      res.status(500).send("There was an error retrieving parking");
      return;

    }

  });

  for (let i = 0; i < carPark.length; i++) {
    let currPark = carPark[i];
    if (!(currPark.isComplete)) {
      res.status(400).send("There is already a parking session ongoing with this car");
      return;

    }

  }

  const startTime = new Date();
  let stringTime = newParking.endTime;

  newParking.endTime = new Date(stringTime);
  newParking.startTime = startTime;
  newParking.owner = user;
  newParking.isComplete = false;

  const query = new Parking(newParking);
  await query.save((err, p) => {
    if (err) {
      res.status(500).send("Unable to create a parking");
      return;
    }

    res.status(201).set("Content-Type", "application/json").json(p);

  });

  // starting twilio notfication timer
  sms.Start(newParking.endTime, user.phonenumber, query._id, currCar);

};

// getParkingHandler returns all parking sessions the user current has ongoing
const getParkingHandler = async (req, res, { Parking, user }) => {
  try {
    const parking = await Parking.find();
    let filteredPark = parking.filter((p) => {
      return p.owner._id == user._id;
    });
    res.set("Content-Type", "application/json").json(filteredPark);
  } catch (e) {
    res.status(500).send("There was an issue getting Parking");
  }
};

// getSpecParkingHandler returns the information on parking session with the given parkid
const getSpecParkingHandler = async (req, res, { Parking, user }) => {
  try {
    const parking = await Parking.findOne({ _id: req.params.parkid });
    res.set("Content-Type", "application/json").json(parking);
  } catch (e) {
    res.status(500).send("There was an issue getting Parking");
  }
};

// patchSpecParkingHandler marks the parking session with the given parkid complete, stopping the twilio timer
const patchSpecParkingHandler = async (req, res, { Parking, user, sms }) => {
  try {
    const parking = await Parking.find({ _id: req.params.parkid });

    if (
      !(parking[0].owner._id == user._id && parking[0].owner.email == user.email) //
    ) {
      res.status(403).send("You did not create this parking");
      return;
    }

    if (req.get("Content-Type") != "application/json") {
      res.status(415).send("Must be JSON");
    }

    let newPark = req.body;

    if (!newPark) {
      res.status(400).send("No updates to make");
      return;
    }

    if (newPark.isComplete) parking.isComplete = newPark.isComplete;

    let pk = await Parking.findByIdAndUpdate(req.params.parkid, { isComplete: req.body.isComplete }, function (err, docs) {
      if (err) {
        res.status(500).send("There was an error making parking changes");
        return;

      }

    });

    // clearing timeout
    sms.Stop(req.params.parkid);

    res.status(200).set("Content-Type", "application/json").json(pk);

  } catch (e) {
    res.status(404).send("Could not find Parking with id " + req.params.parkid);
    return;
  }
};

// deleteSpecParkingHandler removes the specific parking session with given parkid from the parking schema
const deleteSpecParkingHandler = async (req, res, { Parking, user }) => {
  try {
    const parking = await Parking.findOne({ _id: req.params.parkid });

    if (
      !(parking[0].owner._id == user._id && parking[0].owner.email == user.email)
    ) {
      res.status(403).send("You did not create this parking");
      return;
    }

    await Parking.deleteOne({ _id: req.params.parkid });
  } catch (e) {
    res.status(404).send("Could not find Parking with id " + req.params.parkid);
    return;
  }
};

// invalidMethod handles all inappropriate requests of unsupported/invalid methods
const invalidMethod = async (req, res, {}) => {
  res.status(405).send("This method is not allowed");
  return;
};

module.exports = {
  postParkingHandler,
  getParkingHandler,
  getSpecParkingHandler,
  patchSpecParkingHandler,
  deleteSpecParkingHandler,
  invalidMethod,
};