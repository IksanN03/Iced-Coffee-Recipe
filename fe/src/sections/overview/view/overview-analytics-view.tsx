import { Typography } from '@mui/material';
import { DashboardContent } from 'src/layouts/dashboard';
import { getEmailFromToken } from 'src/utils/auth';

export function OverviewAnalyticsView() {
  const userEmail = getEmailFromToken();

  return (
    <DashboardContent maxWidth="xl">
      <Typography variant="h4" sx={{ mb: { xs: 3, md: 5 } }}>
        Welcome, {userEmail}!
      </Typography>
    </DashboardContent>
  );
}
