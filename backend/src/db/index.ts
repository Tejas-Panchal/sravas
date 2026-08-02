import mongoose from "mongoose";
import { DB_NAME } from "../constants.js";

const connectDB = async (): Promise<void> => {
  try {
    const connectionInstance: mongoose.Mongoose = await mongoose.connect(
      `${process.env.MONGODB_URI}/${DB_NAME}`
    );
    console.log(
      `mongoDB connected. DB host: ${connectionInstance.connection.host}`
    );
  } catch (error) {
    console.log(
      "mongoDB connection error:",
      error instanceof Error ? error.message : error
    );
    process.exit(1);
  }
};

export default connectDB;
