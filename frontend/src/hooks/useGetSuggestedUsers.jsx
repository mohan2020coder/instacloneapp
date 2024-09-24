import { setSuggestedUsers } from "@/redux/authSlice";
import axios from "axios";
import { useEffect } from "react";
import { useDispatch } from "react-redux";
import {getToken} from "../lib/utils"; // Assuming getToken retrieves the token

const useGetSuggestedUsers = () => {
  const dispatch = useDispatch();

  useEffect(() => {
    const fetchSuggestedUsers = async () => {
      try {
        // Retrieve the access token from storage
        const token = getToken();

        if (!token) {
          console.warn("No access token found. User might not be authenticated.");
          // Handle the case where the token is not available
          return; // Exit if no token
        }

        const response = await axios.get(`${process.env.NEXT_PUBLIC_API_URL}/user/suggested`, {
          withCredentials: true,
          // Include the token in the Authorization header
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        if (response.data.success) {
          dispatch(setSuggestedUsers(response.data.users));
        } else {
          console.error("Failed to fetch suggested users:", response.data.message || "Unknown error");
        }
      } catch (error) {
        console.error("Error fetching suggested users:", error);
      }
    };

    fetchSuggestedUsers();
  }, [dispatch]); // Include dispatch in the dependency array

  return; // No need to return anything from a custom hook
};

export default useGetSuggestedUsers;