// components/SocketProvider.js
import { useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { io } from 'socket.io-client';
import { setSocket } from '../redux/socketSlice';
import { setOnlineUsers } from '../redux/chatSlice';
import { setLikeNotification } from '../redux/rtnSlice';

const SocketProvider = ({ children }) => {
  const user = useSelector((state) => state.auth.user);
  const socket = useSelector((state) => state.socketio.socket);
  const dispatch = useDispatch();

  useEffect(() => {
    if (user) {
      const socketio = io(process.env.NEXT_PUBLIC_SOCKET_URL, {
        query: {
          userId: user._id,
        },
        transports: ['websocket'],
      });

      dispatch(setSocket(socketio));

      socketio.on('getOnlineUsers', (onlineUsers) => {
        dispatch(setOnlineUsers(onlineUsers));
      });

      socketio.on('notification', (notification) => {
        dispatch(setLikeNotification(notification));
      });

      return () => {
        socketio.close();
        dispatch(setSocket(null));
      };
    } else if (socket) {
      socket.close();
      dispatch(setSocket(null));
    }
  }, [user, socket, dispatch]);

  return children;
};

export default SocketProvider;
