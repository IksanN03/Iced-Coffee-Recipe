import { Label } from 'src/components/label';
import { SvgColor } from 'src/components/svg-color';

// ----------------------------------------------------------------------

const icon = (name: string) => (
  <SvgColor width="100%" height="100%" src={`/assets/icons/navbar/${name}.svg`} />
);

export const navData = [
  {
    title: 'Dashboard',
    path: '/',
    icon: icon('ic-analytics'),
  },
  {
    title: 'Inventory',
    path: '/inventory',
    icon: icon('ic-user'),
  },
  {
    title: 'Recipe',
    path: '/recipe',
    icon: icon('ic-user'),
  },


];
