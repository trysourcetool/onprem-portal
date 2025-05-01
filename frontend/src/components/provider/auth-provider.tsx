import { createContext, useContext, useEffect, useRef, useState } from 'react';
import { useLocation, useNavigate } from '@tanstack/react-router';
import { useQuery } from '@tanstack/react-query';
import type { FC, ReactNode } from 'react';
import type { User } from '@/api/modules/users';
import { api } from '@/api';

type AuthContextType = {
  isAuthorized: boolean;
  account: User | null;
};

export const authContext = createContext<AuthContextType>({
  isAuthorized: false,
  account: null,
});

export const useAuth = () => {
  const { isAuthorized, account } = useContext(authContext);
  return { isAuthorized, account };
};

export const AuthProvider: FC<{ children: ReactNode }> = (props) => {
  const isChecking = useRef<boolean>(false);
  const initialLoading = useRef(false);
  const [isAuthorized, setIsAuthorized] = useState(false);
  const pathname = useLocation().pathname;
  const navigate = useNavigate();

  const checkComplete = () => {
    setTimeout(() => {
      isChecking.current = false;
    }, 0);
  };

  const { data: account } = useQuery({
    enabled: isAuthorized,
    queryKey: ['user'],
    queryFn: async () => {
      const res = await api.users.getMe();
      return res;
    },
  });

  const handleAuthorizedRoute = () => {
    if (pathname.startsWith('/login')) {
      navigate({ to: '/' });
      checkComplete();
      return;
    }
  };

  const handleUnauthorizedRoute = () => {
    if (
      pathname === '/' ||
      (pathname.startsWith('/users') &&
        !pathname.startsWith('/users/oauth/google'))
    ) {
      navigate({ to: '/login' });
      checkComplete();
      return;
    }
  };

  useEffect(() => {
    if (isChecking.current || isAuthorized) {
      return;
    }
    isChecking.current = true;

    if (account) {
      handleAuthorizedRoute();
    } else {
      handleUnauthorizedRoute();
    }
  }, [account, isAuthorized]);

  useEffect(() => {
    if (initialLoading.current) return;
    initialLoading.current = true;
    (async () => {
      const res = await api.auth.refreshToken();
      try {
        if (res.expiresAt) {
          setIsAuthorized(true);
        } else {
          setIsAuthorized(false);
        }
      } catch (error) {
        console.error(error);
        setIsAuthorized(false);
      } finally {
        initialLoading.current = false;
      }
    })();
  }, []);

  return (
    <authContext.Provider
      value={{
        isAuthorized,
        account: account?.user ?? null,
      }}
    >
      {props.children}
    </authContext.Provider>
  );
};
