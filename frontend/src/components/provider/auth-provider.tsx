import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useRef,
  useState,
} from 'react';
import { useLocation, useNavigate } from '@tanstack/react-router';
import { useMutation, useQuery } from '@tanstack/react-query';
import { toast } from 'sonner';
import type { FC, ReactNode } from 'react';
import type { User } from '@/api/modules/users';
import { api } from '@/api';

type AuthContextType = {
  isAuthorized: boolean;
  account: User | null;
  handleLogout: () => void;
  handleAuthorized: () => void;
};

export const authContext = createContext<AuthContextType>({
  isAuthorized: false,
  account: null,
  handleLogout: () => {},
  handleAuthorized: () => {},
});

export const useAuth = () => {
  const { isAuthorized, account, handleLogout, handleAuthorized } =
    useContext(authContext);
  return { isAuthorized, account, handleLogout, handleAuthorized };
};

export const AuthProvider: FC<{ children: ReactNode }> = (props) => {
  const isChecking = useRef<boolean>(false);
  const initialLoading = useRef(false);
  const [isAuthChecked, setIsAuthChecked] = useState(false);
  const [isAuthorized, setIsAuthorized] = useState(false);
  const pathname = useLocation().pathname;
  const navigate = useNavigate();

  const checkComplete = useCallback(() => {
    setTimeout(() => {
      isChecking.current = false;
    }, 0);
  }, []);

  const { data: account, isLoading: isAccountLoading } = useQuery({
    enabled: isAuthorized,
    queryKey: ['user'],
    queryFn: async () => {
      const res = await api.users.getMe();
      return res;
    },
  });

  const logout = useMutation({
    mutationFn: async () => {
      await api.auth.logout();
    },
    onSuccess: () => {
      window.location.href = '/login';
    },
    onError: () => {
      toast('Failed to logout - Please try again');
    },
  });

  const handleAuthorized = () => {
    setIsAuthorized(true);
  };

  const handleAuthorizedRoute = useCallback(() => {
    if (pathname.startsWith('/login')) {
      navigate({ to: '/' });
      checkComplete();
      return;
    }
  }, [pathname, navigate, checkComplete]);

  const handleUnauthorizedRoute = useCallback(() => {
    if (
      pathname === '/' ||
      pathname === '/billing' ||
      (pathname.startsWith('/users') &&
        !pathname.startsWith('/users/oauth/google'))
    ) {
      navigate({ to: '/login' });
      checkComplete();
      return;
    }
  }, [pathname, navigate, checkComplete]);

  useEffect(() => {
    if (isChecking.current || !isAuthChecked || isAccountLoading) {
      return;
    }
    isChecking.current = true;

    if (account) {
      handleAuthorizedRoute();
    } else {
      handleUnauthorizedRoute();
    }
  }, [
    account,
    isAuthChecked,
    handleAuthorizedRoute,
    handleUnauthorizedRoute,
    isAccountLoading,
  ]);

  useEffect(() => {
    if (initialLoading.current) return;
    initialLoading.current = true;
    (async () => {
      try {
        const res = await api.auth.refreshToken();
        if (res.expiresAt) {
          setIsAuthorized(true);
        } else {
          setIsAuthorized(false);
        }
      } catch (error) {
        console.error(error);
        setIsAuthorized(false);
      } finally {
        setIsAuthChecked(true);
        initialLoading.current = false;
      }
    })();
  }, []);

  return (
    <authContext.Provider
      value={{
        isAuthorized,
        account: account?.user ?? null,
        handleLogout: logout.mutate,
        handleAuthorized,
      }}
    >
      {isAuthChecked && !isAccountLoading ? props.children : null}
    </authContext.Provider>
  );
};
