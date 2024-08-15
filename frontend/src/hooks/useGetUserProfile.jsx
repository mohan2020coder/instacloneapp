import { setUserProfile } from "@/redux/authSlice";
import axios from "axios";
import { useEffect } from "react";
import { useDispatch } from "react-redux";

const useGetUserProfile = (userId) => {
    const dispatch = useDispatch();

    useEffect(() => {
        const fetchUserProfile = async () => {
            if (!userId) return; // Handle case where userId is not provided

            try {
                const res = await axios.get(`${process.env.NEXT_PUBLIC_API_URL}/user/${userId}/profile`, { withCredentials: true });
                if (res.data.success) {
                    dispatch(setUserProfile(res.data.user));
                }
            } catch (error) {
                console.error(error); // Use console.error for error logging
            }
        };

        fetchUserProfile();
    }, [userId, dispatch]); // Include dispatch in dependencies

};

export default useGetUserProfile;
