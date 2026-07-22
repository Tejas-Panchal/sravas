import mongoose from "mongoose";
import { Video } from "../models/video.model.js";
// import { Subscription } from "../models/subscription.model.js";
// import { Like } from "../models/like.model.js";
// import { ApiError } from "../utils/ApiError.js";
import { ApiResponse } from "../utils/ApiResponse.js";
import { asyncHandler } from "../utils/asyncHandler.js";

const getChannelStats = asyncHandler(async (req, res) => {
  const userId = req.user._id;

  const stats = await Video.aggregate([
    {
      $match: {
        owner: new mongoose.Types.ObjectId(userId),
      },
    },
    {
      $lookup: {
        from: "subscriptions",
        localField: "owner",
        foreignField: "channel",
        as: "subscribers",
      },
    },
    {
      $lookup: {
        from: "likes",
        localField: "_id",
        foreignField: "video",
        as: "likes",
      },
    },
    {
      $group: {
        _id: null,
        totalVideos: { $sum: 1 },
        totalViews: { $sum: "$views" },
        totalSubscribers: { $first: { $size: "$subscribers" } },
        totalLikes: { $sum: { $size: "$likes" } },
      },
    },
    {
      $project: {
        _id: 0,
        totalVideos: 1,
        totalViews: 1,
        totalSubscribers: 1,
        totalLikes: 1,
      },
    },
  ]);

  if (!stats.length) {
    return res.status(200).json(
      new ApiResponse(
        200,
        {
          totalVideos: 0,
          totalViews: 0,
          totalSubscribers: 0,
          totalLikes: 0,
        },
        "channel stats fetched successfully"
      )
    );
  }

  return res
    .status(200)
    .json(new ApiResponse(200, stats[0], "channel stats fetched successfully"));
});

const getChannelVideos = asyncHandler(async (req, res) => {
  const { page = 1, limit = 10 } = req.query;

  const options = {
    page: parseInt(page) || 1,
    limit: parseInt(limit) || 10,
    sort: { createdAt: -1 },
  };

  const videos = Video.aggregate([
    { $match: { owner: new mongoose.Types.ObjectId(req.user._id) } },
  ]);

  const result = await Video.aggregatePaginate(videos, options);

  return res
    .status(200)
    .json(new ApiResponse(200, result, "channel videos fetched successfully"));
});

export { getChannelStats, getChannelVideos };
