import axios from "axios";

const axiosInstance = axios.create({
  baseURL: import.meta.env.VITE_BACKEND_URL,
  withCredentials: true,
  validateStatus: (status) => {
    return status >= 200 && status < 500;
  },
});

// TODO: handle 401 error

export { axiosInstance };
