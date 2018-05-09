const mongodb = require('mongodb');
const express = require('express');

const MongoClient = mongodb.MongoClient;
const app = express();

const url = "app-database";

app.use((req, res, next) => {
  res.header("Access-Control-Allow-Origin", "*");
  res.header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept");
  next();
});

MongoClient.connect(url, (err, client) => {
  if (err) {
    return console.log(err.name + ': ' + err.message);
  }

  var db = client.db("main");

  app.get('/api/site', (req, res) => {
    db.collection("site").findOne({}, function(err, document) {
      if (err) throw err;
      res.send({
        profile: document.profile,
        experience: document.experience,
        skills: document.skills,
      });
    });
  });

  app.listen(8081, () => {
    console.log('Listening on port 8081...');
  });
});
