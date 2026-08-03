import { Router } from "express";
import {
  createPost,
  deletePost,
  getUserPosts,
  updatePost,
} from "../controllers/post.controller.ts";
import { verifyJWT } from "../middlewares/auth.ts";

const router: Router = Router();
router.use(verifyJWT);

router.route("/").post(createPost);
router.route("/user/:userId").get(getUserPosts);
router.route("/:postId").patch(updatePost).delete(deletePost);

export default router;
