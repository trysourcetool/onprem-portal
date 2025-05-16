import { createFileRoute } from '@tanstack/react-router';
import { Copy } from 'lucide-react';
import { toast } from 'sonner';
import { useEffect } from 'react';
import { useAuth } from '@/components/provider/auth-provider';
import { Button } from '@/components/ui/button';
import { PageHeader } from '@/components/common/page-header';
import { useBreadcrumbs } from '@/hooks/use-breadcrumbs';
import { checkTrialExpiredDays } from '@/lib/account';
import { useSubscription } from '@/components/provider/subscription-provider';

export default function Index() {
  const { account } = useAuth();
  const { setBreadcrumbsState } = useBreadcrumbs();
  const { subscription } = useSubscription();
  const onCopy = async (value: string) => {
    try {
      await navigator.clipboard.writeText(value);
      toast('License Key copied to clipboard');
    } catch (error) {
      console.error(error);
      toast('Failed to copy license key');
    }
  };

  useEffect(() => {
    setBreadcrumbsState?.([{ label: 'License Key', to: '/' }]);
  }, [setBreadcrumbsState]);

  return (
    <div>
      <PageHeader label="License Key" />
      <div className="flex w-screen flex-col gap-4 px-4 py-6 md:w-auto md:gap-6 md:px-6">
        <div className="bg-muted flex flex-col gap-4 px-6 py-4">
          <div className="flex flex-col gap-1">
            <p className="text-lg font-bold">License Key</p>
            {subscription?.status === 'trial' && (
              <p className="text-muted-foreground text-sm">
                Trial license:{' '}
                {account && checkTrialExpiredDays(subscription).expiredDays}/14
                days remaining
              </p>
            )}
          </div>
          <div className="flex gap-3">
            <code className="bg-input flex-1 rounded-md px-3 py-2.5 text-sm">
              {account?.license?.key}
            </code>
            <Button
              variant="default"
              size="icon"
              onClick={() => onCopy(account?.license?.key ?? '')}
            >
              <Copy className="size-4" />
            </Button>
          </div>

          <p className="text-sm">
            For license key setup instructions, see{' '}
            <a
              href="https://docs.trysourcetool.com/docs/getting-started/deployment"
              target="_blank"
              className="font-bold"
            >
              our documentation
            </a>
          </p>
        </div>
      </div>
    </div>
  );
}

export const Route = createFileRoute('/_default/')({
  component: Index,
});
