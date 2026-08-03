import mongoose, { Schema } from "mongoose";
import mongooseAggregatePaginate from "mongoose-aggregate-paginate-v2";

export interface IComment {
  content: string;
  video: mongoose.Types.ObjectId;
  owner: mongoose.Types.ObjectId;
}

export type CommentModel = mongoose.AggregatePaginateModel<IComment>;
export type ICommentDocument = mongoose.HydratedDocument<IComment>;

const commentSchema = new Schema<IComment, CommentModel>(
  {
    content: {
      type: String,
      required: true,
    },
    video: {
      type: mongoose.Types.ObjectId,
      ref: "Video",
    },
    owner: {
      type: mongoose.Types.ObjectId,
      ref: "User",
    },
  },
  { timestamps: true }
);

commentSchema.plugin(mongooseAggregatePaginate);

export const Comment = mongoose.model<IComment, CommentModel>(
  "Comment",
  commentSchema
);
