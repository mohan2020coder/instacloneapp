import { setPosts } from "@/redux/postSlice";
import axios from "axios";
import { useEffect } from "react";
import { useDispatch } from "react-redux";
import  {getToken} from "../lib/utils";

const useGetAllPosts = () => { // Use plural for consistency
  const dispatch = useDispatch();

  useEffect(() => {
    const fetchAllPosts = async () => {
      try {
        const response = await axios.get(`${process.env.NEXT_PUBLIC_API_URL}/post/all`, {
          withCredentials: true,
          // Include headers for token if necessary
          headers: {
            Authorization: `Bearer ${getToken()}`, // Replace with your token retrieval function
          },
        });

        if (response.data.success) {
          console.log(response.data.posts);
          dispatch(setPosts(response.data.posts));
        } else {
          // Handle unsuccessful response (e.g., display error message)
          console.error("Failed to fetch posts:", response.data.message || "Unknown error");
        }
      } catch (error) {
        console.error("Error fetching posts:", error);
      }
    };

    fetchAllPosts();
  }, [dispatch]); // Add dispatch to the dependency array

  return; // No need to return anything from a custom hook
};

export default useGetAllPosts;

