import dotenv from "dotenv";
import { Sequelize } from "sequelize";
dotenv.config();

export const sequelize = new Sequelize(
  process.env.MYSQL_DB || "identityservice",
  process.env.MYSQL_USER || "root",
  process.env.MYSQL_PASSWORD || "123456",
  {
    host: process.env.MYSQL_HOST || "localhost",
    dialect: "mysql",
    logging: false,
  }
);
