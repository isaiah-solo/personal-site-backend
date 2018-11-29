const express = require('express');
const mysql = require('mysql');

const app = express();
const pool = mysql.createPool(
  {
    connectionLimit: 20,
    database: 'site',
    host: 'localhost',
    password: 'admin',
    user: 'root',
  }
);

app.use((req, res, next) => {
  res.header('Access-Control-Allow-Origin', '*');
  res.header('Access-Control-Allow-Headers', 'Origin, X-Requested-With, Content-Type, Accept');
  next();
});

app.get('/api/static/profile', (req, res) => {
  pool.query('select * from profile', (error, rows) => {
    if (error) {
      throw error;
    }
    if (rows.length < 1) {
      throw 'profile table should not be empty.';
    }
    const {headline, id, name} = rows[0];
    pool.query(
      'select * from profile_icon where profile_id = ?',
      [id],
      (error, rows) => {
        if (error) {
          throw error;
        }
        res.send({
          profile: {
            headline,
            icons: rows.map(({name, website}) => (
              {name, website}
            )),
            name,
          }
        });
      }
    );
  });
});

app.get('/api/static/about', (req, res) => {
  pool.query(
    'select j.company as company, j.end_date as endDate, j.id as id, j.position as position, j.start_date as startDate, j.website as website, s.label as skillLabel, s.link as skillLink, jd.text as detailText from job j, job_to_skill js, skill s, job_detail jd where j.id = js.job_id and s.id = js.skill_id and j.id = jd.job_id',
    (error, rows) => {
      if (error) {
        throw error;
      }
      if (rows.length < 1) {
        throw 'profile table should not be empty.';
      }
      let jobMap = {};
      for (const row of rows) {
        const {
          company,
          detailText,
          endDate,
          id,
          position,
          skillLabel,
          skillLink,
          startDate,
          website
        } = row;
        if (id in jobMap) {
          jobMap[id].details = Array.from(
            new Set([
              ...jobMap[id].details,
              detailText
            ])
          );
          if (
            !jobMap[id].skills.some(({label}) => (
              label === skillLabel
            ))
          ) {
            jobMap[id].skills = [
              ...jobMap[id].skills,
              {label: skillLabel, link: skillLink}
            ];
          }
        } else {
          jobMap[id] = {
            company,
            details: [
              detailText
            ],
            endDate,
            position,
            skills: [
              {
                label: skillLabel,
                link: skillLink
              }
            ],
            startDate,
            website,
          };
        }
      }
      res.send({
        about: {
          jobs: Object.keys(jobMap).map(jobKey => (
            jobMap[jobKey]
          )),
        },
      });
    }
  );
});

app.listen(8081, () => {
  console.log('Listening on port 8081...');
});

