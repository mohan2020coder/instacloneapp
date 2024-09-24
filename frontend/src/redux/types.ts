// redux/types.ts
import { Socket } from 'socket.io-client';

// Define your state shape
export interface AuthState {
  user: { _id: string } | null;
}

export interface SocketState {
  socket: Socket | null;
}

export interface ChatState {
  onlineUsers: any[]; // Update this type as needed
}

export interface RtnState {
  notifications: any[]; // Update this type as needed
}

// RootState type definition
export interface RootState {
  auth: AuthState;
  socketio: SocketState;
  chat: ChatState;
  rtn: RtnState;
}
