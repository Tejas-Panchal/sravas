import { isValidObjectId } from "mongoose";
import { Playlist } from "../models/playlist.model.js";
import { ApiError } from "../utils/ApiError.js";
import { ApiResponse } from "../utils/ApiResponse.js";
import { asyncHandler } from "../utils/asyncHandler.js";

const createPlaylist = asyncHandler(async (req, res) => {
  const { name, description } = req.body;

  if (!name?.trim()) throw new ApiError(400, "playlist name is required");
  if (!description?.trim())
    throw new ApiError(400, "playlist description is required");

  const playlist = await Playlist.create({
    name,
    description,
    owner: req.user._id,
    videos: [],
  });

  if (!playlist)
    throw new ApiError(500, "something went wrong while creating the playlist");

  return res
    .status(201)
    .json(new ApiResponse(201, playlist, "playlist created successfully"));
});

const getUserPlaylists = asyncHandler(async (req, res) => {
  const { userId } = req.params;

  if (!isValidObjectId(userId)) throw new ApiError(400, "invalid user id");

  const playlists = await Playlist.find({ owner: userId }).sort({
    createdAt: -1,
  });

  return res
    .status(200)
    .json(
      new ApiResponse(200, playlists, "user playlists fetched successfully")
    );
});

const getPlaylistById = asyncHandler(async (req, res) => {
  const { playlistId } = req.params;

  if (!isValidObjectId(playlistId))
    throw new ApiError(400, "invalid playlist id");

  const playlist = await Playlist.findById(playlistId);

  if (!playlist) throw new ApiError(404, "playlist not found");

  return res
    .status(200)
    .json(new ApiResponse(200, playlist, "playlist fetched successfully"));
});

const addVideoToPlaylist = asyncHandler(async (req, res) => {
  const { playlistId, videoId } = req.params;

  if (!isValidObjectId(playlistId))
    throw new ApiError(400, "invalid playlist id");
  if (!isValidObjectId(videoId)) throw new ApiError(400, "invalid video id");

  const playlist = await Playlist.findById(playlistId);

  if (!playlist) throw new ApiError(404, "playlist not found");

  if (playlist.owner.toString() !== req.user._id.toString())
    throw new ApiError(403, "you can only add to your own playlists");

  if (playlist.videos.includes(videoId))
    throw new ApiError(400, "video already in playlist");

  playlist.videos.push(videoId);
  await playlist.save();

  return res
    .status(200)
    .json(
      new ApiResponse(200, playlist, "video added to playlist successfully")
    );
});

const removeVideoFromPlaylist = asyncHandler(async (req, res) => {
  const { playlistId, videoId } = req.params;

  if (!isValidObjectId(playlistId))
    throw new ApiError(400, "invalid playlist id");
  if (!isValidObjectId(videoId)) throw new ApiError(400, "invalid video id");

  const playlist = await Playlist.findById(playlistId);

  if (!playlist) throw new ApiError(404, "playlist not found");

  if (playlist.owner.toString() !== req.user._id.toString())
    throw new ApiError(403, "you can only remove from your own playlists");

  if (!playlist.videos.includes(videoId))
    throw new ApiError(404, "video not found in playlist");

  playlist.videos.pull(videoId);
  await playlist.save();

  return res
    .status(200)
    .json(
      new ApiResponse(200, playlist, "video removed from playlist successfully")
    );
});

const deletePlaylist = asyncHandler(async (req, res) => {
  const { playlistId } = req.params;

  if (!isValidObjectId(playlistId))
    throw new ApiError(400, "invalid playlist id");

  const playlist = await Playlist.findById(playlistId);

  if (!playlist) throw new ApiError(404, "playlist not found");

  if (playlist.owner.toString() !== req.user._id.toString())
    throw new ApiError(403, "you can only delete your own playlists");

  await Playlist.findByIdAndDelete(playlistId);

  return res
    .status(200)
    .json(new ApiResponse(200, {}, "playlist deleted successfully"));
});

const updatePlaylist = asyncHandler(async (req, res) => {
  const { playlistId } = req.params;
  const { name, description } = req.body;

  if (!isValidObjectId(playlistId))
    throw new ApiError(400, "invalid playlist id");

  if (!name?.trim() && !description?.trim())
    throw new ApiError(
      400,
      "at least one field (name or description) is required"
    );

  const playlist = await Playlist.findById(playlistId);

  if (!playlist) throw new ApiError(404, "playlist not found");

  if (playlist.owner.toString() !== req.user._id.toString())
    throw new ApiError(403, "you can only update your own playlists");

  if (name) playlist.name = name;
  if (description) playlist.description = description;
  await playlist.save();

  return res
    .status(200)
    .json(new ApiResponse(200, playlist, "playlist updated successfully"));
});

export {
  createPlaylist,
  getUserPlaylists,
  getPlaylistById,
  addVideoToPlaylist,
  removeVideoFromPlaylist,
  deletePlaylist,
  updatePlaylist,
};
