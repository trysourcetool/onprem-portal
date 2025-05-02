import {
  index,
  layout,
  physical,
  rootRoute,
} from '@tanstack/virtual-file-routes';

export default rootRoute('root.tsx', [
  layout('default', 'layout-default.tsx', [
    index('index.tsx'),
    physical('/login', 'login'),
    physical('/auth/magic/authenticate', 'auth/magic/authenticate'),
    physical('/signup/followup', 'signup/followup'),
  ]),
]);
