import { setMessages } from "@/redux/chatSlice";
import axios from "axios";
import { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";

const useGetAllMessage = () => {
    const dispatch = useDispatch();
    const { selectedUser } = useSelector(store => store.auth);

    useEffect(() => {
        const fetchAllMessage = async () => {
            try {
                const res = await axios.get(`${process.env.NEXT_PUBLIC_API_URL}/message/all/${selectedUser?._id}`, { withCredentials: true });
                if (res.data.success) {  
                    dispatch(setMessages(res.data.messages));
                }
            } catch (error) {
                console.log(error);
            }
        };

        if (selectedUser?._id) { // Ensure selectedUser exists before making the request
            fetchAllMessage();
        }
    }, [selectedUser, dispatch]); // Add dispatch to dependencies if it is used in the effect

    // Optionally, you can return something if needed
};

export default useGetAllMessage;
