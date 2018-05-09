var MongoClient = require('mongodb').MongoClient;
var url = "mongodb://localhost:27017/";

const express = require('express');
const app = express();

app.use((req, res, next) => {
  res.header("Access-Control-Allow-Origin", "*");
  res.header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept");
  next();
});

MongoClient.connect(url, (err, db) => {
  try {
    if (err) throw err;
    var dbo = db.db("main");

    app.get('/api/site', (req, res) => {
      dbo.collection("site").findOne({}, function(err, document) {
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
  } catch (e) {
    db.close();
  }
});
