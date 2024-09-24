import { Heart, Home, LogOut, MessageCircle, PlusSquare, Search, TrendingUp } from 'lucide-react';
import React, { useState } from 'react';
import { Avatar, AvatarFallback, AvatarImage } from './ui/avatar';
import { toast } from 'sonner';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import { useDispatch, useSelector } from 'react-redux';
import { setAuthUser } from '@/redux/authSlice';
import CreatePost from './CreatePost';
import { setPosts, setSelectedPost } from '@/redux/postSlice';
import { Popover, PopoverContent, PopoverTrigger } from './ui/popover';
import { Button } from './ui/button';

const LeftSidebar = () => {
    const navigate = useNavigate();
    const { user } = useSelector(store => store.auth);
    const { likeNotification } = useSelector(store => store.realTimeNotification);
    const dispatch = useDispatch();
    const [open, setOpen] = useState(false);

    const logoutHandler = async () => {
        try {
            const res = await axios.get(`${process.env.NEXT_PUBLIC_API_URL}/user/logout`, { withCredentials: true });
            if (res.data.success) {
                dispatch(setAuthUser(null));
                dispatch(setSelectedPost(null));
                dispatch(setPosts([]));
                navigate("/login");
                toast.success(res.data.message);
            }
        } catch (error) {
            toast.error(error.response?.data?.message || 'An error occurred.');
        }
    };

    const sidebarHandler = (textType) => {
        switch (textType) {
            case 'Logout':
                logoutHandler();
                break;
            case 'Create':
                setOpen(true);
                break;
            case 'Profile':
                navigate(`/profile/${user?._id}`);
                break;
            case 'Home':
                navigate("/");
                break;
            case 'Messages':
                navigate("/chat");
                break;
            default:
                break;
        }
    };

    const sidebarItems = [
        { icon: <Home aria-label="Home" />, text: "Home" },
        { icon: <Search aria-label="Search" />, text: "Search" },
        { icon: <TrendingUp aria-label="Explore" />, text: "Explore" },
        { icon: <MessageCircle aria-label="Messages" />, text: "Messages" },
        { icon: <Heart aria-label="Notifications" />, text: "Notifications" },
        { icon: <PlusSquare aria-label="Create Post" />, text: "Create" },
        {
            icon: (
                <Avatar className='w-6 h-6'>
                    <AvatarImage src={user?.profilePicture} alt={user?.username} />
                    <AvatarFallback>{user?.username?.[0]}</AvatarFallback>
                </Avatar>
            ),
            text: "Profile"
        },
        { icon: <LogOut aria-label="Logout" />, text: "Logout" },
    ];

    return (
        <div className='fixed top-0 z-10 left-0 px-4 border-r border-gray-300 w-[16%] h-screen'>
            <div className='flex flex-col'>
                <h1 className='my-8 pl-3 font-bold text-xl'>LOGO</h1>
                <div>
                    {sidebarItems.map((item, index) => (
                        <div
                            key={index}
                            onClick={() => sidebarHandler(item.text)}
                            className='flex items-center gap-3 relative hover:bg-gray-100 cursor-pointer rounded-lg p-3 my-3'
                        >
                            {item.icon}
                            <span>{item.text}</span>
                            {item.text === "Notifications" && likeNotification.length > 0 && (
                                <Popover>
                                    <PopoverTrigger asChild>
                                        <Button
                                            size='icon'
                                            className="rounded-full h-5 w-5 bg-red-600 hover:bg-red-600 absolute bottom-6 left-6"
                                            aria-label={`Notifications (${likeNotification.length})`}
                                        >
                                            {likeNotification.length}
                                        </Button>
                                    </PopoverTrigger>
                                    <PopoverContent>
                                        <div>
                                            {likeNotification.length === 0 ? (
                                                <p>No new notifications</p>
                                            ) : (
                                                likeNotification.map((notification) => (
                                                    <div key={notification.userId} className='flex items-center gap-2 my-2'>
                                                        <Avatar>
                                                            <AvatarImage src={notification.userDetails?.profilePicture} />
                                                            <AvatarFallback>{notification.userDetails?.username?.[0]}</AvatarFallback>
                                                        </Avatar>
                                                        <p className='text-sm'>
                                                            <span className='font-bold'>{notification.userDetails?.username}</span> liked your post
                                                        </p>
                                                    </div>
                                                ))
                                            )}
                                        </div>
                                    </PopoverContent>
                                </Popover>
                            )}
                        </div>
                    ))}
                </div>
            </div>
            <CreatePost open={open} setOpen={setOpen} />
        </div>
    );
};

export default LeftSidebar;
