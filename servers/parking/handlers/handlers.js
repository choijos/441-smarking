// Parking
const postParkingHandler = async (req, res, { Parking, user, smsNotif }) => {
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

  const startTime = new Date();
  // where is endtime coming from?

  newParking.startTime = startTime;
  newParking.owner = user;
  newParking.isComplete = false;

  //let stop = smsNotif(10, user.phonenumber); // call stop() to stop the timer?

  const query = new Parking(newParking);
  query.save((err, p) => {
    if (err) {
      res.status(500).send("Unable to create a parking");
      return;
    }

    res.status(201).set("Content-Type", "application/json").json(p);
  });
};

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

// Spec. Parking

// Spec. Channels
const getSpecParkingHandler = async (req, res, { Parking, user }) => {
  try {
    const parking = await Parking.findOne({ _id: req.params.parkid });
    res.set("Content-Type", "application/json").json(parking);
  } catch (e) {
    res.status(500).send("There was an issue getting Parking");
  }
};

const patchSpecParkingHandler = async (req, res, { Parking, user }) => {
  try {
    const parking = await Parking.find({ _id: req.params.parkid });

    if (
      !(parking[0].owner._id == user._id && parking[0].owner.email == user.email)
    ) {
      res.status(403).send("You did not create this parking");
      return;
    }

    if (req.get("Content-Type") != "application/json") {
      res.status(415).send("Must be JSON");
    }

    let newPark = req.body;

    if (!newPark) {
      res.status(400).send("Must provide a new channel");
      return;
    }

    // if (newPark.endTime) parking.endTime = newPark.endTime;
    // if (newPark.notes) parking.notes = newPark.notes;
    if (newPark.isComplete) parking.isComplete = newPark.isComplete;

    let pk = await Parking.findByIdAndUpdate(req.params.parkid, { isComplete: req.body.isComplete }, function (err, docs) {
      if (err) {
        res.status(500).send("There was an error making parking changes");
        return;

      }

    });

    res.status(200).set("Content-Type", "application/json").json(pk);

  } catch (e) {
    res.status(404).send("Could not find Parking with id " + req.params.parkid);
    return;
  }
};

const deleteSpecParkingHandler = async (req, res, { Parking, user }) => {
  try {
    const parking = await Parking.find({ _id: req.params.parkid });

    if (
      !(parking.owner[0]._id == user._id && parking.owner[0].email == user.email)
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
