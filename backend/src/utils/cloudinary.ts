import { v2 as cloudinary, type UploadApiResponse } from "cloudinary";
import fs from "fs";

cloudinary.config({
  cloud_name: process.env.CLOUDINARY_CLOUD_NAME,
  api_key: process.env.CLOUDINARY_API_KEY,
  api_secret: process.env.CLOUDINARY_API_SECRET,
});

const uploadOnCloudinary = async (
  localFilePath: string | undefined
): Promise<any | null> => {
  try {
    if (!localFilePath) return null;
    const response: UploadApiResponse = await cloudinary.uploader.upload(
      localFilePath,
      {
        resource_type: "auto",
      }
    );
    fs.unlinkSync(localFilePath);
    return response;
  } catch (error) {
    fs.unlinkSync(localFilePath as string);
    return null;
  }
};

const deleteFromCloudinary = async (url: string | undefined) => {
  if (!url) return null;
  try {
    const publicId = url.split("/").pop()?.split(".")[0];
    if (!publicId) return null;
    return await cloudinary.uploader.destroy(publicId);
  } catch {
    return null;
  }
};

export { uploadOnCloudinary, deleteFromCloudinary };
