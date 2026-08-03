import dotenv from "dotenv";
import connectDB from "./db/index.ts";
import { app } from "./app.ts";

dotenv.config({ path: "./.env" });

connectDB()
  .then(() => {
    app.listen(process.env.PORT || 8000, () => {
      console.log(`Server is running at port ${process.env.PORT || 8000}`);
    });
  })
  .catch((error) => {
    console.error(
      "mongoDB connection failed!",
      error instanceof Error ? error.message : error
    );
  });
