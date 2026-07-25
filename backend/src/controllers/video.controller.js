import mongoose, { isValidObjectId } from "mongoose";
import { Video } from "../models/video.model.js";
import { User } from "../models/user.model.js";
import { ApiError } from "../utils/ApiError.js";
import { ApiResponse } from "../utils/ApiResponse.js";
import { asyncHandler } from "../utils/asyncHandler.js";
import { uploadOnCloudinary } from "../utils/cloudinary.js";

const getAllVideos = asyncHandler(async (req, res) => {
  const { page = 1, limit = 10, query, sortBy, sortType, userId } = req.query;

  const pipeline = [];

  const match = {};
  if (query) {
    match.$or = [
      { title: { $regex: query, $options: "i" } },
      { description: { $regex: query, $options: "i" } },
    ];
  }

  if (userId && isValidObjectId(userId)) {
    match.owner = new mongoose.Types.ObjectId(userId);
  }

  if (!userId || userId !== req.user._id.toString()) {
    match.isPublished = true;
  }

  pipeline.push({ $match: match });
  pipeline.push({
    $lookup: {
      from: "users",
      localField: "owner",
      foreignField: "_id",
      as: "owner",
      pipeline: [{ $project: { _id: 1, username: 1, fullName: 1, avatar: 1 } }],
    },
  });
  pipeline.push({ $addFields: { owner: { $arrayElemAt: ["$owner", 0] } } });

  const sortField = sortBy || "createdAt";
  const sortOrder = sortType === "asc" ? 1 : -1;
  pipeline.push({ $sort: { [sortField]: sortOrder } });

  pipeline.push({
    $project: {
      videoFile: 1,
      thumbnail: 1,
      title: 1,
      description: 1,
      duration: 1,
      views: 1,
      isPublished: 1,
      owner: 1,
      createdAt: 1,
      updatedAt: 1,
    },
  });

  const options = {
    page: parseInt(page) || 1,
    limit: parseInt(limit) || 10,
  };

  const aggregate = Video.aggregate(pipeline);
  const result = await Video.aggregatePaginate(aggregate, options);

  return res
    .status(200)
    .json(new ApiResponse(200, result, "videos fetched successfully"));
});

const publishAVideo = asyncHandler(async (req, res) => {
  const { title, description } = req.body;

  if (!title?.trim()) throw new ApiError(400, "video title is required");
  if (!description?.trim())
    throw new ApiError(400, "video description is required");

  const videoLocalPath = req.files?.videoFile?.[0]?.path;
  const thumbnailLocalPath = req.files?.thumbnail?.[0]?.path;

  if (!videoLocalPath) throw new ApiError(400, "video file is required");
  if (!thumbnailLocalPath) throw new ApiError(400, "thumbnail is required");

  const videoUpload = await uploadOnCloudinary(videoLocalPath);
  const thumbnailUpload = await uploadOnCloudinary(thumbnailLocalPath);

  if (!videoUpload) throw new ApiError(500, "failed to upload video");
  if (!thumbnailUpload) throw new ApiError(500, "failed to upload thumbnail");

  const video = await Video.create({
    videoFile: videoUpload.url,
    thumbnail: thumbnailUpload.url,
    title,
    description,
    duration: videoUpload.duration,
    owner: req.user._id,
  });

  if (!video)
    throw new ApiError(500, "something went wrong while publishing the video");

  return res
    .status(201)
    .json(new ApiResponse(201, video, "video published successfully"));
});

const getVideoById = asyncHandler(async (req, res) => {
  const { videoId } = req.params;

  if (!isValidObjectId(videoId)) throw new ApiError(400, "invalid video id");

  const video = await Video.findById(videoId).populate(
    "owner",
    "_id username fullName avatar"
  );

  if (!video) throw new ApiError(404, "video not found");

  video.views += 1;
  await video.save({ validateBeforeSave: false });

  return res
    .status(200)
    .json(new ApiResponse(200, video, "video fetched successfully"));
});

const updateVideo = asyncHandler(async (req, res) => {
  const { videoId } = req.params;
  const { title, description } = req.body;

  if (!title?.trim() && !description?.trim() && !req.file)
    throw new ApiError(400, "at least one field to update is required");

  if (!isValidObjectId(videoId)) throw new ApiError(400, "invalid video id");

  const video = await Video.findById(videoId);

  if (!video) throw new ApiError(404, "video not found");

  if (video.owner.toString() !== req.user._id.toString())
    throw new ApiError(403, "you can only update your own videos");

  if (title) video.title = title;
  if (description) video.description = description;

  if (req.file) {
    const thumbnailUpload = await uploadOnCloudinary(req.file.path);
    if (!thumbnailUpload) throw new ApiError(500, "failed to upload thumbnail");
    video.thumbnail = thumbnailUpload.url;
  }

  await video.save({ validateBeforeSave: false });

  return res
    .status(200)
    .json(new ApiResponse(200, video, "video updated successfully"));
});

const deleteVideo = asyncHandler(async (req, res) => {
  const { videoId } = req.params;

  if (!isValidObjectId(videoId)) throw new ApiError(400, "invalid video id");

  const video = await Video.findById(videoId);

  if (!video) throw new ApiError(404, "video not found");

  if (video.owner.toString() !== req.user._id.toString())
    throw new ApiError(403, "you can only delete your own videos");

  await Video.findByIdAndDelete(videoId);

  return res
    .status(200)
    .json(new ApiResponse(200, {}, "video deleted successfully"));
});

const togglePublishStatus = asyncHandler(async (req, res) => {
  const { videoId } = req.params;

  if (!isValidObjectId(videoId)) throw new ApiError(400, "invalid video id");

  const video = await Video.findById(videoId);

  if (!video) throw new ApiError(404, "video not found");

  if (video.owner.toString() !== req.user._id.toString())
    throw new ApiError(403, "you can only toggle your own videos");

  video.isPublished = !video.isPublished;
  await video.save({ validateBeforeSave: false });

  return res
    .status(200)
    .json(
      new ApiResponse(
        200,
        { isPublished: video.isPublished },
        "publish status toggled successfully"
      )
    );
});

export {
  getAllVideos,
  publishAVideo,
  getVideoById,
  updateVideo,
  deleteVideo,
  togglePublishStatus,
};
