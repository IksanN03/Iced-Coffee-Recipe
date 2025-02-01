import { Helmet } from 'react-helmet-async';

import { CONFIG } from 'src/config-global';

import { RecipeView } from 'src/sections/recipe/view';

// ----------------------------------------------------------------------

export default function Page() {
  return (
    <>
      <Helmet>
        <title> {`Recipe - ${CONFIG.appName}`}</title>
      </Helmet>

      <RecipeView />
    </>
  );
}
