import { clsx } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs) {
  return twMerge(clsx(inputs))
}

export const readFileAsDataURL = (file) => {
  return new Promise((resolve) => {
    const reader = new FileReader();
    reader.onloadend = () => {
      if (typeof reader.result === 'string') resolve(reader.result);
    }
    reader.readAsDataURL(file);
  })
}


// Function to retrieve token (implementation details depend on your authentication mechanism)
export const  getToken = () =>{
  // Replace with your logic to retrieve the access token from storage (e.g., localStorage)
  const token = localStorage.getItem("accessToken");
  if (!token) {
    console.warn("No access token found. Authentication might be required.");
  }
  return token;
}