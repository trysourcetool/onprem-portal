import {
  createFileRoute,
  useNavigate,
  useSearch,
} from '@tanstack/react-router';
import { ArrowLeft, Mail } from 'lucide-react';
import { zodValidator } from '@tanstack/zod-adapter';
import { object, string } from 'zod';
import { toast } from 'sonner';
import { useMutation } from '@tanstack/react-query';
import { Button } from '@/components/ui/button';
import {
  Card,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { api } from '@/api';

export default function EmailSent() {
  const search = useSearch({ from: '/_default/login/emailSent/' });
  const email = search.email;
  const navigate = useNavigate();

  const mutation = useMutation({
    mutationFn: async (params: { email: string }) => {
      await api.auth.requestMagicLink({
        data: params,
      });
    },
    onSuccess: () => {
      toast('Email sent - A new verification link has been sent to your email');
    },
    onError: () => {
      toast('Login failed - Please check your email');
    },
  });

  const handleResendEmail = () => {
    if (!email || mutation.isPending) {
      return;
    }

    mutation.mutate({ email });
  };

  return (
    <div className="m-auto flex w-full items-center justify-center">
      <Card className="flex w-full max-w-[384px] flex-col gap-4 p-6">
        <CardHeader className="space-y-6 p-0">
          <CardTitle className="text-foreground text-2xl font-semibold">
            Verify email address
          </CardTitle>
          <div className="border-border flex items-center gap-3 rounded-md border p-3">
            <Mail className="h-5 w-5" />
            <CardDescription className="text-muted-foreground flex-1 text-sm">
              To continue, click the link sent to{' '}
              <span className="text-foreground font-medium">{email}</span>
            </CardDescription>
          </div>
        </CardHeader>

        <p className="text-muted-foreground text-center text-xs font-normal">
          Not seeing the email in your inbox?{' '}
          <button
            type="button"
            onClick={handleResendEmail}
            disabled={mutation.isPending}
            className="cursor-pointer underline"
          >
            Try sending again
          </button>
        </p>
      </Card>

      <div className="fixed bottom-8">
        <Button
          variant="secondary"
          className="cursor-pointer"
          onClick={() => navigate({ to: '/login' })}
        >
          <ArrowLeft className="h-4 w-4" />
          Go back
        </Button>
      </div>
    </div>
  );
}

export const Route = createFileRoute('/_default/login/emailSent/')({
  component: EmailSent,
  validateSearch: zodValidator(
    object({
      email: string(),
    }),
  ),
});
