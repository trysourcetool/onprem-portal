import { useEffect, useRef } from 'react';
import {
  createFileRoute,
  useNavigate,
  useSearch,
} from '@tanstack/react-router';
import { Loader2 } from 'lucide-react';
import { zodValidator } from '@tanstack/zod-adapter';
import { object, string } from 'zod';
import { useMutation } from '@tanstack/react-query';
import { toast } from 'sonner';
import { api } from '@/api';

export default function MagicLinkAuth() {
  const isInitialLoading = useRef(false);
  const search = useSearch({ from: '/_default/auth/magic/authenticate/' });
  const token = search.token;
  const navigate = useNavigate();

  const mutation = useMutation({
    mutationFn: async (params: { token: string }) => {
      const result = await api.auth.authenticateWithMagicLink({
        data: params,
      });

      return result;
    },
    onSuccess: (data) => {
      navigate({
        to: '/signup/followup',
        search: { token: data.registrationToken },
      });
    },
    onError: () => {
      toast('Failed to authenticate - Please try again');
      navigate({ to: '/login' });
    },
  });

  useEffect(() => {
    if (isInitialLoading.current) {
      return;
    }
    isInitialLoading.current = true;

    if (!token) {
      toast('Invalid token - Please try again');
      navigate({ to: '/login' });
      return;
    }

    mutation.mutate({ token });
  }, [navigate, mutation, token]);

  return (
    <div className="m-auto flex items-center justify-center">
      <Loader2 className="size-8 animate-spin" />
    </div>
  );
}

export const Route = createFileRoute('/_default/auth/magic/authenticate/')({
  component: MagicLinkAuth,
  validateSearch: zodValidator(
    object({
      token: string(),
    }),
  ),
});
