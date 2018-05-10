const mongodb = require('mongodb');
const express = require('express');

const MongoClient = mongodb.MongoClient;
const app = express();

const url = "mongodb://localhost:27017/";

app.use((req, res, next) => {
  res.header("Access-Control-Allow-Origin", "*");
  res.header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept");
  next();
});

app.get('/api/site', (req, res) => {
  MongoClient.connect(url, (err, client) => {
    if (err) {
      return console.log(err.name + ': ' + err.message);
    }

    var db = client.db("main");

    db.collection("site").findOne({}, function(err, document) {
      if (err) throw err;
      res.send({
        profile: document.profile,
        experience: document.experience,
        skills: document.skills,
      });
      db.close();
    });
  });
});

app.listen(8081, () => {
  console.log('Listening on port 8081...');
});

