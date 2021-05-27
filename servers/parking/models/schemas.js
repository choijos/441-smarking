const Schema = require("mongoose").Schema;

const parkingSchema = new Schema({
  carID: { type: Number, required: true, unique: true },
  owner: { _id: { type: Number }, email: { type: String } },
  photoURL: { type: String, required: false, unique: false },
  startTime: { type: Date, required: true, unique: false },
  endTime: { type: Date, required: false, unique: false },
  notes: { type: String, required: false, unique: false },
});

module.exports = { parkingSchema };