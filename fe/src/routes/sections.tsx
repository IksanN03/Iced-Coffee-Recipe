import { lazy, Suspense, useEffect, useRef, useState } from 'react';
import { Outlet, Navigate, useRoutes,  useSearchParams, useNavigate } from 'react-router-dom';
import { isAuthenticated } from 'src/utils/auth';
import Box from '@mui/material/Box';
import LinearProgress, { linearProgressClasses } from '@mui/material/LinearProgress';
import { varAlpha } from 'src/theme/styles';
import { AuthLayout } from 'src/layouts/auth';
import { DashboardLayout } from 'src/layouts/dashboard';
import { Alert, AlertColor, Snackbar } from '@mui/material';

export const HomePage = lazy(() => import('src/pages/home'));
export const BlogPage = lazy(() => import('src/pages/blog'));
export const InventoryPage = lazy(() => import('src/pages/inventory'));
export const RecipePage = lazy(() => import('src/pages/recipe'));
export const SignInPage = lazy(() => import('src/pages/sign-in'));
export const ProductsPage = lazy(() => import('src/pages/products'));
export const Page404 = lazy(() => import('src/pages/page-not-found'));

const PrivateRoute = ({ children }: { children: React.ReactNode }) => {
  const [searchParams] = useSearchParams();
  const token = searchParams.get('token');
  
  if (!token) {
    return isAuthenticated() ? children : <Navigate to="/sign-in" replace />;
  }
  
  return null;

};



const renderFallback = (
  <Box display="flex" alignItems="center" justifyContent="center" flex="1 1 auto">
    <LinearProgress
      sx={{
        width: 1,
        maxWidth: 320,
        bgcolor: (theme) => varAlpha(theme.vars.palette.text.primaryChannel, 0.16),
        [`& .${linearProgressClasses.bar}`]: { bgcolor: 'text.primary' },
      }}
    />
  </Box>
);

export function Router() {
    const [searchParams] = useSearchParams();
    const navigate = useNavigate();
    const apiCalledRef = useRef(false);
    const [toast, setToast] = useState<{
      open: boolean;
      type: AlertColor;
      message: string;
    }>({
      open: false,
      type: 'info',
      message: ''
    });
  
    const handleCloseToast = () => {
      setToast({ ...toast, open: false });
    };
  
    useEffect(() => {
      const token = searchParams.get('token');
      if (token && !apiCalledRef.current) {
        apiCalledRef.current = true;
        fetch(`${import.meta.env.VITE_BASE_URL_BACKEND}/auth/magic-link?token=${token}`)
          .then(async (response) => {
            const data = await response.json();
            
            if (response.ok) {
              localStorage.setItem('token', data.data.access_token);
              navigate('/', { replace: true });
              return;
            }
            
            const messageType = data.message.danger ? 'error' 
              : data.message.warning ? 'warning' 
              : 'success' as AlertColor;
            
            const message = data.message.danger || data.message.warning || data.message.success;
            
            setToast({ 
              open: true, 
              type: messageType, 
              message 
            });
            
            setTimeout(() => {
              navigate('/sign-in', { replace: true });
            }, 3000);
          });
      }
    }, [searchParams, navigate]);
  

  

  return useRoutes([
    {
        element: (
          <>
            <Snackbar 
            open={toast.open}
            autoHideDuration={6000}
            onClose={handleCloseToast}
            anchorOrigin={{ vertical: 'top', horizontal: 'right' }}
          >
            <Alert 
              elevation={6} 
              variant="filled" 
              severity={toast.type}
              onClose={handleCloseToast}
            >
              {toast.message}
            </Alert>
          </Snackbar>

            <PrivateRoute>
              <DashboardLayout>
                <Suspense fallback={renderFallback}>
                  <Outlet />
                </Suspense>
              </DashboardLayout>
            </PrivateRoute>
          </>
        ),
      children: [
        { element: <HomePage />, index: true },
        { path: 'inventory', element: <InventoryPage /> },
        { path: 'recipe', element: <RecipePage /> },
        { path: 'products', element: <ProductsPage /> },
        { path: 'blog', element: <BlogPage /> },
      ],
    },
    {
      path: 'sign-in',
      element: (
        isAuthenticated() ? (
          <Navigate to="/" replace />
        ) : (
          <AuthLayout>
            <SignInPage />
          </AuthLayout>
        )
      ),
    },
    
    {
      path: '404',
      element: <Page404 />,
    },
    {
      path: '*',
      element: <Navigate to="/404" replace />,
    },
  ]);
}
