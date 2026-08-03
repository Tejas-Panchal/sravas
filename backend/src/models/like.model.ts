import mongoose, { Schema } from "mongoose";

export interface ILike {
  video?: mongoose.Types.ObjectId;
  comment?: mongoose.Types.ObjectId;
  post?: mongoose.Types.ObjectId;
  likedBy: mongoose.Types.ObjectId;
}

export type ILikeDocument = mongoose.HydratedDocument<ILike>;

const likeSchema = new Schema<ILike>(
  {
    video: { type: Schema.Types.ObjectId, ref: "Video" },
    comment: { type: Schema.Types.ObjectId, ref: "Comment" },
    post: { type: Schema.Types.ObjectId, ref: "Post" },
    likedBy: { type: Schema.Types.ObjectId, ref: "User" },
  },
  { timestamps: true }
);

export const Like = mongoose.model<ILike>("Like", likeSchema);
