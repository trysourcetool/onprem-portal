import { index, layout, rootRoute, route } from '@tanstack/virtual-file-routes';

export default rootRoute('root.tsx', [
  layout('default', 'layout-default.tsx', [
    index('index.tsx'),
    route('login', 'login/index.tsx'),
  ]),
]);
