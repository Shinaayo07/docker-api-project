const fs = require('fs');

require('dotenv').config();
const { Pool } = require('pg');

const pool = new Pool({
  connectionString: process.env.DATABASE_URL,
});


// const { Pool } = require('pg');

// require('dotenv').config();  // Load env vars from .env

// console.log(process.env.DATABASE_URL);  // For testing

let databaseUrl;

if (process.env.DATABASE_URL) {
  databaseUrl = process.env.DATABASE_URL;
} else if (process.env.DATABASE_URL_FILE) {
  databaseUrl = fs.readFileSync(process.env.DATABASE_URL_FILE, 'utf8').trim();
} else {
  throw new Error('DATABASE_URL or DATABASE_URL_FILE environment variable must be set');
}

// databaseUrl =
//   process.env.DATABASE_URL ||
//   fs.readFileSync(process.env.DATABASE_URL_FILE, 'utf8');

// const pool = new Pool({
//   connectionString: databaseUrl,
// });

// the pool will emit an error on behalf of any idle clients
// it contains if a backend error or network partition happens
pool.on('error', (err, client) => {
  console.error('Unexpected error on idle client', err);
  process.exit(-1);
});

// async/await - check out a client
const getDateTime = async () => {
  const client = await pool.connect();
  try {
    const res = await client.query('SELECT NOW() as now;');
    return res.rows[0];
  } catch (err) {
    console.log(err.stack);
  } finally {
    client.release();
  }
};

module.exports = { getDateTime };
