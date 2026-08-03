import { Router } from "express";
import {
  createPost,
  deletePost,
  getUserPosts,
  updatePost,
} from "../controllers/post.controller.js";
import { verifyJWT } from "../middlewares/auth.js";

const router: Router = Router();
router.use(verifyJWT);

router.route("/").post(createPost);
router.route("/user/:userId").get(getUserPosts);
router.route("/:postId").patch(updatePost).delete(deletePost);

export default router;
