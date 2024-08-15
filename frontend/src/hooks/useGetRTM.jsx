import { setMessages } from "@/redux/chatSlice";
import { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";

const useGetRTM = () => {
    const dispatch = useDispatch();
    const { socket } = useSelector(store => store.socketio);
    const { messages } = useSelector(store => store.chat);

    useEffect(() => {
        const handleNewMessage = (newMessage) => {
            dispatch(setMessages([...messages, newMessage]));
        };

        if (socket) {
            socket.on('newMessage', handleNewMessage);
        }

        return () => {
            if (socket) {
                socket.off('newMessage', handleNewMessage);
            }
        };
    }, [socket, messages, dispatch]); // Include all necessary dependencies

};

export default useGetRTM;
