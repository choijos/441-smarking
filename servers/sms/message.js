// const accountSid = process.env.TWILIO_ACCOUNT_SID;
// const authToken = process.env.TWILIO_AUTH_TOKEN;
const accountSid = 'ACbf391b4d4234222e10ca877559fc0838'; // Your Account SID from www.twilio.com/console
const authToken = 'd07e3e6bfe47f2c0c579cfe2db2d13b8';   // Your Auth Token from www.twilio.com/console

const client = require('twilio')(accountSid, authToken);

function startTimer(sec, phone) {
  timer = setTimeout(() => { // So this will have to be some sort of callback function for event handlers - when the client presses the start button or something, this function should start so maybe shouldn't be it's own endpoint/microservice?
    let msgBody = secs + " seconds have elapsed";
    client.messages
      .create({
        body: msgBody,
        from: '+12512734782', // set as environment variable later?
        to: phone // get from db
      })
      .then(message => console.log(message.sid));
  }, (sec * 1000));

  return stop;

  function stop() {
    if (timer) {
      clearTimeout(timer);
      timer = 0;

    }

  }

};

module.exports = startTimer;