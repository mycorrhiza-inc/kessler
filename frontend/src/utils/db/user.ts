import pg from "pg";

export interface User {
  id: number;
  name: string;
  email: string;
	stripe_id: string;
  createdAt: Date;
  updatedAt: Date;
}

async function AddUserToDB(user: User) {
  const { Client } = pg;
  const client = new Client();
  await client.connect();
	const query = "INSERT INTO users (id, name, email, created_at, ubbpdated_at) VALUES ($1, $2, $3, NOW(), NOW())";
  const res = await client.query(" $1::text as message", ["Hello world!"]);
  console.log(res.rows[0].message); // Hello world!
  await client.end();
  // Add user to database
}
