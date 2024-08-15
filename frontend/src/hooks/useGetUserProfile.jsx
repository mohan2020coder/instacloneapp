import { setUserProfile } from "@/redux/authSlice";
import axios from "axios";
import { useEffect } from "react";
import { useDispatch } from "react-redux";
import {getToken} from "../lib/utils"; // Assuming getToken retrieves the token

const useGetUserProfile = (userId) => {
  const dispatch = useDispatch();

  useEffect(() => {
    const fetchUserProfile = async () => {
      if (!userId) return; // Handle case where userId is not provided

      try {
        // Retrieve the access token from storage
        const token = getToken();

        if (!token) {
          console.warn("No access token found. User profile might not be accessible.");
          return; // Exit if no token
        }

        const response = await axios.get(`${process.env.NEXT_PUBLIC_API_URL}/user/${userId}/profile`, {
          withCredentials: true,
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        if (response.data.success) {
          dispatch(setUserProfile(response.data.user));
        } else {
          console.error("Failed to fetch user profile:", response.data.message || "Unknown error");
        }
      } catch (error) {
        console.error("Error fetching user profile:", error);
      }
    };

    fetchUserProfile();
  }, [userId, dispatch]); // Include all necessary dependencies

  return; // No need to return anything from a custom hook
};

export default useGetUserProfile;