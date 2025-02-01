import { useState, useCallback, useEffect } from 'react';
import Box from '@mui/material/Box';
import Link from '@mui/material/Link';
import TextField from '@mui/material/TextField';
import Typography from '@mui/material/Typography';
import LoadingButton from '@mui/lab/LoadingButton';
import Alert from '@mui/material/Alert';
import { useRouter } from 'src/routes/hooks';
import { useSearchParams } from 'react-router-dom';

export function SignInView() {
  const router = useRouter();
const [searchParams] = useSearchParams();
  const [email, setEmail] = useState('');
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState(false);
  const [error, setError] = useState('');

  // Handle magic link token on component mount
  useEffect(() => {
    const token = searchParams.get('token');
    if (token) {
      // Store token in localStorage
      localStorage.setItem('authToken', token);
      // Redirect to protected page
      router.push('/');
    }
  }, [searchParams, router]);

  const handleSubmitEmail = async () => {
    // Add validation check
    if (!email) {
      setError('Email is required');
      return;
    }

    try {
      setLoading(true);
      setError('');
      
      const response = await fetch(`${import.meta.env.VITE_BASE_URL_BACKEND}/auth/submit-email`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email }),
      });

      const data = await response.json();

      if (response.ok) {
        setSuccess(true);
      } else {
        setError(data.message || 'Something went wrong');
      }
    } catch (err) {
      setError('Failed to submit email');
    } finally {
      setLoading(false);
    }
  };


  const renderForm = (
    <Box display="flex" flexDirection="column" alignItems="flex-end">
      {error && (
        <Alert severity="error" sx={{ mb: 2, width: '100%' }}>
          {error}
        </Alert>
      )}
      
      {success ? (
        <Alert severity="success" sx={{ mb: 2, width: '100%' }}>
          Magic link has been sent to your email. Please check your inbox.
        </Alert>
      ) : (
        <>
                <TextField
          fullWidth
          required
          name="email"
          type="email"
          label="Email address"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          InputLabelProps={{ shrink: true }}
          sx={{ mb: 3 }}
        />


          <LoadingButton
            fullWidth
            size="large"
            type="submit"
            color="inherit"
            variant="contained"
            loading={loading}
            onClick={handleSubmitEmail}
          >
            Send Magic Link
          </LoadingButton>
        </>
      )}
    </Box>
  );

  return (
    <>
      <Box gap={1.5} display="flex" flexDirection="column" alignItems="center" sx={{ mb: 5 }}>
        <Typography variant="h5">Sign in with Magic Link</Typography>
        <Typography variant="body2" color="text.secondary">
          Enter your email to receive a magic link
        </Typography>
      </Box>

      {renderForm}
    </>
  );
}
