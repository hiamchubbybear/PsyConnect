import { DataTypes } from "sequelize";
import { sequelize } from "../config/db.js";

export const Notification = sequelize.define(
  "notifications",
  {
    id: { type: DataTypes.INTEGER, autoIncrement: true, primaryKey: true },
    userId: { type: DataTypes.STRING, allowNull: false },
    title: { type: DataTypes.STRING, allowNull: false },
    body: { type: DataTypes.TEXT, allowNull: false },
    type: { type: DataTypes.STRING, allowNull: false }, // e.g., 'like', 'comment', 'follow'
    metadata: { type: DataTypes.JSON, allowNull: true }, // e.g., { postId: '...', authorId: '...' }
    isRead: { type: DataTypes.BOOLEAN, defaultValue: false },
  },
  {
    timestamps: true,
    tableName: "notifications",
  }
);
