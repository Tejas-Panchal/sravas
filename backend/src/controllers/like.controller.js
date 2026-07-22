import { isValidObjectId } from "mongoose";
import { Like } from "../models/like.model.js";
import { ApiError } from "../utils/ApiError.js";
import { ApiResponse } from "../utils/ApiResponse.js";
import { asyncHandler } from "../utils/asyncHandler.js";

const toggleVideoLike = asyncHandler(async (req, res) => {
  const { videoId } = req.params;

  if (!isValidObjectId(videoId)) throw new ApiError(400, "invalid video id");

  const existingLike = await Like.findOne({
    video: videoId,
    likedBy: req.user._id,
  });

  if (existingLike) {
    await Like.findByIdAndDelete(existingLike._id);
    return res
      .status(200)
      .json(new ApiResponse(200, { isLiked: false }, "video unliked"));
  }

  const like = await Like.create({
    video: videoId,
    likedBy: req.user._id,
  });

  if (!like) throw new ApiError(500, "something went wrong while liking video");

  return res
    .status(201)
    .json(new ApiResponse(201, { isLiked: true }, "video liked"));
});

const toggleCommentLike = asyncHandler(async (req, res) => {
  const { commentId } = req.params;

  if (!isValidObjectId(commentId))
    throw new ApiError(400, "invalid comment id");

  const existingLike = await Like.findOne({
    comment: commentId,
    likedBy: req.user._id,
  });

  if (existingLike) {
    await Like.findByIdAndDelete(existingLike._id);
    return res
      .status(200)
      .json(new ApiResponse(200, { isLiked: false }, "comment unliked"));
  }

  const like = await Like.create({
    comment: commentId,
    likedBy: req.user._id,
  });

  if (!like)
    throw new ApiError(500, "something went wrong while liking comment");

  return res
    .status(201)
    .json(new ApiResponse(201, { isLiked: true }, "comment liked"));
});

const toggleTweetLike = asyncHandler(async (req, res) => {
  const { tweetId } = req.params;
  if (!isValidObjectId(tweetId)) throw new ApiError(400, "invalid tweet id");

  const existingLike = await Like.findOne({
    tweet: tweetId,
    likedBy: req.user._id,
  });

  if (existingLike) {
    await Like.findByIdAndDelete(existingLike._id);
    return res
      .status(200)
      .json(new ApiResponse(200, { isLiked: false }, "tweet unliked"));
  }

  const like = await Like.create({
    tweet: tweetId,
    likedBy: req.user._id,
  });

  if (!like) throw new ApiError(500, "something went wrong while liking tweet");

  return res
    .status(201)
    .json(new ApiResponse(201, { isLiked: true }, "tweet liked"));
});

const getLikedVideos = asyncHandler(async (req, res) => {
  const likes = await Like.find({
    likedBy: req.user._id,
    video: { $ne: null },
  })
    .populate("video")
    .sort({ createdAt: -1 });

  const videos = likes.map((like) => like.video).filter(Boolean);

  return res
    .status(200)
    .json(new ApiResponse(200, videos, "liked videos fetched successfully"));
});

export { toggleCommentLike, toggleTweetLike, toggleVideoLike, getLikedVideos };

