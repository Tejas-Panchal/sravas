import { ApiResponse } from "../utils/ApiResponse.ts";
import { asyncHandler } from "../utils/asyncHandler.ts";

const healthcheck = asyncHandler(async (_req, res) => {
  return res
    .status(200)
    .json(new ApiResponse(200, { status: "OK" }, "healthcheck passed"));
});

export { healthcheck };
