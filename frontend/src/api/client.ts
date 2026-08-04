import axios from "axios";
import config from "../config/config.ts";

const client = axios.create({
  baseURL: config.apiBaseUrl,
  withCredentials: true,
});

export default client;
