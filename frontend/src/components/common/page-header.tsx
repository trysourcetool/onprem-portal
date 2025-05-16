import { AlertCircle, Clock } from 'lucide-react';
import clsx from 'clsx';
import { Link } from '@tanstack/react-router';
import { Alert, AlertDescription, AlertTitle } from '../ui/alert';
import { useAuth } from '../provider/auth-provider';
import { useSubscription } from '../provider/subscription-provider';
import type { FC } from 'react';
import { cn } from '@/lib/utils';
import { checkAccountExpiredDays, checkTrialExpiredDays } from '@/lib/account';

export const PageHeader: FC<{
  label: string;
  description?: string;
  border?: boolean;
}> = ({ label, description, border = true }) => {
  const { account } = useAuth();
  const { subscription } = useSubscription();

  const { expiredDays, isTrialEnd } = checkTrialExpiredDays(subscription);

  return (
    <>
      <div
        className={cn('flex flex-col gap-2', border && 'border-b p-4 md:p-6')}
      >
        <h1 className="text-foreground text-3xl font-bold">{label}</h1>
        {description && (
          <p className="text-muted-foreground text-base font-normal">
            {description}
          </p>
        )}
      </div>
      {(subscription?.status === 'trial' ||
        (subscription?.status === 'canceled' && isTrialEnd)) && (
        <Alert
          variant={isTrialEnd ? 'destructive' : 'default'}
          className={clsx(
            'rounded-none border-none',
            isTrialEnd ? 'bg-destructive/10' : 'bg-muted',
          )}
        >
          {!isTrialEnd && <Clock className="h-4 w-4" />}
          {isTrialEnd && <AlertCircle className="h-4 w-4" />}
          <AlertTitle>
            {isTrialEnd
              ? 'WARNING: ACCOUNT DELETION'
              : 'You are currently on a free trial.'}
          </AlertTitle>
          <AlertDescription>
            {isTrialEnd ? (
              `Your free account will be permanently deleted after ${checkAccountExpiredDays(
                account,
              )} days - upgrade now to prevent deletion.`
            ) : (
              <span>
                Your trial will end in {expiredDays} days. To continue using all
                features after day 14, please upgrade.{' '}
                <Link className="font-bold" to={'/settings/billing'}>
                  Learn more about pricing here.
                </Link>
              </span>
            )}
          </AlertDescription>
        </Alert>
      )}
    </>
  );
};
