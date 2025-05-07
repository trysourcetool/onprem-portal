import { useEffect, useRef } from 'react';
import {
  createFileRoute,
  useNavigate,
  useSearch,
} from '@tanstack/react-router';
import { zodValidator } from '@tanstack/zod-adapter';
import { object, string } from 'zod';
import { toast } from 'sonner';
import { Loader2 } from 'lucide-react';
import { useMutation } from '@tanstack/react-query';
import { api } from '@/api';
import { useAuth } from '@/components/provider/auth-provider';

export default function GoogleAuthenticate() {
  const isProcessing = useRef(false);
  const search = useSearch({
    from: '/_default/auth/google/callback/',
  });
  const code = search.code;
  const state = search.state;
  const navigate = useNavigate();
  const { handleAuthorized } = useAuth();

  const mutation = useMutation({
    mutationFn: async (params: { code: string; state: string }) => {
      const result = await api.auth.authenticateWithGoogle({
        data: params,
      });

      return result;
    },
    onSuccess: (data) => {
      if (data.isNewUser) {
        navigate({
          to: '/signup/followup',
          search: { token: data.registrationToken },
        });
      } else {
        handleAuthorized();
        navigate({ to: '/' });
      }
    },
    onError: () => {
      toast('Failed to authenticate - Please try again');
      navigate({ to: '/login' });
    },
  });

  useEffect(() => {
    if (isProcessing.current) {
      return;
    }
    isProcessing.current = true;

    if (!code || !state) {
      toast('Invalid token - Please try again');
      navigate({ to: '/login' });
      return;
    }

    mutation.mutate({ code, state });
  }, [navigate, mutation, code, state]);

  return (
    // Fixed escaped quotes
    <div className="m-auto flex items-center justify-center p-10">
      <Loader2 className="size-8 animate-spin" />
    </div>
  );
}

export const Route = createFileRoute('/_default/auth/google/callback/')({
  component: GoogleAuthenticate,
  validateSearch: zodValidator(
    object({
      code: string(),
      state: string(),
    }),
  ),
});
