import { useForm } from 'react-hook-form';
import { object, string } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import {
  createFileRoute,
  useNavigate,
  useSearch,
} from '@tanstack/react-router';
import { zodValidator } from '@tanstack/zod-adapter';
import { toast } from 'sonner';
import { useMutation } from '@tanstack/react-query';
import type { z } from 'zod';
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form';
import { Button } from '@/components/ui/button';
import { Card, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { api } from '@/api';
import { useAuth } from '@/components/provider/auth-provider';

export default function Followup() {
  const search = useSearch({
    from: '/_default/signup/followup/',
  });
  const token = search.token;
  const navigate = useNavigate();
  const { handleAuthorized } = useAuth();

  const mutation = useMutation({
    mutationFn: async (params: {
      token: string;
      firstName: string;
      lastName: string;
    }) => {
      const result = await api.auth.registerWithMagicLink({
        data: params,
      });

      return result;
    },
    onSuccess: () => {
      handleAuthorized();
      navigate({ to: '/' });
      toast('Signup success - Next, create an organization');
    },
    onError: () => {
      toast('Failed to register - Please try again');
    },
  });

  const schema = object({
    firstName: string({
      required_error: 'First name is required',
    }),
    lastName: string({
      required_error: 'Last name is required',
    }),
  });

  type Schema = z.infer<typeof schema>;

  const form = useForm<Schema>({
    resolver: zodResolver(schema),
  });

  const onSubmit = form.handleSubmit((data) => {
    if (!token) {
      toast('Invalid token - Please try again');
      return;
    }
    mutation.mutate({
      token,
      ...data,
    });
  });

  return (
    <div className="m-auto flex w-full items-center justify-center">
      <Form {...form}>
        <Card className="flex w-full max-w-sm flex-col gap-6 p-6">
          <CardHeader className="p-0">
            <CardTitle>What's your name?</CardTitle>
          </CardHeader>
          <form onSubmit={onSubmit} className="flex flex-col gap-4">
            <div className="flex items-start gap-3">
              <FormField
                control={form.control}
                name="firstName"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>First Name</FormLabel>
                    <FormControl>
                      <Input placeholder="Rachel" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="lastName"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Last Name</FormLabel>
                    <FormControl>
                      <Input placeholder="Bennett" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>

            <Button type="submit" className="w-full cursor-pointer">
              Continue
            </Button>
          </form>
        </Card>
      </Form>
    </div>
  );
}

export const Route = createFileRoute('/_default/signup/followup/')({
  component: Followup,
  validateSearch: zodValidator(
    object({
      token: string(),
    }),
  ),
});
