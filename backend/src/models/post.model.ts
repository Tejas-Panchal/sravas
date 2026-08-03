import mongoose, { Schema } from "mongoose";

export interface IPost {
  content: string;
  owner: mongoose.Types.ObjectId;
}

export type IPostDocument = mongoose.HydratedDocument<IPost>;

const postSchema = new Schema<IPost>(
  {
    content: {
      type: String,
      required: true,
    },
    owner: {
      type: mongoose.Types.ObjectId,
      ref: "User",
    },
  },
  { timestamps: true }
);

export const Post = mongoose.model<IPost>("Post", postSchema);
