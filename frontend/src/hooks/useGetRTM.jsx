import { setMessages } from "@/redux/chatSlice";
import { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";

const useGetRTM = () => {
  const dispatch = useDispatch();
  const { socket, isAuthenticated } = useSelector(store => store.socketio);
  const { messages } = useSelector(store => store.chat);

  useEffect(() => {
    const handleNewMessage = (newMessage) => {
      dispatch(setMessages([...messages, newMessage]));
    };

    // Check authentication before subscribing to events
    if (socket && isAuthenticated) {
      socket.on('newMessage', handleNewMessage);
    }

    return () => {
      if (socket) {
        socket.off('newMessage', handleNewMessage);
      }
    };
  }, [socket, messages, dispatch, isAuthenticated]); // Include all necessary dependencies

  return; // No need to return anything from a custom hook
};

export default useGetRTM;