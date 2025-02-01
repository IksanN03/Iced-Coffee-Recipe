import { jwtDecode } from "jwt-decode";

export const setToken = (token: string) => localStorage.setItem('token', token);

export const getToken = () => localStorage.getItem('token');

export const removeToken = () => localStorage.removeItem('token');

export const isAuthenticated = () => !!getToken();

export const getEmailFromToken = () => {
  const token = getToken();
  if (token) {
    const decoded = jwtDecode<{ email: string }>(token);
    return decoded.email;
  }
  return '';
};
