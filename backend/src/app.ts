import express, { type Express } from "express";
import cors from "cors";
import cookieParser from "cookie-parser";

import userRouter from "./routes/user.route.ts";
import videoRouter from "./routes/video.route.ts";
import subscriptionRouter from "./routes/subscription.route.ts";
import commentRouter from "./routes/comment.route.ts";
import likeRouter from "./routes/like.route.ts";
import postRouter from "./routes/post.route.ts";
import playlistRouter from "./routes/playlist.route.ts";
import dashboardRouter from "./routes/dashboard.route.ts";
import healthcheckRouter from "./routes/healthcheck.route.ts";

const app: Express = express();

app.use(cors({ origin: process.env.CORS_ORIGIN, credentials: true }));
app.use(express.json({ limit: "16kb" }));
app.use(express.urlencoded({ extended: true, limit: "16kb" }));
app.use(express.static("public"));
app.use(cookieParser());

app.use("/api/v1/users", userRouter);
app.use("/api/v1/videos", videoRouter);
app.use("/api/v1/subscriptions", subscriptionRouter);
app.use("/api/v1/comments", commentRouter);
app.use("/api/v1/likes", likeRouter);
app.use("/api/v1/posts", postRouter);
app.use("/api/v1/playlist", playlistRouter);
app.use("/api/v1/dashboard", dashboardRouter);
app.use("/api/v1/healthcheck", healthcheckRouter);

export { app };
