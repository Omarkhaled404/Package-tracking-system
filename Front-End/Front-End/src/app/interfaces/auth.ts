export interface RegisterPostData {
    name: string;
    email: string;
    phone: string;
    password: string;
    role: string;

}
export interface LoginPostData {
    email: string;
    password: string;
  }
  
// Define the User interface
export interface User {
    id?: number;       // Optional user ID (or any other unique identifier)
    name: string;      // Full name of the user
    email: string;     // Email address
    phone?: string;    // Optional phone number
    role: string;      // Role of the user, e.g., "admin" or "user"
  }
  
