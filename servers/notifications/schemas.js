const Schema = require("mongoose").Schema;

const Device = new Schema({
  deviceId: { type: String, required: true},
  platform: { type: String, required: true}

});

module.exports = {
  Device
};