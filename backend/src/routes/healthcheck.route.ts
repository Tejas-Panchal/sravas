import { Router } from "express";
import { healthcheck } from "../controllers/healthcheck.controller.ts";

const router: Router = Router();

router.route("/").get(healthcheck);

export default router;
