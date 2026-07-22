import mongoose, { isValidObjectId } from "mongoose";
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
  // TODO: add video to playlist
});

const removeVideoFromPlaylist = asyncHandler(async (req, res) => {
  const { playlistId, videoId } = req.params;
  // TODO: remove video from playlist
});

const deletePlaylist = asyncHandler(async (req, res) => {
  const { playlistId } = req.params;
  // TODO: delete playlist
});

const updatePlaylist = asyncHandler(async (req, res) => {
  const { playlistId } = req.params;
  const { name, description } = req.body;
  //TODO: update playlist
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
