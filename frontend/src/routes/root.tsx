import {
  Link,
  Outlet,
  createRootRoute,
  useNavigate,
} from '@tanstack/react-router';
import { ThemeProvider } from 'next-themes';
import { TanStackRouterDevtools } from '@tanstack/react-router-devtools';
import type { ErrorComponentProps } from '@tanstack/react-router';
import { CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { AuthProvider } from '@/components/provider/auth-provider';

function Fallback(props: ErrorComponentProps) {
  const navigate = useNavigate();
  return (
    <div className="m-auto flex items-center justify-center">
      <div className="flex max-w-[374px] flex-col gap-6 p-6">
        <CardHeader className="p-0">
          <CardTitle>Error: {props.error.name}</CardTitle>
          <CardDescription>{props.error.message}</CardDescription>
        </CardHeader>
        <Button
          onClick={() => {
            props.reset();
            navigate({ to: '/' });
          }}
        >
          Back to home
        </Button>
      </div>
    </div>
  );
}

export default function App() {
  return (
    <>
      <ThemeProvider enableSystem={false} attribute="class">
        <AuthProvider>
          <Outlet />
          <TanStackRouterDevtools position="bottom-right" />
        </AuthProvider>
      </ThemeProvider>
    </>
  );
}

export const Route = createRootRoute({
  component: App,
  errorComponent: Fallback,
  notFoundComponent: () => {
    return (
      <div className="m-auto flex items-center justify-center">
        <div className="flex max-w-[374px] flex-col gap-6 p-6">
          <CardHeader className="p-0">
            <CardTitle>Page not found</CardTitle>
            <CardDescription>
              The page you are looking for does not exist.
            </CardDescription>
          </CardHeader>
          <Button asChild>
            <Link to="/">Back to home</Link>
          </Button>
        </div>
      </div>
    );
  },
});
