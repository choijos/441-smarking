// Parking
const postParkingHandler = async (req, res, { Parking, user }) => {
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

  newParking.startTime = startTime;
  newParking.owner = user;

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
    const parking = await Parking.findOne({ _id: req.params.parkid });

    if (
      !(parking.owner._id == user._id && parking.creator.email == user.email)
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

    if (newPark.endTime) parking.endTime = newPark.endTime;
    if (newPark.notes) parking.notes = newPark.notes;

    await parking.save();
    res.status(200).set("Content-Type", "application/json").json(channel);
  } catch (e) {
    res.status(404).send("Could not find Parking with id " + req.params.parkid);
    return;
  }
};

const deleteSpecParkingHandler = async (req, res, { Parking, user }) => {
  try {
    const parking = await Parking.findOne({ _id: req.params.parkid });

    if (
      !(parking.owner._id == user._id && parking.creator.email == user.email)
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

module.exports = {
  postParkingHandler,
  getParkingHandler,
  getSpecParkingHandler,
  patchSpecParkingHandler,
  deleteSpecParkingHandler,
};