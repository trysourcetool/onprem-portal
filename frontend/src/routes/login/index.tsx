import { createFileRoute, useNavigate } from '@tanstack/react-router';
import { object, string } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import { useForm } from 'react-hook-form';
import { useState } from 'react';
import { toast } from 'sonner';
import type { z } from 'zod';
import { SocialButtonGoogle } from '@/components/common/social-button-google';
import { Button } from '@/components/ui/button';
import {
  Card,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormMessage,
} from '@/components/ui/form';
import { api } from '@/api';

export default function Login() {
  const navigate = useNavigate();
  const [isRequestMagicLinkWaiting, setIsRequestMagicLinkWaiting] =
    useState(false);
  const [isOauthGoogleAuthWaiting, setIsOauthGoogleAuthWaiting] =
    useState(false);
  const schema = object({
    email: string({
      required_error: 'Email is required',
    }).email('Email is invalid'),
  });

  type Schema = z.infer<typeof schema>;

  const form = useForm<Schema>({
    resolver: zodResolver(schema),
  });

  const onSubmit = form.handleSubmit(async (data) => {
    if (isOauthGoogleAuthWaiting || isRequestMagicLinkWaiting) {
      return;
    }
    setIsRequestMagicLinkWaiting(true);
    try {
      const result = await api.auth.requestMagicLink({
        data: { email: data.email },
      });

      navigate({ to: '/login/emailSent', search: { email: result.email } });
    } catch (error) {
      console.error(error);
      toast('Login failed - Please check your email');
    } finally {
      setIsRequestMagicLinkWaiting(false);
    }
  });

  const handleGoogleAuth = async () => {
    if (isOauthGoogleAuthWaiting || isRequestMagicLinkWaiting) {
      return;
    }
    setIsOauthGoogleAuthWaiting(true);
    try {
      const result = await api.auth.requestGoogleAuthLink();
      window.location.href = result.authUrl;
    } catch (error) {
      console.error(error);
      toast('Failed to retrieve Url - Please try again');
    } finally {
      setIsOauthGoogleAuthWaiting(false);
    }
  };

  return (
    <div className="flex flex-1 items-center justify-center">
      <Form {...form}>
        <Card className="flex w-full max-w-[384px] flex-col gap-6 p-6">
          <CardHeader>
            <CardTitle className="text-2xl">Welcome to Sourcetool.</CardTitle>
            <CardDescription>Log in to build your first app.</CardDescription>
          </CardHeader>

          <form onSubmit={onSubmit} className="flex flex-col gap-4">
            <SocialButtonGoogle
              onClick={handleGoogleAuth}
              label="Continue with Google"
            />

            <div className="relative flex items-center justify-center">
              <span className="text-foreground text-sm font-medium">or</span>
            </div>
            <FormField
              control={form.control}
              name="email"
              render={({ field }) => (
                <FormItem>
                  <FormControl>
                    <Input
                      id="email"
                      type="email"
                      placeholder="Enter your personal or work email"
                      required
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <Button type="submit" className="w-full">
              Continue with email
            </Button>
            <p className="text-muted-foreground text-center text-xs">
              By continuing, you agree to Sourcetool's Consumer Terms and Usage
              Policy, and acknowledge their Privacy Policy.
            </p>
          </form>
        </Card>
      </Form>
    </div>
  );
}

export const Route = createFileRoute('/_default/login')({
  component: Login,
});
