import mongoose, { isValidObjectId } from "mongoose";
import { Tweet } from "../models/tweet.model.js";
import { User } from "../models/user.model.js";
import { ApiError } from "../utils/ApiError.js";
import { ApiResponse } from "../utils/ApiResponse.js";
import { asyncHandler } from "../utils/asyncHandler.js";

const createTweet = asyncHandler(async (req, res) => {
  const { content } = req.body;

  if (!content?.trim()) throw new ApiError(400, "tweet content is required");

  const tweet = await Tweet.create({
    content,
    owner: req.user._id,
  });

  if (!tweet)
    throw new ApiError(500, "something went wrong while creating the tweet");

  return res
    .status(201)
    .json(new ApiResponse(201, tweet, "tweet created successfully"));
});

const getUserTweets = asyncHandler(async (req, res) => {
  const { userId } = req.params;

  if (!isValidObjectId(userId)) throw new ApiError(400, "invalid user id");

  const tweets = await Tweet.find({ owner: userId }).sort({ createdAt: -1 });

  if (!tweets)
    throw new ApiError(500, "something went wrong while fetching tweets");

  return res
    .status(200)
    .json(new ApiResponse(200, tweets, "user tweets fetched successfully"));
});

const updateTweet = asyncHandler(async (req, res) => {
  const { tweetId } = req.params;
  const { content } = req.body;

  if (!isValidObjectId(tweetId)) throw new ApiError(400, "invalid tweet id");
  if (!content?.trim()) throw new ApiError(400, "tweet content is required");

  const tweet = await Tweet.findById(tweetId);

  if (!tweet) throw new ApiError(404, "tweet not found");
  if (tweet.owner.toString() !== req.user._id.toString())
    throw new ApiError(403, "you can only update your own tweets");

  tweet.content = content;
  await tweet.save();

  return res
    .status(200)
    .json(new ApiResponse(200, tweet, "tweet updated successfully"));
});

const deleteTweet = asyncHandler(async (req, res) => {
  const { tweetId } = req.params;

  if (!isValidObjectId(tweetId)) throw new ApiError(400, "invalid tweet id");

  const tweet = await Tweet.findById(tweetId);

  if (!tweet) throw new ApiError(404, "tweet not found");
  if (tweet.owner.toString() !== req.user._id.toString())
    throw new ApiError(403, "you can only delete your own tweets");

  await Tweet.findByIdAndDelete(tweetId);

  return res
    .status(200)
    .json(new ApiResponse(200, {}, "tweet deleted successfully"));
});

export { createTweet, getUserTweets, updateTweet, deleteTweet };
