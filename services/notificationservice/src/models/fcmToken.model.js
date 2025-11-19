import { DataTypes } from "sequelize";
import { sequelize } from "../config/db.js";

export const FcmToken = sequelize.define(
  "fcm_tokens",
  {
    id: { type: DataTypes.INTEGER, autoIncrement: true, primaryKey: true },
    userId: { type: DataTypes.INTEGER, allowNull: false },
    token: { type: DataTypes.STRING, allowNull: false },
  },
  {
    timestamps: true,
    tableName: "fcm_tokens",
  }
);
