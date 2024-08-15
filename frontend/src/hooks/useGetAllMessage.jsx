import { setMessages } from "@/redux/chatSlice";
import axios from "axios";
import { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import {getToken} from "../lib/utils"; // Assuming getToken retrieves the token

const useGetAllMessage = () => {
  const dispatch = useDispatch();
  const { selectedUser } = useSelector(store => store.auth);

  useEffect(() => {
    const fetchAllMessage = async () => {
      if (!selectedUser?._id) return; // Handle case where selectedUser is not available

      try {
        // Retrieve the access token from storage
        const token = getToken();

        if (!token) {
          console.warn("No access token found. Messages might not be accessible.");
          return; // Exit if no token
        }

        const response = await axios.get(`${process.env.NEXT_PUBLIC_API_URL}/message/all/${selectedUser?._id}`, {
          withCredentials: true,
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        if (response.data.success) {
          dispatch(setMessages(response.data.messages));
        } else {
          console.error("Failed to fetch messages:", response.data.message || "Unknown error");
        }
      } catch (error) {
        console.error("Error fetching messages:", error);
      }
    };

    fetchAllMessage();
  }, [selectedUser, dispatch]); // Include all necessary dependencies

  return; // No need to return anything from a custom hook
};

export default useGetAllMessage;