import mongoose, { isValidObjectId } from "mongoose";
import { Comment } from "../models/comment.model.js";
import { ApiError } from "../utils/ApiError.js";
import { ApiResponse } from "../utils/ApiResponse.js";
import { asyncHandler } from "../utils/asyncHandler.js";

const getVideoComments = asyncHandler(async (req, res) => {
  const { videoId } = req.params;
  const { page = 1, limit = 10 } = req.query;

  if (!videoId) throw new ApiError(400, "video id missing");

  if (!isValidObjectId(videoId)) throw new ApiError(400, "invalid video id");

  const comments = Comment.aggregate([
    {
      $match: {
        video: new mongoose.Types.ObjectId(videoId),
      },
    },
    {
      $lookup: {
        from: "users",
        localField: "owner",
        foreignField: "_id",
        as: "owner",
        pipeline: [
          {
            $project: {
              _id: 1,
              username: 1,
              fullName: 1,
              avatar: 1,
            },
          },
        ],
      },
    },
    {
      $addFields: {
        owner: { $arrayElemAt: ["$owner", 0] },
      },
    },
    {
      $project: {
        content: 1,
        video: 1,
        owner: 1,
        createdAt: 1,
        updatedAt: 1,
      },
    },
  ]);

  const options = {
    page: parseInt(page) || 1,
    limit: parseInt(limit) || 10,
    sort: { createdAt: -1 },
  };

  const result = await Comment.aggregatePaginate(comments, options);

  return res
    .status(200)
    .json(new ApiResponse(200, result, "video comments fetched successfully"));
});

const addComment = asyncHandler(async (req, res) => {
  const { videoId } = req.params;
  const { content } = req.body;

  if (!content?.trim()) throw new ApiError(400, "comment content is required");

  if (!isValidObjectId(videoId)) throw new ApiError(400, "invalid video id");

  const comment = await Comment.create({
    content,
    video: videoId,
    owner: req.user._id,
  });

  const createdComment = await Comment.findById(comment._id).populate(
    "owner",
    "_id username fullName avatar"
  );

  if (!createdComment)
    throw new ApiError(500, "somthing went wrong while adding the comment");

  return res
    .status(201)
    .json(new ApiResponse(201, createdComment, "comment added successfully"));
});

const updateComment = asyncHandler(async (req, res) => {
  const { commentId } = req.params;
  const { content } = req.body;

  if (!isValidObjectId(commentId))
    throw new ApiError(400, "invalid comment id");

  if (!content?.trim()) throw new ApiError(400, "content is required");

  const comment = await Comment.findById(commentId);

  if (!comment) throw new ApiError(404, "comment not found");

  if (comment.owner.toString() !== req.user._id.toString())
    throw new ApiError(403, "you can only update your own comments");

  comment.content = content;
  await comment.save();

  const updatedComment = await Comment.findById(commentId).populate(
    "owner",
    "_id username fullName avatar"
  );

  return res
    .status(200)
    .json(new ApiResponse(200, updatedComment, "comment updated successfully"));
});

const deleteComment = asyncHandler(async (req, res) => {
  const { commentId } = req.params;

  if (!isValidObjectId(commentId))
    throw new ApiError(400, "invalid comment id");

  const comment = await Comment.findById(commentId);

  if (!comment) throw new ApiError(404, "comment not found");

  if (comment.owner.toString() !== req.user._id.toString())
    throw new ApiError(403, "you can only delete your own comments");

  await Comment.findByIdAndDelete(commentId);

  return res
    .status(200)
    .json(new ApiResponse(200, {}, "comment deleted successfully"));
});

export { getVideoComments, addComment, updateComment, deleteComment };
