import { isValidObjectId } from "mongoose";
import { Post } from "../models/post.model.ts";
import { ApiError } from "../utils/ApiError.ts";
import { ApiResponse } from "../utils/ApiResponse.ts";
import { asyncHandler } from "../utils/asyncHandler.ts";

const createPost = asyncHandler(async (req, res) => {
  const { content } = req.body;

  if (!content?.trim()) throw new ApiError(400, "post content is required");

  const post = await Post.create({
    content,
    owner: req.user._id,
  });

  if (!post)
    throw new ApiError(500, "something went wrong while creating the post");

  return res
    .status(201)
    .json(new ApiResponse(201, post, "post created successfully"));
});

const getUserPosts = asyncHandler(async (req, res) => {
  const { userId } = req.params as { userId: string };

  if (!isValidObjectId(userId)) throw new ApiError(400, "invalid user id");

  const posts = await Post.find({ owner: userId }).sort({ createdAt: -1 });

  if (!posts)
    throw new ApiError(500, "something went wrong while fetching posts");

  return res
    .status(200)
    .json(new ApiResponse(200, posts, "user posts fetched successfully"));
});

const updatePost = asyncHandler(async (req, res) => {
  const { postId } = req.params as { postId: string };
  const { content } = req.body;

  if (!isValidObjectId(postId)) throw new ApiError(400, "invalid post id");
  if (!content?.trim()) throw new ApiError(400, "post content is required");

  const post = await Post.findById(postId);
  if (!post) throw new ApiError(404, "post not found");

  if (post.owner.toString() !== req.user._id.toString())
    throw new ApiError(403, "you can only update your own posts");

  post.content = content;
  await post.save();

  return res
    .status(200)
    .json(new ApiResponse(200, post, "post updated successfully"));
});

const deletePost = asyncHandler(async (req, res) => {
  const { postId } = req.params as { postId: string };

  if (!isValidObjectId(postId)) throw new ApiError(400, "invalid post id");

  const post = await Post.findById(postId);
  if (!post) throw new ApiError(404, "post not found");

  if (post.owner.toString() !== req.user._id.toString())
    throw new ApiError(403, "you can only delete your own posts");

  await Post.findByIdAndDelete(postId);

  return res
    .status(200)
    .json(new ApiResponse(200, {}, "post deleted successfully"));
});

export { createPost, getUserPosts, updatePost, deletePost };
